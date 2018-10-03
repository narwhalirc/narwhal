package tusk

import (
	"github.com/JoshStrobl/trunk"
	"github.com/lrstanley/girc"
	"strings"
)

// NarwhalAutoKicker is our autokicker plugin
var NarwhalAutoKicker NarwhalAutoKickerPlugin

// NarwhallAutoKickerConfig is our configuration for the Narwhal autokicker
type NarwhalAutoKickerConfig struct {
	// Enabled determines whether to enable this functionality
	Enabled bool

	// Hosts to kick. Matches from end.
	Hosts []string

	// MessageMatches is a list of messages that will result in kicks
	MessageMatches []string

	// Users to kick. Matches from beginning.
	Users []string
}

// NarwhalAutoKickerPlugin is our Autokick plugin
type NarwhalAutoKickerPlugin struct{}

func (autokicker *NarwhalAutoKickerPlugin) Parse(c *girc.Client, e girc.Event) {
	msg := strings.TrimSpace(e.Trailing)
	user := e.Source.Name
	host := e.Source.Host

	var userShouldBeKicked bool

	// #region Hosts Kick List Check

	if len(Config.Commands.AutoKick.Hosts) > 0 { // If we have a Hosts list
		for _, kickHost := range Config.Commands.AutoKick.Hosts {
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
		if len(Config.Commands.AutoKick.MessageMatches) > 0 { // If we have a Messages list
			for _, match := range Config.Commands.AutoKick.MessageMatches {
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
		if len(Config.Commands.AutoKick.Users) > 0 { // If we have a Kick list
			for _, kickUser := range Config.Commands.AutoKick.Users {
				if strings.HasSuffix(kickUser, "*") { // If we should not be doing exact match
					if strings.HasPrefix(user, strings.Replace(kickUser, "*", "", -1)) { // If the username begins with this kickUser
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
		for _, channel := range Config.Channels { // For each channel
			KickUser(c, channel, user)
		}
	}
}

// AddUsers will add the specified users to the AutoKick Users list, if they aren't already added
func (autokicker *NarwhalAutoKickerPlugin) AddUsers(users []string) {
	for _, requestedAddUser := range users {
		Config.Commands.AutoKick.Users = append(Config.Commands.AutoKick.Users, requestedAddUser)
	}

	Config.Commands.AutoKick.Users = DeduplicateList(Config.Commands.AutoKick.Users) // Deduplicate users and set to AutoKick Users

	if saveErr := SaveConfig(); saveErr != nil {
		trunk.LogWarn("Failed to update the configuration: " + saveErr.Error())
	}
}

// RemoveUsers will remove the specified users from the AutoKick Users list, if they are added
func (autokicker *NarwhalAutoKickerPlugin) RemoveUsers(users []string) {
	var usersList = make(map[string]bool) // Map of users and their add / remove state
	newUsers := []string{}                // Users to retain

	for _, user := range Config.Commands.AutoKick.Users { // For each user in Users
		for _, userToRemove := range users { // Users we're wanting to remove
			if userToRemove == user { // If this blacklist user matches the user we're wanting to remove
				trunk.LogWarn("Match user: " + userToRemove + " by " + user)
				usersList[userToRemove] = true // Should remove the user
				break
			}
		}

		if _, exists := usersList[user]; !exists { // User shouldn't be removed
			trunk.LogWarn("Should not remove user, therefore appending: " + user)
			newUsers = append(newUsers, user)
		}
	}

	Config.Commands.AutoKick.Users = DeduplicateList(newUsers) // Deduplicate users and set to AutoKick Users

	if saveErr := SaveConfig(); saveErr != nil {
		trunk.LogWarn("Failed to update the configuration: " + saveErr.Error())
	}
}
