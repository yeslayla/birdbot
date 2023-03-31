package app

import (
	"log"
	"os"
	"plugin"

	"github.com/yeslayla/birdbot/common"
)

type PluginLoader struct{}

func NewPluginLoader() PluginLoader {
	return PluginLoader{}
}

func (loader PluginLoader) LoadPlugin(pluginPath string) common.Component {
	plug, err := plugin.Open(pluginPath)
	if err != nil {
		log.Printf("Failed to load plugin '%s': %s", pluginPath, err)
		return nil
	}

	sym, err := plug.Lookup("Component")
	if err != nil {
		log.Printf("Failed to load plugin '%s': failed to get Component: %s", pluginPath, err)
		return nil
	}

	var component common.Component
	component, ok := sym.(common.Component)
	if !ok {
		log.Printf("Failed to load plugin '%s': Plugin component does not properly implement interface!", pluginPath)
	}

	return component
}

func (loader PluginLoader) LoadPlugins(directory string) []common.Component {

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

		if comp := loader.LoadPlugin(path.Name()); comp != nil {
			components = append(components, comp)
		}
	}

	return components
}
