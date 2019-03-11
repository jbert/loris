# Loris

This is a simple set/get redis-a-like implemented in golang as an experiment
to better understand the performance characteristics of systems like this (and
to have fun in general).

The goal is to:

- gain more experience benchmarking and profiling software under load
- explore the tradeoffs of different ways of sharing data on a multi-threaded
  sysstem

Naively, one might imagine (I did :-) that redis throughput is ultimately
limited by its single threaded nature. However, coming up with a load which
demonstrates that looks like it might be challenging.

## Usage

You can specify different storage approaches on the command line:

```
$ loris -store mutex:map 
```

Will run a single map, protected by a mutex:

```
$ loris -store sharded:mutex:map 
```

Will run a number of maps (16), each protected by a mutex.
