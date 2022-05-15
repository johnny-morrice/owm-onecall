package owmonecall

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

// Use the OneCall API
//
// * https://openweathermap.org/api/one-call-api
// * https://api.openweathermap.org/data/2.5/onecall?lat={lat}&lon={lon}&exclude={part}&appid={API key}

func OneCall(lat, long decimal.Decimal, appId string, optionals ...OptionalParameter) (*OneCallResponse, error) {
	latText := url.PathEscape(lat.String())
	longText := url.PathEscape(long.String())
	optQuery := ""
	for _, opt := range optionals {
		optQuery += "&"
		optQuery += fmt.Sprintf("%s=%s", url.PathEscape(opt.Name), url.PathEscape(opt.Value))
	}
	url := fmt.Sprintf("https://api.openweathermap.org/data/2.5/onecall?lat=%s&lon=%s&appid=%s", latText, longText, appId)
	url += optQuery
	client := &http.Client{}
	resp, err := client.Get(url)
	if err != nil {
		return nil, errors.Wrap(err, "API call failed")
	}
	out := &OneCallResponse{}
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(out)
	if err != nil {
		return nil, errors.Wrap(err, "JSON decoding failed")
	}
	return out, nil
}

type OptionalParameter struct {
	Name  string
	Value string
}

func Lang(lang string) OptionalParameter {
	return OptionalParameter{
		Name:  "lang",
		Value: lang,
	}
}

func Exclude(excludes []string) OptionalParameter {
	joined := strings.Join(excludes, ",")
	return OptionalParameter{
		Name:  "excludes",
		Value: joined,
	}
}

func ImperialUnits() OptionalParameter {
	return OptionalParameter{
		Name:  "units",
		Value: "imperial",
	}
}

func MetricUnits() OptionalParameter {
	return OptionalParameter{
		Name:  "units",
		Value: "metric",
	}
}

func StandardUnits() OptionalParameter {
	return OptionalParameter{
		Name:  "units",
		Value: "standard",
	}
}

type Weather struct {
	ID          int64  `json:"id"`
	Main        string `json:"main"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
}

type CurrentWeather struct {
	Dt         int64                      `json:"dt"`
	Sunrise    int64                      `json:"sunrise"`
	Sunset     int64                      `json:"sunset"`
	Temp       decimal.Decimal            `json:"temp"`
	FeelsLike  decimal.Decimal            `json:"feels_like"`
	Pressure   int64                      `json:"pressure"`
	Humidity   int64                      `json:"humidity"`
	DewPoint   decimal.Decimal            `json:"dew_point"`
	Uvi        decimal.Decimal            `json:"uvi"`
	Clouds     int64                      `json:"clouds"`
	Visibility int64                      `json:"visibility"`
	WindSpeed  decimal.Decimal            `json:"wind_speed"`
	WindDeg    int64                      `json:"wind_deg"`
	Weather    []Weather                  `json:"weather"`
	Rain       map[string]decimal.Decimal `json:"rain"`
}

type Minutely struct {
	Dt            int64           `json:"dt"`
	Precipitation decimal.Decimal `json:"precipitation"`
}

type Hourly struct {
	Dt         int64           `json:"dt"`
	Temp       decimal.Decimal `json:"temp"`
	FeelsLike  decimal.Decimal `json:"feels_like"`
	Pressure   int64           `json:"pressure"`
	Humidity   int64           `json:"humidity"`
	DewPoint   decimal.Decimal `json:"dew_point"`
	Uvi        decimal.Decimal `json:"uvi"`
	Clouds     int64           `json:"clouds"`
	Visibility int64           `json:"visibility"`
	WindSpeed  decimal.Decimal `json:"wind_speed"`
	WindDeg    int64           `json:"wind_deg"`
	WindGust   decimal.Decimal `json:"wind_gust"`
	Weather    []Weather       `json:"weather"`
	Pop        decimal.Decimal `json:"pop"`
}

type DayTemp struct {
	Day   decimal.Decimal `json:"day"`
	Min   decimal.Decimal `json:"min"`
	Max   decimal.Decimal `json:"max"`
	Night decimal.Decimal `json:"night"`
	Eve   decimal.Decimal `json:"eve"`
	Morn  decimal.Decimal `json:"morn"`
}

type DayFeelsLike struct {
	Day   decimal.Decimal `json:"day"`
	Night decimal.Decimal `json:"night"`
	Eve   decimal.Decimal `json:"eve"`
	Morn  decimal.Decimal `json:"morn"`
}

type Daily struct {
	Dt        int64           `json:"dt"`
	Sunrise   int64           `json:"sunrise"`
	Sunset    int64           `json:"sunset"`
	Moonrise  int64           `json:"moonrise"`
	Moonset   int64           `json:"moonset"`
	MoonPhase decimal.Decimal `json:"moon_phase"`
	Temp      DayTemp         `json:"temp"`
	FeelsLike DayFeelsLike    `json:"feels_like"`
	Pressure  int64           `json:"pressure"`
	Humidity  int64           `json:"humidity"`
	DewPoint  decimal.Decimal `json:"dew_point"`
	WindSpeed decimal.Decimal `json:"wind_speed"`
	WindDeg   decimal.Decimal `json:"wind_deg"`
	Weather   []Weather       `json:"weather"`
	Clouds    int64           `json:"clouds"`
	Pop       decimal.Decimal `json:"pop"`
	Rain      decimal.Decimal `json:"rain"`
	Uvi       decimal.Decimal `json:"uvi"`
}

type Alert struct {
	SenderName string   `json:"sender_name"`
	Event      string   `json:"event"`
	Start      int64    `json:"start"`
	End        int64    `json:"end"`
	Decription string   `json:"description"`
	Tags       []string `json:"tags"`
}

type OneCallResponse struct {
	Lat            decimal.Decimal `json:"lat"`
	Lon            decimal.Decimal `json:"lon"`
	Timezone       string          `json:"timezone"`
	TimezoneOffset int             `json:"timezone_offset"`
	Current        CurrentWeather  `json:"current"`
	Minutely       []Minutely      `json:"minutely"`
	Hourly         []Hourly        `json:"hourly"`
	Daily          []Daily         `json:"daily"`
	Alerts         []Alert         `json:"alerts"`
}
