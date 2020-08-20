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

