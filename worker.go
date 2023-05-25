/*
Copyright © 2023 Patrick Hermann patrick.hermann@sva.de
*/

package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/stuttgart-things/yacht-worker/worker"

	// "github.com/stuttgart-things/yacht-worker/worker"

	"github.com/fatih/color"
	sthingsBase "github.com/stuttgart-things/sthingsBase"
	sthingsCli "github.com/stuttgart-things/sthingsCli"

	goVersion "go.hein.dev/go-version"
)

var (
	logfilePath   = "yw.log"
	shortened     = false
	version       = "unset"
	date          = "unknown"
	commit        = "unknown"
	output        = "yaml"
	log           = sthingsBase.StdOutFileLogger(logfilePath, "2006-01-02 15:04:05", 50, 3, 28)
	redisClient   = sthingsCli.CreateRedisClient(redisAddress+":"+redisPort, redisPassword)
	redisAddress  = os.Getenv("REDIS_SERVER")
	redisPort     = os.Getenv("REDIS_PORT")
	redisPassword = os.Getenv("REDIS_PASSWORD")
	prRanges      = os.Getenv("PR_RANGES")
	revisionRunID = os.Getenv("REVISION_RUN_ID")
	prs           = make(map[int][]string)
	prCount       = strings.Split(prRanges, ";")
)

const banner = `
___    ___ ___       __
|\  \  /  /|\  \     |\  \
\ \  \/  / \ \  \    \ \  \
 \ \    / / \ \  \  __\ \  \
  \/  /  /   \ \  \|\__\_\  \
__/  / /      \ \____________\
|\___/ /        \|____________|
\|___|/

`

func main() {

	// STARTUP / STATUS OUTPUT
	color.Cyan(banner)
	color.Cyan("YACHT WORKER")
	resp := goVersion.FuncWithOutput(shortened, version, commit, date, output)
	color.Cyan(resp + "\n")
	log.Info("YW server started")

	// GETTING REVISIONRUN DATA FROM REDIS
	pipelineRunData := redisClient.HGetAll(revisionRunID).Val()

	// CREATE PRS + DEBUG OUTPUT
	for i := 0; i < len(prCount); i++ {
		fmt.Println("STAGE", i)
		intVar, _ := strconv.Atoi(prCount[i])

		for j := 0; j < intVar; j++ {
			prs[i] = append(prs[i], pipelineRunData[strconv.Itoa(i)+":"+strconv.Itoa(j)])
			fmt.Println("PIPELINE", j)
			fmt.Println(pipelineRunData[strconv.Itoa(i)+"."+strconv.Itoa(j)])
		}
	}

	// DEBUG OUTPUT OF ALL PRS
	fmt.Println(prs)

	// CREATING AND WATCHING RUNS
	worker.ConsumeRevisionRun(prs)
	log.Warn("YW stopped")
	os.Exit(0)

}
