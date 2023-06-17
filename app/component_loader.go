package app

import (
	"log"

	"github.com/yeslayla/birdbot/common"
)

type ComponentLoader struct {
	bot *Bot
}

func NewComponentLoader(bot *Bot) *ComponentLoader {
	return &ComponentLoader{
		bot: bot,
	}
}

func (loader *ComponentLoader) LoadComponent(component common.Module) {
	if err := component.Initialize(loader); err != nil {
		log.Print("Failed to load component: ", err)
	}
}

func (loader *ComponentLoader) OnReady(handler func() error) error {
	loader.bot.onReadyHandlers = append(loader.bot.onReadyHandlers, handler)
	return nil
}

func (loader *ComponentLoader) OnNotify(handler func(string) error) error {
	loader.bot.onNotifyHandlers = append(loader.bot.onNotifyHandlers, handler)
	return nil
}

func (loader *ComponentLoader) OnEventCreate(handler func(common.Event) error) error {
	loader.bot.onEventCreatedHandlers = append(loader.bot.onEventCreatedHandlers, handler)
	return nil
}
func (loader *ComponentLoader) OnEventDelete(handler func(common.Event) error) error {
	loader.bot.onEventDeletedHandlers = append(loader.bot.onEventDeletedHandlers, handler)
	return nil
}
func (loader *ComponentLoader) OnEventUpdate(handler func(common.Event) error) error {
	loader.bot.onEventUpdatedHandlers = append(loader.bot.onEventUpdatedHandlers, handler)
	return nil
}
func (loader *ComponentLoader) OnEventComplete(handler func(common.Event) error) error {
	loader.bot.onEventCompletedHandlers = append(loader.bot.onEventCompletedHandlers, handler)
	return nil
}

func (loader *ComponentLoader) RegisterExternalChat(channelID string, chat common.ExternalChatModule) error {
	if _, ok := loader.bot.channelChats[channelID]; !ok {
		loader.bot.channelChats[channelID] = []common.ExternalChatModule{}
	}

	loader.bot.channelChats[channelID] = append(loader.bot.channelChats[channelID], chat)

	return nil
}

func (loader *ComponentLoader) CreateEvent(event common.Event) error {
	return loader.bot.Session.CreateEvent(event)
}

func (loader *ComponentLoader) Notify(message string) error {
	loader.bot.Notify(message)
	return nil
}

func (loader *ComponentLoader) RegisterCommand(name string, config common.ChatCommandConfiguration, handler func(common.User, map[string]any) string) {
	loader.bot.Session.RegisterCommand(name, config, handler)
}
