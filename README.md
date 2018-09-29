
# aryadb

[![Build Status](https://travis-ci.com/gushitong/aryadb.svg?branch=master)](https://travis-ci.com/gushitong/aryadb)
[![NoSQL](https://img.shields.io/badge/db-NoSQL-blue.svg)](https://github.com/gushitong/aryadb)
[![Go Report Card](https://goreportcard.com/badge/github.com/gushitong/aryadb?service=github)](https://goreportcard.com/report/github.com/gushitong/aryadb)
[![License](https://img.shields.io/badge/License-Apache-green.svg)]((https://github.com/gushitong/aryadb))

aryadb is a high performance no-sql database build on [BadgerDB](https://github.com/dgraph-io/badger) with redis protocol
support. It meant to provide a key-value store alternative to redis.

* Pure `golang` implementation, no `c/c++` dependency
* Compatible with redis protocol, Redis client are supported
* Persistent all data to the disk

## Example
```bash
go get github.com/gushitong/aryadb
```
    
```bash
$ aryadb 
2018/09/29 12:05:53 started server at :6380    
```    

```bash
$ redis-cli -p 6380
127.0.0.1:6380> SET k 1
OK
127.0.0.1:6380> GET k
"1"
127.0.0.1:6380> HSET hash k v
(integer) 1
127.0.0.1:6380> HGET hash k
"v"
127.0.0.1:6380> PING
PONG
```    
    
## Redis Command Support

|  Strings   | Lists    | Hashes    | Sets      | Sorted Sets   |
|:----------:|:--------:|:---------:|:---------:|:-------------:|
| `append`   | `lindex` | `hdel`    | `sadd`    | `zadd`        |
| `decr`     | `llen`   | `hexists` | `scard`   | `zcard`       |
| `decrby`   | `lpop`   | `hget`    | `sismember`| `zcount`     |
| `get`      | `lpush`  | `hgetall` | `smembers` | `zincrby`    |
| `getbit`   | `lpushx` | `hincrby` | `spop`    | `zpopmax`     |
| `getrange` | `lrange` | `hincrbyfloat` |      | `zpopmin`     |
| `getset`   | `lset`   | `hkeys`   |           | `zrange`      |
| `incr`     |          | `hlen`    |           | `zrangebyscore`|   
| `incrby`   |          | `hmget`   |           | `zrank`       |
| `incrfloat`|          | `hmset`   |           | `zrevrange`   |
| `mget`     |          | `hscan`   |           | `zrevrangebyscore`|
| `mset`     |          | `hset`    |           | `zrevrank`    |
| `msetnx`   |          | `hsetnx`  |           | `zscore`      |  
| `set`      |          | `hstrlen` |           |               |
| `setbit`   |          | `hvals`   |           |               |
| `setex`    |          |           |           |               |
| `setnx`    |          |           |           |               |
| `setrange` |          |           |           |               |
| `strlen`   |          |           |           |               |

## Benchmark

This benchmark running on my local mac, aryadb has better performance on SSD.

* redis benchmark:

```bash
    $ redis-benchmark -p 6379 -c 50 -n 10000 -q 
    PING_INLINE: 56497.18 requests per second
    PING_BULK: 56818.18 requests per second
    SET: 58823.53 requests per second
    GET: 57142.86 requests per second
    INCR: 59171.60 requests per second
    LPUSH: 59523.81 requests per second
    RPUSH: 57142.86 requests per second
    LPOP: 59880.24 requests per second
    RPOP: 52631.58 requests per second
    SADD: 61728.39 requests per second
    SPOP: 56497.18 requests per second
    LPUSH (needed to benchmark LRANGE): 59880.24 requests per second
    LRANGE_100 (first 100 elements): 19047.62 requests per second
    LRANGE_300 (first 300 elements): 9074.41 requests per second
    LRANGE_500 (first 450 elements): 6501.95 requests per second
    LRANGE_600 (first 600 elements): 4933.40 requests per second
    MSET (10 keys): 45871.56 requests per second
```

* aryadb benchmark:

```bash
    $ redis-benchmark -p 6380 -c 50 -n 10000 -q 
    PING_INLINE: 45662.10 requests per second
    PING_BULK: 48543.69 requests per second
    SET: 26881.72 requests per second
    GET: 32051.28 requests per second
    INCR: 30769.23 requests per second
    LPUSH: 24937.66 requests per second
    RPUSH: 38461.54 requests per second
    LPOP: 24691.36 requests per second
    RPOP: 39215.69 requests per second
    SADD: 36496.35 requests per second
    SPOP: 30120.48 requests per second
    LPUSH (needed to benchmark LRANGE): 26954.18 requests per second
    LRANGE_100 (first 100 elements): 8960.57 requests per second
    LRANGE_300 (first 300 elements): 8045.05 requests per second
    LRANGE_500 (first 450 elements): 8312.55 requests per second
    LRANGE_600 (first 600 elements): 9049.77 requests per second
    MSET (10 keys): 33670.04 requests per second
```