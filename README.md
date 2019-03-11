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
