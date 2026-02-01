package app

import (
	"fmt"
	"strings"
)

func collectSrtServerSummaryMessage() string {
	srtClient := GetSrtClientApi()
	s, err := srtClient.GetSummary()
	if err != nil {
		return "Failed to get SRT server summary: " + err.Error()
	}
	sys := s.Data.System
	self := s.Data.Self

	cpuLv := levelByPercent(sys.CPUPercent)
	ramLv := levelByPercent(sys.MemRamPercent)
	diskLv := levelByPercent(sys.DiskBusyPercent)
	loadLv := levelByLoad(sys.Load1m, sys.CPUs)

	srsMemMB := self.MemKB / 1024
	srsMemLv := OK
	if srsMemMB > 2048 {
		srsMemLv = CRITICAL
	} else if srsMemMB > 1024 {
		srsMemLv = WARNING
	}

	var b strings.Builder
	b.WriteString("SRT Server Summary Report\n")
	b.WriteString("```")
	b.WriteString("**📡 SRS SUMMARY REPORT**\n")
	b.WriteString(fmt.Sprintf("> Server: `%s`\n", s.Server))
	b.WriteString(fmt.Sprintf("> Version: `%s`\n", self.Version))
	b.WriteString(fmt.Sprintf("> Uptime: `%ds`\n\n", self.SrsUptime))
	b.WriteString("**🎥 SRS PROCESS**\n")
	b.WriteString(fmt.Sprintf("- CPU: `%.2f%%` → %s\n", self.CPUPercent, levelByPercent(self.CPUPercent)))
	b.WriteString(fmt.Sprintf("- Memory: `%d MB` → %s\n", srsMemMB, srsMemLv))
	b.WriteString(fmt.Sprintf("- Connections: `%d`\n\n", sys.ConnSrs))
	b.WriteString("**🖥 SYSTEM**\n")
	b.WriteString(fmt.Sprintf("- CPU: (%d) `%.2f%%` → %s\n", sys.CPUs, sys.CPUPercent, cpuLv))
	b.WriteString(fmt.Sprintf("- RAM: `%.2f%%` → %s\n", sys.MemRamPercent, ramLv))
	b.WriteString(fmt.Sprintf("- Disk Busy: `%.2f%%` → %s\n", sys.DiskBusyPercent, diskLv))
	b.WriteString(fmt.Sprintf("- Load(1m): `%.2f` / CPU(%d) → %s\n", sys.Load1m, sys.CPUs, loadLv))
	b.WriteString(fmt.Sprintf("- Connections: `%d` (UDP: %d)\n", sys.ConnSys, sys.ConnSysUDP))
	b.WriteString("```")
	return b.String()
}
