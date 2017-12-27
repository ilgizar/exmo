package trader

import (
    "fmt"

    "github.com/ilgizar/exmo/api"
)

func Watch(pair PairStruct) {
    res := api.Ticker(pair.Pair)
    if ticker, ok := res[pair.Pair]; ok {
        if ticker.BuyPrice <= float64(pair.Floor) {
            fmt.Printf("time to buy: %+v\n", ticker.BuyPrice)
        } else if ticker.SellPrice >= float64(pair.Roof) {
            fmt.Printf("time to sell: %+v\n", ticker.SellPrice)
        } else {
            fmt.Printf("%+v\n", ticker)
        }
    }
}
