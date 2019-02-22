package tusk

// This file contains the types pertaining to the Slap plugin

// NarwhalSlapConfig is our configuration for the Narwhal autokicker
type NarwhalSlapConfig struct {
	// CustomActions is a list of custom actions on how to slap a user
	CustomActions []string
}

// NarwhalSlapPlugin is our slap plugin
type NarwhalSlapPlugin struct {
	Objects []string
}