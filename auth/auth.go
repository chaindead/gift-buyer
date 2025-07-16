package auth

import (
	"github.com/amarnathcjd/gogram/telegram"
	"github.com/chaindead/gift-buyer/config"
	"github.com/rs/zerolog/log"
)

func Auth(cfg *config.Config) {
	client, err := telegram.NewClient(telegram.ClientConfig{
		AppID:   cfg.AppID,
		AppHash: cfg.AppHash,
	})
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create telegram client")
	}

	if _, err = client.Conn(); err != nil {
		log.Fatal().Err(err).Msg("Failed to conn telegram client")
	}

	if err = client.AuthPrompt(); err != nil {
		log.Fatal().Err(err).Msg("Failed to auth telegram client")
	}

	log.Info().Str("session", client.ExportSession()).Msg("Successfully authenticated telegram session, for env")
}
