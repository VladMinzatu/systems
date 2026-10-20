package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	tea "charm.land/bubbletea/v2"
)

const cpuSampleInterval = time.Second

// cpuStat holds the aggregate CPU time counters from the "cpu" line of
// /proc/stat, in USER_HZ ticks.
type cpuStat struct {
	idle  uint64
	total uint64
}

type cpuSampleMsg struct {
	stat cpuStat
	err  error
}

// sampleCPU returns a command that reads /proc/stat after cpuSampleInterval.
func sampleCPU() tea.Cmd {
	return tea.Tick(cpuSampleInterval, func(time.Time) tea.Msg {
		stat, err := readCPUStat()
		return cpuSampleMsg{stat: stat, err: err}
	})
}

// readCPUStat parses the aggregate "cpu" line of /proc/stat:
//
//	cpu  user nice system idle iowait irq softirq steal guest guest_nice
func readCPUStat() (cpuStat, error) {
	f, err := os.Open("/proc/stat")
	if err != nil {
		return cpuStat{}, err
	}
	defer f.Close()

	sc := bufio.NewScanner(f)
	for sc.Scan() {
		fields := strings.Fields(sc.Text())
		if len(fields) < 5 || fields[0] != "cpu" {
			continue
		}
		var s cpuStat
		// guest and guest_nice are already included in user and nice.
		for i, field := range fields[1:min(len(fields), 9)] {
			v, err := strconv.ParseUint(field, 10, 64)
			if err != nil {
				return cpuStat{}, fmt.Errorf("parse /proc/stat: %w", err)
			}
			s.total += v
			if i == 3 || i == 4 { // idle, iowait
				s.idle += v
			}
		}
		return s, nil
	}
	if err := sc.Err(); err != nil {
		return cpuStat{}, err
	}
	return cpuStat{}, fmt.Errorf("no cpu line in /proc/stat")
}

// cpuPercent returns the share of non-idle time between two samples, 0-100.
func cpuPercent(prev, cur cpuStat) float64 {
	total := cur.total - prev.total
	if total == 0 {
		return 0
	}
	idle := cur.idle - prev.idle
	return 100 * float64(total-idle) / float64(total)
}
