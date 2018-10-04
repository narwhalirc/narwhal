package tusk

import (
	"github.com/lrstanley/girc"
	"sort"
	"strings"
)

// This file contains misc. utilities

// BanUser will ban the specified user from a channel
func BanUser(c *girc.Client, channel string, user string) {
	c.Cmd.Ban(channel, user)
}

// BanUsers will ban multiple users from a channel
func BanUsers(c *girc.Client, channel string, users []string) {
	for _, user := range users { // For each user
		BanUser(c, channel, user) // Issue a BanUser
	}
}

// DeduplicateList will eliminate duplicates from a list
func DeduplicateList(list []string) []string {
	var itemsInList = make(map[string]bool) // Define itemsInList as a list of items. Makes it easy to determine that we've already added an item
	newList := []string{}

	for _, entry := range list { // For each entry in list
		if _, exists := itemsInList[entry]; !exists { // entry not in list
			itemsInList[entry] = true
			newList = append(newList, entry) // Add the entry
		}
	}

	sort.Strings(newList) // Sort our entries
	return newList
}

// KickUser will kick the specified user from a channel
func KickUser(c *girc.Client, channel string, user string) {
	c.Cmd.Kick(channel, user, "Detected by this Narwhal for kick approval. Kicking.")
}

// KickUsers will kick multiple users from a channel
func KickUsers(c *girc.Client, channel string, users []string) {
	for _, user := range users { // For each user
		KickUser(c, channel, user) // Issue a KickUser
	}
}

// ParseMessage will parse an event and return a NarwhalMessage
func ParseMessage(e girc.Event) NarwhalMessage {
	var channel string
	var command string
	var params []string
	user := e.Source.Name

	if user == "" { // User is somehow empty
		user = e.Source.Ident // Change to using Ident
	}

	if e.IsFromChannel() { // If this is from a channel
		channel = e.Params[0] // Channel is first param
	}

	message := strings.TrimSpace(e.Trailing)
	msgSplit := strings.Split(message, " ")                 // Split on whitespace
	command = strings.Replace(msgSplit[:1][0], ".", "", -1) // Get the first item, remove .

	if len(msgSplit) > 1 {
		params = msgSplit[1:]
	}

	return NarwhalMessage{
		Channel: channel,
		Command: command,
		Host:    e.Source.Host,
		Issuer:  user,
		Message: e.Trailing,
		Params:  params,
	}
}

// UnbanUser will unban the specified user from a channel
func UnbanUser(c *girc.Client, channel string, user string) {
	c.Cmd.Unban(channel, user)
}

// UnbanUsers will unban multiple users from a channel
func UnbanUsers(c *girc.Client, channel string, users []string) {
	for _, user := range users {
		UnbanUser(c, channel, user) // Issue an UnbanUser
	}
}
