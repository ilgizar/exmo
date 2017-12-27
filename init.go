package main

import (
    "flag"
    "os"
    "os/signal"
    "syscall"
)

var configFile string

func init() {
    flag.StringVar(&configFile, "config",       "exmo.conf",    "path to config file")
}

func initHUP() {
    c := make(chan os.Signal, 1)
    signal.Notify(c, syscall.SIGHUP)

    go func(){
        for _ = range c {
            reloadConfig()
        }
    }()
}
