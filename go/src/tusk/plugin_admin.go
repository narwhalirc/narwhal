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

func (adminmanager *NarwhalAdminPlugin) Parse(c *girc.Client, e girc.Event, m NarwhalMessage) {
	if len(Config.Users.Admins) > 0 { // If there are any admins set
		var userIsAdmin bool

		for _, admin := range Config.Users.Admins { // For each listed admin
			if m.Issuer == admin { // If this is a match
				userIsAdmin = true
				break
			}
		}

		if userIsAdmin { // If the user issuing a command is an admin
			adminmanager.CommandIssuer(c, e, m) // Pass along to our command issuer
		}
	}
}

// CommandIssuer is our primary function for command and param handling
func (adminmanager *NarwhalAdminPlugin) CommandIssuer(c *girc.Client, e girc.Event, m NarwhalMessage) {
	eventChannel := m.Channel
	cmd := m.Command
	params := m.Params
	hasGlobal := strings.HasPrefix(cmd, "global")

	// #region Global commands (not channel specific)

	if hasGlobal {
		nonGlobalCommand := strings.Replace(cmd, "global", "", -1) // Get the non-global equivelant for when we do per-user action across multiple channels

		for _, channel := range Config.Channels { // For each channel the bot is in
			adminmanager.CommandIssuer(c, e, NarwhalMessage{ // Issue a non-global command against this user for this specific command
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
		case "addkicker": // Add Kicker without kick attempt
			NarwhalAutoKicker.AddUsers(params) // Add the users to Autokick
			break
		case "addmsg": // Add Message to our MessageMatches
			msg := strings.Replace(m.Message, "."+m.Command, " ", -1) // Remove command from our whole message to get entire message to add
			NarwhalAutoKicker.AddMessage(msg)                         // Add the message to Autokick MessageMatches
			break
		case "ban": // Ban
			KickUsers(c, eventChannel, params) // Kick the users before issuing ban
			NarwhalAutoKicker.AddUsers(params) // Add the users to Autokick
			BanUsers(c, eventChannel, params)  // Ban the users
			break
		case "kick": // Kick
			KickUsers(c, eventChannel, params) // Kick the users
			NarwhalAutoKicker.AddUsers(params) // Add the users to Autokick
			break
		case "proclaim": // Proclamation
			proclamationMessage := "Behold, I am your robot narwhal overlord. Bow before me, puny hoooomans, or I shall unleash source code upon you."
			c.Cmd.Reply(e, proclamationMessage)
			c.Cmd.Action(m.Channel, "means to say to visit https://github.com/JoshStrobl/narwhal")
			break
		case "removemsg": // Remove Message from our MessageMatches
			msg := strings.Replace(m.Message, "."+m.Command, " ", -1) // Remove command from our whole message to get entire message to add
			NarwhalAutoKicker.RemoveMessage(msg)                      // Remove the message from Autokick MessageMatches
			break
		case "unban": // Unban
			NarwhalAutoKicker.RemoveUsers(params) // Remove the users from Autokick
			UnbanUsers(c, eventChannel, params)   // Unban the users
			break
		case "unkick": // Unkick
			NarwhalAutoKicker.RemoveUsers(params) // Remove the users from Autokick
			break
		default:
		}
	}

	// #endregion
}
