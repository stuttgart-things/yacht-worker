/*
Copyright © 2023 Patrick Hermann patrick.hermann@sva.de
*/

package worker

import (
	"context"
	"fmt"
	"os"

	redis "github.com/redis/go-redis/v9"

	sthingsBase "github.com/stuttgart-things/sthingsBase"
)

var (
	// elasticSearchUrl   = os.Getenv("ELASTICSEARCH_URL")
	// elasticSearchIndex = os.Getenv("ELASTICSEARCH_STATUS_INDEX")
	sendToElastic = os.Getenv("STATUS_TO_ELASTICSEARCH")
	// logfilePathStatus  = "yaw-status.log"
	redisAddress  = os.Getenv("REDIS_SERVER")
	redisPort     = os.Getenv("REDIS_PORT")
	redisPassword = os.Getenv("REDIS_PASSWORD")
	ctx           = context.Background()
)

type PipelineRunStatus struct {
	Name       string `mapstructure:"name"`
	Status     string `mapstructure:"status"`
	Completed  int    `mapstructure:"completed"`
	Failed     int    `mapstructure:"failed"`
	Canceled   int    `mapstructure:"canceled"`
	Incomplete int    `mapstructure:"incomplete"`
	Skipped    int    `mapstructure:"skipped"`
}

func GetPipelineRunStatus(prStatus string) (pipelineRunStatus map[string]string) {

	pipelineRunStatus = make(map[string]string)

	pipelineRunStatus["status"], _ = sthingsBase.GetRegexSubMatch(prStatus, `Succeeded\s(\w+)`)
	pipelineRunStatus["completed"], _ = sthingsBase.GetRegexSubMatch(prStatus, `Completed:\s(\w+)`)
	pipelineRunStatus["failed"], _ = sthingsBase.GetRegexSubMatch(prStatus, `Failed:\s(\w+)`)
	pipelineRunStatus["canceled"], _ = sthingsBase.GetRegexSubMatch(prStatus, `Cancelled\s(\w+)`)
	pipelineRunStatus["incomplete"], _ = sthingsBase.GetRegexSubMatch(prStatus, `Incomplete:\s(\w+)`)
	pipelineRunStatus["skipped"], _ = sthingsBase.GetRegexSubMatch(prStatus, `Skipped:\s(\w+)`)

	return pipelineRunStatus

}

func GetRevisionRunStatus(prName string) int {

	redisClient := redis.NewClient(&redis.Options{
		Addr:     redisAddress + ":" + redisPort,
		Password: redisPassword, // no password set
		DB:       0,             // use default DB
	})

	err := redisClient.Set(ctx, prName, "DONE", 0).Err()
	if err != nil {
		panic(err)
	}

	// DECREASE TOTAL COUNT
	statusValue, err := redisClient.Get(ctx, "countPipelineRuns").Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("DECREASING FOR "+prName, statusValue)
	redisClient.Decr(ctx, "countPipelineRuns")
	statusValue, err = redisClient.Get(ctx, "countPipelineRuns").Result()
	fmt.Println("AFTER DECREASING..", statusValue)

	if err != nil {
		panic(err)
	}

	return sthingsBase.ConvertStringToInteger(statusValue)

}
