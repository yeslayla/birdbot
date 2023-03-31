package app

import (
	"fmt"
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

func (loader *ComponentLoader) LoadComponent(component common.Component) {
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

func (loader *ComponentLoader) RegisterGameModule(ID string, plugin common.GameModule) error {
	return fmt.Errorf("unimplemented")
}

func (loader *ComponentLoader) CreateEvent(event common.Event) error {
	return loader.bot.Session.CreateEvent(event)
}

func (loader *ComponentLoader) Notify(message string) error {
	loader.bot.Notify(message)
	return nil
}
