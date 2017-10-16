package apiconfig

import (
	"encoding/json"
	"os"

	"github.com/pkg/errors"
)

// DBConfig stores the database info
type DBConfig struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Database string `json:"database"`
	Host     string `json:"host"`
}

// APIConfig stores the API fields loaded from the JSON file
type APIConfig struct {
	Database DBConfig `json:"database"`
	Giphy    Giphy    `json:"giphy"`
}

var config *APIConfig

// Parse reads and parses the API config from the specified path.
func Parse(path string) (*APIConfig, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, errors.Wrap(err, "Error opening the API config file")
	}
	defer f.Close()

	config = new(APIConfig)
	err = json.NewDecoder(f).Decode(config)
	if err != nil {
		return nil, errors.Wrap(err, "Error decoding the API config file")
	}
	return config, nil
}

// Get returns the parsed config, or a new APIConfig{} if it's nil.
func Get() *APIConfig {
	if config == nil {
		config = new(APIConfig)
	}
	return config
}
