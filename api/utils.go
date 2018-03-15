package api

import (
    "fmt"
    "strconv"
)

func Interface2Int(value interface{}) int {
    res, _ := strconv.ParseInt(fmt.Sprintf("%v", value), 10, 64)

    return int(res)
}

func Interface2Float(value interface{}) float64 {
    res, _ := strconv.ParseFloat(fmt.Sprintf("%v", value), 64)

    return res
}

func Interface2Map(value interface{}) map[string]float64 {
    res := make(map[string]float64)

    for key, val := range value.(map[string]interface{}) {
        res[key], _ = strconv.ParseFloat(val.(string), 64)
    }

    return res
}
