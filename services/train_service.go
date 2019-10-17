package services

import "time"

type TrainService struct {
}

type TrainRoute struct {
	Origin      string
	Destination string
	Departure   time.Time
	Duration    int
	Footprint   int
}

func NewTrainService() *TrainService {
	return &TrainService{}
}

func (ts *TrainService) GetRoutes(origin string, destination string, departure time.Time) []TrainRoute {
	return []TrainRoute{
		{
			Origin:      origin + " a-reeno",
			Destination: destination + " a-who",
			Departure:   departure.Add(1*time.Hour + 7*time.Minute + 11*time.Second),
			Duration:    2*60*60 + 31*60,
			Footprint:   123,
		},
	}
}
