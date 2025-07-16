package main

import (
	"fmt"

	"github.com/amarnathcjd/gogram/telegram"
)

func NewClientWithSession(appID int32, appHash, session string) (*telegram.Client, error) {
	client, err := telegram.NewClient(telegram.ClientConfig{
		AppID:   appID,
		AppHash: appHash,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create telegram client: %w", err)
	}

	if _, err = client.Conn(); err != nil {
		return nil, fmt.Errorf("failed to conn telegram client: %w", err)
	}

	if _, err = client.ImportSession(session); err != nil {
		return nil, fmt.Errorf("failed to import telegram session: %w", err)
	}

	return client, nil
}

func SendNewGift(c *telegram.Client, toPeer any, giftId int64, upgradable bool, message ...string) (telegram.PaymentsPaymentResult, error) {
	userPeer, err := c.ResolvePeer(toPeer)
	if err != nil {
		return nil, err
	}

	inv := &telegram.InputInvoiceStarGift{
		Peer:           userPeer,
		GiftID:         giftId,
		IncludeUpgrade: upgradable,
		HideName:       false,
	}

	if len(message) > 0 {
		entites, textPart := c.FormatMessage(message[0], c.ParseMode())
		inv.Message = &telegram.TextWithEntities{
			Text:     textPart,
			Entities: entites,
		}
	}

	form, err := c.PaymentsGetPaymentForm(inv, &telegram.DataJson{})
	if err != nil {
		return nil, err
	}

	return c.PaymentsSendStarsForm(form.(*telegram.PaymentsPaymentFormStarGift).FormID, inv)
}
