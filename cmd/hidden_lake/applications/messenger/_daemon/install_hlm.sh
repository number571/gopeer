#!/bin/bash

# root mode
echo "
[Unit]
Description=HiddenLakeMessenger

[Service]
ExecStart=/root/hlm_amd64_linux -path=/root -pasw=/root/pasw.key
Restart=always
RestartSec=10

[Install]
WantedBy=multi-user.target
" > /etc/systemd/system/hidden_lake_messenger.service

cd /root && \
    rm -f hlm_amd64_linux && \
    wget https://github.com/number571/go-peer/releases/latest/download/hlm_amd64_linux && \
    chmod +x hlm_amd64_linux

systemctl daemon-reload
systemctl enable hidden_lake_messenger.service
systemctl restart hidden_lake_messenger.service
