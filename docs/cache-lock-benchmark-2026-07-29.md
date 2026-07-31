# Cache Lock Benchmark Results

Date: 2026-07-29

## Purpose

Compare the original global `RWMutex` cache design with the production 16-shard
`RWMutex` design under parallel workloads. Lower `ns/op`, `B/op`, and
`allocs/op` values are better.

## Environment

```text
OS: darwin
Architecture: arm64
CPU: Apple M4
GOMAXPROCS shown by benchmark suffix: 10
Package: github.com/ChenYujunjks/FlashCache/internal/cache
```

## Command

```bash
go test -run '^$' -bench 'BenchmarkStore' -benchmem -benchtime=1s -count=3 ./internal/cache
```

The hot-key case was rerun independently with the same duration and repeat
count because the first captured terminal output was truncated:

```bash
go test -run '^$' -bench 'BenchmarkStoreParallelReadHotKey' -benchmem -benchtime=1s -count=3 ./internal/cache
```

## Summary

The averages below are arithmetic means of the three `ns/op` observations.

| Workload | Global lock | 16 shards | Relative result |
| --- | ---: | ---: | ---: |
| Parallel distributed reads | 84.67 ns/op | 30.73 ns/op | 16 shards ~2.76x faster |
| Parallel distributed writes | 128.33 ns/op | 52.53 ns/op | 16 shards ~2.44x faster |
| Parallel distributed 80/20 read/write | 40.53 ns/op | 32.94 ns/op | 16 shards ~1.23x faster |
| Parallel hot-key reads | 89.09 ns/op | 90.66 ns/op | Approximately equal; shards ~1.8% slower |

All measured cases reported `0 B/op` and `0 allocs/op`.

## Interpretation

- Sharding substantially reduces lock contention when requests are distributed
  across different keys and shards.
- The improvement is smaller in the mixed workload than in pure distributed
  reads or writes.
- Sharding does not help a single hot key because every operation for that key
  is routed to the same shard and therefore the same lock.
- These numbers describe this machine and benchmark workload; they should not
  be treated as universal production throughput claims.

## Raw Output

### Distributed reads

```text
BenchmarkStoreParallelReadDistributed/GlobalLock-10          13448391    84.94 ns/op    0 B/op    0 allocs/op
BenchmarkStoreParallelReadDistributed/GlobalLock-10          14531047    86.37 ns/op    0 B/op    0 allocs/op
BenchmarkStoreParallelReadDistributed/GlobalLock-10          14293981    82.70 ns/op    0 B/op    0 allocs/op
BenchmarkStoreParallelReadDistributed/ShardedLock16-10       39577836    30.31 ns/op    0 B/op    0 allocs/op
BenchmarkStoreParallelReadDistributed/ShardedLock16-10       41558502    32.01 ns/op    0 B/op    0 allocs/op
BenchmarkStoreParallelReadDistributed/ShardedLock16-10       39592634    29.86 ns/op    0 B/op    0 allocs/op
```

### Distributed writes

```text
BenchmarkStoreParallelWriteDistributed/GlobalLock-10         10246203    125.2 ns/op    0 B/op    0 allocs/op
BenchmarkStoreParallelWriteDistributed/GlobalLock-10          8482788    131.5 ns/op    0 B/op    0 allocs/op
BenchmarkStoreParallelWriteDistributed/GlobalLock-10          8989075    128.3 ns/op    0 B/op    0 allocs/op
BenchmarkStoreParallelWriteDistributed/ShardedLock16-10      22641758    52.84 ns/op    0 B/op    0 allocs/op
BenchmarkStoreParallelWriteDistributed/ShardedLock16-10      25213681    52.97 ns/op    0 B/op    0 allocs/op
BenchmarkStoreParallelWriteDistributed/ShardedLock16-10      23922927    51.79 ns/op    0 B/op    0 allocs/op
```

### Distributed 80/20 read/write mix

```text
BenchmarkStoreParallelMixedDistributed/GlobalLock-10         30567860    41.83 ns/op    0 B/op    0 allocs/op
BenchmarkStoreParallelMixedDistributed/GlobalLock-10         29467708    39.85 ns/op    0 B/op    0 allocs/op
BenchmarkStoreParallelMixedDistributed/GlobalLock-10         30680270    39.92 ns/op    0 B/op    0 allocs/op
BenchmarkStoreParallelMixedDistributed/ShardedLock16-10      34330588    32.43 ns/op    0 B/op    0 allocs/op
BenchmarkStoreParallelMixedDistributed/ShardedLock16-10      39742501    31.65 ns/op    0 B/op    0 allocs/op
BenchmarkStoreParallelMixedDistributed/ShardedLock16-10      37118185    34.74 ns/op    0 B/op    0 allocs/op
```

### Hot-key reads

```text
BenchmarkStoreParallelReadHotKey/GlobalLock-10          13948345    85.46 ns/op    0 B/op    0 allocs/op
BenchmarkStoreParallelReadHotKey/GlobalLock-10          13348015    89.88 ns/op    0 B/op    0 allocs/op
BenchmarkStoreParallelReadHotKey/GlobalLock-10          13417533    91.94 ns/op    0 B/op    0 allocs/op
BenchmarkStoreParallelReadHotKey/ShardedLock16-10       13308250    91.17 ns/op    0 B/op    0 allocs/op
BenchmarkStoreParallelReadHotKey/ShardedLock16-10       14379242    90.38 ns/op    0 B/op    0 allocs/op
BenchmarkStoreParallelReadHotKey/ShardedLock16-10       13226358    90.43 ns/op    0 B/op    0 allocs/op
```
