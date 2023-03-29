# WireGuard Exporter

## This is a WireGuard exporter for Prometheus. It exposes metrics about the amount of data transmitted through WireGuard interfaces.
Flags

    `-a`: URL path for surfacing collected metrics. Default is /metrics.

    `-p`: Address for WireGuard exporter. Default is :9586.

    `-c`: Path to file with name:key value. Default is /etc/wireguard/configs/clients.txt.
## Metrics

This exporter exposes the following metrics:

    ### wireguard_bytes_received: Total number of bytes received.
    ### wireguard_bytes_sent: Total number of bytes sent.
    ### wireguard_last_handshake: UNIX timestamp seconds of the last handshake.
    ### wireguard_counter_config: Configuration counter.
    ### 

## Installation

To use this exporter, you need to have a working installation of WireGuard and Go. Then, follow these steps:

    Clone the repository: git clone https://github.com/1kurops/wireguard-go-exporter.git
    Navigate to the directory: cd wireguard-exporter
    Compile the binary: go build -o wireguard-exporter
    Run the binary: ./wireguard-exporter
