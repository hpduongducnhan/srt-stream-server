package app

import (
	"fmt"
	"strings"
)

func collectSrtServerClientsMessage() string {
	srtClient := GetSrtClientApi()
	s, err := srtClient.GetClients()
	if err != nil {
		return "Failed to get SRT server clients: " + err.Error()
	}
	b := strings.Builder{}
	b.WriteString("SRT Server Clients Overview\n")
	b.WriteString("```")
	b.WriteString("**👥 SRS CLIENTS OVERVIEW**\n")
	b.WriteString(fmt.Sprintf("> Total Clients: `%d`\n", len(s.Clients)))
	b.WriteString("\n")
	for i, c := range s.Clients {
		uptimeStr := formatUptime(int64(c.Alive))
		b.WriteString(fmt.Sprintf("- **%d. ID:** `%s`\n", i+1, c.ID))
		b.WriteString(fmt.Sprintf("    - Type: `%s`\n", c.Type))
		b.WriteString(fmt.Sprintf("    - Stream: `%s`\n", c.Stream))
		b.WriteString(fmt.Sprintf("    - VHost: `%s`\n", c.Vhost))
		b.WriteString(fmt.Sprintf("    - IP: `%s`\n", c.IP))
		b.WriteString(fmt.Sprintf("    - Uptime: `%s`\n", uptimeStr))
		b.WriteString(fmt.Sprintf("    - SendBytes: `%d`\n", c.SendBytes))
		b.WriteString(fmt.Sprintf("    - RecvBytes: `%d`\n", c.RecvBytes))
	}
	b.WriteString("```")
	return b.String()
}

func collectSrtServerClientDetailMessage(clientID string) string {
	srtClient := GetSrtClientApi()
	s, err := srtClient.GetClient(clientID)
	if err != nil {
		return "Failed to get SRT client detail: " + err.Error()
	}
	c := s.Client
	b := strings.Builder{}
	b.WriteString("SRT Server Client Detail\n")
	b.WriteString("```")
	b.WriteString("**👤 CLIENT DETAIL**\n")
	b.WriteString(fmt.Sprintf("> ID: `%s`\n", c.ID))
	b.WriteString(fmt.Sprintf("> Type: `%s`\n", c.Type))
	b.WriteString(fmt.Sprintf("> Stream: `%s`\n", c.Stream))
	b.WriteString(fmt.Sprintf("> VHost: `%s`\n", c.Vhost))
	b.WriteString(fmt.Sprintf("> Name: `%s`\n", c.Name))
	b.WriteString(fmt.Sprintf("> IP: `%s`\n", c.IP))
	b.WriteString(fmt.Sprintf("> Uptime: `%.0fs`\n", c.Alive))
	b.WriteString(fmt.Sprintf("> SendBytes: `%d`\n", c.SendBytes))
	b.WriteString(fmt.Sprintf("> RecvBytes: `%d`\n", c.RecvBytes))
	b.WriteString("\n**📈 KBPS**\n")
	b.WriteString(fmt.Sprintf("- Recv: `%d`\n", c.Kbps.Recv30s))
	b.WriteString(fmt.Sprintf("- Send: `%d`\n", c.Kbps.Send30s))
	b.WriteString("```")
	return b.String()
}
