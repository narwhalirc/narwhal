package tusk

import (
	"github.com/lrstanley/girc"
	"strings"
)

// NarwhalAdminManager is our admin management plugin
var NarwhalAdminManager NarwhalAdminPlugin

// NarwhalAdminConfig is our configuration for the Narwhal admin plugin
type NarwhalAdminConfig struct {
	// Enabled determines whether to enable this functionality
	Enabled bool
}

// NarwhalAdminPlugin is our Admin plugin
type NarwhalAdminPlugin struct{}

func (adminmanager *NarwhalAdminPlugin) Parse(c *girc.Client, e girc.Event) {
	if len(Config.Users.Admins) > 0 { // If there are any admins set
		var userIsAdmin bool
		user := e.Source.Name

		if user == "" { // User is somehow empty
			user = e.Source.Ident // Change to using Ident
		}

		for _, admin := range Config.Users.Admins { // For each listed admin
			if user == admin { // If this is a match
				userIsAdmin = true
				break
			}
		}

		if userIsAdmin { // If the user issuing a command is an admin
			narwhalMessage := ParseMessage(e)
			adminmanager.CommandIssuer(c, narwhalMessage) // Pass along to our command issuer
		}
	}
}

// CommandIssuer is our primary function for command and param handling
func (adminmanager *NarwhalAdminPlugin) CommandIssuer(c *girc.Client, m NarwhalMessage) {
	eventChannel := m.Channel
	cmd := m.Command
	params := m.Params
	hasGlobal := strings.HasPrefix(cmd, "global")

	// #region Global commands (not channel specific)

	if hasGlobal {
		nonGlobalCommand := strings.Replace(cmd, "global", "", -1) // Get the non-global equivelant for when we do per-user action across multiple channels

		for _, channel := range Config.Channels { // For each channel the bot is in
			adminmanager.CommandIssuer(c, NarwhalMessage{ // Issue a non-global command against this user for this specific command
				Channel: channel, // Change our channel to this one
				Command: nonGlobalCommand,
				Issuer:  m.Issuer,
				Params:  m.Params,
			})
		}
	}

	// #endregion

	// #region Channel-specific commands

	if !hasGlobal {
		switch cmd {
		case "ban": // Ban
			KickUsers(c, eventChannel, params) // Kick the users before issuing ban
			NarwhalAutoKicker.AddUsers(params) // Add the users to Autokick
			BanUsers(c, eventChannel, params)  // Ban the users
			break
		case "kick": // Kick
			KickUsers(c, eventChannel, params) // Kick the users
			NarwhalAutoKicker.AddUsers(params) // Add the users to Autokick
			break
		case "unban": // Unban
			NarwhalAutoKicker.RemoveUsers(params) // Remove the users from Autokick
			UnbanUsers(c, eventChannel, params)   // Unban the users
			break
		case "unkick": // Unkick
			NarwhalAutoKicker.RemoveUsers(params) // Remove the users from Autokick
			break
		case "welcome": // Welcome message

		default:
		}
	}

	// #endregion
}
