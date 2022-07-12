/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package main

import (
	"github.com/hamster-shared/hamster-provider/cmd"
	golog "github.com/ipfs/go-log/v2"
	log "github.com/sirupsen/logrus"
	"os"
)

func main() {
	golog.SetAllLoggers(golog.LevelInfo) // Change to INFO for extra info
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)
	log.Info("hello logrus")
	cmd.Execute()
}
