package utils

import "testing"

func TestParseConfig(t *testing.T) {
	ParseConfig("../config.ini")
	t.Log(Conf)
}
