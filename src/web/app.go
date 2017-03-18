package main

import (
    "net/http"
    . "myrouter"
    . "mymodel"
    . "myutil"
    "strconv"
)

type MyResponse struct {
    status bool
    data interface{}
}

func main() {
    myHandler := MyHandler{}

    myHandler.HandleFunc("^/orders$", []string{"GET"}, func(w http.ResponseWriter, r *http.Request, slugs *map[string]string) {
        query := r.URL.Query()
        search := query["search"]
        page := query["page"]

        pageInt := 1
        if len(page) > 0 {
            pageInt, err := strconv.Atoi(page[0])
            if pageInt < 1 || err != nil {
                pageInt = 1
            }
        }

        var searchStr string
        if len(search) > 0 {
            searchStr = search[0]
        }

        myResponse := MyResponse{false, nil}
        orders, err := Order.Get(searchStr, pageInt)
        if err == nil {
            myResponse.status = true
            myResponse.data = map[string]interface{}{
                "orders": orders,
                "page": pageInt,
            }

            nextQuery := r.URL.Query()
            nextQuery.Set("page", string(pageInt + 1))
            nextParams := [][]string{}
            for k, v := range nextQuery {
                for _, _v := range v {
                    nextParams = append(nextParams, []string{k, _v})
                }
            }
            if url, err := ConstructUrl(r.URL.Path, nextParams); err == nil {
                _map := &myResponse.data.(map[string]interface{})
                (*_map)["next"] = url
            }

            prevQuery := r.URL.Query()
            if pageInt > 1 {
                prevQuery.Set("page", string(pageInt - 1))
            }
            prevParam := [][]string{}
            for k, v := range prevQuery {
                for _, _v := range v {
                    prevParam = append(prevParam, []string{k, _v})
                }
            }
            if url, err := ConstructUrl(r.URL.Path, prevParam); err == nil {
                myResponse.data["prev"] = url
            }
        }

        jsonStr, _ := json.Marshal(myResponse)
        w.Header().Set("Content-Type", "application/json")
        fmt.Fprintf(w, string(jsonStr))
    })

    myHandler.HandleFunc("^/orders$", []string{"POST"}, func(w http.ResponseWriter, r *http.Request, slugs *map[string]string) {
        
    })

    myHandler.HandleFunc("^/orders/(?P<id>.*)/(?P<action>.*)$",
        []string{"PUT": true}, func(w http.ResponseWriter, r *http.Request, slugs *map[string]string) {
        action, ok := (*slugs)["action"]
        if !ok {
            http.NotFound(w, r)
        }

        switch action {
        case "process":

        case "cancel":

        case "finish":

        default:
            http.NotFound(w, r)
        }
    })

    http.ListenAndServe(":9000", &myHandler)
}