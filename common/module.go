package common

type Module interface {
	Initialize(birdbot ModuleManager) error
}

// ModuleManager is the primary way for a module to interact with BirdBot
// by listening to events and committing actions
type ModuleManager interface {
	OnReady(func() error) error

	OnNotify(func(string) error) error

	// Event events
	OnEventCreate(func(Event) error) error
	OnEventDelete(func(Event) error) error
	OnEventUpdate(func(Event) error) error
	OnEventComplete(func(Event) error) error

	// Actions
	CreateEvent(event Event) error
	Notify(message string) error

	// Commands
	RegisterCommand(string, ChatCommandConfiguration, func(User, map[string]any) string)

	// Submodules
	RegisterExternalChat(channelID string, chat ExternalChatModule) error
}
