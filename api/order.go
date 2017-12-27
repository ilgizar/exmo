package api

import (
    "fmt"
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


func Open(pair string, quantity, price float64, type_name string) OrderResultStruct {
    return orderResult(Query("order_create", Params{"pair": pair, "quantity": fmt.Sprintf("%f", quantity), "price": fmt.Sprintf("%f", price), "type": type_name}, true))
}

func Cancel(id int) OrderResultStruct {
    return orderResult(Query("order_cancel", Params{"order_id": string(id)}, true))
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

            order.Ask.Quantity = interface2float(v["ask_quantity"])
            order.Ask.Amount   = interface2float(v["ask_amount"])
            order.Ask.Top      = interface2float(v["ask_top"])
            order.Ask.Glass    = ab2slice(v["ask"])

            order.Bid.Quantity = interface2float(v["bid_quantity"])
            order.Bid.Amount   = interface2float(v["bid_amount"])
            order.Bid.Top      = interface2float(v["bid_top"])
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
            res.OrderID = interface2int(val)
        }
    } else {
        res.Error = err.Error()
    }

    return res
}

func orderList(opened bool, resp Response, err error) OrdersStruct {
    var orders OrdersStruct
    if err == nil {
        for pair, val := range resp {
            var lots OrderSlice
            for _, item := range val.([]interface{}) {
                t := item.(map[string]interface{})
                var lot OrderStruct

                lot.OrderID  = interface2int(t["order_id"])
                if opened {
                    lot.Date = interface2int(t["created"])
                    lot.Type = t["type"].(string)
                } else {
                    lot.Date = interface2int(t["date"])
                    lot.Type = t["order_type"].(string)
                }
                lot.Pair     = t["pair"].(string)
                lot.Quantity = t["quantity"].(float64)
                lot.Price    = t["price"].(float64)
                lot.Amount   = t["amount"].(float64)

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
        ab.Quantity = interface2float(v[0])
        ab.Amount   = interface2float(v[1])
        ab.Price    = interface2float(v[2])

        result = append(result, ab)
    }

    return result
}
