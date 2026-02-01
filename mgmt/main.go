package main

import (
	"ndd/srt/app"
)

func main() {
	logger := app.GetLogger()
	botSession := app.GetDiscordSession()
	err := botSession.Open()
	if err != nil {
		logger.Error().Err(err).Msg("failed to open communication session")
	}
	defer botSession.Close()

	addr := ":10081"

	httpServer := app.NewHttpServer(addr)
	go func() {
		logger.Info().Msg("starting API server on " + addr)
		if err := httpServer.ListenAndServe(); err != nil {
			logger.Error().Err(err).Msg("API server stopped")
		}
	}()

	logger.Info().Msg("server is now running. Press CTRL-C to exit.")
	select {}
}
