package gcal

import (
	"github.com/BurntSushi/toml"
)

// Config contains user calendar config.
type Config struct {
	CalendarID string `toml:"calendar_id"`
	Credential string `toml:"credential_file"`
}

// LoadConfig returns user google calendar config.
func LoadConfig(fp string) (Config, error) {
	var c Config
	var err error
	if _, err = toml.DecodeFile(fp, &c); err != nil {
		return c, err
	}
	return c, nil
}
