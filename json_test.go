package owmonecall

import (
	"encoding/json"
	"os"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/shopspring/decimal"
)

func TestDeserialiseResponse(t *testing.T) {
	file, err := os.Open("testdata/sample.json")
	if err != nil {
		t.Fatal("failed to open sample file")
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	actual := &OneCallResponse{}
	err = decoder.Decode(actual)
	if err != nil {
		t.Fatal("failed to decode sample file JSON")
	}
	expected := expectedSampleJson()
	assertEqual(t, expected, actual)
}

func expectedSampleJson() *OneCallResponse {
	return &OneCallResponse{
		Lat:            parseDecimal("33.44"),
		Lon:            parseDecimal("-94.04"),
		Timezone:       "America/Chicago",
		TimezoneOffset: -21600,
		Current: CurrentWeather{
			Dt:         1618317040,
			Sunrise:    1618282134,
			Sunset:     1618333901,
			Temp:       parseDecimal("284.07"),
			FeelsLike:  parseDecimal("282.84"),
			Pressure:   1019,
			Humidity:   62,
			DewPoint:   parseDecimal("277.08"),
			Uvi:        parseDecimal("0.89"),
			Visibility: 10000,
			WindSpeed:  parseDecimal("6"),
			WindDeg:    300,
			Weather:    []Weather{{ID: 500, Main: "Rain", Description: "light rain", Icon: "10d"}},
			Rain:       map[string]decimal.Decimal{"1h": parseDecimal("0.21")},
		},
		Minutely: []Minutely{{Dt: 1618317060, Precipitation: parseDecimal("0.205")}},
		Hourly: []Hourly{
			{
				Dt:         1618315200,
				Temp:       parseDecimal("282.58"),
				FeelsLike:  parseDecimal("280.4"),
				Pressure:   1019,
				Humidity:   68,
				DewPoint:   parseDecimal("276.98"),
				Uvi:        parseDecimal("1.4"),
				Clouds:     19,
				Visibility: 306,
				WindSpeed:  parseDecimal("4.12"),
				WindDeg:    296,
				WindGust:   parseDecimal("7.33"),
				Weather: []Weather{
					{
						ID:          801,
						Main:        "Clouds",
						Description: "few clouds",
						Icon:        "02d",
					}},
				Pop: parseDecimal("1"),
			},
		},
		Daily: []Daily{
			{
				Dt:        1618308000,
				Sunrise:   1618282134,
				Sunset:    1618333901,
				Moonrise:  1618284960,
				Moonset:   1618339740,
				MoonPhase: parseDecimal("0.04"),
				Temp: DayTemp{
					Day:   parseDecimal("279.79"),
					Min:   parseDecimal("275.09"),
					Max:   parseDecimal("284.07"),
					Night: parseDecimal("275.09"),
					Eve:   parseDecimal("279.21"),
					Morn:  parseDecimal("278.49"),
				},
				FeelsLike: DayFeelsLike{
					Day:   parseDecimal("277.59"),
					Night: parseDecimal("276.27"),
					Eve:   parseDecimal("276.49"),
					Morn:  parseDecimal("276.27"),
				},
				Pressure:  1020,
				Humidity:  81,
				DewPoint:  parseDecimal("276.77"),
				WindSpeed: parseDecimal("3.06"),
				WindDeg:   parseDecimal("294"),
				Weather: []Weather{
					{
						ID:          500,
						Main:        "Rain",
						Description: "light rain",
						Icon:        "10d",
					},
				},
				Clouds: 56,
				Pop:    parseDecimal("0.2"),
				Rain:   parseDecimal("0.62"),
				Uvi:    parseDecimal("1.93"),
			},
		},
		Alerts: []Alert{
			{
				SenderName: "NWS Tulsa",
				Event:      "Heat Advisory",
				Start:      1597341600,
				End:        1597366800,
				Decription: strings.Join([]string{
					"...HEAT ADVISORY REMAINS IN EFFECT FROM 1 PM THIS AFTERNOON TO",
					"8 PM CDT THIS EVENING...",
					"* WHAT...Heat index values of 105 to 109 degrees expected.",
					"* WHERE...Creek, Okfuskee, Okmulgee, McIntosh, Pittsburg,",
					"Latimer, Pushmataha, and Choctaw Counties.",
					"* WHEN...From 1 PM to 8 PM CDT Thursday.",
					"* IMPACTS...The combination of hot temperatures and high",
					"humidity will combine to create a dangerous situation in which",
					"heat illnesses are possible.",
				}, "\n"),
				Tags: []string{"Extreme temperature value"},
			},
		},
	}
}

func parseDecimal(num string) decimal.Decimal {
	d, err := decimal.NewFromString(num)
	if err != nil {
		panic("invalid decimal: " + num)
	}
	return d
}

func assertEqual(t *testing.T, expected interface{}, actual interface{}) {
	t.Helper()
	diff := cmp.Diff(expected, actual)
	if diff != "" {
		t.Logf("expected %v but received %v", expected, actual)
		t.Fatal(diff)
	}
}
