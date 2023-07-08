---
title: Performance
---

# Estimated performance

According to some simple benchmark, for a simple event, it takes:

* ~60us for parsing
* ~9us for converting to Prometheus metric

So with a security margin, we can estimate an event should take 100us maximum.

Roughly, the exporter should be able to handle about 10kQps.

For a base of 1000 Salt minions, it should be able to sustain 10 jobs per minion per second, which is a quite high for Salt.
