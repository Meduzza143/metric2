package handlers

// import (
// 	"net/http"
// 	"net/http/httptest"
// 	"testing"

// 	"github.com/gorilla/mux"
// 	"github.com/stretchr/testify/assert"
// )

// func tester() {

// }

// func TestUpdateHandle(t *testing.T) {
// 	testCases := []struct {
// 		method   string
// 		request  string
// 		wantCode int
// 		varName  string
// 		varType  string
// 		varValue string
// 		testName string
// 	}{
// 		{
// 			testName: "test 1 (OK)",
// 			method:   http.MethodPost,
// 			request:  "/update/",
// 			wantCode: 200,
// 			varName:  "var_1",
// 			varType:  "gauge",
// 			varValue: "123",
// 		},
// 		{
// 			testName: "test 2 (400)",
// 			method:   http.MethodPost,
// 			request:  "/update/",
// 			wantCode: 400,
// 			varName:  "var_2",
// 			varType:  "wrong_type",
// 			varValue: "123",
// 		},
// 	}

// 	for _, tc := range testCases {
// 		t.Run(tc.testName, func(t *testing.T) {
// 			req, err := http.NewRequest(tc.method, tc.request, nil)
// 			if err != nil {
// 				t.Fatal(err)
// 			}
// 			vars := map[string]string{
// 				"name":  tc.varName,
// 				"value": tc.varValue,
// 				"type":  tc.varType,
// 			}
// 			req = mux.SetURLVars(req, vars)
// 			w := httptest.NewRecorder()
// 			handler := http.HandlerFunc(UpdateHandle)
// 			handler.ServeHTTP(w, req)

// 			assert.Equal(t, tc.wantCode, w.Code)
// 		})
// 	}
// }
