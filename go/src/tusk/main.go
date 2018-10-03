package tusk

import (
	"github.com/JoshStrobl/trunk"
	"github.com/lrstanley/girc"
	"os"
	"os/user"
)

// A Narwhal is no Narwhal without their tusk!

// Config is our Narwhal Config
var Config NarwhalConfig

// HomeDir is the current user's home directory
var HomeDir string

// WorkDir is our current working directory
var WorkDir string

func init() {
	var getUserErr error
	var getWdErr error
	var currentUser *user.User

	currentUser, getUserErr = user.Current() // Attempt to get the current user

	if getUserErr == nil { // If we successfully got the user
		HomeDir = currentUser.HomeDir
	} else {
		trunk.LogFatal("Failed to get the current user: " + getUserErr.Error())
	}

	WorkDir, getWdErr = os.Getwd() // Get the current working directory

	if getWdErr != nil { // If we failed to get the current working dir
		trunk.LogFatal("Failed to get the current working directory: " + getWdErr.Error())
	}
}

// NewTusk will create a new tusk for our Narwhal, but only one tusk is allowed at a time.
func NewTusk() {
	var newTuskErr error

	Config, newTuskErr = ReadConfig()

	if newTuskErr == nil { // Read our config
		ircConfig := girc.Config{
			Server: Config.Network,
			Port:   Config.Port,
			Name:   Config.Name,
			Nick:   Config.User,
			User:   Config.User,
			SASL: &girc.SASLPlain{
				User: Config.User,
				Pass: Config.Password,
			},
		}

		client := girc.New(ircConfig)
		client.Handlers.Add(girc.CONNECTED, OnConnected) // On CONNECTED, trigger OnConnected
		client.Handlers.Add(girc.JOIN, Parser)           // On JOIN, trigger our Parser
		client.Handlers.Add(girc.INVITE, OnInvite)       // On INVITE, trigger OnInvite
		client.Handlers.Add(girc.PRIVMSG, Parser)        // On PRIVMSG, trigger our Parser

		if newTuskErr = client.Connect(); newTuskErr != nil { // Failed during run
			trunk.LogFatal("Failed to run client: " + newTuskErr.Error())
		}
	} else {
		trunk.LogFatal("Failed to read or parse config: " + newTuskErr.Error())
	}
}
