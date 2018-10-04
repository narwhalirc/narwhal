package tusk

import (
	"github.com/lrstanley/girc"
)

// Our structs and interfaces

// NarwhalConfig is our primary Narwhal configuration
type NarwhalConfig struct {
	// Network is the IRC network to connection to.
	Network string

	// Port is the port on the network we're connecting to. Likely 6667.
	Port int

	// User is the IRC Bot username
	User string

	// Name is the IRC Bot name
	Name string

	// FallbackNick is the IRC bot fallback nickname if first nick is registered to someone else
	FallbackNick string `toml:"FallbackNick,omitempty"`

	// Password is the IRC bot password for authentication
	Password string

	// Plugins is a list of plugin configurations
	Plugins NarwhalPluginsConfig `toml:"Plugins,omitempty"`

	// Channels is a list of channels to join
	Channels []string

	// Users is our users configuration
	Users NarwhalUsersConfig `toml:"Users,omitempty"`
}

// NarwhalMessage is a custom message
type NarwhalMessage struct {
	Command string
	Issuer  string
	Params  []string
}

// NarwhalPlugin is a plugin interface
type NarwhalPlugin interface {
	Parse(c *girc.Client, e girc.Event)
}

// NarwhalUsersConfig is our configuration for blacklisting users, administrative users, and autokicking
type NarwhalUsersConfig struct {
	// Admins is an array of users authorized to perform admin actions
	Admins []string

	// Blacklist is an array of users blacklisted from performing Plugins
	Blacklist []string
}

// NarwhalPluginsConfig is a list of command configurations
type NarwhalPluginsConfig struct {
	Admin    NarwhalAdminConfig
	AutoKick NarwhalAutoKickerConfig
	Song     NarwhalSongConfig
}
