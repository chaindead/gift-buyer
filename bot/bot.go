package main

import (
	"flag"
	"fmt"
	"sort"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/chaindead/gift-buyer/auth"
	. "github.com/chaindead/gift-buyer/config"
	_ "github.com/chaindead/gift-buyer/log"

	"github.com/amarnathcjd/gogram/telegram"
	"github.com/rs/zerolog/log"
)

var knownGifts = make(map[int64]bool)
var printInitialGifts sync.Once
var activityCounter int64

var needAuth = flag.Bool("auth", false, "get session string for bot to run")

// crashWatcher monitors activity and crashes app if no activity for 1 minute
func crashWatcher() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	lastCount := atomic.LoadInt64(&activityCounter)

	for {
		select {
		case <-ticker.C:
			currentCount := atomic.LoadInt64(&activityCounter)
			if currentCount == lastCount {
				log.Fatal().Int64("last_count", lastCount).
					Int64("current_count", currentCount).
					Msg("No activity detected for 1 minute - crashing app")
			}
			lastCount = currentCount
			log.Debug().Int64("activity_count", currentCount).Msg("Activity check passed")
		}
	}
}

func main() {
	cfg, err := LoadConfig()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to load configuration")
	}

	if *needAuth {
		auth.Auth(cfg)

		return
	}

	go crashWatcher()

	client, err := NewClientWithSession(cfg.AppID, cfg.AppHash, cfg.Session)
	if err != nil {
		log.Fatal().Err(err).Send()
	}

	mainClient := client
	adminPeer, err := client.ResolvePeer(cfg.Admin)
	if err != nil {
		log.Fatal().Err(err).Str("admin", cfg.Admin).Msg("Failed to resolve admin peer")
	}

	log.Info().Str("admin", cfg.Admin).Dur("poll_interval", cfg.PollInterval).Msg("Bot started, monitoring gifts")

	var previousHash int32 = 0
	for ; ; time.Sleep(cfg.PollInterval) {
		activityCounter++
		// Get available gifts
		availGifts, err := client.PaymentsGetStarGifts(previousHash)
		if err != nil {
			log.Err(err).Msg("Failed to get star gifts")
			continue
		}

		// Handle different response types
		switch response := availGifts.(type) {
		case *telegram.PaymentsStarGiftsNotModified:
			// No changes detected
			log.Debug().
				Int32("hash", previousHash).
				Int64("cnt", activityCounter).
				Msg("No gift changes detected")
		case *telegram.PaymentsStarGiftsObj:
			// Print initial gifts info on first start
			printInitialGifts.Do(func() {
				printLimitedGiftsInfo(response.Gifts)
			})

			// Try to auto-buy gifts
			autoBuyGifts(mainClient, response.Gifts)

			// Gifts have been modified
			currentHash := response.Hash

			var newExist bool
			for _, gift := range response.Gifts {
				if _, ok := knownGifts[gift.(*telegram.StarGiftObj).ID]; ok {
					continue
				} else {
					newExist = true
					knownGifts[gift.(*telegram.StarGiftObj).ID] = true
				}
			}

			log.Info().Int32("previous_hash", previousHash).Int32("current_hash", currentHash).
				Int("gifts_count", len(response.Gifts)).Msg("Gift changes detected")

			if !newExist {
				log.Info().Int64("cnt", activityCounter).Msg("No new-gift changes detected")
				previousHash = currentHash

				continue
			}

			message := formatGiftUpdateMessage(response.Gifts)
			if message == "" {
				log.Info().Msg("Empty message skip")
				continue
			}
			_, err = client.SendMessage(adminPeer, message, &telegram.SendOptions{
				ParseMode: "HTML",
			})
			if err != nil {
				log.Err(err).Msg("Failed to send notification to admin")
				continue
			}

			log.Info().Msg("Notification sent to admin")
			previousHash = currentHash

		default:
			log.Warn().Interface("response_type", response).Msg("Unexpected response type")
		}
	}
}

func autoBuyGifts(client *telegram.Client, gifts []telegram.StarGift) {
	//const maxTotalCount = 500000
	const maxTotalCount = 50000
	const maxBuyCount = 100

	// Find limited gifts with total count <= 50,000
	type targetGift struct {
		ID          int64
		Updgradable bool
		Total       int32
	}

	var targetGifts []targetGift
	for _, gift := range gifts {
		gft := gift.(*telegram.StarGiftObj)
		if gft.AvailabilityTotal > 0 && gft.AvailabilityTotal <= maxTotalCount && gft.AvailabilityRemains > 0 {
			targetGifts = append(targetGifts, targetGift{
				ID:          gft.ID,
				Total:       gft.AvailabilityTotal,
				Updgradable: gft.CanUpgrade,
			})
		}
	}

	if len(targetGifts) == 0 {
		return
	}

	// Sort gifts by availability total (ascending - start with smallest)
	sort.Slice(targetGifts, func(i, j int) bool {
		return targetGifts[i].Total < targetGifts[j].Total
	})

	// Buy each gift sequentially until error, then move to next
	for _, gift := range targetGifts {
		log.Info().Int64("gift_id", gift.ID).Int32("total", gift.Total).Msg("Starting to buy gift")

		for range maxBuyCount {
			_, err := SendNewGift(client, client.Me().Username, gift.ID, false, "auto-buy")
			if err != nil {
				log.Err(err).Int64("gift_id", gift.ID).Msg("Failed to buy gift, moving to next")
				break
			}

			//go upgradeGift(client)

			log.Info().Int64("gift_id", gift.ID).Msg("Successfully bought gift")
			time.Sleep(time.Second)
		}
	}
}

func printLimitedGiftsInfo(gifts []telegram.StarGift) {
	type giftInfo struct {
		ID        int64
		Available int32
		Total     int32
	}

	var limitedGifts []giftInfo
	for _, gift := range gifts {
		gft := gift.(*telegram.StarGiftObj)
		if gft.AvailabilityTotal > 0 {
			limitedGifts = append(limitedGifts, giftInfo{
				ID:        gft.ID,
				Available: gft.AvailabilityRemains,
				Total:     gft.AvailabilityTotal,
			})
		}
	}

	// Sort by total count
	sort.Slice(limitedGifts, func(i, j int) bool {
		return limitedGifts[i].Total < limitedGifts[j].Total
	})

	fmt.Println("Limited gifts info:")
	for _, gift := range limitedGifts {
		fmt.Printf("Gift %d: %d/%d\n", gift.ID, gift.Available, gift.Total)
	}
}

func formatGiftUpdateMessage(gifts []telegram.StarGift) string {
	var availableGifts []string
	var soldOutGifts []string

	for _, gift := range gifts {
		gft := gift.(*telegram.StarGiftObj)

		if gft.AvailabilityRemains == 0 {
			continue
		}

		giftInfo := fmt.Sprintf("🎁 Gift %d\n"+
			"⭐ Stars: %d\n"+
			"📦 Available: %d/%d",
			gft.ID, gft.Stars, gft.AvailabilityRemains, gft.AvailabilityTotal)
		availableGifts = append(availableGifts, giftInfo)
	}

	if len(availableGifts) == 0 && len(soldOutGifts) == 0 {
		return ""
	}

	var message strings.Builder
	if len(availableGifts) > 0 {
		message.WriteString("✅ <b>Available Limited Gifts:</b>\n")
		for _, gift := range availableGifts {
			message.WriteString(gift)
			message.WriteString("\n\n")
		}
	}

	message.WriteString(fmt.Sprintf("⏰ <i>Updated at: %s</i>", time.Now().Format("15:04:05")))

	return message.String()
}
