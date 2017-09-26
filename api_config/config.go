package api_config

import (
	"encoding/json"
	"os"

	"github.com/pkg/errors"
)

type ApiConfig struct {
	Giphy Giphy `json:"giphy"`
}

var config *ApiConfig

// Parse reads and parses the API config from the specified path.
func Parse(path string) (*ApiConfig, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, errors.Wrap(err, "Error opening the API config file")
	}
	defer f.Close()

	config = new(ApiConfig)
	err = json.NewDecoder(f).Decode(config)
	if err != nil {
		return nil, errors.Wrap(err, "Error decoding the API config file")
	}
	return config, nil
}

// Get returns the parsed config, or a new ApiConfig{} if it's nil.
func Get() *ApiConfig {
	if config == nil {
		config = new(ApiConfig)
	}
	return config
}
