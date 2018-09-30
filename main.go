package main

import (
	"os"
	"runtime"

	"github.com/mrl-athomelab/website/application/logger"
	"github.com/mrl-athomelab/website/application/server"
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {
	configFile := "config.yaml"
	if len(os.Args) > 1 {
		configFile = os.Args[1]
	}

	s, err := server.Prepare(configFile)
	if err != nil {
		logger.Error("Error on preparing server , %v", err)
		return
	}
	logger.Error("Serve error %v", s.Run())
}
