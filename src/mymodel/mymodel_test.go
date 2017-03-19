package mymodel

import (
    "testing"
)

func TestMockDatastore(t *testing.T) {
    datastore := MockDatastore{}
    err := datastore.init()
    if err != nil {
        t.Fatal(err)
    }

    // datastore.getOrderCollection("", 1)
    // test if changing returned value change underlying data
    // t.Errorf("Handler for %v %v returned wrong HTTP status code: got %v expected %v", testCase.method, testCase.path, status, testCase.expectedStatus)
}