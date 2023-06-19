package app

import (
	"log"
	"path/filepath"

	"github.com/yeslayla/birdbot-common/common"
)

type ComponentLoader struct {
	bot       *Bot
	configDir string
}

func NewComponentLoader(bot *Bot, configDir string) *ComponentLoader {
	return &ComponentLoader{
		bot:       bot,
		configDir: configDir,
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

func (loader *ComponentLoader) RegisterExternalChat(ID string, chat common.ExternalChatModule) error {
	loader.bot.chatHandlers[ID] = chat
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

func (loader *ComponentLoader) GetConfigPath(fileName string) string {
	return filepath.Join(loader.configDir, "birdbot", fileName)
}
