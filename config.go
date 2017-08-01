package main

import (
	"fmt"
	"log"

	"github.com/BurntSushi/toml"
)

const (
	cfgpath = "/etc/pod-killer/config.toml"
)

// Config structure
type Config struct {
	General   GeneralConfig
	Scheduler SchedulerConfig
	Killer    KillerConfig
	Influx    InfluxConfig
}

// GeneralConfig structure
type GeneralConfig struct {
	Port int
}

// SchedulerConfig structure
type SchedulerConfig struct {
	Timezone string
	Crontime string `toml:"crontime"`
}

// KillerConfig structure
type KillerConfig struct {
	Selector            string
	NamespaceDenyPolicy bool     `toml:"namespace_deny_policy"`
	NamespaceList       []string `toml:"namespace_list"`
}

// InfluxConfig structure
type InfluxConfig struct {
	Hostname string
	Port     int
}

func initDefault() *Config {
	log.Printf("Load default values...")
	const layout = "2017-07-27 15:00:46"
	config := &Config{
		General: GeneralConfig{
			Port: 8081,
		},
		Scheduler: SchedulerConfig{
			Timezone: "Europe/Kiev",
			Crontime: "0 s p * * *",
		},
		Killer: KillerConfig{
			Selector:            "destiny in (doomed)",
			NamespaceDenyPolicy: true,
			NamespaceList:       []string{"kube-system"},
		},
		Influx: InfluxConfig{
			Hostname: "localhost",
			Port:     8086,
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

	return *config, nil
}
