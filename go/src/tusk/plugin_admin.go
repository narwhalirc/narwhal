package tusk

import (
	"github.com/lrstanley/girc"
	"strings"
)

// NarwhalAdminManager is our admin management plugin
var NarwhalAdminManager NarwhalAdminPlugin

// NarwhalAdminPlugin is our configuration for the Narwhal admin plugin
type NarwhalAdminPlugin struct {
	// Enabled determines whether to enable this functionality
	Enabled bool
}

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
			narwhalMessage := ParseMessage(user, strings.TrimSpace(e.Trailing))
			adminmanager.CommandIssuer(c, e, narwhalMessage) // Pass along to our command issuer
		}
	}
}

// CommandIssuer is our primary function for command and param handling
func (adminmanager *NarwhalAdminPlugin) CommandIssuer(c *girc.Client, e girc.Event, narwhalMessage NarwhalMessage) {
	cmd := narwhalMessage.Command
	params := narwhalMessage.Params
	hasGlobal := strings.HasPrefix(cmd, "global")

	// #region Global commands (not channel specific)

	if hasGlobal {
		nonGlobalCommand := strings.Replace(cmd, "global", "", -1) // Get the non-global equivelant for when we do per-user action across multiple channels

		for _, channel := range Config.Channels { // For each channel the bot is in
			specifiedKickParams := []string{channel}
			specifiedKickParams = append(specifiedKickParams, params...)

			adminmanager.CommandIssuer(c, e, NarwhalMessage{ // Issue a non-global command against this user for this specific command
				Command: nonGlobalCommand,
				Issuer:  narwhalMessage.Issuer,
				Params:  specifiedKickParams,
			})
		}
	}

	// #endregion

	// #region Channel-specific commands

	if !hasGlobal {
		if len(params) == 1 {
			c.Cmd.ReplyTo(e, "You must pass a #channel with this command.")
		} else if len(params) > 1 { // Multiple params
			channel := params[:1][0] // Designate the first argument as the channel
			users := params[1:]      // Designate all arguments after as users

			switch cmd {
			case "ban": // Ban
				KickUsers(c, channel, users)      // Kick the users before issuing ban
				NarwhalAutoKicker.AddUsers(users) // Add the users to Autokick
				BanUsers(c, channel, users)       // Ban the users
				break
			case "kick": // Kick
				KickUsers(c, channel, users)      // Kick the users
				NarwhalAutoKicker.AddUsers(users) // Add the users to Autokick
				break
			case "unban": // Unban
				NarwhalAutoKicker.RemoveUsers(users) // Remove the users from Autokick
				UnbanUsers(c, channel, users)        // Unban the users
				break
			case "unkick": // Unkick
				NarwhalAutoKicker.RemoveUsers(users) // Remove the users from Autokick
				break
			default:
			}
		}
	}

	// #endregion
}
