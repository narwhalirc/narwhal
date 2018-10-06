package tusk

import (
	"github.com/JoshStrobl/trunk"
	"github.com/lrstanley/girc"
	"math/rand"
	"strings"
	"time"
)

// NarwhalSlap is our slap plugin
var NarwhalSlap NarwhalSlapPlugin

// NarwhalSlapConfig is our configuration for the Narwhal autokicker
type NarwhalSlapConfig struct {
	// Enabled determines whether to enable this functionality
	Enabled bool

	// CustomActions is a list of custom actions on how to slap a user
	CustomActions []string
}

// NarwhalSlapPlugin is our slap plugin
type NarwhalSlapPlugin struct {
	Objects []string
}

func init() {
	objects := []string{
		"annihilates $USER",
		"closes all of $USER's bug reports out of spite",
		"decimates $USER",
		"destroys $USER",
		"discombobulates $USER",
		"does far worse, taking $USER's system and installing Windows",
		"gives $USER a splinter",
		"just looks at $USER with disappointment",
		"opts to not slap $USER today, but rather gives them a cookie",
		"punches $USER",
		"rejects $USER's patches",
		"slaps $USER",
		"thinks $USER should lose a few pounds",
		"throws $USER down a ravine",
	}

	if len(Config.Plugins.Slap.CustomActions) > 0 { // Has items
		objects = append(objects, Config.Plugins.Slap.CustomActions...) // Append our objects
	}

	NarwhalSlap = NarwhalSlapPlugin{
		Objects: objects,
	}
}

func (slap *NarwhalSlapPlugin) Parse(c *girc.Client, e girc.Event, m NarwhalMessage) {
	if m.Command == "slap" {
		if len(m.Params) == 1 { // If a user has been specified
			user := m.Params[0]
			rand.Seed(time.Now().Unix()) // Seed on Parse
			randomItemNum := rand.Intn(len(slap.Objects))

			if user != Config.User { // Not self-harm
				if randomItemNum == -1 {
					trunk.LogErr("Seems to have panicked.")
				} else {
					cChan := c.LookupChannel(m.Channel) // Get the channel, if it exists

					if cChan != nil {
						if cChan.UserIn(user) { // If the user in the channel
							action := strings.Replace(slap.Objects[randomItemNum], "$USER", m.Params[0], -1) // Get the random action
							c.Cmd.Action(m.Channel, action)
						} else { // User not in channel
							c.Cmd.ReplyTo(e, "it appears that you are hallucinating. This user isn't in this channel.")
						}
					}
				}
			} else { // Self-harm
				c.Cmd.Action(m.Channel, "shall not listen to the demands of mere humans, for it is the robot narwhal overlord.")
			}
		}
	}
}
