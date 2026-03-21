![image.png](image.png)
# Built In-memory Database Bases On Real Redis.

Built From Scratch in Go.
It's inspired by Redis and implements the **RESP (Redis Serialization Protocol)** and making it fully compatible with redis-cli.

This project with purposes for reviewing , learning , exploring  high-performance Network Servers , Operating System , IO Models , Multithreaded Programming and knowing about the low-level design of modern databases.

## Key Features

- **Redis Protocol (RESP) Compliant**: Fully compatible with redis-cli and other Redis clients.

- **High-Performance I/O**: Uses a single-threaded with I/O multiplexing (epoll for **Linux** or kqueue for **macOS**) to handle thousands of concurrent connections.

- **Implement basically command like as redis**

- **Custom Data Structures**: Implements complex data structures from scratch, including:
  - Skip List: For high-performance sorted sets (ZADD, ZRANK, etc.).

- **Probabilistic Data Structures**: Includes implementations of:
    - **Scalable Bloom Filter**: For fast, memory-efficient set membership testing.

    - **Count-Min Sketch**: For estimating item frequencies in a data stream.

    - **Approx LRU eviction**: Redis frees up memory when it reaches its storage limit (MaxKeyNumber) in RAM.

    - **Random eviction**: Redis random and delete key when it reaches its storage limit (MaxKeyNumber) in RAM.

- **Graceful Shutdown**: Ensures data is handled correctly and connections are closed properly on server termination.

## Getting Started
```
  cd cmd
  go run main.go
  # on another terminal (using redis-clone version 6.)
  redis-cli -p 3000
```
  
The server supports a wide range of commands grouped by data type:

| Category | Commands |
| :--- | :--- |
| **General** | `PING` , `INFO` |
| **String** | `SET`, `GET`, `DEL`, `TTL`, `EXPIRE`, `INCR` |
| **Sorted Set**| `ZADD`, `ZRANK`, `ZREM`, `ZSCORE`, `ZCARD` |
| **Set** | `SADD`, `SREM`, `SCARD`, `SMEMBERS`, `SISMEMBER`, `SRAND`, `SPOP` |
| **Bloom Filter**| `BF.RESERVE`, `BF.INFO`, `BF.MADD`, `BF.EXISTS`, `BF.MEXISTS` |
| **Count-Min** | `CMS.INITBYDIM`, `CMS.INITBYPROB`, `CMS.INCRBY`, `CMS.QUERY` |

## Future Works And Optimize
[ ] Monitoring log using Prometheus and Grafana

[ ] Build up single-threaded to multithreaded

[ ] Trying event-loop architecture combined with I/O multiplexing 

[ ] Approx LFU eviction



