---
title: Performance
---

# Estimated performance

According to a simple benchmark, for a single event it takes:

* ~60µs for parsing
* ~9µs for converting to Prometheus metric

With a security margin, we can estimate processing an event should take 100µs maximum.

Roughly, the exporter should be able to handle about 10kQps.

For a base of 1000 Salt minions, it should be able to sustain 10 jobs per minion per second, which is quite high for Salt.
