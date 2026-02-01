package app

import (
	"fmt"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/mem"
)

type CPUInfo struct {
	ModelName string  `json:"model_name"`
	Cores     int32   `json:"cores"`
	Mhz       float64 `json:"mhz"`
}

type RAMInfo struct {
	Total uint64 `json:"total_bytes"`
}

type DiskInfo struct {
	Total     uint64 `json:"total_bytes"`
	Available uint64 `json:"available_bytes"`
}

type SystemInfo struct {
	CPU  []CPUInfo `json:"cpu"`
	RAM  RAMInfo   `json:"ram"`
	Disk DiskInfo  `json:"disk"`
}

func formatBytes(bytes uint64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := uint64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

func GetSystemInfo() (*SystemInfo, error) {
	// CPU
	cpuStats, err := cpu.Info()
	if err != nil {
		return nil, err
	}

	cpus := make([]CPUInfo, 0, len(cpuStats))
	for _, c := range cpuStats {
		cpus = append(cpus, CPUInfo{
			ModelName: c.ModelName,
			Cores:     c.Cores,
			Mhz:       c.Mhz,
		})
	}

	// RAM
	vm, err := mem.VirtualMemory()
	if err != nil {
		return nil, err
	}

	// Disk: chỉ lấy root `/`
	du, err := disk.Usage("/")
	if err != nil {
		return nil, err
	}

	return &SystemInfo{
		CPU: cpus,
		RAM: RAMInfo{
			Total: vm.Total,
		},
		Disk: DiskInfo{
			Total:     du.Total,
			Available: du.Free, // hoặc du.Available nếu muốn exclude reserved blocks
		},
	}, nil
}
