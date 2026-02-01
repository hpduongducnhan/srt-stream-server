package app

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

func logActionToDiscord(msg string, hook HookRequest) {
	if discordSession == nil {
		return
	}
	toSendDiscordMessage := fmt.Sprintf(
		"action %s\n client_id=%s\n ip=%s\n vhost=%s\n app=%s\n stream=%s\n param=%s",
		hook.Action, hook.ClientID, hook.IP, hook.Vhost, hook.App, hook.Stream, hook.Param,
	)
	_, err := discordSession.ChannelMessageSend(
		"1467104920552345756",
		fmt.Sprintf(
			"Current time %s\n%s\n```%s```",
			time.Now().Format(time.RFC1123),
			msg, toSendDiscordMessage,
		),
	)
	if err != nil {
		logger.Error().Err(err).Msg("failed to send action log to discord")
	}
}

func isEventAllowed(hook HookRequest) bool {
	// If stream filter is disabled, allow all
	if !GetStreamFilterEnabled() {
		return true
	}

	// Step 1: Check IP allowlist
	allowedIPList := GetAllowedIPs()
	if len(allowedIPList) > 0 {
		// IP list has values -> check if client IP is in the list
		ipAllowed := false
		for _, allowedIP := range allowedIPList {
			if hook.IP == allowedIP {
				ipAllowed = true
				break
			}
		}
		if !ipAllowed {
			logger.Warn().Str("ip", hook.IP).Msg("IP not in allowed list")
			return false
		}
	}
	// If IP list is empty -> allow by default (skip IP check)

	// Step 2: Check app/stream allowlist
	allowedStreamList := GetAllowedStreams()
	if len(allowedStreamList) > 0 {
		// Stream list has values -> check if app/stream pair is in the list
		streamAllowed := false
		for _, allowed := range allowedStreamList {
			if hook.App == allowed.App && hook.Stream == allowed.Stream {
				streamAllowed = true
				break
			}
		}
		if !streamAllowed {
			logger.Warn().
				Str("app", hook.App).
				Str("stream", hook.Stream).
				Msg("app/stream not in allowed list")
			return false
		}
	}
	// If stream list is empty -> allow by default (skip stream check)

	return true
}

func handleSrtHook(w http.ResponseWriter, r *http.Request) {
	var hook HookRequest
	err := json.NewDecoder(r.Body).Decode(&hook)
	if err != nil {
		logger.Error().Err(err).Msg("failed to decode hook request")
		sendJSON(w, 400, "invalid request")
		return
	}

	logger.Info().Str("action", hook.Action).
		Str("client_id", hook.ClientID).
		Str("ip", hook.IP).
		Str("vhost", hook.Vhost).
		Str("app", hook.App).
		Str("stream", hook.Stream).
		Str("param", hook.Param).
		Msg("received hook")

	logActionToDiscord("received hook", hook)

	var respCode int
	var respMsg string

	switch hook.Action {
	case "on_connect":
		if !isEventAllowed(hook) {
			logger.Warn().Str("client_id", hook.ClientID).Msg("connection not allowed, rejecting")
			logActionToDiscord("Client connection not allowed, rejecting", hook)
			respCode = 403
			respMsg = "connection not allowed"
			break
		}
		logActionToDiscord("Client connected", hook)
		respCode = 0
		respMsg = "ok"
	case "check_publish":
		if !isEventAllowed(hook) {
			logger.Warn().Str("client_id", hook.ClientID).Msg("publish not allowed, rejecting")
			logActionToDiscord("Stream publish not allowed, rejecting", hook)
			respCode = 403
			respMsg = "publish not allowed"
			break
		}
		logActionToDiscord("Stream published", hook)
		respCode = 0
		respMsg = "ok"
	case "check_play":
		if !isEventAllowed(hook) {
			logger.Warn().Str("client_id", hook.ClientID).Msg("play not allowed, rejecting")
			logActionToDiscord("Stream play not allowed, rejecting", hook)
			respCode = 403
			respMsg = "play not allowed"
			break
		}
		logActionToDiscord("Stream played", hook)
		respCode = 0
		respMsg = "ok"
	case "check_close":
		logActionToDiscord("Client disconnected", hook)
		respCode = 0
		respMsg = "ok"
	// Bổ sung các action SRT có thể gọi nhưng không cần kiểm soát đặc biệt
	case "on_play":
		logActionToDiscord("on_play", hook)
		respCode = 0
		respMsg = "ok"
	case "on_publish":
		logActionToDiscord("on_publish", hook)
		respCode = 0
		respMsg = "ok"
	case "on_unpublish":
		logActionToDiscord("on_unpublish", hook)
		respCode = 0
		respMsg = "ok"
	case "on_stop":
		logActionToDiscord("on_stop", hook)
		respCode = 0
		respMsg = "ok"
	case "on_close":
		logActionToDiscord("on_close", hook)
		respCode = 0
		respMsg = "ok"
	case "on_dvr":
		logActionToDiscord("on_dvr", hook)
		respCode = 0
		respMsg = "ok"
	case "on_hls":
		logActionToDiscord("on_hls", hook)
		respCode = 0
		respMsg = "ok"
	default:
		logger.Warn().Str("action", hook.Action).Msg("unknown action")
		logActionToDiscord("Unknown", hook)
		respCode = 400
		respMsg = "unknown action"
	}

	logger.Info().Str("action", hook.Action).Msg("action handled, sending response")
	sendJSON(w, respCode, respMsg)
}
