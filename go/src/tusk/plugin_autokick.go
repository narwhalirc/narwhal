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

// NarwhallAutoKickerConfig is our configuration for the Narwhal autokicker
type NarwhalAutoKickerConfig struct {
	// Enabled determines whether to enable this functionality
	Enabled bool

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

func (autokicker *NarwhalAutoKickerPlugin) Parse(c *girc.Client, e girc.Event) {
	msg := strings.TrimSpace(e.Trailing)
	user := e.Source.Name
	host := e.Source.Host

	if user == "" { // User is somehow empty
		user = e.Source.Ident // Change to using Ident
	}

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
				userIdent := e.Source.Ident + "@" + host

				if userIdent == host { // If the user ident matches host
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
				if msg == match { // If this is an exact match
					userShouldBeKicked = true
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
				if strings.HasSuffix(kickUser, "*") { // If we should not be doing exact match
					kickUserWithoutSuffix := strings.Replace(kickUser, "*", "", -1)
					if strings.HasPrefix(user, kickUserWithoutSuffix) { // If the username begins with this kickUser
						userShouldBeKicked = true
						break
					} else if user == kickUserWithoutSuffix { // Identical match
						userShouldBeKicked = true
						break
					}
				} else { // If we should be doing an exact match
					if user == kickUser { // If the user should be kicked
						userShouldBeKicked = true
						break
					}
				}
			}
		}
	}

	// #endregion

	if userShouldBeKicked {
		trunk.LogInfo("AutoKick triggered. Kicking " + user)
		kickCount := 0

		if Config.Plugins.AutoKick.EnabledAutoban { // If we've enabled autoban
			var exists bool
			kickCount, exists = autokicker.Tracker[user] // Get the current kickCount if it exists

			if exists {
				kickCount++ // Increment the counter
			} else {
				kickCount = 1 // Set to 1
			}

			autokicker.Tracker[user] = kickCount // Update our tracker
		}

		for _, channel := range Config.Channels { // For each channel
			KickUser(c, channel, user)

			if Config.Plugins.AutoKick.EnabledAutoban && (kickCount > Config.Plugins.AutoKick.MinimumKickToBanCount) { // User has been kicked more than our minimum
				BanUser(c, channel, user)
			}
		}
	}
}

// AddUsers will add the specified users to the AutoKick Users list, if they aren't already added
func (autokicker *NarwhalAutoKickerPlugin) AddUsers(users []string) {
	for _, requestedAddUser := range users {
		Config.Plugins.AutoKick.Users = append(Config.Plugins.AutoKick.Users, requestedAddUser)
	}

	Config.Plugins.AutoKick.Users = DeduplicateList(Config.Plugins.AutoKick.Users) // Deduplicate users and set to AutoKick Users

	if saveErr := SaveConfig(); saveErr != nil {
		trunk.LogWarn("Failed to update the configuration: " + saveErr.Error())
	}
}

// RemoveUsers will remove the specified users from the AutoKick Users list, if they are added
func (autokicker *NarwhalAutoKickerPlugin) RemoveUsers(users []string) {
	var usersList = make(map[string]bool) // Map of users and their add / remove state
	newUsers := []string{}                // Users to retain

	for _, user := range Config.Plugins.AutoKick.Users { // For each user in Users
		for _, userToRemove := range users { // Users we're wanting to remove
			if userToRemove == user { // If this blacklist user matches the user we're wanting to remove
				usersList[userToRemove] = true   // Should remove the user
				delete(autokicker.Tracker, user) // Delete user from Tracker
				break
			}
		}

		if _, exists := usersList[user]; !exists { // User shouldn't be removed
			newUsers = append(newUsers, user)
		}
	}

	Config.Plugins.AutoKick.Users = DeduplicateList(newUsers) // Deduplicate users and set to AutoKick Users

	if saveErr := SaveConfig(); saveErr != nil {
		trunk.LogWarn("Failed to update the configuration: " + saveErr.Error())
	}
}
