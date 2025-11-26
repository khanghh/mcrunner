package sysmetrics

import (
	"errors"
	"net"
	"os"
	"runtime"
	"strconv"
	"strings"
)

// ResourceUsage holds container resource info
type ResourceUsage struct {
	MemoryUsage uint64  // current memory usage in bytes
	MemoryLimit uint64  // max allowed memory in bytes (0 = unlimited)
	CPUUsage    float64 // current CPU usage percent
	CPULimit    float64 // max CPUs allowed
	DiskUsage   uint64  // current disk usage in bytes
	DiskSize    uint64  // disk size in bytes
}

// GetMemoryUsageBytes returns current memory usage in bytes
func GetMemoryUsageBytes() (uint64, error) {
	// cgroup v2 path first, then v1 fallback
	candidates := []string{
		"/sys/fs/cgroup/memory.current",               // cgroup v2
		"/sys/fs/cgroup/memory/memory.usage_in_bytes", // cgroup v1
	}
	data, err := readFirstExisting(candidates)
	if err != nil {
		return 0, err
	}
	return strconv.ParseUint(strings.TrimSpace(string(data)), 10, 64)
}

// GetMemoryLimitBytes returns max memory limit in bytes (0 = unlimited)
func GetMemoryLimitBytes() (uint64, error) {
	// cgroup v2 path first, then v1 fallback
	candidates := []string{
		"/sys/fs/cgroup/memory.max",                   // cgroup v2
		"/sys/fs/cgroup/memory/memory.limit_in_bytes", // cgroup v1
	}
	data, err := readFirstExisting(candidates)
	if err != nil {
		return 0, err
	}
	str := strings.TrimSpace(string(data))
	if str == "max" || str == "-1" { // -1 for some v1 unlimited cases
		return 0, nil
	}
	return strconv.ParseUint(str, 10, 64)
}

// GetCPULimit returns maximum allowed CPU cores (can be fractional)
func GetCPULimit() (float64, error) {
	// Try cgroup v2 cpu.max
	if data, err := os.ReadFile("/sys/fs/cgroup/cpu.max"); err == nil {
		parts := strings.Fields(strings.TrimSpace(string(data)))
		if len(parts) >= 2 {
			if parts[0] == "max" {
				return float64(runtime.NumCPU()), nil
			}
			quota, _ := strconv.ParseFloat(parts[0], 64)
			period, _ := strconv.ParseFloat(parts[1], 64)
			if period == 0 {
				return float64(runtime.NumCPU()), nil
			}
			return quota / period, nil
		}
	}

	// Fallback to cgroup v1 cpu.cfs_* values
	quotaBytes, qErr := os.ReadFile("/sys/fs/cgroup/cpu/cpu.cfs_quota_us")
	periodBytes, pErr := os.ReadFile("/sys/fs/cgroup/cpu/cpu.cfs_period_us")
	if qErr == nil && pErr == nil {
		qStr := strings.TrimSpace(string(quotaBytes))
		pStr := strings.TrimSpace(string(periodBytes))
		if qStr == "-1" { // unlimited
			return float64(runtime.NumCPU()), nil
		}
		quota, _ := strconv.ParseFloat(qStr, 64)
		period, _ := strconv.ParseFloat(pStr, 64)
		if period == 0 {
			return float64(runtime.NumCPU()), nil
		}
		return quota / period, nil
	}

	// As a last resort, return host CPUs
	return float64(runtime.NumCPU()), nil
}

// readCPUUsage reads nanoseconds used by container CPU
func readCPUUsage() (uint64, error) {
	// cgroup v2: cpu.stat with usage_usec
	if data, err := os.ReadFile("/sys/fs/cgroup/cpu.stat"); err == nil {
		for _, line := range strings.Split(string(data), "\n") {
			if strings.HasPrefix(strings.TrimSpace(line), "usage_usec") {
				fields := strings.Fields(line)
				if len(fields) >= 2 {
					val, _ := strconv.ParseUint(fields[1], 10, 64)
					return val * 1000, nil // microseconds -> nanoseconds
				}
			}
		}
	}

	// cgroup v1: cpuacct.usage (already in nanoseconds)
	if data, err := os.ReadFile("/sys/fs/cgroup/cpuacct/cpuacct.usage"); err == nil {
		valStr := strings.TrimSpace(string(data))
		val, err := strconv.ParseUint(valStr, 10, 64)
		if err == nil {
			return val, nil
		}
	}
	return 0, errors.New("unable to read cpu usage from cgroup")
}

// readFirstExisting returns the content of the first existing file from paths
func readFirstExisting(paths []string) ([]byte, error) {
	for _, p := range paths {
		if _, err := os.Stat(p); err == nil {
			return os.ReadFile(p)
		}
	}
	return nil, errors.New("none of the provided cgroup files exist")
}

func GetOutboundIP() (net.IP, error) {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP, nil
}

func GetResourceUsage() *ResourceUsage {
	return getMonitor().GetResourceUsage()
}
