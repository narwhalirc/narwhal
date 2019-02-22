package tusk

// This file contains the types pertaining to the Autokick plugin

// NarwhallAutoKickerConfig is our configuration for the Narwhal autokicker
type NarwhalAutoKickerConfig struct {
	// EnabledAutoban determines whether to enable the automatic banning of users which exceed our MinimumKickToBanCount
	EnabledAutoban bool `json:",omitempty"`

	// Hosts to kick. Matches from end.
	Hosts []string

	// MessageMatches is a list of messages that will result in kicks
	MessageMatches []string

	// MinimumKickToBanCount is a minimum amount of times a user should be kicked before being automatically banned. Only enforced when EnabledAutoban is set
	MinimumKickToBanCount int `json:",omitempty"`

	// Users to kick. Matches from beginning.
	Users []string
}

// NarwhalAutoKickerPlugin is our Autokick plugin
type NarwhalAutoKickerPlugin struct {
	// Tracker is a map of usernames to the amount of times they've been kicked
	Tracker map[string]int
}