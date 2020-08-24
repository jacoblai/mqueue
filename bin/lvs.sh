ifconfig eth0:0 192.168.11.10 broadcast 192.168.11.10 netmask 255.255.255.255 up
route add -host 192.168.11.10 dev eth0:0
echo "1" > /proc/sys/net/ipv4/ip_forward
ipvsadm -C
ipvsadm -A -t 192.168.11.10:80 -s rr
ipvsadm -a -t 192.168.11.10:80 -r 192.168.11.21:80 -g
ipvsadm -a -t 192.168.11.10:80 -r 192.168.11.22:80 -g
ipvsadm
sysctl -p