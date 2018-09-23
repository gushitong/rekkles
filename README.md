
# aryadb

[![Build Status](https://travis-ci.com/gushitong/aryadb.svg?branch=master)](https://travis-ci.com/gushitong/aryadb)

aryadb is a high performance no-sql database build on [BadgerDB](https://github.com/dgraph-io/badger) with redis protocol
support. It meant to provide a key-value store alternative to redis.

## Install

    go get github.com/gushitong/aryadb
    
## Redis Command Support

* String:

    append decr decrby get getbit getrange getset incr incrby mget mset msetnx set setbit setex setnx setrange strlen

## Benchmark

This benchmark running on my local mac, aryadb has better performance on SSD.

* redis benchmark:

```bash
    $ redis-benchmark -p 6379 -t get,set -n 50000 -r 50000  -e
    ====== SET ======
      50000 requests completed in 0.90 seconds
      50 parallel clients
      3 bytes payload
      keep alive: 1
    
    92.45% <= 1 milliseconds
    99.14% <= 2 milliseconds
    99.52% <= 3 milliseconds
    99.85% <= 4 milliseconds
    99.86% <= 5 milliseconds
    99.86% <= 6 milliseconds
    99.90% <= 28 milliseconds
    99.91% <= 29 milliseconds
    100.00% <= 30 milliseconds
    55309.73 requests per second
    
    ====== GET ======
      50000 requests completed in 0.85 seconds
      50 parallel clients
      3 bytes payload
      keep alive: 1
    
    98.82% <= 1 milliseconds
    99.85% <= 2 milliseconds
    99.98% <= 3 milliseconds
    100.00% <= 3 milliseconds
    59031.88 requests per second
```

* aryadb benchmark:

```bash
    $ redis-benchmark -p 6380 -t get,set -n 50000 -r 50000  -e
    ====== SET ======
      50000 requests completed in 1.62 seconds
      50 parallel clients
      3 bytes payload
      keep alive: 1
    
    42.17% <= 1 milliseconds
    98.73% <= 2 milliseconds
    99.86% <= 3 milliseconds
    100.00% <= 3 milliseconds
    30883.26 requests per second
    
    ====== GET ======
      50000 requests completed in 1.46 seconds
      50 parallel clients
      3 bytes payload
      keep alive: 1
    
    84.07% <= 1 milliseconds
    99.88% <= 2 milliseconds
    100.00% <= 2 milliseconds
    34293.55 requests per second
```