package main

import (
	"sync"
	"time"

	"github.com/amarnathcjd/gogram/telegram"
	"github.com/rs/zerolog/log"
)

var upgradeMutex sync.Mutex

func upgradeGift(client *telegram.Client) {
	upgradeMutex.Lock()
	defer upgradeMutex.Unlock()

	log.Info().Msg("Starting gift upgrade watcher")

	timeout := time.After(5 * time.Minute)
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()
	//A wait of 3 seconds is required before calling the method (code 420)
	time.Sleep(3 * time.Second)

loop:
	for {
		select {
		case <-timeout:
			log.Info().Msg("Gift upgrade watcher timeout (5 minutes), stopping")
			return

		case <-ticker.C:
			log.Debug().Msg("Checking for upgradable star gifts...")

			savedGifts, err := client.PaymentsGetSavedStarGifts(&telegram.PaymentsGetSavedStarGiftsParams{Peer: &telegram.InputPeerSelf{}})
			if err != nil {
				log.Err(err).Msg("Failed to get saved star gifts")
				continue
			}

			for _, gift := range savedGifts.Gifts {
				if !gift.CanUpgrade {
					continue
				}

				inputGift := &telegram.InputSavedStarGiftUser{MsgID: gift.MsgID}

				_, err = client.PaymentsUpgradeStarGift(false, inputGift)
				if err != nil {
					log.Err(err).Msg("Failed to upgrade star gift")
					continue
				}
				log.Info().Msg("Successfully upgraded star gift")

				break loop
			}

			log.Debug().Msg("No upgradable gifts found")
		}
	}
}
