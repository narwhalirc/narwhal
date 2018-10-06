package tusk

import (
	"github.com/lrstanley/girc"
	"strings"
)

// NarwhalAdminManager is our admin management plugin
var NarwhalAdminManager NarwhalAdminPlugin

// NarwhalAdminConfig is our configuration for the Narwhal admin plugin
type NarwhalAdminConfig struct {
	// DisabledCommands is a list of admin commands to disable
	DisabledCommands []string
	// Enabled determines whether to enable this functionality
	Enabled bool
}

// NarwhalAdminPlugin is our Admin plugin
type NarwhalAdminPlugin struct{}

func (adminmanager *NarwhalAdminPlugin) Parse(c *girc.Client, e girc.Event, m NarwhalMessage) {
	if len(Config.Users.Admins) > 0 { // If there are any admins set
		var userIsAdmin bool

		for _, admin := range Config.Users.Admins { // For each listed admin
			userIsAdmin = Matches(admin, m.Issuer) // Check for a match against the username

			if !userIsAdmin { // User not an admin by nick
				userIsAdmin = Matches(admin, m.Host) // Check for a match against the host (more secure in some cases)
			}

			if userIsAdmin { // If this is a match
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

	if !hasGlobal && !IsInStringArr(Config.Plugins.Admin.DisabledCommands, cmd) { // Not a global command and not disabled
		switch cmd {
		case "addhost": // Add Host to AutoKick Hosts
			NarwhalAutoKicker.AddHost(m.MessageNoCmd) // Add the host
			break
		case "addkicker": // Add Kicker without kick attempt
			NarwhalAutoKicker.AddUsers(params) // Add the users to Autokick
			break
		case "addmsg": // Add Message to our MessageMatches
			NarwhalAutoKicker.AddMessage(m.MessageNoCmd) // Add the message to Autokick MessageMatches
			break
		case "ban": // Ban
			KickUsers(c, eventChannel, params) // Kick the users before issuing ban
			NarwhalAutoKicker.AddUsers(params) // Add the users to Autokick
			BanUsers(c, eventChannel, params)  // Ban the users
			break
		case "blacklist": // Blacklist
			adminmanager.Blacklist(params) // Blacklist the user(s)
			break
		case "leave": // Leave a channel
			c.Cmd.Action(m.Channel, "has far more important things to attend do!")
			c.Cmd.Part(m.Channel)
		case "kick": // Kick
			KickUsers(c, eventChannel, params) // Kick the users
			NarwhalAutoKicker.AddUsers(params) // Add the users to Autokick
			break
		case "proclaim": // Proclamation
			proclamationMessage := "Behold, I am your robot narwhal overlord. Bow before me, puny hoooomans, or I shall unleash source code upon you."
			c.Cmd.Reply(e, proclamationMessage)
			c.Cmd.Action(m.Channel, "means to say to visit https://github.com/JoshStrobl/narwhal")
			break
		case "removehost": // Remove Host from AutoKick Hosts
			NarwhalAutoKicker.RemoveHost(m.MessageNoCmd) // Remove the host
			break
		case "removemsg": // Remove Message from our MessageMatches
			NarwhalAutoKicker.RemoveMessage(m.MessageNoCmd) // Remove the message from Autokick MessageMatches
			break
		case "rmblacklist": // Remove user(s) from Blacklist
			adminmanager.RemoveFromBlacklist(params)
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

// Blacklist will add users to the blacklist
func (adminmanager *NarwhalAdminPlugin) Blacklist(users []string) {
	Config.Users.Blacklist = append(Config.Users.Blacklist, users...) // Add users
	Config.Users.Blacklist = DeduplicateList(Config.Users.Blacklist)
	SaveConfig()
}

// RemoveFromBlacklist will remove users from the blacklist
func (adminmanager *NarwhalAdminPlugin) RemoveFromBlacklist(users []string) {
	Config.Users.Blacklist = RemoveFromStringArr(Config.Users.Blacklist, users) // Remove the user
	SaveConfig()
}
