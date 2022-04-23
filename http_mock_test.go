package owmonecall

import (
	"io"
	"os"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/shopspring/decimal"
)

func TestWithDefaultParameters(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	// Exact URL match
	httpmock.RegisterResponder(
		"GET",
		"https://api.openweathermap.org/data/2.5/onecall?lat=3.14&lon=5.67&appid=c8b58ab0-1360-4a3a-9f70-3742e48ad2fe",
		httpmock.NewStringResponder(200, readTextFile(t, "testdata/sample.json")),
	)

	// do stuff that makes a request to articles
	lat, err := decimal.NewFromString("3.14")
	if err != nil {
		t.Fatal("failed to parse decimal")
	}
	long, err := decimal.NewFromString("5.67")
	if err != nil {
		t.Fatal("failed to parse decimal")
	}
	appId := "c8b58ab0-1360-4a3a-9f70-3742e48ad2fe"
	actual, err := OneCall(lat, long, appId)
	if err != nil {
		t.Fatalf("API call failure: %v", err.Error())
	}

	assertEqual(t, expectedSampleJson(), actual)

	// get count info
	callCount := httpmock.GetTotalCallCount()
	assertEqual(t, 1, callCount)
}

func TestWithOptionalParameters(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	// Exact URL match
	httpmock.RegisterResponder(
		"GET",
		"https://api.openweathermap.org/data/2.5/onecall?lat=3.14&lon=5.67&appid=c8b58ab0-1360-4a3a-9f70-3742e48ad2fe&excludes=foo,bar&units=metric&lang=en",
		httpmock.NewStringResponder(200, readTextFile(t, "testdata/sample.json")),
	)

	// do stuff that makes a request to articles
	lat, err := decimal.NewFromString("3.14")
	if err != nil {
		t.Fatal("failed to parse decimal")
	}
	long, err := decimal.NewFromString("5.67")
	if err != nil {
		t.Fatal("failed to parse decimal")
	}
	appId := "c8b58ab0-1360-4a3a-9f70-3742e48ad2fe"
	actual, err := OneCall(lat, long, appId, Exclude([]string{"foo", "bar"}), MetricUnits(), Lang("en"))
	if err != nil {
		t.Fatalf("API call failure: %v", err.Error())
	}

	assertEqual(t, expectedSampleJson(), actual)

	// get count info
	callCount := httpmock.GetTotalCallCount()
	assertEqual(t, 1, callCount)
}

func readTextFile(t *testing.T, path string) string {
	t.Helper()
	file, err := os.Open(path)
	if err != nil {
		panic("failed to open file")
	}
	defer file.Close()
	bytes, err := io.ReadAll(file)
	if err != nil {
		panic("failed to read file")
	}
	return string(bytes)
}
