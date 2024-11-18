package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

type RequestData struct {
	CurrentTimezone     string `json:"currentTimezone"`
	DestinationTimeZone string `json:"destinationTimezone"`
}

type ResponseData struct {
	Status  string `json:"status"`
	Time    string `json:"time_difference"`
	Message string `json:"message"`
}

func handleRequest(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var requestData RequestData

	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		response := ResponseData{
			Status:  "error",
			Message: "Invalid JSON format",
			Time:    "",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode((response))
		return
	}

	timeDifference := convert(requestData.CurrentTimezone, requestData.DestinationTimeZone)
	response := ResponseData{
		Status:  "success",
		Message: "converting from " + requestData.CurrentTimezone + " to " + requestData.DestinationTimeZone,
		Time:    timeDifference,
	}

	if timeDifference == "invalid timezone" || timeDifference == "error converting the timezone" {
		response.Status = "error"
		response.Message = timeDifference
		response.Time = ""
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)

}

func convert(c, d string) string {

	timeZones := map[string]string{
		"EST": "-5",
		"CST": "-6",
		"MST": "-7",
		"PST": "-8",
		"AKT": "-9",
		"HST": "-10",
	}

	currentTime, ok1 := timeZones[c]
	destinationTime, ok2 := timeZones[d]

	if !ok1 || !ok2 {
		return "invalid timezone"
	}

	currentTimeOffset, err1 := strconv.Atoi(currentTime)
	destinationTimeOffset, err2 := strconv.Atoi(destinationTime)

	if err1 != nil || err2 != nil {
		fmt.Println("errror converting the timezone")
		return "error"
	}

	delta := destinationTimeOffset - currentTimeOffset

	if delta > 0 {
		return "+" + strconv.Itoa(delta)
	} else {
		return strconv.Itoa(delta)
	}

}

func main() {

	http.HandleFunc("/time", handleRequest)
	fmt.Println("server is listening on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
