/*
Copyright (C) 2016  Eric Ziscky

    This program is free software; you can redistribute it and/or modify
    it under the terms of the GNU General Public License as published by
    the Free Software Foundation; either version 2 of the License, or
    (at your option) any later version.

    This program is distributed in the hope that it will be useful,
    but WITHOUT ANY WARRANTY; without even the implied warranty of
    MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
    GNU General Public License for more details.

    You should have received a copy of the GNU General Public License along
    with this program; if not, write to the Free Software Foundation, Inc.,
    51 Franklin Street, Fifth Floor, Boston, MA 02110-1301 USA.
*/

package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/BurntSushi/toml"

	"github.com/ziscky/mock-pesa/c2b"
	"github.com/ziscky/mock-pesa/common"
)

//api: ensure all mock api implementations satisfy this interface
type api interface {
	Start()
	Stop()
}

//parseConf parses the config from the path given
//or returns the default settings
func parseConf(path string) common.Config {
	var conf common.Config
	merchantID := os.Getenv("MERCHANT_ID")
	passkey := os.Getenv("PASSKEY")
	if merchantID == "" {
		merchantID = "12345"
	}
	if passkey == "" {
		passkey = "54321"
	}
	if _, err := toml.DecodeFile(path, &conf); err != nil {
		fmt.Println("Failed loading conf file,using defaults.")
		return common.Config{
			MaxAmount:                    70000,
			MinAmount:                    10,
			MerchantID:                   merchantID,
			CallBackDelay:                0,
			SAGPasskey:                   passkey,
			MaxCustomerTransactionPerDay: 150000,
			EnabledAPIS:                  []string{"c2b"},
		}
	}
	return conf

}

//startAPIS starts the apis
func startAPIS(apis ...api) {
	for _, api := range apis {
		api.Start()
	}
}

//stopAPIS stops the apis
func stopAPIS(apis ...api) {
	for _, api := range apis {
		api.Stop()
	}
}

func main() {
	c2bPort := flag.String("-c2b", "7000", "-c2b=portno.Default=7000")
	conf := flag.String("-conf", "config", "-conf=/path/to/conf. Default=./config")
	flag.Parse()

	config := parseConf(*conf)
	fmt.Println("Config:\n", config.ToString())

	var enabledAPIS []api
	for _, v := range config.EnabledAPIS {
		if v == "c2b" {
			enabledAPIS = append(enabledAPIS, c2b.NewAPI(*c2bPort, config))
		}
	}

	fmt.Println("Starting:", config.EnabledAPIS)
	startAPIS(enabledAPIS...)

	serve := make(chan os.Signal)
	signal.Notify(serve, syscall.SIGINT, syscall.SIGTERM) //Submit to user demands to quit, but gracefully yee sir.

	<-serve
	//graceful
	stopAPIS(enabledAPIS...)

}
