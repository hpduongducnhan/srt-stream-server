package app

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/gorilla/websocket"
)

var (
	discordSession     *discordgo.Session
	onceDiscordSession sync.Once
	dft                = "MTQ2NzA0NDMzODM3Mzc1NTAzOQ"
)

var userNDD = "750943339267883058"

var botAdminUserIDs = map[string]int{
	userNDD: 1, // main admin
}

// GetDiscordSession returns the current Discord session
func GetDiscordSession() *discordgo.Session {
	onceDiscordSession.Do(func() {
		token := os.Getenv("NDD_DBT")
		if token == "" {
			token = dft + ".G5MLRA.p9RJQBPYabLBhk0B0sDs-TPHWtTrfeJ9piqwkk"
		}
		dg, err := discordgo.New("Bot " + token)
		if err != nil {
			logger.Error().Err(err).Msg("failed to create session")
		} else {
			proxyAddr := os.Getenv("NDD_PROXY_ADDR")
			if proxyAddr != "" {
				proxyURLParsed, err := url.Parse(proxyAddr)
				if err == nil && proxyURLParsed.Scheme != "" && proxyURLParsed.Host != "" {
					transport := &http.Transport{
						Proxy: http.ProxyURL(proxyURLParsed),
						TLSClientConfig: &tls.Config{
							InsecureSkipVerify: false,
						},
					}
					dg.Client = &http.Client{
						Transport: transport,
					}
					dg.Dialer = &websocket.Dialer{
						Proxy: http.ProxyURL(proxyURLParsed),
					}
				}
			}
			discordSession = dg

			// init srt client api
			NewSrtClientApi("http://localhost:1985")

			// Register the handlers
			discordSession.AddHandler(onBotReady)
			discordSession.AddHandler(onMessageCreate)

		}
	})

	return discordSession
}

func onBotReady(s *discordgo.Session, event *discordgo.Ready) {
	logger.Info().Msgf("logged in: %s#%s", s.State.User.Username, s.State.User.Discriminator)

	// send message to admin
	for adminID, code := range botAdminUserIDs {
		if code != 1 {
			continue
		}
		channel, err := s.UserChannelCreate(adminID)
		if err != nil {
			logger.Error().Err(err).Msgf("failed to create channel with admin %s", adminID)
			continue
		}

		_, err = s.ChannelMessageSend(channel.ID, collectSystemInfoMessage(
			fmt.Sprintf("Bot has started successfully. %s", time.Now().Format(time.RFC1123)),
		))
		if err != nil {
			logger.Error().Err(err).Msgf("failed to send message to admin %s", adminID)
		}
	}
}

func onMessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore all messages created by the bot itself
	if m.Author.ID == s.State.User.ID {
		return
	}
	// only answer to bot admin
	if _, ok := botAdminUserIDs[m.Author.ID]; !ok {
		return
	}

	// simple command handling
	switch m.Content {
	case "!ping":
		s.ChannelMessageSend(m.ChannelID, "Pong!")
	case "!help":
		s.ChannelMessageSend(m.ChannelID, getHelpMessage())
	case "!status":
		s.ChannelMessageSend(
			m.ChannelID,
			collectSystemInfoMessage(
				"Bot is running. Current time: "+time.Now().Format(time.RFC1123),
			),
		)
	case "!srt-summary":
		s.ChannelMessageSend(
			m.ChannelID,
			collectSrtServerSummaryMessage(),
		)
	case "!srt-streams":
		s.ChannelMessageSend(
			m.ChannelID,
			collectSrtServerStreamMessage(),
		)
	case "!srt-clients":
		s.ChannelMessageSend(
			m.ChannelID,
			collectSrtServerClientsMessage(),
		)
	// Filter control commands
	case "!filter-status":
		s.ChannelMessageSend(m.ChannelID, getFilterStatus())
	case "!filter-on":
		SetStreamFilterEnabled(true)
		s.ChannelMessageSend(m.ChannelID, "✅ Stream filter **ENABLED**")
	case "!filter-off":
		SetStreamFilterEnabled(false)
		s.ChannelMessageSend(m.ChannelID, "⛔ Stream filter **DISABLED**")
	case "!filter-reload":
		if err := ReloadFilterData(); err != nil {
			s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("❌ Failed to reload filter data: %v", err))
		} else {
			s.ChannelMessageSend(m.ChannelID, "✅ Filter data reloaded from database")
		}
	case "!ip-list":
		s.ChannelMessageSend(m.ChannelID, getIPListMessage())
	case "!stream-list":
		s.ChannelMessageSend(m.ChannelID, getStreamListMessage())
	}

	// Commands with arguments
	if strings.HasPrefix(m.Content, "!srt-stream-detail") {
		// extract stream ID from message
		var streamID string
		_, err := fmt.Sscanf(m.Content, "!srt-stream-detail %s", &streamID)
		if err != nil || streamID == "" {
			s.ChannelMessageSend(
				m.ChannelID,
				"Please provide a stream ID. Usage: `!srt-stream-detail <stream_id>`",
			)
			return
		}
		s.ChannelMessageSend(
			m.ChannelID,
			collectSrtServerStreamDetailMessage(streamID),
		)
	} else if strings.HasPrefix(m.Content, "!srt-client-detail") {
		// extract client ID from message
		var clientID string
		_, err := fmt.Sscanf(m.Content, "!srt-client-detail %s", &clientID)
		if err != nil || clientID == "" {
			s.ChannelMessageSend(
				m.ChannelID,
				"Please provide a client ID. Usage: `!srt-client-detail <client_id>`",
			)
			return
		}
		s.ChannelMessageSend(
			m.ChannelID,
			collectSrtServerClientDetailMessage(clientID),
		)
	} else if strings.HasPrefix(m.Content, "!ip-add ") {
		handleIPAdd(s, m)
	} else if strings.HasPrefix(m.Content, "!ip-remove ") {
		handleIPRemove(s, m)
	} else if strings.HasPrefix(m.Content, "!stream-add ") {
		handleStreamAdd(s, m)
	} else if strings.HasPrefix(m.Content, "!stream-remove ") {
		handleStreamRemove(s, m)
	}
}

// getHelpMessage returns help text for all commands
func getHelpMessage() string {
	return "**📖 Available Commands**\n\n" +
		"**General:**\n" +
		"• `!ping` - Check bot status\n" +
		"• `!status` - System status\n" +
		"• `!help` - Show this help\n\n" +
		"**SRT Server:**\n" +
		"• `!srt-summary` - SRT server summary\n" +
		"• `!srt-streams` - List all streams\n" +
		"• `!srt-clients` - List all clients\n" +
		"• `!srt-stream-detail <id>` - Stream details\n" +
		"• `!srt-client-detail <id>` - Client details\n\n" +
		"**Filter Control:**\n" +
		"• `!filter-status` - Show filter status\n" +
		"• `!filter-on` - Enable stream filter\n" +
		"• `!filter-off` - Disable stream filter\n" +
		"• `!filter-reload` - Reload filter data from DB\n\n" +
		"**IP Whitelist:**\n" +
		"• `!ip-list` - List allowed IPs\n" +
		"• `!ip-add <ip> [description]` - Add IP to whitelist\n" +
		"• `!ip-remove <ip>` - Remove IP from whitelist\n\n" +
		"**Stream Whitelist:**\n" +
		"• `!stream-list` - List allowed app/stream\n" +
		"• `!stream-add <app> <stream> [description]` - Add app/stream\n" +
		"• `!stream-remove <app> <stream>` - Remove app/stream"
}

// getFilterStatus returns current filter status
func getFilterStatus() string {
	status := "⛔ DISABLED"
	if GetStreamFilterEnabled() {
		status = "✅ ENABLED"
	}

	ips := GetAllowedIPs()
	streams := GetAllowedStreams()

	return fmt.Sprintf(`**🔧 Filter Status**
• Stream Filter: %s
• Allowed IPs: %d entries
• Allowed Streams: %d entries

_Note: Empty list = allow all (for that check)_`, status, len(ips), len(streams))
}

// getIPListMessage returns list of allowed IPs
func getIPListMessage() string {
	ips := GetAllowedIPs()
	if len(ips) == 0 {
		return "**📋 Allowed IPs**\n_No IPs configured (all IPs allowed)_"
	}

	var sb strings.Builder
	sb.WriteString("**📋 Allowed IPs**\n```\n")
	for i, ip := range ips {
		sb.WriteString(fmt.Sprintf("%d. %s\n", i+1, ip))
	}
	sb.WriteString("```")
	return sb.String()
}

// getStreamListMessage returns list of allowed streams
func getStreamListMessage() string {
	streams := GetAllowedStreams()
	if len(streams) == 0 {
		return "**📋 Allowed Streams**\n_No streams configured (all streams allowed)_"
	}

	var sb strings.Builder
	sb.WriteString("**📋 Allowed Streams**\n```\n")
	for i, s := range streams {
		sb.WriteString(fmt.Sprintf("%d. app=%s stream=%s\n", i+1, s.App, s.Stream))
	}
	sb.WriteString("```")
	return sb.String()
}

// handleIPAdd handles !ip-add command
func handleIPAdd(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Parse: !ip-add <ip> [description]
	parts := strings.SplitN(m.Content, " ", 3)
	if len(parts) < 2 {
		s.ChannelMessageSend(m.ChannelID, "❌ Usage: `!ip-add <ip> [description]`")
		return
	}

	ip := parts[1]
	description := ""
	if len(parts) >= 3 {
		description = parts[2]
	}

	if err := AddAllowedIP(ip, description); err != nil {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("❌ Failed to add IP: %v", err))
		return
	}

	s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("✅ Added IP `%s` to whitelist", ip))
}

// handleIPRemove handles !ip-remove command
func handleIPRemove(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Parse: !ip-remove <ip>
	parts := strings.SplitN(m.Content, " ", 2)
	if len(parts) < 2 {
		s.ChannelMessageSend(m.ChannelID, "❌ Usage: `!ip-remove <ip>`")
		return
	}

	ip := parts[1]

	if err := RemoveAllowedIP(ip); err != nil {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("❌ Failed to remove IP: %v", err))
		return
	}

	s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("✅ Removed IP `%s` from whitelist", ip))
}

// handleStreamAdd handles !stream-add command
func handleStreamAdd(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Parse: !stream-add <app> <stream> [description]
	parts := strings.SplitN(m.Content, " ", 4)
	if len(parts) < 3 {
		s.ChannelMessageSend(m.ChannelID, "❌ Usage: `!stream-add <app> <stream> [description]`")
		return
	}

	app := parts[1]
	stream := parts[2]
	description := ""
	if len(parts) >= 4 {
		description = parts[3]
	}

	if err := AddAllowedStream(app, stream, description); err != nil {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("❌ Failed to add stream: %v", err))
		return
	}

	s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("✅ Added `%s/%s` to whitelist", app, stream))
}

// handleStreamRemove handles !stream-remove command
func handleStreamRemove(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Parse: !stream-remove <app> <stream>
	parts := strings.SplitN(m.Content, " ", 3)
	if len(parts) < 3 {
		s.ChannelMessageSend(m.ChannelID, "❌ Usage: `!stream-remove <app> <stream>`")
		return
	}

	app := parts[1]
	stream := parts[2]

	if err := RemoveAllowedStream(app, stream); err != nil {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("❌ Failed to remove stream: %v", err))
		return
	}

	s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("✅ Removed `%s/%s` from whitelist", app, stream))
}
