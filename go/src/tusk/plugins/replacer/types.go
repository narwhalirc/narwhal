package tusk

// NarwhalReplacerConfig is our configuration for the Narwhal replacer plugin
type NarwhalReplacerConfig struct {
	// CachedMessageLimit is our limit of how many messages to cache
	CachedMessageLimit int
}

// NarwhalReplacerPlugin is our Replacer plugin
type NarwhalReplacerPlugin struct{}