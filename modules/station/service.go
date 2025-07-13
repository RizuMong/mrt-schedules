package station

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
	"time"
)

type Service interface {
	GetAllStations() (response []StationResponse, err error)
	CheckSchedulesByStation(id string) (response []ScheduleResponse, err error)
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

func (s *service) CheckSchedulesByStation(id string) (response []ScheduleResponse, err error) {
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

	var schedule []Schedule
	err = json.Unmarshal(byteResponse, &schedule)

	if err != nil {
		return nil, err
	}

	// schedule selected by station id
	var scheduleSelected Schedule

	for _, item := range schedule {
		if item.StationId == id {
			scheduleSelected = item
			break
		}
	}

	if scheduleSelected.StationId == "" {
		err = errors.New("Station not found")
		return nil, err
	}

	response, err = ConvertDataToResponse(scheduleSelected)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func ConvertDataToResponse(schedule Schedule) (response []ScheduleResponse, err error) {
	var (
		LebakBulusTripName = "Station Lebak Bulus Grab"
		BundaranHITripName = "Station Bundaran HI Bank DKI"
	)

	scheduleLebakBulus := schedule.ScheduleLebakBulus
	scheduleBundaranHI := schedule.ScheduleBundaranHI

	scheduleLebakBulusParsed, err := ConvertScheduleToTimeFormat(scheduleLebakBulus)
	if err != nil {
		return
	}

	scheduleBundaranHIParsed, err := ConvertScheduleToTimeFormat(scheduleBundaranHI)
	if err != nil {
		return
	}

	// convert response
	for _, item := range scheduleLebakBulusParsed {
		if item.Format("15:04") > time.Now().Format("15:04") {
			response = append(response, ScheduleResponse{
				StationName: LebakBulusTripName,
				Time:        item.Format("15:04"),
			})
		}
	}

	for _, item := range scheduleBundaranHIParsed {
		if item.Format("15:04") > time.Now().Format("15:04") {
			response = append(response, ScheduleResponse{
				StationName: BundaranHITripName,
				Time:        item.Format("15:04"),
			})
		}
	}

	return response, nil
}

func ConvertScheduleToTimeFormat(schedule string) (response []time.Time, err error) {
	var (
		parsedTime time.Time
		schedules  = strings.Split(schedule, ",")
	)

	for _, item := range schedules {
		trimmedTime := strings.TrimSpace(item)

		if trimmedTime == "" {
			continue
		}

		// Handle time ranges like "16:46: 16:51" by splitting on ": " and taking individual times
		if strings.Contains(trimmedTime, ": ") {
			times := strings.Split(trimmedTime, ": ")
			for _, timeStr := range times {
				timeStr = strings.TrimSpace(timeStr)
				if timeStr == "" {
					continue
				}

				// Validate that timeStr looks like a time (HH:MM format)
				if len(timeStr) != 5 || timeStr[2] != ':' {
					continue // Skip invalid time strings
				}

				parsedTime, err = time.Parse("15:04", timeStr)
				if err != nil {
					continue // Skip invalid times instead of returning error
				}
				response = append(response, parsedTime)
			}
		} else {
			// Validate that trimmedTime looks like a time (HH:MM format)
			if len(trimmedTime) != 5 || trimmedTime[2] != ':' {
				continue // Skip invalid time strings
			}

			parsedTime, err = time.Parse("15:04", trimmedTime)
			if err != nil {
				continue // Skip invalid times instead of returning error
			}
			response = append(response, parsedTime)
		}
	}

	return
}
