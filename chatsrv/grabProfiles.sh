#!/usr/bin/env bash

curl -o heap.pprof http://localhost:6060/debug/pprof/heap

curl -o cpu.pprof http://localhost:6060/debug/pprof/profile?seconds=30

curl -o block.pprof http://localhost:6060/debug/pprof/block

curl -o mutex.pprof http://localhost:6060/debug/pprof/mutex

curl -o trace.out http://localhost:6060/debug/pprof/trace?seconds=5
#go tool trace trace.out

# http://localhost:6060/debug/pprof/ 

# https://blog.golang.org/2011/06/profiling-go-programs.html