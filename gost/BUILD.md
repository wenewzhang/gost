Gost with tls
```
./gost -L tls://199.15.77.130:443
iptables -t nat -I SS-TCP -d 199.15.77.130  -j RETURN  
/usr/bin/gost-arm7-client -D  -L red://192.168.2.1:10060 -F=tls://199.15.77.130:443

iptables -t nat -I SS-TCP -d 199.15.77.130  -j RETURN  
./gost -D  -L red://192.168.2.1:10060 -F=tls://199.15.77.130:443


./gost-arm7-client -D  -L red://192.168.3.1:10060 -F=tls://199.15.77.130:443



./gost -L tls://173.82.56.214:445

iptables -t nat -I SS-TCP -d 173.82.56.214  -j RETURN  
/usr/bin/gost-arm7-client -D  -L red://192.168.3.1:10060 -F=tls://173.82.56.214:445


```

Gost with tls & auth
```
./gost -L tls://user:pass@173.82.56.214:8080

iptables -t nat -I SS-TCP -d 173.82.56.214  -j RETURN  
/usr/bin/gost-arm7-client -D  -L red://192.168.3.1:10060 -F=tls://user:pass@173.82.56.214:8080

```
```


Show CPU info
```
cat /proc/cpuinfo
```

Python
```
python -m http.server --bind 192.168.2.108 --directory  ./
wget http://192.168.2.108:8000/gost-arm7-client
chmod +x gost-arm7-client

wget http://192.168.3.108:8000/run.sh
chmod +x run.sh

```

Debug go app with args
```
 /home/jimmy/go/bin/dlv debug --  -L red://192.168.3.1:10060 -F=tls://199.15.77.130:443
```

Upgrade Go modules for go.mod

List all 
```
go list -u -m all
```

Upgrade go.mod
```
go get -u ./...
```

Upgrade a package
```
get foo@v1.6.2
go get foo@e3702bed2

```


Check real memory
```
ps -e -o 'pid,comm,args,pcpu,rsz,vsz,stime,user,uid' |grep gost
``` 