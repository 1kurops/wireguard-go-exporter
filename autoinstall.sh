#!/bin/bash

go mod init main.go
go mod tidy
go build -o wireguard-go-exporter .

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
    echo "ExecStart=/opt/wireguard-go-exporter/wireguard-go-exporter -c=/etc/wireguard/${1}_wireguard/${1}_wg0.conf -i ${1}_wg0"
    echo
    echo "[Install]"
    echo "WantedBy=multi-user.target"
} >/etc/systemd/system/wireguard-go-exporter.service

systemctl daemon-reload
systemctl start wireguard-go-exporter && systemctl enable wireguard-go-exporter
systemctl restart wireguard-go-exporter
