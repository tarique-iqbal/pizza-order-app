package geocoder

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"restaurant-service/internal/domain/restaurant"
)

type openCageGeocoder struct {
	apiKey string
}

func NewOpenCageGeocoder(apiKey string) restaurant.Geocoder {
	return &openCageGeocoder{apiKey: apiKey}
}

func (s *openCageGeocoder) GeocodeAddress(addr restaurant.RestaurantAddress) (float64, float64, error) {
	query := fmt.Sprintf("%s %s, %s %s", addr.House, addr.Street, addr.PostalCode, addr.City)
	endpoint := fmt.Sprintf("https://api.opencagedata.com/geocode/v1/json?q=%s&key=%s", url.QueryEscape(query), s.apiKey)

	resp, err := http.Get(endpoint)
	if err != nil {
		return 0, 0, err
	}
	defer resp.Body.Close()

	var data struct {
		Results []struct {
			Geometry struct {
				Lat float64 `json:"lat"`
				Lng float64 `json:"lng"`
			} `json:"geometry"`
		} `json:"results"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return 0, 0, err
	}

	if len(data.Results) == 0 {
		return 0, 0, fmt.Errorf("no geocoding results found")
	}

	return data.Results[0].Geometry.Lat, data.Results[0].Geometry.Lng, nil
}
