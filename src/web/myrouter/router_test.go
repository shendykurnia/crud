package myrouter

import (
    "testing"
    "net/http"
    "net/http/httptest"
    "fmt"
)

func Test(t *testing.T) {
    myHandler := MyHandler{}

    myHandler.HandleFunc("^/a$", map[string]interface{}{"GET": true}, func(w http.ResponseWriter, r *http.Request, slugs *map[string]string) {
        fmt.Fprintf(w, "1")
    })

    myHandler.HandleFunc("^/a$", map[string]interface{}{"POST": true}, func(w http.ResponseWriter, r *http.Request, slugs *map[string]string) {
        fmt.Fprintf(w, "2")
    })

    myHandler.HandleFunc("^/a/{id}/b$", map[string]interface{}{"PUT": true}, func(w http.ResponseWriter, r *http.Request, slugs *map[string]string) {
        fmt.Fprintf(w, "3")
    })

    for _, testCase := range []struct {
        method string
        path string
        expectedStatus int
        expectedBody string
    }{
        {"GET", "/a", http.StatusOK, "1"},
        {"POST", "/a", http.StatusOK, "2"},
        {"PUT", "/a/{id}/b", http.StatusOK, "3"},
        {"GET", "/a/b", http.StatusNotFound, ""},
        {"PUT", "/a/3/b", http.StatusOK, "3"},
    } {
        req, err := http.NewRequest(testCase.method, testCase.path, nil)
        if err != nil {
            t.Fatal(err)
        }

        rr := httptest.NewRecorder()

        myHandler.ServeHTTP(rr, req)

        if status := rr.Code; status != testCase.expectedStatus {
            t.Errorf("Handler for %v %v returned wrong HTTP status code: got %v expected %v", testCase.method, testCase.path, status, testCase.expectedStatus)
        }

        if testCase.expectedBody != "" && rr.Body.String() != testCase.expectedBody {
            t.Errorf("Handler for %v %v returned unexpected body: got %v expected %v", testCase.method, testCase.path, rr.Body.String(), testCase.expectedBody)
        }
    }
}