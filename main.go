package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
	namespace = "wireguard"
)

var (
	addr    = flag.String("a", "/metrics", "URL path for surfacing collected metrics")
	port    = flag.String("p", ":9586", "address for WireGuard exporter")
	clients = flag.String("c", "/etc/wireguard/configs/clients.txt", "Path to file with name:key value")
)

type collector struct {
	bytesReceived *prometheus.Desc
	bytesSent     *prometheus.Desc
	lasthandshake *prometheus.Desc
	counterconfig *prometheus.Desc
}

func newCollector() *collector {
	flag.Parse()
	return &collector{
		bytesReceived: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "bytes_received"),
			"Total number of bytes received.",
			[]string{"interface", "public_key", "name"},
			nil,
		),
		bytesSent: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "bytes_sent"),
			"Total number of bytes sent.",
			[]string{"interface", "public_key", "name"},
			nil,
		),
		lasthandshake: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "last_handshake"),
			"UNIX timestamp seconds of the last handshake",
			[]string{"interface", "public_key", "name"},
			nil,
		),
		counterconfig: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "counter_config"),
			"Configuration counter.",
			[]string{"interface"},
			nil,
		),
	}
}

func (c *collector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.bytesReceived
	ch <- c.bytesSent
	ch <- c.lasthandshake
	ch <- c.counterconfig
}

func (c *collector) Collect(ch chan<- prometheus.Metric) {
	flag.Parse()
	file, err := os.Open(*clients)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	keyToName := make(map[string]string)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)

		key := fields[1]
		user_name := fields[0]

		keyToName[key] = user_name
	}
	if err := scanner.Err(); err != nil {
		panic(err)
	}

	cmd := exec.Command("wg", "show", "all", "dump")
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("Error running command: %v\n", err)
		return
	}

	dump := strings.Split(string(output), "\n")
	count := 0
	for _, line := range dump[1:] {
		if line == "" {
			continue
		}
		count++
		fields := strings.Fields(line)
		interfaceName := fields[0]
		publicKey := fields[1]
		lasthandshake, _ := strconv.ParseFloat(fields[5], 64)
		user_name, ok := keyToName[publicKey]
		if !ok {
			fmt.Println("error user name", publicKey)
			return
		}
		bytesReceived, _ := strconv.ParseFloat(fields[6], 64)
		bytesSent, _ := strconv.ParseFloat(fields[7], 64)

		ch <- prometheus.MustNewConstMetric(
			c.bytesReceived,
			prometheus.CounterValue,
			bytesReceived,
			interfaceName, publicKey, user_name,
		)

		ch <- prometheus.MustNewConstMetric(
			c.bytesSent,
			prometheus.CounterValue,
			bytesSent,
			interfaceName, publicKey, user_name,
		)

		ch <- prometheus.MustNewConstMetric(
			c.lasthandshake,
			prometheus.GaugeValue,
			float64(lasthandshake),
			interfaceName, publicKey, user_name,
		)
	}
	ch <- prometheus.MustNewConstMetric(
		c.counterconfig,
		prometheus.GaugeValue,
		float64(count),
		"ru_wg0",
	)
}

func main() {
	flag.Parse()
	collector := newCollector()
	prometheus.MustRegister(collector)

	endpoint := http.NewServeMux()
	endpoint.Handle(*addr, promhttp.Handler())

	log.Printf("starting WireGuard exporter on %q", *port, *addr)
	log.Printf("clients path is :", *clients)
	s := &http.Server{
		Addr:         *port,
		Handler:      endpoint,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
	}
	log.Fatal(s.ListenAndServe())
}
