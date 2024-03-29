# WireGuard Exporter

## WireGuard exporter is a Prometheus exporter for WireGuard VPN. 
It collects data from the wg show command and exports it in the Prometheus format, which can be scraped by Prometheus server to collect the metrics and store them in the time-series database.
    Flags:
|flag|description|Default|
|----------------|----------------|----------------|
|`-a`| URL path for surfacing collected metrics.| Default is /metrics.|
|`-p`| Address for WireGuard exporter.| Default is :9586.|
|`-c`| Path to file with name:key value.| Default is /etc/wireguard/wg0.conf.|
|`-i`| Wireguard interface | Default is wg0.|
## Metrics

This exporter exposes the following metrics:
|metrics|description|
|-------------------------|-------------------------|
|`wireguard_bytes_received:`| Total number of bytes received.|
|`wireguard_bytes_sent:`| Total number of bytes sent.|
|`wireguard_last_handshake:`| UNIX timestamp seconds of the last handshake.|
|`wireguard_counter_config:`| Configuration counter.|

## Installation

### Dependencies

Before installing WireGuard Exporter, make sure the following dependencies are installed on the system.

- Go programming language
- Prometheus client library for Go
- WireGuard VPN

To use this exporter, follow these steps:

1. Clone the repository: `git clone https://github.com/1kurops/wireguard-go-exporter.git`

2. Navigate to the directory: `cd wireguard-exporter`

3. Compile the binary: `go build -o wireguard-exporter`

4. Run the binary: `./wireguard-exporter`
