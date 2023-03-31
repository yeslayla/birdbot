package common

type Component interface {
	Initialize(birdbot ComponentManager) error
}

// ComponentManager is the primary way for a component to interact with BirdBot
// by listening to events and committing actions
type ComponentManager interface {
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

	RegisterGameModule(ID string, plugin GameModule) error
}
