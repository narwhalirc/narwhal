package tusk

import (
	"github.com/lrstanley/girc"
	"strings"
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

func (song *NarwhalSongPlugin) Parse(c *girc.Client, e girc.Event) {
	msg := strings.TrimSpace(e.Trailing)

	for _, query := range []string{"!song", ".song"} {
		if msg == query {
			c.Cmd.Reply(e, "NARWHALS NARWHALS: https://www.youtube.com/watch?v=ykwqXuMPsoc")
			break
		}
	}
}
