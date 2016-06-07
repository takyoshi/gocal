package gocal

import (
	"testing"
)

func TestLoadConfig(t *testing.T) {
	conf, err := LoadConfig("test/calendar.toml")
	if err != nil {
		t.Errorf("%v", err)
	}

	if conf.CalendarID != "primary" {
		t.Errorf("Wrong calendar_id value: got(%v), expected(primary)", conf.CalendarID)
	}
	if conf.Credential != "no file" {
		t.Errorf("Wrong credential_file value: got(%v), expected(no file)", conf.Credential)
	}
}

func TestErrLoadConfig(t *testing.T) {
	conf, err := LoadConfig("nothing_file")
	if err == nil {
		t.Errorf("err should not nil. %v", conf)
	}
}
