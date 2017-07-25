package main

import (
	"log"

	"fmt"

	"github.com/BurntSushi/toml"
)

const (
	cfgpath = "/etc/pod-killer/config.toml"
)

// Configuration structure
type Config struct {
	Title string
}

func getConfig() (Config, error) {
	log.Printf("Initialize config...")

	var conf Config
	if _, err := toml.DecodeFile(cfgpath, &conf); err != nil {
		log.Fatal(fmt.Sprintf("Unable to load config %s\n", cfgpath), err.Error())
	}

	// SetDefaults()

	//if err := viper.ReadInConfig(); err != nil {
	//	return err
	//}

	//ValidateConfigs()
	//setupWatch()
	return conf, nil
}
