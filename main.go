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
	addr      = flag.String("a", "/metrics", "URL path for surfacing collected metrics")
	port      = flag.String("p", ":9586", "address for WireGuard exporter")
	config    = flag.String("c", "/etc/wireguard/wg0.conf", "Path to main file config")
	Interface = flag.String("i", "wg0", "Wireguard interface")
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
	file, err := os.Open(*config)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	keyToName := make(map[string]string)

	scanner := bufio.NewScanner(file)
	var currentBlock string
	inBlock := false

	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "### begin ") && strings.HasSuffix(line, " ###") {
			currentBlock = strings.TrimPrefix(line, "### begin ")
			currentBlock = strings.TrimSuffix(currentBlock, " ###")
			inBlock = true
		} else if inBlock && strings.HasPrefix(line, "### end ") && strings.HasSuffix(line, " ###") {
			inBlock = false
			currentBlock = ""
		} else if inBlock && strings.HasPrefix(line, "PublicKey = ") {
			publicKey := strings.TrimPrefix(line, "PublicKey = ")
			keyToName[publicKey] = currentBlock
		}
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
			lasthandshake,
			interfaceName, publicKey, user_name,
		)
	}
	ch <- prometheus.MustNewConstMetric(
		c.counterconfig,
		prometheus.GaugeValue,
		float64(count),
		*Interface,
	)
}

func main() {
	flag.Parse()
	collector := newCollector()
	prometheus.MustRegister(collector)

	endpoint := http.NewServeMux()
	endpoint.Handle(*addr, promhttp.Handler())

	log.Printf("starting WireGuard exporter on %q", *port, *addr)
	log.Printf("Config path is :", *config)
	log.Printf("Interface exporting is :", *Interface)
	s := &http.Server{
		Addr:         *port,
		Handler:      endpoint,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
	}
	log.Fatal(s.ListenAndServe())
}
