package tusk

import (
	"github.com/lrstanley/girc"
)

// NarwhalSong is our song plugin
var NarwhalSong NarwhalSongPlugin

// NarwhalSongConfig is our configuration for the Narwhal song
type NarwhalSongConfig struct {
	// Enabled determines whether to enable this functionality
	Enabled bool
}

// NarwhalSong is our Song plugin
type NarwhalSongPlugin struct{}

func (song *NarwhalSongPlugin) Parse(c *girc.Client, e girc.Event, m NarwhalMessage) {
	if m.Command == "song" {
		c.Cmd.Reply(e, "NARWHALS NARWHALS: https://www.youtube.com/watch?v=ykwqXuMPsoc")
	}
}
