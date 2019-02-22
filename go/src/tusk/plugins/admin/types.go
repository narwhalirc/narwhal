package tusk

// This file contains the types pertaining to the Admin plugin

// NarwhalAdminConfig is our configuration for the Narwhal admin plugin
type NarwhalAdminConfig struct {
	// DisabledCommands is a list of admin commands to disable
	DisabledCommands []string
}

// NarwhalAdminPlugin is our Admin plugin
type NarwhalAdminPlugin struct{}