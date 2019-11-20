#!/bin/bash

sudo apt-get update
sudo apt-get install -y hostapd dnsmasq

sudo systemctl stop hostapd
sudo systemctl stop dnsmasq


config="interface wlan0
static ip_address=192.168.0.10/24
denyinterfaces eth0
denyinterfaces wlan0"
echo "${config}" | sudo tee -a /etc/dhcpcd.conf


sudo mv /etc/dnsmasq.conf /etc/dnsmasq.conf.orig
dnsmasq="interface=wlan0
  dhcp-range=192.168.0.11,192.168.0.30,255.255.255.0,24h"
echo "${dnsmasq}" | sudo tee -a /etc/dnsmasq.conf



hostapd="interface=wlan0
bridge=br0
hw_mode=g
channel=7
wmm_enabled=0
macaddr_acl=0
auth_algs=1
ignore_broadcast_ssid=0
wpa=2
wpa_key_mgmt=WPA-PSK
wpa_pairwise=TKIP
rsn_pairwise=CCMP
ssid=NETWORK
wpa_passphrase=PASSWORD"

echo "${hostapd}" | sudo tee -a /etc/hostapd/hostapd.conf

echo 'DAEMON_CONF="/etc/hostapd/hostapd.conf"' | sudo tee -a /etc/default/hostapd


echo 'net.ipv4.ip_forward=1' | sudo tee -a /etc/sysctl.conf

sudo iptables -t nat -A POSTROUTING -o eth0 -j MASQUERADE
sudo sh -c "iptables-save > /etc/iptables.ipv4.nat"
iptables-restore < /etc/iptables.ipv4.nat
sudo sed -i 's/exit\ 0//g' /etc/rc.local
echo 'iptables-restore < /etc/iptables.ipv4.nat' | sudo tee -a /etc/rc.local
echo 'exit 0' | sudo tee -a /etc/rc.local

sudo systemctl unmask hostapd
sudo systemctl enable hostapd
sudo systemctl start hostapd
sudo systemctl start dnsmasq
