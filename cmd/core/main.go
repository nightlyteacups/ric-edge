/**
 * Copyright 2019 Rightech IoT. All Rights Reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package main

import (
	"flag"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	"github.com/Rightech/ric-edge/internal/app/core/config"
	"github.com/Rightech/ric-edge/internal/app/core/entrypoint"
	"github.com/Rightech/ric-edge/pkg/update"
)

// set at build time via ldflags
var version string // nolint: gochecknoglobals

const name = "core-" + runtime.GOOS + "-" + runtime.GOARCH

func main() {
	printCfg := flag.Bool("default-config", false, "print default configuration")
	printMinCfg := flag.Bool("min-config", false, "print minimal configuration")
	flag.Parse()

	if printCfg != nil && *printCfg {
		printConfig("default-config")
		return
	}

	if printMinCfg != nil && *printMinCfg {
		printConfig("min-config")
		return
	}

	config.Setup(version)

	log.Info("Version: ", viper.GetString("version"))

	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM)

	if viper.GetBool("check_updates") {
		res := update.Check(viper.GetString("version"), name)
		if res != "" {
			log.Info("New version available. Download it: ", res)

			if viper.GetBool("auto_download_updates") {
				update.Download(res)
				return
			}
		}
	}

	err := entrypoint.Start(signalCh)
	if err != nil {
		log.Fatal(err)
	}
}
