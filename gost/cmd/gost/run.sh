#!/bin/sh
rm gost-$1-client
wget http://192.168.2.108:8000/gost-$1-client
chmod +x gost-$1-client
iptables -t nat -I SS-TCP -d 173.82.56.214  -j RETURN  
./gost-$1-client -D -L red://192.168.2.1:10060 -F=tls://user:pass@199.15.77.130:443


