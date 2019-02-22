package tusk

import (
	"github.com/JoshStrobl/trunk"
	"github.com/lrstanley/girc"
	"strings"
)

// NarwhalAutoKicker is our autokicker plugin
var NarwhalAutoKicker NarwhalAutoKickerPlugin

func init() {
	NarwhalAutoKicker = NarwhalAutoKickerPlugin{
		Tracker: make(map[string]int),
	}
}

func (autokicker *NarwhalAutoKickerPlugin) Parse(c *girc.Client, e girc.Event, m NarwhalMessage) {
	var userShouldBeKicked bool

	// #region Hosts Kick List Check

	if len(Config.Plugins.AutoKick.Hosts) > 0 { // If we have a Hosts list
		for _, kickHost := range Config.Plugins.AutoKick.Hosts {
			if strings.HasPrefix(kickHost, "*!*@") { // If we're checking for any username with this host
				reducedHost := strings.Replace(kickHost, "*!*@", "", -1)

				if reducedHost == kickHost { // Exact host match
					userShouldBeKicked = true
					break
				}
			} else { // If we're looking for ident+host
				userIdent := m.Issuer + "@" + m.Host

				if userIdent == m.Host { // If the user ident matches host
					userShouldBeKicked = true
					break
				}
			}
		}
	}

	// #endregion

	// #region Message Matching

	if !userShouldBeKicked { // If we haven't yet determined to kick
		if len(Config.Plugins.AutoKick.MessageMatches) > 0 { // If we have a Messages list
			for _, match := range Config.Plugins.AutoKick.MessageMatches {
				userShouldBeKicked = Matches(match, m.Message) // Check if string matches our requirements

				if userShouldBeKicked {
					break
				}
			}
		}
	}

	// #endregion

	// #region Users Kick List Check

	if !userShouldBeKicked { // If we haven't yet determined to kick
		if len(Config.Plugins.AutoKick.Users) > 0 { // If we have a Kick list
			for _, kickUser := range Config.Plugins.AutoKick.Users {
				userShouldBeKicked = Matches(kickUser, m.Issuer) // Check if string matches our requirements

				if userShouldBeKicked {
					break
				}
			}
		}
	}

	// #endregion

	if userShouldBeKicked && !IsAdmin(m) { // If the user is indicated to be kicked and is not an admin
		trunk.LogInfo("AutoKick triggered. Kicking " + m.Issuer)
		kickCount := 0

		if Config.Plugins.AutoKick.EnabledAutoban { // If we've enabled autoban
			var exists bool
			kickCount, exists = autokicker.Tracker[m.Issuer] // Get the current kickCount if it exists

			if exists {
				kickCount++ // Increment the counter
			} else {
				kickCount = 1 // Set to 1
			}

			autokicker.Tracker[m.Issuer] = kickCount // Update our tracker
		}

		for _, channel := range Config.Channels { // For each channel
			KickUser(c, channel, m.Issuer)

			if Config.Plugins.AutoKick.EnabledAutoban && (kickCount >= Config.Plugins.AutoKick.MinimumKickToBanCount) { // User has been kicked more than our minimum
				BanUser(c, channel, m.Issuer)
			}
		}
	}
}

// AddHost will add the specified host to the AutoKick Hosts list, if they aren't already added
func (autokicker *NarwhalAutoKickerPlugin) AddHost(host string) {
	Config.Plugins.AutoKick.Hosts = append(Config.Plugins.AutoKick.Hosts, host)
	Config.Plugins.AutoKick.Hosts = DeduplicateList(Config.Plugins.AutoKick.Hosts)
	SaveConfig()
}

// AddMessage will add the specified message to the AutoKick MessageMatches list, if they aren't already added
func (autokicker *NarwhalAutoKickerPlugin) AddMessage(message string) {
	Config.Plugins.AutoKick.MessageMatches = append(Config.Plugins.AutoKick.MessageMatches, message) // Add the msg
	Config.Plugins.AutoKick.MessageMatches = DeduplicateList(Config.Plugins.AutoKick.MessageMatches) // Deduplicate messages and set to MessageMatches
	SaveConfig()
}

// AddUsers will add the specified users to the AutoKick Users list, if they aren't already added
func (autokicker *NarwhalAutoKickerPlugin) AddUsers(users []string) {
	for _, requestedAddUser := range users {
		Config.Plugins.AutoKick.Users = append(Config.Plugins.AutoKick.Users, requestedAddUser)
	}

	Config.Plugins.AutoKick.Users = DeduplicateList(Config.Plugins.AutoKick.Users) // Deduplicate users and set to AutoKick Users
	SaveConfig()
}

// RemoveHost will remove the specified host from the AutoKick Hosts list, if they are added
func (autokicker *NarwhalAutoKickerPlugin) RemoveHost(host string) {
	hosts := []string{host}
	Config.Plugins.AutoKick.Hosts = RemoveFromStringArr(Config.Plugins.AutoKick.Hosts, hosts)
	Config.Plugins.AutoKick.Hosts = DeduplicateList(Config.Plugins.AutoKick.Hosts)
	SaveConfig()
}

// RemoveMessage will remove the specified message from the AutoKick MessageMatches list, if they are added
func (autokicker *NarwhalAutoKickerPlugin) RemoveMessage(message string) {
	message = strings.TrimSpace(message)
	messages := []string{message}
	Config.Plugins.AutoKick.MessageMatches = RemoveFromStringArr(Config.Plugins.AutoKick.MessageMatches, messages) // Remove the specified items from the string array
	Config.Plugins.AutoKick.MessageMatches = DeduplicateList(Config.Plugins.AutoKick.MessageMatches)               // Deduplicate MessageMatches
	SaveConfig()
}

// RemoveUsers will remove the specified users from the AutoKick Users list, if they are added
func (autokicker *NarwhalAutoKickerPlugin) RemoveUsers(users []string) {
	for _, user := range users { // Users we're wanting to remove
		delete(autokicker.Tracker, user) // Delete user from Tracker if they exist
	}

	Config.Plugins.AutoKick.Users = RemoveFromStringArr(Config.Plugins.AutoKick.Users, users) // Remove the specified items from the string array
	Config.Plugins.AutoKick.Users = DeduplicateList(Config.Plugins.AutoKick.Users)            // Deduplicate users and set to AutoKick Users
	SaveConfig()
}
