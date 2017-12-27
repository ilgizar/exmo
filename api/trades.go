package api

type TradeStruct struct {
    TradeID  int
    Date     int
    Type     string
    Pair     string
    OrderID  int
    Quantity float64
    Price    float64
    Amount   float64
}

type TradeSlice []TradeStruct

type TradesStruct map[string]TradeSlice

type InOutTradesStruct struct {
    Currency string
    Amount   float64
}

type OrderTradesStruct struct {
    Type   string
    In     InOutTradesStruct
    Out    InOutTradesStruct
    Trades TradeSlice
}


func Trades(pair string, limit int, offset int, user bool) TradesStruct {
    var resp Response
    var err error
    if user {
        resp, err = Query("user_trades", Params{"pair": pair, "limit": string(limit), "offset": string(offset)}, true)
    } else {
        resp, err = Query("trades", Params{"pair": pair}, false)
    }

    trades := make(TradesStruct)
    if err == nil {
        for pair, val := range resp {
            trades[pair] = trade2struct(val)
        }
    }

    return trades
}

func OrderTrades(id int) OrderTradesStruct {
    resp, err := Query("order_trades", Params{"order_id": string(id)}, true)

    var result OrderTradesStruct
    if err == nil {
        result.Type = resp["type"].(string)

        result.In.Currency = resp["in_currency"].(string)
        result.In.Amount = resp["in_amount"].(float64)

        result.Out.Currency = resp["out_currency"].(string)
        result.Out.Amount = resp["out_amount"].(float64)

        result.Trades = trade2struct(resp["trades"])
    }

    return result
}

func trade2struct(val interface{}) TradeSlice {
    var lots TradeSlice
    for _, trade := range val.([]interface{}) {
        t := trade.(map[string]interface{})
        var lot TradeStruct

        lot.TradeID = interface2int(t["trade_id"])
        lot.Date = interface2int(t["date"])
        lot.Type = t["type"].(string)
        lot.Quantity = interface2float(t["quantity"])
        lot.Amount = interface2float(t["amount"])

        if _, ok := t["order_id"]; ok {
            lot.Pair = t["pair"].(string)
            lot.OrderID = interface2int(t["order_id"])
            lot.Price = t["price"].(float64)
        }

        lots = append(lots, lot)
    }

    return lots
}
