package api

import (
    "fmt"
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

func Info() UserInfoStruct {
    resp, err := Query("user_info", nil, true)

    var info UserInfoStruct
    if err == nil {
        info.UID      = interface2int(resp["uid"])
        info.Date     = interface2int(resp["server_date"])
        info.Balances = interface2map(resp["balances"])
        info.Reserved = interface2map(resp["reserved"])
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
    resp, err := Query("ticker", nil, false)

    tickers := TickersStruct{}
    if err == nil {
        for pair, value := range resp {
            if pairName == "" || pairName == pair {
                val := value.(map[string]interface{})
                ticker := TickerStruct{}

                ticker.BuyPrice  = interface2float(val["buy_price"])
                ticker.SellPrice = interface2float(val["sell_price"])
                ticker.LastTrade = interface2float(val["last_trade"])
                ticker.High      = interface2float(val["high"])
                ticker.Low       = interface2float(val["low"])
                ticker.Avg       = interface2float(val["avg"])
                ticker.Vol       = interface2float(val["vol"])
                ticker.VolCurr   = interface2float(val["vol_curr"])
                ticker.Updated   = interface2int(val["updated"])

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

                p.Quantity.Min = interface2float(val["min_quantity"])
                p.Quantity.Max = interface2float(val["max_quantity"])
                p.Price.Min    = interface2float(val["min_price"])
                p.Price.Max    = interface2float(val["max_price"])
                p.Amount.Min   = interface2float(val["min_amount"])
                p.Amount.Max   = interface2float(val["max_amount"])

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
