
# aryadb

[![Build Status](https://travis-ci.com/gushitong/aryadb.svg?branch=master)](https://travis-ci.com/gushitong/aryadb)

aryadb is a high performance no-sql database build on [BadgerDB](https://github.com/dgraph-io/badger) with redis protocol
support. It meant to provide a key-value store alternative to redis.

## Install

    go get github.com/gushitong/aryadb
    
## Redis Command Support

* String Command

```bash
append decr decrby get getbit getrange getset incr incrby mget mset msetnx set setbit setex setnx setrange strlen
```

* Hash Command

```bash
hdel hexists hget hgetall hincrby hincrbyfloat hkeys hlen hmget hmset hscan hset hsetnx hstrlen hvals
```

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