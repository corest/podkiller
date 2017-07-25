package main

import (
	"log"
	"strconv"
	"fmt"

	"github.com/BurntSushi/toml"
)

const (
	cfgpath = "/etc/pod-killer/config.toml"
)

// Configuration structure
type Config struct {
	general generalConfig
	scheduler schedulerConfig
}

type generalConfig struct {
	port int
}

type schedulerConfig struct {
	timezone string
    running_interval int
    allowed_period [2]string
}

func initDefault() *Config {
	log.Printf("Load default values...")
	config := &Config{
		general: generalConfig{
			port: 8081,
		},
		scheduler: schedulerConfig {
			timezone: "Europe/Kiev",
			running_interval: -1,
			allowed_period: [2]string{"00.00 Mon", "23.59 Sun"},
		},
	}
	return config
}

func getConfig() (Config, error) {
	log.Printf("Initialize config...")

    config := initDefault()

	if _, err := toml.DecodeFile(cfgpath, config); err != nil {
		log.Fatal(fmt.Sprintf("Unable to load config %s\n", cfgpath), err.Error())
	} 
	
	log.Printf(fmt.Sprintf("Pod-killer uses config from %s\n", cfgpath))
	

	//if err := viper.ReadInConfig(); err != nil {
	//	return err
	//}

	//ValidateConfigs()
	//setupWatch()
	return *config, nil
}
