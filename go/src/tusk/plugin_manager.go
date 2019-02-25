package tusk

import (
	"github.com/JoshStrobl/trunk"
	"os"
	"path/filepath"
	"plugin"
	"strings"
)

// This file contains our currently basic plugin management functionality

// NarwhalPluginManager is our plugin manager
type NarwhalPluginManager struct {
	Modules map[string]plugin.Symbol
}

// IsEnabled will return if the plugin specified is enabled
func (pm *NarwhalPluginManager) IsEnabled(pluginName string) bool {
	return IsInStringArr(Config.Plugins.Enabled, pluginName)
}

// LoadPlugins is responsible for loading any plugins from our modules directory
func (pm *NarwhalPluginManager) LoadPlugins() error {
	var loadPluginsErr error

	for _, paths := range Paths { // For each path in paths
		modulesPath := filepath.Join(paths, "modules") // Add modules dir

		if modDir, modOpenErr := os.Open(modulesPath); modOpenErr == nil { // If there wasn't an error opening this directory
			if directoryItems, dirReadErr := modDir.Readdir(-1); dirReadErr == nil { // If we successfully read the contents of the directory
				if len(directoryItems) > 0 { // Have content
					for _, dirItem := range directoryItems { // For each directory item
						fileName := dirItem.Name()
						fileExt := filepath.Ext(fileName)

						if !dirItem.IsDir() && (fileExt == ".so") { // If this may be a .so file
							pluginName := strings.Replace(fileName, fileExt, "", -1) // Remove .so

							if _, alreadyAdded := pm.Modules[pluginName]; !alreadyAdded { // If we haven't already added this plugin
								if plugin, pluginOpenErr := plugin.Open(filepath.Join(modulesPath, fileName)); pluginOpenErr == nil { // Attempt file open
									if parseFunc, lookupErr := plugin.Lookup("Parse"); lookupErr == nil { // If we successfully looked up the Parse func symbol for this plugin
										trunk.LogSuccess("Added plugin: " + pluginName)
										pm.Modules[pluginName] = parseFunc
									} else {
										loadPluginsErr = lookupErr
										break
									}
								} else {
									loadPluginsErr = pluginOpenErr
									break
								}
							}
						}
					}
				}
			} else {
				loadPluginsErr = dirReadErr
			}
		}

		if loadPluginsErr != nil {
			break
		}
	}

	return loadPluginsErr
}