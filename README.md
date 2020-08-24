# mqueue
高性能http队列 每秒TPS达7万以上 （分时落盘，防丢失）

* wrk测试
```
wrk -t 16 -c 100 -d 30s --latency --timeout 5s -s post.lua http://localhost:8088/api/queue
```
* 测试结果
```
Running 30s test @ http://localhost:8088/api/queue
  16 threads and 100 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency     1.39ms    1.60ms  45.07ms   99.29%
    Req/Sec     4.61k   612.62    30.27k    95.13%
  Latency Distribution
     50%    1.30ms
     75%    1.40ms
     90%    1.51ms
     99%    2.59ms
  2204626 requests in 30.10s, 681.21MB read
Requests/sec:  73235.09
Transfer/sec:     22.63MB
```

* 支持自动TTL删除功能，通过Query参数指定ttl，值为数值型/秒，到达指定时间后记录会自动被清除。
```
http://127.0.0.1:8088/api/queue?ttl=10
```

* lvs 安装
```
$sudo apt-get install ipvsadm
$sudo apt-get install keepalived
$sudo vim /etc/keepalived/keepalived.conf //参考bin目录
$sudo ifconfig eth0:0 192.168.101.210 netmask 255.255.255.0 broadcast 192.168.101.210
$sudo route add -host 192.168.101.210 dev eth0:0
$sudo echo "1" > /proc/sys/net/ipv4/ip_forward
$sudo ipvsadm -C
$sudo ipvsadm -A -t 192.168.101.210:8088 -s rr
$sudo ipvsadm -a -t 192.168.101.210:8088 -r 192.168.101.68:8088 -g
$sudo ipvsadm -a -t 192.168.101.210:8088 -r 192.168.101.69:8088 -g
$sudo ipvsadm
$sudo sysctl -p
```

* lvs 节点配置
```
$sudo ifconfig lo:0 192.168.101.210 broadcast 192.168.101.210 netmask 255.255.255.255 up
$sudo route add -host 192.168.101.210 dev lo:0
$sudo echo "1"> /proc/sys/net/ipv4/conf/lo/arp_ignore
$sudo echo "2"> /proc/sys/net/ipv4/conf/lo/arp_announce
$sudo echo "1" > /proc/sys/net/ipv4/conf/all/arp_ignore
$sudo echo "2" > /proc/sys/net/ipv4/conf/all/arp_announce
$sudo sysctl -p
```

* lvs 双节点负载均衡压测数据
``` 
Running 30s test @ http://192.168.101.210:8088/api/db
  16 threads and 1000 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency     8.47ms    4.10ms 217.94ms   77.87%
    Req/Sec     7.00k     1.11k   14.38k    71.97%
  Latency Distribution
     50%    7.83ms
     75%   10.27ms
     90%   13.27ms
     99%   21.82ms
  3352874 requests in 30.08s, 1.01GB read
Requests/sec: 111452.24
Transfer/sec:     34.44MB
```

