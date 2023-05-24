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

	"github.com/fatih/color"
	sthingsBase "github.com/stuttgart-things/sthingsBase"
	sthingsCli "github.com/stuttgart-things/sthingsCli"

	goVersion "go.hein.dev/go-version"
)

var (
	logfilePath   = "yaw.log"
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

	color.Cyan(banner)
	color.Cyan("YACHT WORKER")
	resp := goVersion.FuncWithOutput(shortened, version, commit, date, output)
	color.Cyan(resp + "\n")

	log.Info("YW server started")

	prs := make(map[int][]string)

	customerData := redisClient.HGetAll(revisionRunID).Val()

	s := strings.Split(prRanges, ";")

	for i := 0; i < len(s); i++ {
		fmt.Println("STAGE", i)

		intVar, _ := strconv.Atoi(s[i])

		for j := 0; j < intVar; j++ {

			prs[i] = append(prs[i], customerData[strconv.Itoa(i)+":"+strconv.Itoa(j)])
			fmt.Println("PIPELINE", j)

			fmt.Println(customerData[strconv.Itoa(i)+"."+strconv.Itoa(j)])

		}
	}

	fmt.Println(prs)
	worker.ConsumeRevisionRun(prs)
	log.Warn("YW stopped")
	os.Exit(0)

}
