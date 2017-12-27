package main

import (
    "flag"
    "time"

    "github.com/ilgizar/exmo/api"
    "github.com/ilgizar/exmo/trader"
)

func main() {
    flag.Parse()

    readConfig()

    initHUP()

    api.Init(config.Keys.Public, config.Keys.Secret)

    info()

    for _, pair := range config.Pairs {
        go func() {
            trader.Init(pair.Pair)
            c := time.Tick(time.Second)
            for _ = range c {
                trader.Watch(pair)
            }
        }()
    }

    c := time.Tick(time.Second)
    for _ = range c {}
}
