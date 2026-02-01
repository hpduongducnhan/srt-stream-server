package app

import (
	"fmt"
	"strings"
)

func formatUptime(seconds int64) string {
	if seconds <= 0 {
		return "0s"
	}
	h := seconds / 3600
	m := (seconds % 3600) / 60
	s := seconds % 60
	var parts []string
	if h > 0 {
		parts = append(parts, fmt.Sprintf("%dh", h))
	}
	if m > 0 {
		parts = append(parts, fmt.Sprintf("%dm", m))
	}
	if s > 0 || len(parts) == 0 {
		parts = append(parts, fmt.Sprintf("%ds", s))
	}
	return strings.Join(parts, "")
}

func collectSrtServerStreamMessage() string {
	srtClient := GetSrtClientApi()
	s, err := srtClient.GetStreams()
	if err != nil {
		return "Failed to get SRT server streams: " + err.Error()
	}
	b := strings.Builder{}
	b.WriteString("SRT Server Streams Overview\n")
	b.WriteString("```")
	b.WriteString("**📡 SRS STREAMS OVERVIEW**\n")
	b.WriteString(fmt.Sprintf("> Total Streams: `%d`\n", len(s.Streams)))
	b.WriteString("\n")
	for i, stream := range s.Streams {
		uptimeStr := formatUptime(stream.LiveMS / 1000)
		videoTick := "❌"
		if stream.Video != nil {
			videoTick = "✅"
		}
		audioTick := "❌"
		if stream.Audio != nil {
			audioTick = "✅"
		}
		b.WriteString(fmt.Sprintf("- **%d. %s**\n", i+1, stream.Name))
		b.WriteString(fmt.Sprintf("    - ID: `%s`\n", stream.ID))
		b.WriteString(fmt.Sprintf("    - App: `%s`\n", stream.App))
		b.WriteString(fmt.Sprintf("    - VHost: `%s`\n", stream.Vhost))
		b.WriteString(fmt.Sprintf("    - Url: `%s`\n", stream.URL))
		b.WriteString(fmt.Sprintf("    - Clients: `%d`\n", stream.Clients))
		b.WriteString(fmt.Sprintf("    - Uptime: `%s`\n", uptimeStr))
		b.WriteString(fmt.Sprintf("    - Video: %s\n", videoTick))
		b.WriteString(fmt.Sprintf("    - Audio: %s\n", audioTick))
	}

	b.WriteString("```")
	return b.String()
}

func collectSrtServerStreamDetailMessage(streamID string) string {
	srtClient := GetSrtClientApi()
	s, err := srtClient.GetStream(streamID)
	if err != nil {
		return "Failed to get SRT stream detail: " + err.Error()
	}
	stream := s.Stream
	// Lấy danh sách client
	clientsResp, err := srtClient.GetClients()
	var clients []Client
	if err == nil {
		for _, c := range clientsResp.Clients {
			if c.Stream == stream.ID && c.Vhost == stream.Vhost {
				clients = append(clients, c)
			}
		}
	}
	pubCount := 0
	playCount := 0
	for _, c := range clients {
		if c.Type == "srt-publish" {
			pubCount++
		} else if c.Type == "srt-play" {
			playCount++
		}
	}
	b := strings.Builder{}
	b.WriteString("SRT Server Stream Detail\n")
	b.WriteString("```")
	b.WriteString("**🎬 STREAM DETAIL**\n")
	b.WriteString(fmt.Sprintf("> Name: `%s`\n", stream.Name))
	b.WriteString(fmt.Sprintf("> App: `%s`\n", stream.App))
	b.WriteString(fmt.Sprintf("> VHost: `%s`\n", stream.Vhost))
	b.WriteString(fmt.Sprintf("> Stream ID: `%s`\n", stream.ID))
	b.WriteString(fmt.Sprintf("> Uptime: `%d ms`\n", stream.LiveMS))
	b.WriteString(fmt.Sprintf("> Clients: `%d` (Pub: %d, Play: %d)\n", stream.Clients, pubCount, playCount))
	b.WriteString(fmt.Sprintf("> Frames: `%d`\n", stream.Frames))
	b.WriteString(fmt.Sprintf("> SendBytes: `%d`\n", stream.SendBytes))
	b.WriteString(fmt.Sprintf("> RecvBytes: `%d`\n", stream.RecvBytes))
	b.WriteString("\n")
	b.WriteString("**👥 CLIENTS**\n")
	if len(clients) == 0 {
		b.WriteString("- No clients found for this stream.\n")
	} else {
		for i, c := range clients {
			b.WriteString(fmt.Sprintf("- %d. ID: `%s`\n", i+1, c.ID))
			b.WriteString(fmt.Sprintf("    - Type: `%s`\n", c.Type))
			b.WriteString(fmt.Sprintf("    - IP: `%s`\n", c.IP))
			b.WriteString(fmt.Sprintf("    - Uptime: `%.0fs`\n", c.Alive))
			b.WriteString(fmt.Sprintf("    - SendBytes: `%d`\n", c.SendBytes))
			b.WriteString(fmt.Sprintf("    - RecvBytes: `%d`\n", c.RecvBytes))
		}
	}
	b.WriteString("\n**📺 VIDEO**\n")
	if stream.Video != nil {
		b.WriteString(fmt.Sprintf("- Codec: `%s`\n", stream.Video.Codec))
		b.WriteString(fmt.Sprintf("- Profile: `%s`\n", stream.Video.Profile))
		b.WriteString(fmt.Sprintf("- Level: `%s`\n", stream.Video.Level))
		b.WriteString(fmt.Sprintf("- Resolution: `%dx%d`\n", stream.Video.Width, stream.Video.Height))
	} else {
		b.WriteString("- No video info\n")
	}
	b.WriteString("\n**🔊 AUDIO**\n")
	if stream.Audio != nil {
		b.WriteString(fmt.Sprintf("- Codec: `%s`\n", stream.Audio.Codec))
		b.WriteString(fmt.Sprintf("- Profile: `%s`\n", stream.Audio.Profile))
		b.WriteString(fmt.Sprintf("- SampleRate: `%d`\n", stream.Audio.SampleRate))
		b.WriteString(fmt.Sprintf("- Channel: `%d`\n", stream.Audio.Channel))
	} else {
		b.WriteString("- No audio info\n")
	}
	b.WriteString("\n**📈 KBPS**\n")
	b.WriteString(fmt.Sprintf("- Recv: `%d`\n", stream.Kbps.Recv30s))
	b.WriteString(fmt.Sprintf("- Send: `%d`\n", stream.Kbps.Send30s))
	b.WriteString("\n**📝 PUBLISH**\n")
	b.WriteString(fmt.Sprintf("- Active: `%v`\n", stream.Publish.Active))
	b.WriteString(fmt.Sprintf("- CID: `%s`\n", stream.Publish.CID))
	b.WriteString("```")
	return b.String()
}
