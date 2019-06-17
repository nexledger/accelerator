/*
 *    Copyright 2019 Samsung SDS
 *
 *    Licensed under the Apache License, Version 2.0 (the "License");
 *    you may not use this file except in compliance with the License.
 *    You may obtain a copy of the License at
 *
 *        http://www.apache.org/licenses/LICENSE-2.0
 *
 *    Unless required by applicable law or agreed to in writing, software
 *    distributed under the License is distributed on an "AS IS" BASIS,
 *    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *    See the License for the specific language governing permissions and
 *    limitations under the License.
 */

package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/nexledger/accelerator/pkg/server"
)

var (
	configPath string
	flagSet    *flag.FlagSet
	logger     *zap.Logger
)

func main() {
	initializeArguments()
	startLoggers()

	s, err := server.New(configPath)
	if err != nil {
		logger.Error("Failed to create server: " + err.Error())
		return
	}

	failure := s.Serve()

	awaitTermination(failure)

	logger.Info("Shutting down the server......")
	closeLoggers()
	s.Stop()
	logger.Info("Stopped.")
}

func awaitTermination(failures ...chan error) {
	terminated := make(chan bool)

	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit
		terminated <- true
	}()

	for _, failure := range failures {
		go func() {
			<-failure
			terminated <- true
		}()
	}

	<-terminated
}

func initializeArguments() {
	flagSet = flag.NewFlagSet("server", flag.ExitOnError)
	flagSet.StringVar(&configPath, "f", "accelerator.yaml", "-f <configFilePath> : config file path")
	flagSet.Parse(os.Args[1:])
}

func startLoggers() {
	config := zap.NewDevelopmentConfig()
	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	l, err := config.Build()
	if err != nil {
		logger = zap.NewNop()
	}
	logger = l
	zap.ReplaceGlobals(logger)
	defer logger.Sync()
}

func closeLoggers() {
	logger.Sync()
}
