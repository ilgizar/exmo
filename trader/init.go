package trader

import (
    "github.com/ilgizar/exmo/api"
)

var settings api.PairStruct

func Init(pair string) {
    res := api.Pairs(pair)
    settings = res[pair]
}
