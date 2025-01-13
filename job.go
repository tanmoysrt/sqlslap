package main

import (
	"fmt"
	"sync"
	"time"
)

type RunnerStats struct {
	Total  uint
	Failed uint
}

type Job struct {
	Title       string
	Runner      func(int, *Job) // Modified to accept runner ID
	Stats       map[int]*RunnerStats
	StartTime   time.Time
	StopChannel chan struct{}
	RunnerCount int
	WaitGroup   sync.WaitGroup
	Lock        sync.Mutex
}

func NewJob(title string, runnerCount int, runner func(int, *Job)) Job {
	stats := make(map[int]*RunnerStats)
	for i := 0; i < runnerCount; i++ {
		stats[i] = &RunnerStats{Total: 0, Failed: 0}
	}

	return Job{
		Title:       title,
		Runner:      runner,
		Stats:       stats,
		StartTime:   time.Now(),
		RunnerCount: runnerCount,
		StopChannel: make(chan struct{}),
		WaitGroup:   sync.WaitGroup{},
	}
}

func (op *Job) Start() {
	for i := 0; i < op.RunnerCount; i++ {
		go op.Runner(i, op) // Pass runner ID
	}
}

func (op *Job) WaitAndLog() {
	go func() {
		for {
			time.Sleep(1 * time.Second)
			op.Log()
		}
	}()

	fmt.Scanln()
	op.Stop()
	op.WaitGroup.Wait()
}

func (op *Job) Stop() {
	close(op.StopChannel)
}

func (op *Job) Log() {
	elapsed := time.Since(op.StartTime)
	elapsed = elapsed.Round(time.Second)

	var totalOps, totalFailed uint

	// clear screen
	fmt.Print("\033[H\033[2J")
	fmt.Printf("Job: %s | Runners : %d\n", op.Title, op.RunnerCount)
	fmt.Println("┌────────┬──────────┬─────────┬──────────────┐")
	fmt.Printf("│ Runner │  Total   │ Failed  │ Op/s         │\n")
	fmt.Println("├────────┼──────────┼─────────┼──────────────┤")

	for i := 0; i < op.RunnerCount; i++ {
		stats := op.Stats[i]
		throughput := float64(stats.Total) / elapsed.Seconds()
		throughput = float64(int(throughput*100)) / 100
		if stats.Total == 0 {
			throughput = 0
		}
		fmt.Printf("│ %-6d │ %-8d │ %-7d │ %-10.2f   │\n", i, stats.Total, stats.Failed, throughput)
		totalOps += stats.Total
		totalFailed += stats.Failed
	}

	fmt.Println("├────────┴──────────┴─────────┴──────────────┤")
	totalThroughput := float64(totalOps) / elapsed.Seconds()
	totalThroughput = float64(int(totalThroughput*100)) / 100
	if totalOps == 0 {
		totalThroughput = 0
	}
	fmt.Printf("│ %-6s │ %-8d │ %-7d │ %-10.1f   │\n", "Total", totalOps, totalFailed, totalThroughput)
	fmt.Println("└────────────────────────────────────────────┘")
}

func (op *Job) AddToTotal(runnerId int, count uint) {
	op.Lock.Lock()
	defer op.Lock.Unlock()
	op.Stats[runnerId].Total += count
}

func (op *Job) AddToFailed(runnerId int, count uint) {
	op.Lock.Lock()
	defer op.Lock.Unlock()
	op.Stats[runnerId].Failed += count
}
