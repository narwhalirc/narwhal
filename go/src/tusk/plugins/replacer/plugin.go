package tusk

import (
	"fmt"
	"github.com/lrstanley/girc"
	"strings"
)

// NarwhalReplacer is our string replacer plugin
var NarwhalReplacer NarwhalReplacerPlugin

// cachedMessages is an array of NarwhalMessage structs
var cachedMessages []NarwhalMessage

func init() {
	cachedMessages = []NarwhalMessage{} // Create a new slice
}

func (replacer *NarwhalReplacerPlugin) Parse(c *girc.Client, e girc.Event, m NarwhalMessage) {
	if m.Command != "r" { // Not replacer
		return
	}

	m.MessageNoCmd = strings.Replace(m.MessageNoCmd, ", ", ",", -1)
	potentialUserReplacerArr := strings.SplitAfterN(m.MessageNoCmd, ",", 2) // Attempt to split the replaces

	replacerMessage := m.MessageNoCmd // Default to message being without .r command
	username := m.Issuer // Default to person which wrote the message

	if len(potentialUserReplacerArr) > 1 { // Username specified
		username = potentialUserReplacerArr[0]
		username = strings.Replace(username, ",", "", -1) // Remove , reference
		replacerMessage = strings.Join(potentialUserReplacerArr[1:], "")
	}

	for _, message := range cachedMessages { // For each message
		if (message.Issuer == username) && message.Message != m.Message { // If this username matches and not this most recent message
			replacerMessage = strings.Replace(replacerMessage, "s/", "", -1)
			messageReplacement := strings.Split(replacerMessage, "/") // Split on / (so hello/world/ is [s hello world])

			if len(messageReplacement) == 3 { // Must be exact match
				searchWord := messageReplacement[0]
				replaceWord := messageReplacement[1]

				newMessage := strings.Replace(message.MessageNoCmd, searchWord, replaceWord, -1) // Replace all instances
				newMessage = strings.Replace(newMessage, "  ", " ", -1) // Replace double whitespace with single space

				c.Cmd.Action(m.Channel, fmt.Sprintf("thinks %s meant to say: %s", message.Issuer, newMessage))
			}

			break
		}
	}
}

// AddToCache will add the provided messaged to our cached list and remove the first entry if we get above our limit
func (replacer *NarwhalReplacerPlugin) AddToCache(m NarwhalMessage) {
	limit := Config.Plugins.Replacer.CachedMessageLimit
	newList := append([]NarwhalMessage{m}, cachedMessages...) // Prepend to a new list
	cachedMessages = newList

	cacheLen := len(cachedMessages)

	if cacheLen> limit { // If this is above our limit
		cachedMessages = cachedMessages[:(cacheLen - 1)]
	}
}