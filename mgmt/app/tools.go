package app

func collectSystemInfoMessage(msg string) string {
	message := msg + "\n"
	publicIP, _ := getPublicIp()
	clientIP, _ := getClientIp()
	systemInfo, _ := GetSystemInfo()

	message += "```"
	message += "Public IP: " + publicIP + "\n"
	message += "Client IP: " + clientIP + "\n"
	if systemInfo != nil {
		if len(systemInfo.CPU) > 0 {
			message += "CPU: " + systemInfo.CPU[0].ModelName + " (" + string(systemInfo.CPU[0].Cores) + " cores)\n"
		}
		message += "RAM: " + formatBytes(systemInfo.RAM.Total) + "\n"
		message += "Disk: " + formatBytes(systemInfo.Disk.Total) + " total, " + formatBytes(systemInfo.Disk.Available) + " available\n"
	}
	message += "```\n"
	return message
}
