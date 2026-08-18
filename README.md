# Floodgate

A DDoS mitigation system built from scratch in Go, to deeply understand how rate limiting, distributed state, and traffic anomaly detection actually work — not just to use a library that does it for you.

This is a learning-first project: every algorithm, lock, and Redis script here was written and debugged by hand, then verified under real concurrent load (including Go's race detector). The goal isn't to compete with production-grade infrastructure like Cloudflare or AWS Shield — it's to build a smaller, fully-understood version of the same core ideas those systems rely on.

## What it does

Limits how many requests each IP can make, using a swappable rate-limiting algorithm behind a shared interface. Started in-memory, now moving state into Redis so it can work across multiple server instances.

## Progress

- [x] Token bucket rate limiter (in-memory), with HTTP middleware
- [x] Fixed window and sliding window algorithms, all swappable via a shared `Limiter` interface
- [x] Redis-backed fixed window (atomic `INCR` + `EXPIRE`)
- [x] Redis-backed token bucket (atomic Lua script)
- [x] Redis-backed sliding window
- [x] Traffic anomaly detection
- [x] Auto-blocking
- [ ] Attack simulator / load tester
- [ ] Deployment
