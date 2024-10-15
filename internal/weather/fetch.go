package weather

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

var (
	RequestBase        = RequestBuilder()
	ErrNotExistingCity = errors.New("1006")
)

type Weather struct {
	Location struct {
		Name string `json:"name"`
	}
	Current struct {
		Temp_c    float32 `json:"temp_c"`
		Is_day    int     `json:"is_day"`
		Condition struct {
			Text string `json:"text"`
		}
	}
	Error struct {
		Code int `json:"code"`
	}
}

func RequestBuilder() string {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	return fmt.Sprintf("http://api.weatherapi.com/v1/current.json?key=%s", os.Getenv("WEATHER_TOKEN"))
}

func WeatherRequest(city string) (Weather, error) {
	var res Weather

	resp, err := http.Get(fmt.Sprintf("%s&q=%s", RequestBase, city))
	if err != nil {
		return res, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return res, err
	}

	if err := json.Unmarshal(body, &res); err != nil {
		return res, err
	}

	return res, nil
}

func CheckCityExists(city string) (string, error) {
	resp, err := WeatherRequest(city)
	if err != nil {
		return "", err
	}

	if resp.Error.Code == 1006 {
		return "", ErrNotExistingCity
	}

	return resp.Location.Name, nil
}
