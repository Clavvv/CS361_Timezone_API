package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandleRequest(t *testing.T) {
	/*
		Citation:
		Copied/Adapted from -> https://pkg.go.dev/testing
		Reason: Wanted to use the equivalent of assert for Python testing but for Go.
	*/

	tests := []struct {
		name             string //defining format of structs
		requestBody      RequestData
		expectedStatus   int
		expectedResponse ResponseData
		// test cases
	}{
		{
			name: "Valid Timezones",
			requestBody: RequestData{
				CurrentTimezone:     "EST",
				DestinationTimeZone: "PST",
			},
			expectedStatus: http.StatusOK,
			expectedResponse: ResponseData{
				Status:  "success",
				Message: "converting from EST to PST",
				Time:    "-3",
			},
		},
		{
			name: "Invalid Current Timezone",
			requestBody: RequestData{
				CurrentTimezone:     "XYZ",
				DestinationTimeZone: "PST",
			},
			expectedStatus: http.StatusBadRequest,
			expectedResponse: ResponseData{
				Status:  "error",
				Message: "invalid timezone",
				Time:    "",
			},
		},
		{
			name: "Invalid Destination Timezone",
			requestBody: RequestData{
				CurrentTimezone:     "EST",
				DestinationTimeZone: "XYZ",
			},
			expectedStatus: http.StatusBadRequest,
			expectedResponse: ResponseData{
				Status:  "error",
				Message: "invalid timezone",
				Time:    "",
			},
		},
		{
			name: "Invalid JSON",
			requestBody: RequestData{
				CurrentTimezone:     "EST",
				DestinationTimeZone: "",
			},
			expectedStatus: http.StatusBadRequest,
			expectedResponse: ResponseData{
				Status:  "error",
				Message: "invalid timezone",
				Time:    "",
			},
		},
	}

	for _, tt := range tests {

		/*
			Citation:
			ADAPTED FROM -> https://pkg.go.dev/net/http/httptest
			Reason: Documentation for testing http server via Go's testing features
		*/
		t.Run(tt.name, func(t *testing.T) {
			reqBody, err := json.Marshal(tt.requestBody)
			if err != nil {
				t.Fatalf("could not marshal request body: %v", err) //Marshal is standard term for being able to cast a type as another usually JSON, XML to a data type in Go
			}

			req := httptest.NewRequest(http.MethodPost, "/time", bytes.NewReader(reqBody)) // calls via http request
			w := httptest.NewRecorder()

			handleRequest(w, req)

			if status := w.Result().StatusCode; status != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, status)
			}

			var response ResponseData
			if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
				t.Fatalf("could not decode response body: %v", err)
			}

			if response != tt.expectedResponse {
				t.Errorf("expected response %v, got %v", tt.expectedResponse, response)
			}
		})
	}
}
