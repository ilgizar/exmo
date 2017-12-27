package main

import (
    "github.com/ilgizar/exmo/conf"
    "github.com/ilgizar/exmo/trader"
)

type ConfigStruct struct {
    Keys struct {
        Public  string
        Secret  string
    }

    Pairs []trader.PairStruct
}

var config ConfigStruct

func readConfig() {
    conf.Read(configFile, &config)
}

func reloadConfig() {
    readConfig()
}
