package api

import (
    "fmt"
    "strconv"
)

type OrderResultStruct struct {
    Result  bool
    Error   string
    OrderID int
}

type OrderStruct struct {
    OrderID  int
    Date     int
    Type     string
    Pair     string
    Quantity float64
    Price    float64
    Amount   float64
}

type OrderSlice []OrderStruct

type OrdersStruct map[string]OrderSlice

type GlassStruct struct {
    Quantity float64
    Amount   float64
    Price    float64
}

type GlassesSlice []GlassStruct

type AskBidStruct struct {
    Quantity float64
    Amount   float64
    Top      float64
    Glass    GlassesSlice
}

type OrderBookStruct struct {
    Ask  AskBidStruct
    Bid  AskBidStruct
}

type OrderBooksStruct map[string]OrderBookStruct

var OpenTypes = []string{"buy", "sell", "market_buy", "market_sell", "market_buy_total", "market_sell_total"}


func Open(pair string, quantity, price float64, type_name string) OrderResultStruct {
    return orderResult(Query("order_create", Params{"pair": pair, "quantity": fmt.Sprintf("%f", quantity), "price": fmt.Sprintf("%f", price), "type": type_name}, true))
}

func Cancel(id interface{}) OrderResultStruct {
    return orderResult(Query("order_cancel", Params{"order_id": fmt.Sprintf("%v", id)}, true))
}

func Orders() OrdersStruct {
    resp, err := Query("user_open_orders", nil, true)
    return orderList(true, resp, err)
}

func Cancelled(limit, offset int) OrdersStruct {
    resp, err := Query("user_cancelled_orders", Params{"limit": string(limit), "offset": string(offset)}, true)
    return orderList(false, resp, err)
}

func OrderBook(pair string, limit int) OrderBooksStruct {
    resp, err := Query("order_book", Params{"pair": pair, "limit": string(limit)}, false)

    orders := OrderBooksStruct{}
    if err == nil {
        for pair, val := range resp {
            v := val.(map[string]interface{})

            order := OrderBookStruct{}

            order.Ask.Quantity = Interface2Float(v["ask_quantity"])
            order.Ask.Amount   = Interface2Float(v["ask_amount"])
            order.Ask.Top      = Interface2Float(v["ask_top"])
            order.Ask.Glass    = ab2slice(v["ask"])

            order.Bid.Quantity = Interface2Float(v["bid_quantity"])
            order.Bid.Amount   = Interface2Float(v["bid_amount"])
            order.Bid.Top      = Interface2Float(v["bid_top"])
            order.Bid.Glass    = ab2slice(v["bid"])

            orders[pair] = order
        }
    }

    return orders
}

func orderResult(resp Response, err error) OrderResultStruct {
    var res OrderResultStruct
    if err == nil {
        res.Result      = resp["result"].(bool)
        res.Error       = resp["error"].(string)
        if val, ok := resp["order_id"]; ok {
            res.OrderID = Interface2Int(val)
        }
    } else {
        res.Error = err.Error()
    }

    return res
}

func orderList(opened bool, resp Response, err error) OrdersStruct {
    orders := make(OrdersStruct)
    if err == nil {
        for pair, val := range resp {
            var lots OrderSlice
            for _, item := range val.([]interface{}) {
                t := item.(map[string]interface{})
                var lot OrderStruct

                lot.OrderID  = Interface2Int(t["order_id"])
                if opened {
                    lot.Date = Interface2Int(t["created"])
                    lot.Type = t["type"].(string)
                } else {
                    lot.Date = Interface2Int(t["date"])
                    lot.Type = t["order_type"].(string)
                }
                lot.Pair        = t["pair"].(string)
                lot.Quantity, _ = strconv.ParseFloat(t["quantity"].(string), 64)
                lot.Price, _    = strconv.ParseFloat(t["price"].(string), 64)
                lot.Amount, _   = strconv.ParseFloat(t["amount"].(string), 64)

                lots = append(lots, lot)
            }
            orders[pair] = lots
        }
    }

    return orders
}

func ab2slice(value interface{}) GlassesSlice {
    result := GlassesSlice{}

    for _, val := range value.([]interface{}) {
        v := val.([]interface{})

        ab := GlassStruct{}
        ab.Quantity = Interface2Float(v[0])
        ab.Amount   = Interface2Float(v[1])
        ab.Price    = Interface2Float(v[2])

        result = append(result, ab)
    }

    return result
}
