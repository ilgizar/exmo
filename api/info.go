package api

import (
    "fmt"
    "sync"
    "time"
)

type UserInfoStruct struct {
    UID      int
    Date     int
    Balances map[string]float64
    Reserved map[string]float64
}

type AmountStruct struct {
    Quantity float64
    Amount   float64
    Price    float64
}

type TickerStruct struct {
    BuyPrice  float64
    SellPrice float64
    LastTrade float64
    High      float64
    Low       float64
    Avg       float64
    Vol       float64
    VolCurr   float64
    Updated   int
}

type TickersStruct map[string]TickerStruct

type MinMaxStruct struct {
    Min      float64
    Max      float64
}

type PairStruct struct {
    Quantity MinMaxStruct
    Price    MinMaxStruct
    Amount   MinMaxStruct
}

type PairsStruct map[string]PairStruct

type TickerSharedStruct struct {
    sync.Mutex
    cache    Response
    time     int64
}

const (
    tickerDelay int64 = 60
)

var tickerSharedData TickerSharedStruct


func Info() UserInfoStruct {
    resp, err := Query("user_info", nil, true)

    var info UserInfoStruct
    if err == nil {
        info.UID      = Interface2Int(resp["uid"])
        info.Date     = Interface2Int(resp["server_date"])
        info.Balances = Interface2Map(resp["balances"])
        info.Reserved = Interface2Map(resp["reserved"])
    }

    return info
}

func Amount(pair string, quantity float64) AmountStruct {
    resp, err := Query("required_amount", Params{"pair": pair, "quantity": fmt.Sprintf("%f", quantity)}, true)

    var amount AmountStruct
    if err == nil {
        amount.Quantity = resp["quantity"].(float64)
        amount.Amount   = resp["amount"].(float64)
        amount.Price    = resp["avg_price"].(float64)
    }

    return amount
}

func Ticker(pairName string) TickersStruct {
    now := time.Now().Unix()

    var err error

    tickerSharedData.Lock()
    if tickerSharedData.time < now - tickerDelay {
        tickerSharedData.time = now
        tickerSharedData.cache, err = Query("ticker", nil, false)
    }
    tickerSharedData.Unlock()

    tickers := TickersStruct{}
    if err == nil {
        for pair, value := range tickerSharedData.cache {
            if pairName == "" || pairName == pair {
                val := value.(map[string]interface{})
                ticker := TickerStruct{}

                ticker.BuyPrice  = Interface2Float(val["buy_price"])
                ticker.SellPrice = Interface2Float(val["sell_price"])
                ticker.LastTrade = Interface2Float(val["last_trade"])
                ticker.High      = Interface2Float(val["high"])
                ticker.Low       = Interface2Float(val["low"])
                ticker.Avg       = Interface2Float(val["avg"])
                ticker.Vol       = Interface2Float(val["vol"])
                ticker.VolCurr   = Interface2Float(val["vol_curr"])
                ticker.Updated   = Interface2Int(val["updated"])

                tickers[pair] = ticker
            }
        }
    }

    return tickers
}

func Pairs(pairName string) PairsStruct {
    resp, err := Query("pair_settings", nil, false)

    pairs := PairsStruct{}
    if err == nil {
        for pair, value := range resp {
            if pairName == "" || pairName == pair {
                val := value.(map[string]interface{})
                p := PairStruct{}

                p.Quantity.Min = Interface2Float(val["min_quantity"])
                p.Quantity.Max = Interface2Float(val["max_quantity"])
                p.Price.Min    = Interface2Float(val["min_price"])
                p.Price.Max    = Interface2Float(val["max_price"])
                p.Amount.Min   = Interface2Float(val["min_amount"])
                p.Amount.Max   = Interface2Float(val["max_amount"])

                pairs[pair] = p
            }
        }
    }

    return pairs
}

/*
func Currency() []string {
    resp, err := Query("currency", nil, false)

    var currency []string
    if err == nil {
        currency = resp.([]string)
    }

    return currency
}
*/
