package myutil

import (
    "net/http"
)

func ConstructUrl(url string, params [][]string) (finalUrl string, err error) {
    req, err := http.NewRequest("GET", url, nil)
    if err != nil {
        return url, err
    }

    q := req.URL.Query()
    for _, param := range params {
        if len(param) == 0 {
            continue
        }

        key := param[0]
        val := ""
        if len(param) > 1 {
            val = param[1]
        }
        q.Add(key, val)
    }
    req.URL.RawQuery = q.Encode()

    return req.URL.String(), nil
}