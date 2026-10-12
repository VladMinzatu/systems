# Backlog

---

## Track A: eBPF through applications (ebpf-go)

| # | Project | What you build | Prog / map types exercised | OS fundamental angle |
|---|---|---|---|---|
| A1 | **execsnoop/opensnoop, two ways** | Process and file-open tracer; emit via both perf buffer and ringbuf, load-test both and count dropped events | `tracepoint`, `raw_tp`; `RINGBUF` vs `PERF_EVENT_ARRAY` | Ringbuf reserve/commit design, per-CPU vs shared buffers, ordering |
| A2 | **Syscall latency histograms** | biolatency/funclatency-style log2 histograms per syscall and per cgroup | `fentry`/`fexit` vs `kprobe`/`kretprobe` (measure the overhead difference); `PERCPU_ARRAY`, `HASH` | Per-CPU data to avoid contention, log-bucketed histograms |
| A3 | **Sampling CPU + off-CPU profiler → pprof** | On-CPU sampling plus off-CPU time from `sched_switch`, output as `.pb.gz` | `perf_event`, `tp_btf/sched_switch`; `STACK_TRACE`, `LRU_HASH` | Stack unwinding (frame pointers), ELF symbolization, scheduler states |
| A4 | **Go runtime X-ray** | Trace goroutine create/park/ready and GC phases in an unmodified Go binary | `uprobe` (and why `uretprobe` breaks with Go stack growth); `uprobe.multi` | Go's G/M/P model, reading runtime structs via DWARF offsets |
| A5 | **TLS plaintext sniffer** | Hook `SSL_read`/`SSL_write` in libssl to see HTTPS plaintext from curl/Python | `uprobe`/`uretprobe`; `HASH` keyed on tid to stitch entry and exit | Userspace ABI, calling conventions, shared-library symbol resolution |
| A6 | **XDP firewall + L4 load balancer** | CIDR blocklist, then a Maglev-hashed LB with conntrack | `XDP`; `LPM_TRIE`, `LRU_HASH`, `DEVMAP`, `bpf_redirect_map` | Packet parsing under verifier bounds, consistent hashing, connection tracking |
| A7 | **AF_XDP userspace packet echo/forwarder** | Kernel-bypass packet I/O in Go | `XDP` + `XSKMAP` | UMEM plus fill/completion/RX/TX rings: the textbook shared-memory SPSC ring with memory barriers (next step after the SPSC ringbuf) |
| A8 | **Sockmap TCP proxy** | Proxy that splices sockets in-kernel, benchmarked against a userspace proxy and `splice(2)` | `sk_msg`, `sk_skb`, `sockops`; `SOCKHASH` | Socket layer, zero-copy, the cost of a byte per syscall |
| A9 | **Per-cgroup network policy** | "Container X may only connect to Y", plus transparent connect redirect | `cgroup_skb`, `cgroup/connect4`, `sk_lookup`; `SK_STORAGE`, `CGROUP_STORAGE` | cgroups v2 hierarchy, local storage attached to kernel objects |
| A10 | **Tail-call protocol pipeline: DNS logger** | TC program parses Ethernet → IP → UDP → DNS, dispatching each layer via tail calls | `tcx`/`TC`; `PROG_ARRAY`, `BLOOM_FILTER` (seen-domain dedup) | Layered parsing within verifier instruction and loop limits |
| A11 | **Mini-Tetragon (LSM enforcement)** | Deny exec of binaries outside an allowlist, block writes to protected paths | `lsm`; `TASK_STORAGE`, `INODE_STORAGE` | Security hooks, observing vs enforcing |
| A12 | **Kernel-side token-bucket rate limiter** | Per-source-IP rate limit with periodic refill | `TC`/`XDP`; `bpf_timer`, `PERCPU_HASH`, `bpf_spin_lock` in map values | Token bucket, locking inside the kernel, timers |
| A13 | **Custom `ps`/`lsof` via BPF iterators** | Walk tasks and their open files without parsing `/proc` | `iter/task`, `iter/task_file`; pinned iterators in bpffs | Kernel object graphs; compare with the cost of text parsing in `/proc` |
| A14 | **Toy sched_ext scheduler** | FIFO, then weighted-vtime, then "latency-critical cgroup first". Check ebpf-go struct_ops maturity first; libbpf skeleton is the fallback | `struct_ops` (sched_ext) | Run queues, dispatch queues, fairness, CFS/EEVDF concepts |
| A15 | **TCP congestion control in BPF** | Reno, then something like BBR-lite | `struct_ops` (`tcp_congestion_ops`) | Control loops, AIMD; ties into backpressure work |
| A16 | **Kernel↔user control plane** | Push config/commands from Go into BPF without map polling; stretch: a data structure living in a BPF arena | `USER_RINGBUF`, `QUEUE`/`STACK` maps, `ARENA` (6.9+) | Shared-memory allocators, bidirectional rings |

---

## Track B: OS fundamentals rebuilt in userspace (Go runtime / vLLM style)

| # | Project | What you build | Real-world analog |
|---|---|---|---|
| B1 | **Bounded MPMC queue** | Vyukov MPMC with sequence numbers, benchmarked against channels and a mutex-protected slice | Disruptor, Go `runq` |
| B2 | **Work-stealing scheduler** | Chase-Lev deque plus a mini M:N scheduler with per-worker local queues, a global queue and stealing | Go GMP, Tokio, Rayon |
| B3 | **Locks lab** | Spinlock, ticket lock, MCS lock, futex-based mutex via raw `SYS_FUTEX`; measure fairness/throughput under contention; reproduce Go's mutex starvation mode | Linux qspinlock, `sync.Mutex` |
| B4 | **Size-class allocator** | Thread-cache → central → page-heap allocator over an mmap'd arena, plus a buddy allocator for pages | tcmalloc, Go mcache/mcentral/mheap, Linux buddy + slab |
| B5 | **Paged KV-cache manager** | Block tables, refcounted copy-on-write blocks, prefix sharing via block hashing, LRU eviction | vLLM PagedAttention (really virtual memory) |
| B6 | **Continuous-batching scheduler** | Token-budget request scheduler with preemption (swap vs recompute) on top of B5 | vLLM scheduler, OS preemptive scheduling |
| B7 | **Buffer pool manager** | Fixed frames, pin counts, CLOCK/LRU-K eviction, dirty write-back, `O_DIRECT` + `pread` | Postgres shared buffers, Linux page cache |
| B8 | **WAL + crash recovery** | Checksummed records, group commit, torn-write detection, crash-injection test harness | Postgres/etcd WAL, ext4 journal |
| B9 | **Hierarchical timing wheel** | Timer wheel with cascading, benchmarked against the go-heap heap | Linux timers, Kafka purgatory, Go timers |
| B10 | **RCU / epoch reclamation + seqlock** | Read-mostly concurrent map with epoch-based reclamation | Linux RCU, seqlock for `jiffies` |
| B11 | **Per-CPU overwriting trace buffer** | Lossy flight-recorder ring with sharded buffers and merge-on-read | ftrace, Go flight recorder (extends go-flight-recorder) |
| B12 | **Toy concurrent GC** | Tri-color mark-sweep with a write barrier over a simulated heap | Go GC |
| B13 | **Zero-copy protocol parser** | Resumable state-machine parser for RESP or HTTP/1.1 that handles partial reads | Redis, nginx parser |

---

## Track C: Low-level Linux

| # | Project | What you build | Linux surface |
|---|---|---|---|
| C1 | **Raw io_uring in Go (no liburing)** | mmap the SQ/CQ rings yourself, then build a `cat`/file server | `io_uring_setup`/`io_uring_enter`, shared rings, memory ordering |
| C2 | **epoll reactor / netpoller** | Edge-triggered event loop, compared against Go's netpoller | `epoll`, nonblocking fds, `eventfd`, `timerfd` |
| C3 | **Mini container runtime** | Namespaces, `pivot_root`, cgroups v2 limits; then trace/police it with A9/A11 | `clone` flags, `unshare`, mounts, cgroupfs |
| C4 | **seccomp sandbox supervisor** | Classic-BPF syscall filter plus `SECCOMP_RET_USER_NOTIF` brokering | seccomp, cBPF (eBPF's ancestor), pidfd |
| C5 | **strace clone** | Syscall decoding with argument pretty-printing | `ptrace`, register ABI, reading tracee memory |
| C6 | **Cross-process shared-memory IPC** | SPSC ring in `memfd` plus futex wakeups, benchmarked against pipes and UDS | `memfd_create`, `mmap`, futex, `SCM_RIGHTS` fd passing |
| C7 | **userfaultfd lazy loader / snapshotter** | Demand-page memory from a file and snapshot dirty pages | `userfaultfd`, `mprotect`, soft-dirty bits (Firecracker, CRIU) |
| C8 | **FUSE filesystem** | Log-structured FS with an inode table, then a crash test (combines with B8) | FUSE, VFS concepts |
| C9 | **ELF loader + symbolizer** | Parse ELF/DWARF, load a static binary by hand; reuse in A3 | ELF, `execve` internals, auxv |
| C10 | **Tiny KVM VMM** | Boot real-mode code, emulate a serial port, then go to long mode | `/dev/kvm` ioctls, page tables, VM exits |
| C11 | **Green threads with an asm context switch** | Hand-written stack switching (covers "Assembly projects") | Calling conventions, stacks, `ucontext` |

---

## Projects that build on each other

- **Observability:** A1 → A2 → A3 (+ C9 for symbolization) → A4. Ends as a small Parca/bpftrace.
- **Networking:** A6 → A7 → A10 → A12 → A8. Packet parsing, kernel bypass, policy, proxying.
- **Containers / security:** C3 → C4 → A9 → A11. A runtime with syscall, network and exec enforcement.
- **Inference engine core:** B1 → B4 → B5 → B6. A mini vLLM scheduler with a fake model.
- **Storage:** B13 → B8 → B7 → C8 → C1. A tiny storage engine on io_uring.
- **Scheduling:** B2 → A14 → A15. Userspace, then kernel, then network "scheduling".

Coverage gaps (niche prog types): `flow_dissector`, `lirc`, `cgroup/sysctl`, `cgroup/dev`. Could fold one into A9 for full coverage.
