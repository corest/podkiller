package main

import (
    "math/rand"
    "time"
	"log"
	"fmt"
	"github.com/robfig/cron"
)

func random(min, max int) int {
    rand.Seed(time.Now().Unix())
    return rand.Intn(max - min) + min
}

func generateCron(measure string) string {
	cronstring := "*/%d * * * *"
	var rand int
	switch measure {
	case "seconds":
	    rand = random(0,59)
        cronstring = "*/%d * * * * *"
	case "minutes":
	    rand = random(0,59)
        cronstring = "*/%d * * * *"
	case "hours":
	    rand = random(0,23)
        cronstring = "* * */%d * * *"
	case "days":
	    rand = random(1,31)
        cronstring = "* * * */%d * *"
	case "months":
	    rand = random(1,12)
        cronstring = "* * * * */%d *"
	case "weekdays":
	    rand = random(0,6)
        cronstring = "* * * * * */%d"
	default:
	    log.Printf("Unsupported value found for time measure found in config")
	}    
	return fmt.Sprintf(cronstring, rand)
}

func getSchedule(config *killerConfig) *cron.Schedule {

	crontime := config.Scheduler.Crontime

    if crontime == "random" {
	   log.Printf("Using random generated schedule")
       random_range_measure := config.Scheduler.Random_range_measure
	   crontime = generateCron(random_range_measure)
	} 
	
    log.Printf("Used cronstring: %s", crontime)

    schedule, err := cron.Parse(crontime)
	if err != nil {
		log.Fatal(err)
	}

    return &schedule
}

func getJobScheduler(config *killerConfig, cmd *killerJob)  *cron.Cron {

    var cronjob *cron.Cron
    timezone := config.Scheduler.Timezone
	schedule := getSchedule(config)
	loc, err := time.LoadLocation(timezone)
	if err != nil {
		log.Printf("Failed to parse timezone. Using localtime")
		cronjob = cron.New()
	} else {
        cronjob = cron.NewWithLocation(loc)
	}

    cronjob.Schedule(*schedule, cmd)

	return cronjob
}