# WireGuard Exporter

## WireGuard exporter is a Prometheus exporter for WireGuard VPN. 
It collects data from the wg show command and exports it in the Prometheus format, which can be scraped by Prometheus server to collect the metrics and store them in the time-series database.
    Flags:
|flag|description|Default|
|----------------|----------------|----------------|
|`-a`| URL path for surfacing collected metrics.| Default is /metrics.|
|`-p`| Address for WireGuard exporter.| Default is :9586.|
|`-c`| Path to file with name:key value.| Default is /etc/wireguard/configs/clients.txt.|
## Metrics

This exporter exposes the following metrics:
|metrics|description|
|-------------------------|-------------------------|
|wireguard_bytes_received:| Total number of bytes received.|
|wireguard_bytes_sent:| Total number of bytes sent.|
|wireguard_last_handshake:| UNIX timestamp seconds of the last handshake.|
|wireguard_counter_config:| Configuration counter.|

## Installation

### Dependencies

Before installing WireGuard Exporter, make sure the following dependencies are installed on the system.

- Go programming language
- Prometheus client library for Go
- WireGuard VPN

To use this exporter, follow these steps:

Clone the repository: `git clone https://github.com/1kurops/wireguard-go-exporter.git`

Navigate to the directory: `cd wireguard-exporter`

Compile the binary: `go build -o wireguard-exporter`

Run the binary: `./wireguard-exporter`

### Contributing

If you find a bug or have an idea for a new feature, please open an issue or submit a pull request. We welcome contributions from the community!