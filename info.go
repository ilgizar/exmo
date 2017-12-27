package main

import (
    "fmt"

    "github.com/ilgizar/exmo/api"
)

func info() {
    info := api.Info()
    fmt.Printf("User ID: %+v\n", info.UID)
    for key, val := range info.Balances {
        if val > 0 || info.Reserved[key] > 0 {
            fmt.Printf("%5s: %.2f [%.2f]\n", key, val, info.Reserved[key])
        }
    }
}
