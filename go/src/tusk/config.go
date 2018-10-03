package tusk

// This file contains our configuration logic

import (
	"bufio"
	"bytes"
	"errors"
	"github.com/BurntSushi/toml"
	"io/ioutil"
	"path/filepath"
)

// ConfigPaths is our valid config paths
var ConfigPaths []string

// ConfigFoundPath is the directory we found the config in
var ConfigFoundPath string

func init() {
	ConfigPaths = []string{
		filepath.Join(HomeDir, ".config", "narwhal", "config.toml"), // User's XDG_CONFIG_DIR
		filepath.Join(WorkDir, "config.toml"),                       // Local
		"/etc/narwhal/config.toml",                                  // System-wide
		"/usr/share/defaults/narwhal/config.toml",                   // Stateless vendor config
	}
}

// ReadConfig will read our narwhal configuration, if it exists, and return it
func ReadConfig() (NarwhalConfig, error) {
	var config NarwhalConfig
	var readErr error

	for index, dir := range ConfigPaths { // Search each path
		if configBytes, readErr := ioutil.ReadFile(dir); readErr == nil { // If we successfully read the file
			if len(configBytes) > 0 { // If file is not empty
				decodeErr := toml.Unmarshal(configBytes, &config)

				if decodeErr == nil { // No error during unmarshal
					ConfigFoundPath = dir
					config = SetDefaults(config) // Enforce config defaults
				} else {
					readErr = errors.New("Failed to decode config: " + decodeErr.Error())
				}

				break
			}
		} else { // Failed to read the file
			if index == (len(ConfigPaths) - 1) { // Last file being read
				readErr = errors.New("Failed to find Narwhal's config.toml in any recognized location.")
			}
		}
	}

	return config, readErr
}

// SaveConfig will save the config in our previously recognized location
func SaveConfig() error {
	var saveErr error
	var buffer bytes.Buffer
	writer := bufio.NewWriter(&buffer)
	encoder := toml.NewEncoder(writer) // Create a new toml encoder
	encoder.Indent = "\t"              // Use a tab because we're opinionated

	if saveErr = encoder.Encode(Config); saveErr == nil { // Encode our Config into a buffer
		saveErr = ioutil.WriteFile(ConfigFoundPath, buffer.Bytes(), 0644) // Write the config
	}

	return saveErr
}

// SetDefaults will set the defaults for the provided NarwhalConfig
func SetDefaults(config NarwhalConfig) NarwhalConfig {
	if config.Name == "" {
		config.Name = "Narwhal Bot"
	}

	if config.Network == "" {
		config.Network = "chat.freenode.net:6667" // Default to freenode
	}

	if config.Port == 0 {
		config.Port = 6667 // Default to 6667
	}

	return config
}
