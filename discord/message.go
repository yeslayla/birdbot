package discord

import (
	"log"

	"github.com/bwmarrin/discordgo"
	"github.com/yeslayla/birdbot/core"
	"github.com/yeslayla/birdbot/persistence"
)

func (discord *Discord) WebhookSendMessage(channel *core.Channel, displayName string, message string) {

	webhookData, err := discord.db.GetDiscordWebhook(channel.ID)
	if err != nil {
		log.Printf("Error getting webhook from DB: %s", err)
		return
	}

	if webhookData == nil {
		webhook, err := discord.session.WebhookCreate(channel.ID, "BirdBot", "")
		if err != nil {
			log.Printf("Error creating webhook: %s", err)
			return
		}

		webhookData = &persistence.DBDiscordWebhook{
			ID:    webhook.ID,
			Token: webhook.Token,
		}

		if err := discord.db.SetDiscordWebhook(channel.ID, webhookData); err != nil {
			log.Fatalf("Error failed to store webhook in DB: %s", err)
			return
		}

	}

	if _, err = discord.session.WebhookExecute(webhookData.ID, webhookData.Token, false, &discordgo.WebhookParams{
		Content:  message,
		Username: displayName,
	}); err != nil {
		log.Printf("Failed to send message over webhook: %s", err)
	}

}
