package tusk

import (
	"github.com/JoshStrobl/trunk"
	"github.com/lrstanley/girc"
	"strings"
)

// adminCommands is a list of admin commands
var adminCommands []string

// fixedIgnoreUsers is a list of users we'll ignore messages from no matter what
var fixedIgnoreUsers []string

// fullBlacklist is the combined list of fixedIgnoreUsers and configured blacklist
var fullBlacklist []string

// ignoreCommands is a list of numerical command codes to ignore
var ignoreCommands []string

func init() {
	fixedIgnoreUsers = []string{
		"freenode-connect",
	}

	fullBlacklist = fixedIgnoreUsers

	if len(Config.Users.Blacklist) > 0 { // If there are items in our blacklist
		fullBlacklist = append(fullBlacklist, Config.Users.Blacklist...)
	}

	ignoreCommands = []string{
		"002", // RPL_YOURHOST
		"003", // RPL_CREATED
		"004", // RPL_MYINFO
		"005", // RPL_BOUNCE
		"251", // RPL_LUSERCLIENT
		"252", // RPL_LUSEROP
		"253", // RPL_LUSERUNKNOWN
		"254", // RPL_LUSERCHANNELS
		"255", // RPL_LUSERME
		"265", // RPL_LOCALUSERS
		"266", // RPL_GLOBALUSERS
		"331", // RPL_NOTOPIC
		"332", // RPL_TOPIC
		"333", // RPL_TOPICWHOTIME
		"372", // RPL_MOTD
		"375", // RPL_MOTDSTART
		"376", // RPL_ENDOFMOTD
	}
}

// OnConnected will handle connection to an IRC network
func OnConnected(c *girc.Client, e girc.Event) {
	trunk.LogSuccess("Successfully connected to " + Config.Network + " as " + Config.User)
	if len(Config.Channels) > 0 { // If we have channels set to join
		for _, channel := range Config.Channels { // For each channel to join
			c.Cmd.Join(channel)
			c.Cmd.Mode(channel, "+o") // Attempt to op self
		}
	}
}

// OnInvite will handle a request to invite an IRC channel
func OnInvite(c *girc.Client, e girc.Event) {
	channel := strings.TrimSpace(e.Trailing)
	trunk.LogInfo("Received invite to " + channel + ". Joining.")
	c.Cmd.Join(channel)

	Config.Channels = append(Config.Channels, channel)
	SaveConfig()
}

// Parser will handle the majority of incoming messages, user joins, etc.
func Parser(c *girc.Client, e girc.Event) {
	m := ParseMessage(e)

	var ignoreMessage bool
	command := e.Command

	for _, ignoreCommand := range ignoreCommands { // For each ignore command
		if ignoreCommand == command {
			ignoreMessage = true
			break
		}
	}

	if !ignoreMessage {
		var userInBlacklist bool

		for _, blacklistUser := range fullBlacklist { // For each user
			kickUserWithoutSuffix := strings.Replace(blacklistUser, "*", "", -1)

			if m.Issuer == blacklistUser { // If the user is in the blacklist
				userInBlacklist = true
				break
			} else if strings.HasPrefix(m.Issuer, kickUserWithoutSuffix) { // If the username begins with this kickUser
				userInBlacklist = true
				break
			}
		}

		if Config.Plugins.AutoKick.Enabled { // AutoKick enabled
			NarwhalAutoKicker.Parse(c, e, m) // Run through auto-kicker first
		}

		if !userInBlacklist {
			trunk.LogInfo("Allowed: " + m.Issuer)
			trunk.LogInfo("Received: " + m.Message)
			trunk.LogInfo("Host: " + m.Host)
			trunk.LogInfo("Possible Channel: " + m.Channel)

			if Config.Plugins.Admin.Enabled { // Admin Management enabled
				NarwhalAdminManager.Parse(c, e, m) // Run through management
			}

			if Config.Plugins.Song.Enabled { // Song enabled
				NarwhalSong.Parse(c, e, m) // Run through song
			}

			if Config.Plugins.Slap.Enabled { // Slap enabled
				NarwhalSlap.Parse(c, e, m) // Run through slap
			}
		}
	}
}
