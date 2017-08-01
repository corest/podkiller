package main

import (
	"bytes"
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/robfig/cron"
)

func random(min, max int) int {
	if min == 0 {
		min = 1
	}
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(max-min) + min
}

func parseCronOption(c string, min, max int) string {
	var res string
	rand := strconv.Itoa(random(min, max))
	switch c {
	case "s":
		res = rand
	case "p":
		res = fmt.Sprintf("*/%s", rand)
	default:
		res = c
	}

	return res
}

func getCronTime(config *Config) string {

	cronTemplate := config.Scheduler.Crontime
	var buffer bytes.Buffer
	cronLimits := [6][2]int{{0, 59}, {0, 59}, {0, 23}, {1, 31}, {1, 12}, {0, 6}}

	for i, c := range strings.Split(cronTemplate, " ") {

		if i != 0 {
			buffer.WriteString(" ")
		}
		buffer.WriteString(parseCronOption(c, cronLimits[i][0], cronLimits[i][1]))

	}

	crontime := buffer.String()
	log.Printf("Using cronstring: %s", crontime)

	return crontime
}

func getJobScheduler(config *Config, cmd *KillerJob) (*cron.Cron, error) {

	var cronjob *cron.Cron
	timezone := config.Scheduler.Timezone
	loc, err := time.LoadLocation(timezone)
	if err != nil {
		log.Printf("Failed to parse timezone. Using localtime")
		cronjob = cron.New()
	} else {
		cronjob = cron.NewWithLocation(loc)
	}

	crontime := getCronTime(config)
	schedule, err := cron.Parse(crontime)
	if err != nil {
		log.Fatal(err)
	}

	nextrun := schedule.Next(time.Now())
	log.Printf("Next pod-killer run at: %s", nextrun.String())

	cmd.setSchedule(crontime)
	cronjob.Schedule(schedule, cmd)

	return cronjob, nil
}
