/*
Copyright © 2023 Patrick Hermann patrick.hermann@sva.de
*/
package main

import (
	"github.com/fatih/color"
	sthingsBase "github.com/stuttgart-things/sthingsBase"

	goVersion "go.hein.dev/go-version"
)

var (
	logfilePath = "yaw.log"
	shortened   = false
	version     = "unset"
	date        = "unknown"
	commit      = "unknown"
	output      = "yaml"
	log         = sthingsBase.StdOutFileLogger(logfilePath, "2006-01-02 15:04:05", 50, 3, 28)
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

}
