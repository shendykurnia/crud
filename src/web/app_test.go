package main

import (
    "testing"
    "net/http"
    "net/http/httptest"
    "bytes"
    . "mymodel"
)

func Test(t *testing.T) {
    var datastore Datastore = &MockDatastore{}
    stack := Stack{&datastore, nil, nil}

    myHandler := *initApp(&stack)

    for _, testCase := range []struct {
        method string
        path string
        payload []byte
        expectedStatus int
        expectedBody string
    }{
        {"GET", "/orders", nil, http.StatusOK, `{"data":{"orders":[],"page":1},"status":"success"}`},
        {"POST", "/orders", []byte(`{"shop_id":1,"customer_id":1,"products":[{"id":1},{"id":2}]}`), http.StatusOK, `{"data":{"id":1,"shop_id":1,"customer_id":1,"status":"created","products":[{"id":1,"name":"Table","created":"2017-03-19T00:00:00Z"},{"id":2,"name":"Chair","created":"2017-03-19T00:00:00Z"}],"created":"2017-03-19T00:00:00Z"},"status":"success"}`},
        {"GET", "/orders", nil, http.StatusOK, `{"data":{"orders":[{"id":1,"shop_id":1,"customer_id":1,"status":"created","products":[{"id":1,"name":"Table","created":"2017-03-19T00:00:00Z"},{"id":2,"name":"Chair","created":"2017-03-19T00:00:00Z"}],"created":"2017-03-19T00:00:00Z"}],"page":1},"status":"success"}`},
        {"PUT", "/orders/1/process", nil, http.StatusOK, `{"status":"success"}`},
        {"GET", "/orders", nil, http.StatusOK, `{"data":{"orders":[{"id":1,"shop_id":1,"customer_id":1,"status":"processed","products":[{"id":1,"name":"Table","created":"2017-03-19T00:00:00Z"},{"id":2,"name":"Chair","created":"2017-03-19T00:00:00Z"}],"created":"2017-03-19T00:00:00Z"}],"page":1},"status":"success"}`},
        {"PUT", "/orders/1/cancel", nil, http.StatusOK, `{"status":"success"}`},
        {"GET", "/orders", nil, http.StatusOK, `{"data":{"orders":[{"id":1,"shop_id":1,"customer_id":1,"status":"canceled","products":[{"id":1,"name":"Table","created":"2017-03-19T00:00:00Z"},{"id":2,"name":"Chair","created":"2017-03-19T00:00:00Z"}],"created":"2017-03-19T00:00:00Z"}],"page":1},"status":"success"}`},
        {"PUT", "/orders/1/finish", nil, http.StatusOK, `{"message":"Object error: Invalid status change","status":"error"}`},
        {"POST", "/orders", []byte(`{"shop_id":2,"customer_id":2,"products":[{"id":1}]}`), http.StatusOK, `{"data":{"id":2,"shop_id":2,"customer_id":2,"status":"created","products":[{"id":1,"name":"Table","created":"2017-03-19T00:00:00Z"}],"created":"2017-03-19T00:00:00Z"},"status":"success"}`},
        {"GET", "/orders", nil, http.StatusOK, `{"data":{"orders":[{"id":1,"shop_id":1,"customer_id":1,"status":"canceled","products":[{"id":1,"name":"Table","created":"2017-03-19T00:00:00Z"},{"id":2,"name":"Chair","created":"2017-03-19T00:00:00Z"}],"created":"2017-03-19T00:00:00Z"},{"id":2,"shop_id":2,"customer_id":2,"status":"created","products":[{"id":1,"name":"Table","created":"2017-03-19T00:00:00Z"}],"created":"2017-03-19T00:00:00Z"}],"page":1},"status":"success"}`},
    } {
        req, err := http.NewRequest(testCase.method, testCase.path, bytes.NewBuffer(testCase.payload))
        if err != nil {
            t.Fatal(err)
        }
        if testCase.payload != nil {
            req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
        }

        rr := httptest.NewRecorder()

        myHandler.ServeHTTP(rr, req)

        if status := rr.Code; status != testCase.expectedStatus {
            t.Errorf("Handler for %v %v %v returned wrong HTTP status code: got %v expected %v", testCase.method, testCase.path, testCase.payload, status, testCase.expectedStatus)
        }

        if testCase.expectedBody != "" && rr.Body.String() != testCase.expectedBody {
            t.Errorf("Handler for %v %v %v returned unexpected body: got %v expected %v", testCase.method, testCase.path, string(testCase.payload), rr.Body.String(), testCase.expectedBody)
        }
    }
}