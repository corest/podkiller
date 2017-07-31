package main

import (
	"fmt"
	"log"

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
	Influx    influxConfig
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
	Selector              string
	Namespace_deny_policy bool
	Namespace_list        []string
}

type influxConfig struct {
	Hostname string
	Port     int
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
		Killer: runnerConfig{
			Selector:              "destiny in (doomed)",
			Namespace_deny_policy: true,
			Namespace_list:        []string{"kube-system"},
		},
		Influx: influxConfig{
			Hostname: "localhost",
			Port:     8086,
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
