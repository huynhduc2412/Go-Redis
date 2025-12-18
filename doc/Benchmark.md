### Benchmark commands
```
 src/redis-benchmark -n 10000 -t ping_mbulk -c 200 -h localhost -p 3000
```

### Benchmark result

TCP Server with threads

```
Summary:
  throughput summary: 13262.60 requests per second
  latency summary (msec):
          avg       min       p50       p95       p99       max
       14.479     0.016     1.687   156.287   167.167   170.239

```

TCP Server with one threads handle on by using EPOLL/KQUEUE to monitoring

```
Summary:
  throughput summary: 129870.13 requests per second
  latency summary (msec):
          avg       min       p50       p95       p99       max
        0.777     0.176     0.719     1.207     1.647     2.015
```



