package sysmetrics

import (
	"runtime"
	"sync"
	"syscall"
	"time"
)

var (
	monitorInstance *ResourceMonitor
	monitorOnce     sync.Once
)

// ResourceMonitor tracks resource usage in the background
type ResourceMonitor struct {
	mu            sync.RWMutex
	lastCPUUsage  uint64
	lastCheckTime time.Time
	cpuPercent    float64
	memoryUsage   uint64
	memoryLimit   uint64
	cpuLimit      float64
	diskUsage     uint64
	diskSize      uint64
	stopCh        chan struct{}
	started       bool
}

// Start begins background monitoring of CPU usage
func (rm *ResourceMonitor) Start() {
	rm.mu.Lock()
	if rm.started {
		rm.mu.Unlock()
		return
	}
	rm.started = true

	// Initialize first reading
	usage, err := readCPUUsage()
	if err == nil {
		rm.lastCPUUsage = usage
		rm.lastCheckTime = time.Now()
	}
	rm.mu.Unlock()

	go rm.monitor()
}

// Stop stops the background monitoring
func (rm *ResourceMonitor) Stop() {
	rm.mu.Lock()
	defer rm.mu.Unlock()
	if rm.started {
		close(rm.stopCh)
		rm.started = false
	}
}

// monitor runs in the background and updates all resource usage periodically
func (rm *ResourceMonitor) monitor() {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-rm.stopCh:
			return
		case <-ticker.C:
			rm.updateAllMetrics()
		}
	}
}

// updateAllMetrics updates all resource metrics
func (rm *ResourceMonitor) updateAllMetrics() {
	// Update CPU usage
	usage, err := readCPUUsage()
	if err == nil {
		rm.mu.Lock()
		now := time.Now()
		elapsed := now.Sub(rm.lastCheckTime).Seconds()

		if elapsed > 0 && usage > rm.lastCPUUsage {
			delta := usage - rm.lastCPUUsage // in nanoseconds

			// Normalize by the number of CPUs available to the cgroup (or host fallback)
			cpus, err := GetCPULimit()
			if err != nil || cpus <= 0 {
				cpus = float64(runtime.NumCPU())
			}

			denom := elapsed * cpus * 1e9 // ns worth of CPU time available over interval
			if denom > 0 {
				cpuPercent := (float64(delta) / denom) * 100
				if cpuPercent < 0 {
					cpuPercent = 0
				}
				rm.cpuPercent = cpuPercent
			}
		}

		rm.lastCPUUsage = usage
		rm.lastCheckTime = now
		rm.mu.Unlock()
	}

	// Update memory usage
	if memUsage, err := GetMemoryUsageBytes(); err == nil {
		rm.mu.Lock()
		rm.memoryUsage = memUsage
		rm.mu.Unlock()
	}

	// Update memory limit
	if memLimit, err := GetMemoryLimitBytes(); err == nil {
		rm.mu.Lock()
		rm.memoryLimit = memLimit
		rm.mu.Unlock()
	}

	// Update CPU limit
	if cpuLimit, err := GetCPULimit(); err == nil {
		rm.mu.Lock()
		rm.cpuLimit = cpuLimit
		rm.mu.Unlock()
	}

	// Update disk usage (but not disk size which is cached)
	if diskUsed, diskSize, err := rm.GetDiskStats("/"); err == nil {
		rm.mu.Lock()
		rm.diskUsage = diskUsed
		rm.diskSize = diskSize
		rm.mu.Unlock()
	}
}

// GetCPUPercent returns the current CPU usage percentage from the monitor
func (rm *ResourceMonitor) GetCPUPercent() float64 {
	rm.mu.RLock()
	defer rm.mu.RUnlock()
	return rm.cpuPercent
}

// GetMemoryUsage returns the current memory usage in bytes
func (rm *ResourceMonitor) GetMemoryUsage() uint64 {
	rm.mu.RLock()
	defer rm.mu.RUnlock()
	return rm.memoryUsage
}

// GetMemoryLimit returns the memory limit in bytes
func (rm *ResourceMonitor) GetMemoryLimit() uint64 {
	rm.mu.RLock()
	defer rm.mu.RUnlock()
	return rm.memoryLimit
}

// GetCPULimit returns the CPU limit
func (rm *ResourceMonitor) GetCPULimit() float64 {
	rm.mu.RLock()
	defer rm.mu.RUnlock()
	return rm.cpuLimit
}

// GetDiskUsage returns the current disk usage in bytes
func (rm *ResourceMonitor) GetDiskUsage() uint64 {
	rm.mu.RLock()
	defer rm.mu.RUnlock()
	return rm.diskUsage
}

// GetDiskSize returns the cached disk size in bytes
func (rm *ResourceMonitor) GetDiskSize() uint64 {
	rm.mu.RLock()
	defer rm.mu.RUnlock()
	return rm.diskSize
}

// GetDiskStats returns used and total bytes for the filesystem at path.
// It also caches the total size on first successful read.
func (rm *ResourceMonitor) GetDiskStats(path string) (used uint64, total uint64, err error) {
	var st syscall.Statfs_t
	if err = syscall.Statfs(path, &st); err != nil {
		return 0, 0, err
	}
	// Use Bavail to represent space available to unprivileged user (matches df)
	bs := uint64(st.Bsize)
	total = uint64(st.Blocks) * bs
	freeAvail := uint64(st.Bavail) * bs
	if total >= freeAvail {
		used = total - freeAvail
	} else {
		used = 0
	}
	return
}

// GetResourceUsage returns current container usage snapshot
func (rm *ResourceMonitor) GetResourceUsage() *ResourceUsage {
	return &ResourceUsage{
		MemoryUsage: rm.GetMemoryUsage(),
		MemoryLimit: rm.GetMemoryLimit(),
		CPUUsage:    rm.GetCPUPercent(),
		CPULimit:    rm.GetCPULimit(),
		DiskUsage:   rm.GetDiskUsage(),
		DiskSize:    rm.GetDiskSize(),
	}
}

// getMonitor returns the singleton ResourceMonitor instance
func getMonitor() *ResourceMonitor {
	monitorOnce.Do(func() {
		monitorInstance = &ResourceMonitor{
			stopCh: make(chan struct{}),
		}
		monitorInstance.Start()
	})
	return monitorInstance
}
