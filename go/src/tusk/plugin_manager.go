package tusk

// This file contains our currently basic plugin management functionality

// NarwhalPluginManager is our plugin manager
type NarwhalPluginManager struct{}

// IsEnabled will return if the plugin specified is enabled
func (pm *NarwhalPluginManager) IsEnabled(pluginName string) bool {
	return IsInStringArr(Config.Plugins.Enabled, pluginName)
}
