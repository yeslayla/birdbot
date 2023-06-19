package app

import (
	"github.com/yeslayla/birdbot-common/common"
	"github.com/yeslayla/birdbot/core"
)

type ExternalChatManager struct {
	chat    common.ExternalChatModule
	channel *core.Channel
	bot     *Bot
}

func (manager *ExternalChatManager) SendMessage(user string, message string) {
	manager.bot.Session.WebhookSendMessage(manager.channel, user, message)
}

func (app *Bot) InitalizeExternalChat(channel *core.Channel, chat common.ExternalChatModule) {
	manager := &ExternalChatManager{
		channel: channel,
		chat:    chat,
		bot:     app,
	}

	manager.chat.Initialize(manager)
}
