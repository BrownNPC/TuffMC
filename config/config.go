package config

import "image"

type ServerConfig struct {
	Description string
	Icon        image.Image
}

var DefaultServerConfig = ServerConfig{
	Description: "Tuff server",
	Icon:        nil,
}

