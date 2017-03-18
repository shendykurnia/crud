package main

import (
    "net/http"
    . "myrouter"
    . "mymodel"
    "strconv"
    "encoding/json"
    "fmt"
)

func main() {
    myHandler := MyHandler{}

    // Get and search API
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

        defaultResponse := map[string]interface{}{
            "status": "error",
        }
        myResponse := defaultResponse
        orders, err := GetOrders(searchStr, pageInt)
        if err == nil {
            myResponse = map[string]interface{}{
                "status": "success",
                "data": map[string]interface{}{
                    "orders": orders,
                    "page": pageInt,
                },
            }
        }

        respondWithJson(&w, myResponse)
    })

    // Create API
    myHandler.HandleFunc("^/orders$", []string{"POST"}, func(w http.ResponseWriter, r *http.Request, slugs *map[string]string) {
        defaultResponse := map[string]interface{}{
            "status": "error",
        }
        myResponse := defaultResponse

        r.ParseForm()
        productJson := ProductJson{}
        for key, _ := range r.Form {
            err := json.Unmarshal([]byte(key), &productJson)
            if err != nil {
                myResponse["error"] = "Invalid input"
                respondWithJson(&w, myResponse)
                return
            }
            break
        }

        order, err := CreateOrder(&productJson)
        if err == nil {
            myResponse = map[string]interface{}{
                "status": "success",
                "data": order,
            }
        }

        respondWithJson(&w, myResponse)
    })

    // Process, cancel, and finalize API
    myHandler.HandleFunc("^/orders/(?P<id>.*)/(?P<action>.*)$",
        []string{"PUT"}, func(w http.ResponseWriter, r *http.Request, slugs *map[string]string) {
        action, ok := (*slugs)["action"]
        if !ok {
            http.NotFound(w, r)
        }

        defaultResponse := map[string]interface{}{
            "status": "error",
        }
        myResponse := defaultResponse
        id, ok := (*slugs)["id"]
        if !ok {
            http.NotFound(w, r)
        }

        isValidOrder := false
        idInt, err := strconv.Atoi(id)
        if err == nil {
            if order, err := GetOrder(idInt); order != nil && err == nil {
                isValidOrder = true
            }
        }

        if !isValidOrder {
            myResponse = map[string]interface{}{
                "status": "error",
                "error": "Order not found",
            }
        } else {
            statusMap := map[string]string{
                "process": StatusProcessed,
                "cancel": StatusCanceled,
                "finish": StatusFinished,
            }

            if _, ok := statusMap[action]; !ok {
                http.NotFound(w, r)
                return
            }

            for _action, _status := range statusMap {
                if _action != action {
                    continue
                }

                if err := UpdateOrderStatus(idInt, _status); err == nil {
                    myResponse["status"] = "success"
                    break
                }
            }
        }

        respondWithJson(&w, myResponse)
    })

    http.ListenAndServe(":9000", &myHandler)
}

func respondWithJson(w *http.ResponseWriter, obj interface{}) {
    jsonStr, _ := json.Marshal(obj)
    (*w).Header().Set("Content-Type", "application/json")
    fmt.Fprintf((*w), string(jsonStr))
}