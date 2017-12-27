package api

import (
    "strconv"
)

func interface2int(value interface{}) int {
    var res int
    if val, ok := value.(float64); ok {
        return int(val)
    }

    return res
}

func interface2float(value interface{}) float64 {
    res, _ := strconv.ParseFloat(value.(string), 64)

    return res
}

func interface2map(value interface{}) map[string]float64 {
    res := make(map[string]float64)

    for key, val := range value.(map[string]interface{}) {
        res[key], _ = strconv.ParseFloat(val.(string), 64)
    }

    return res
}
