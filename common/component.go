package common

type Component interface {
	Initialize(birdbot ComponentManager) error
}

type ComponentManager interface {
	OnReady(func() error) error

	OnNotify(func(string) error) error

	// Event events
	OnEventCreate(func(Event) error) error
	OnEventDelete(func(Event) error) error
	OnEventUpdate(func(Event) error) error
	OnEventComplete(func(Event) error) error

	RegisterGameModule(ID string, plugin GameModule) error

	CreateEvent(event Event) error
	Notify(message string) error
}
