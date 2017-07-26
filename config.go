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
type killerConfig struct {
	General generalConfig
	Scheduler schedulerConfig
}

type generalConfig struct {
	Port int
}

type schedulerConfig struct {
	Timezone string
    Crontime string `toml:"crontime"`
    Random_range_measure string
}

func initDefault() *killerConfig {
	log.Printf("Load default values...")
	config := &killerConfig{
		General: generalConfig{
			Port: 8081,
		},
		Scheduler: schedulerConfig {
			Timezone: "Europe/Kiev",
			Crontime: "0 0 * * * *", 
			Random_range_measure: "minutes",
		},
	}
	return config
}

func getConfig() (killerConfig, error) {
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
