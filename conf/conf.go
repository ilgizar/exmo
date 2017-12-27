package conf

import (
    "io/ioutil"
    "os"

    "github.com/influxdata/toml"
)

func Read(filename string, config interface{}) error {
    f, err := os.Open(filename)
    if err != nil {
        return err
    }
    defer f.Close()

    buf, err := ioutil.ReadAll(f)
    if err != nil {
        return err
    }

    err = toml.Unmarshal(buf, config)

    return err
}
