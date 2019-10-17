package server

import (
	"encoding/json"
	"hack_train/services"
	"net/http"
	"strings"
	"time"
)

type TrainServer struct {
}

type routeSuggestion struct {
	Mode        string    `json:"mode"`
	Origin      string    `json:"origin"`
	Destination string    `json:"destination"`
	Departure   time.Time `json:"departure"`
	Duration    int       `json:"duration"`
	Footprint   int       `json:"footprint"`
}

type errorReply struct {
	Status       int    `json:"status"`
	ErrorMessage string `json:"error_message"`
}

const (
	reqOrigin      = "origin"
	reqDestination = "destination"
	reqDeparture   = "departure"
)

func respondWithError(writer http.ResponseWriter, status int, message string) {
	writer.WriteHeader(status)
	errorMessage := errorReply{
		Status:       status,
		ErrorMessage: message,
	}
	jsonMessage, _ := json.Marshal(errorMessage)
	_, _ = writer.Write(jsonMessage)
}

var getRoutes = func(writer http.ResponseWriter, request *http.Request) {
	validationErrors := []string{}
	_ = request.ParseForm()
	origin := request.FormValue(reqOrigin)
	if origin == "" {
		validationErrors = append(validationErrors, "origin is required")
	}

	destination := request.FormValue(reqDestination)
	if destination == "" {
		validationErrors = append(validationErrors, "destination is required")
	}

	departure := request.FormValue(reqDeparture)
	if departure == "" {
		validationErrors = append(validationErrors, "departure is required")
	}

	departureTime, err := time.Parse(time.RFC3339, departure)
	if err != nil {
		validationErrors = append(validationErrors, "unable to parse departure date: "+err.Error())
	}

	if len(validationErrors) > 0 {
		respondWithError(writer, http.StatusBadRequest, strings.Join(validationErrors, "; "))
		return
	}

	suggestion := []routeSuggestion{}
	trainService := services.NewTrainService()
	for _, t := range trainService.GetRoutes(origin, destination, departureTime) {
		suggestion = append(suggestion, routeSuggestion{
			Mode:        "train",
			Origin:      t.Origin,
			Destination: t.Destination,
			Departure:   t.Departure,
			Duration:    t.Duration,
			Footprint:   t.Footprint,
		})
	}
	reply, err := json.Marshal(suggestion)
	if err != nil {
		_, _ = writer.Write([]byte(err.Error()))
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	_, _ = writer.Write(reply)
}

func NewServer() *TrainServer {
	return &TrainServer{}
}

func (t *TrainServer) Run() error {
	http.HandleFunc("/", getRoutes)
	return http.ListenAndServe(":8080", nil)
}
