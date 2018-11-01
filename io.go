package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"time"
)

func readArgs(args []string) {
	if len(os.Args) != 4 {
		println("Usage:    failover <server-address> <server-alias> <settings-path>")
		os.Exit(1)
	}
	varPath = "/tmp/check-" + os.Args[2]
	url = os.Args[1]
	settingsPath = os.Args[3]
}

func logEvent(event string) {
	t := time.Now()
	fmt.Printf("[%d-%02d-%02d %02d:%02d:%02d] %s\n", t.Year(), int(t.Month()), t.Day(), t.Hour(), t.Minute(), t.Second(), event)
}

func readVars() {
	varFile, err := ioutil.ReadFile(varPath)
	if err != nil {
		if os.IsNotExist(err) {
			internetStatus = true
			serverStatus = true
			varFile = []byte{boolToByte([8]bool{internetStatus, serverStatus})}
			writeVars()
			logEvent("Reboot detected!")
		} else {
			panic(err)
		}
	}
	vars := byteToBool(varFile[0])

	internetStatus = vars[0]
	serverStatus = vars[1]
}

func writeVars() {
	err := ioutil.WriteFile(varPath, []byte{boolToByte([8]bool{internetStatus, serverStatus})}, 0600)
	if err != nil {
		panic(err)
	}
}
