#!/bin/bash

GO_VERSION=$(curl -sSL https://golang.org/VERSION?m=text)
GO_URL="https://go.dev/dl/${GO_VERSION}.linux-amd64.tar.gz"
REPO_URL="https://github.com/1kurops/wireguard-go-exporter.git"

if [ -z "$1" ]; then
    echo "Usage: $0 <location>"
    exit 1
fi
locate=$1

if ! wget ${GO_URL}; then
    echo "Error: command failed"
    exit 1
fi

if ! tar -C /usr/local -xzf ${GO_VERSION}.linux-amd64.tar.gz; then
    echo "Error: command failed"
    exit 1
fi

echo 'export PATH=$PATH:/usr/local/go/bin' >>~/.bashrc
source ~/.bashrc
go version
mkdir /opt/goexporter/

git clone ${REPO_URL}
cd wireguard-go-exporter

sed -i "s|ru_wg0|${locate}_wg0|" main.go

go build -o /opt/goexporter/wireguard-go-exporter .
cd /opt/goexporter/
{
    echo "[Unit]"
    echo "Description=Prometheus WireGuard Go Exporter"
    echo "Wants=network-online.target"
    echo "After=network-online.target"
    echo
    echo "[Service]"
    echo "User=root"
    echo "Group=root"
    echo "Type=simple"
    echo "ExecStart=/opt/goexporter/wireguard-go-exporter -p=:$2 -c=/etc/${locate}_wireguard/configs/clients.txt"
    echo
    echo "[Install]"
    echo "WantedBy=multi-user.target"
} >wireguard-go-exporter.service

mv wireguard-go-exporter.service /etc/systemd/system/.
systemctl stop wg_exporter.service
systemctl disable wg_exporter.service
systemctl daemon-reload
systemctl start wireguard-go-exporter && systemctl enable wireguard-go-exporter

rm -r ~/wireguard-go-exporter
