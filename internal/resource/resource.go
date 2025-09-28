package resource

import (
	"runtime"
	"time"
)

type Metadata struct {
	NumCPU            int       `json:"num_cpu"`
	GoMaxProcs        int       `json:"go_max_procs"`
	CPUBenchmarkScore int64     `json:"cpu_benchmark_score"`
	CapturedAt        time.Time `json:"captured_at"`
}

func Capture(numCPU int) Metadata {
	prev := runtime.GOMAXPROCS(0)
	runtime.GOMAXPROCS(numCPU)
	start := time.Now()
	score := runCPUBenchmark()
	return Metadata{
		NumCPU:            numCPU,
		GoMaxProcs:        prev,
		CPUBenchmarkScore: score,
		CapturedAt:        start,
	}
}

func runCPUBenchmark() int64 {
	var ops int64
	end := time.Now().Add(200 * time.Millisecond)
	for time.Now().Before(end) {
		for i := 0; i < 1000; i++ {
			ops += int64(i * i)
		}
	}
	return ops
}