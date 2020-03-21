package utils

import (
	"gopkg.in/ini.v1"
	"log"
)

type Storage struct {
	DSN         string `ini:"dsn"`
	MaxIdleConn int    `ini:"max_idle_conn"`
	MaxOpenConn int    `ini:"max_open_conn"`
}

type Config struct {
	DB         Storage `ini:"db"`
	BaseString string  `ini:"base_string"`
	Host       string  `ini:"host"`
}

var Conf *Config

func ParseConfig(filepath string) {
	cfg, err := ini.Load(filepath)
	if err != nil {
		log.Panicf("fail to read config file: %s due to %v", filepath, err)
	}
	Conf = new(Config)
	err = cfg.MapTo(Conf)
	if err != nil {
		log.Panicf("fail to map config to Config due to %v", err)
	}
}
