package discord

import (
	"log"

	"github.com/bwmarrin/discordgo"
	"github.com/yeslayla/birdbot/core"
	"github.com/yeslayla/birdbot/persistence"
)

const WebhookName = "BirdBot"

// RefreshWebhookState refreshes the state of all webhooks
func (discord *Discord) RefreshWebhookState() {
	channels, err := discord.session.GuildChannels(discord.guildID)
	if err != nil {
		log.Printf("Error getting channels: %s", err)
		return
	}

	for _, channel := range channels {
		webhookData, err := discord.db.GetDiscordWebhook(channel.ID)
		if err != nil {
			log.Printf("Error getting webhook from DB: %s", err)
			return
		}

		if webhookData == nil {
			continue
		}

		_, err = discord.session.WebhookEdit(webhookData.ID, WebhookName, discord.GetAvatarBase64(NewUser(discord.session.State.User)), channel.ID)
		if err != nil {
			log.Printf("Error updating webhook: %s", err)
		}
	}
}

// WebhookSendMessage sends a message to a channel using a webhook
func (discord *Discord) WebhookSendMessage(channel *core.Channel, displayName string, message string) {

	webhookData, err := discord.db.GetDiscordWebhook(channel.ID)
	if err != nil {
		log.Printf("Error getting webhook from DB: %s", err)
		return
	}

	if webhookData == nil {
		webhookAvatar := discord.GetAvatarBase64(NewUser(discord.session.State.User))

		webhook, err := discord.session.WebhookCreate(channel.ID, WebhookName, webhookAvatar)
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
