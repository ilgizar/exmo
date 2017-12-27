package api

import (
    "bytes"
    "crypto/hmac"
    "crypto/sha512"
    "encoding/json"
    "errors"
    "fmt"
    "io/ioutil"
    "net/http"
    "net/url"
    "strconv"
    "time"
)

type Response map[string]interface{}
type Params map[string]string

var apiPublicKey string
var apiSecretKey []byte

func Init(key, secret string) {
    apiPublicKey = key
    apiSecretKey = []byte(secret)
}

func Query(method string, params Params, signed bool) (Response, error) {
    post_params := url.Values{}
    if signed {
        post_params.Add("nonce", nonce())
    }
    if params != nil {
        for key, value := range params {
            post_params.Add(key, value)
        }
    }
    post_content := post_params.Encode()

    req, _ := http.NewRequest("POST", "https://api.exmo.com/v1/" + method, bytes.NewBuffer([]byte(post_content)))
    if signed {
        req.Header.Set("Key", apiPublicKey)
        req.Header.Set("Sign", do_sign(post_content))
    }
    req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
    req.Header.Add("Content-Length", strconv.Itoa(len(post_content)))

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    if resp.Status != "200 OK" {
        return nil, errors.New("http status: " + resp.Status)
    }

    body, err1 := ioutil.ReadAll(resp.Body)
    if err1 != nil {
        return nil, err1
    }

    var dat map[string]interface{}
    err2 := json.Unmarshal([]byte(body), &dat)
    if err2 != nil {
        return nil, err2
    }

    if result, ok := dat["result"]; ok && result.(bool) != true {
        return nil, errors.New(dat["error"].(string))
    }

    return dat, nil
}

func nonce() string {
    return fmt.Sprintf("%d", time.Now().UnixNano())
}

func do_sign(message string) string {
    mac := hmac.New(sha512.New, apiSecretKey)
    mac.Write([]byte(message))
    return fmt.Sprintf("%x", mac.Sum(nil))
}
