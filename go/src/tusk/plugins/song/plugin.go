package tusk

import (
	"github.com/lrstanley/girc"
)

// NarwhalSong is our song plugin
var NarwhalSong NarwhalSongPlugin

func (song *NarwhalSongPlugin) Parse(c *girc.Client, e girc.Event, m NarwhalMessage) {
	if m.Command == "song" {
		c.Cmd.Reply(e, "NARWHALS NARWHALS: https://www.youtube.com/watch?v=ykwqXuMPsoc")
	}
}
