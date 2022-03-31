package main

import (
	"github.com/hortelanobruno/foaas-api/cmd"
	"os"
)

func main() {
	if err := cmd.Cmds().Execute(); err != nil {
		os.Exit(1)
	}
}
