package app

import (
	"log"
	"os"
	"plugin"

	"github.com/yeslayla/birdbot/common"
)

// LoadPlugin loads a plugin and returns its component if successful
func LoadPlugin(pluginPath string) common.Component {

	plug, err := plugin.Open(pluginPath)
	if err != nil {
		log.Printf("Failed to load plugin '%s': %s", pluginPath, err)
		return nil
	}

	// Lookup component symbol
	sym, err := plug.Lookup("Component")
	if err != nil {
		log.Printf("Failed to load plugin '%s': failed to get Component: %s", pluginPath, err)
		return nil
	}

	// Validate component type
	var component common.Component
	component, ok := sym.(common.Component)
	if !ok {
		log.Printf("Failed to load plugin '%s': Plugin component does not properly implement interface!", pluginPath)
	}

	return component
}

// LoadPlugins loads all plugins and componenets in a directory
func LoadPlugins(directory string) []common.Component {

	paths, err := os.ReadDir(directory)
	if err != nil {
		log.Printf("Failed to load plugins: %s", err)
		return []common.Component{}
	}

	var components []common.Component = make([]common.Component, 0)
	for _, path := range paths {
		if path.IsDir() {
			continue
		}

		if comp := LoadPlugin(path.Name()); comp != nil {
			components = append(components, comp)
		}
	}

	return components
}
