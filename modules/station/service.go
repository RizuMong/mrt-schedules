package station

import (
	"encoding/json"
	"io"
	"net/http"
	"time"
)

type Service interface {
	GetAllStations() (response []StationResponse, err error)
}

type service struct {
	client *http.Client
}

func NewService() Service {
	return &service{
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (s *service) GetAllStations() (response []StationResponse, err error) {
	url := "https://jakartamrt.co.id/id/val/stasiuns"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	byteResponse, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var stations []Station
	err = json.Unmarshal(byteResponse, &stations)

	if err != nil {
		return nil, err
	}

	for _, item := range stations {
		response = append(response, StationResponse(item))
	}

	return response, nil
}
