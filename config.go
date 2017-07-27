package main

import (
	"fmt"
	"log"
	"time"

	"github.com/BurntSushi/toml"
)

const (
	cfgpath = "/etc/pod-killer/config.toml"
)

// Configuration structure
type killerConfig struct {
	General   generalConfig
	Scheduler schedulerConfig
	Killer    runnerConfig
}

type generalConfig struct {
	Port int
}

type schedulerConfig struct {
	Timezone             string
	Crontime             string `toml:"crontime"`
	Random_range_measure string
}

type runnerConfig struct {
	Timeslot_deny_policy bool
	Schedule_timeslots   [][2]time.Time
}

func initDefault() *killerConfig {
	log.Printf("Load default values...")
	const layout = "2017-07-27 15:00:46"
	config := &killerConfig{
		General: generalConfig{
			Port: 8081,
		},
		Scheduler: schedulerConfig{
			Timezone: "Europe/Kiev",
			Crontime: "0 s p * * *",
		},
		Killer: runnerConfig{},
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
