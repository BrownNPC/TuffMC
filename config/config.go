// Will be used to configure the server in the future, through toml or json.
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

