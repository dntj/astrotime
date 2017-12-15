package astrotime

import (
	"fmt"
	"testing"
	"time"
)

type place struct {
	lat, lon float64
	times    []sunData
}

type sunData struct {
	day, sunrise, sunset time.Time
}

var (
	northenSummer  = p("2017-07-10T15:04:05Z")
	southernSummer = p("2017-12-29T15:04:05Z")
	midSeason      = p("2017-10-15T15:04:05Z")

	// Various latitudes.
	places = map[string]place{
		"ushuaia": place{
			lat: -54.8019, lon: -68.3030,
			times: []sunData{
				{day: northenSummer, sunrise: p("2017-07-10T12:51:36Z"), sunset: p("2017-07-10T20:26:12Z")},
				{day: southernSummer, sunrise: p("2017-12-29T07:58:08Z"), sunset: p("2017-12-30T01:12:53Z")},
				{day: midSeason, sunrise: p("2017-10-15T09:21:36Z"), sunset: p("2017-10-15T23:17:10Z")},
			},
		},
		"melbourne": place{
			lat: -37.8136, lon: 144.9631,
			times: []sunData{
				{day: northenSummer, sunrise: p("2017-07-09T21:34:30Z"), sunset: p("2017-07-10T07:16:53Z")},
				{day: southernSummer, sunrise: p("2017-12-28T18:59:38Z"), sunset: p("2017-12-29T09:44:58Z")},
				{day: midSeason, sunrise: p("2017-10-14T19:34:24Z"), sunset: p("2017-10-15T08:37:52Z")},
			},
		},
		"manila": place{
			lat: 14.5995, lon: 120.9842,
			times: []sunData{
				{day: northenSummer, sunrise: p("2017-07-09T21:33:22Z"), sunset: p("2017-07-10T10:29:34Z")},
				{day: southernSummer, sunrise: p("2017-12-28T22:20:04Z"), sunset: p("2017-12-29T09:36:37Z")},
				{day: midSeason, sunrise: p("2017-10-14T21:47:25Z"), sunset: p("2017-10-15T09:35:49Z")},
			},
		},
		"ulanBator": place{
			lat: 47.8864, lon: 106.9057,
			times: []sunData{
				{day: northenSummer, sunrise: p("2017-07-09T21:04:35Z"), sunset: p("2017-07-10T12:50:34Z")},
				{day: southernSummer, sunrise: p("2017-12-29T00:41:36Z"), sunset: p("2017-12-29T09:07:50Z")},
				{day: midSeason, sunrise: p("2017-10-14T23:12:04Z"), sunset: p("2017-10-15T10:03:12Z")},
			},
		},
		"reykjavik": place{
			lat: 64.1265, lon: -21.8174,
			times: []sunData{
				{day: northenSummer, sunrise: p("2017-07-10T03:28:45Z"), sunset: p("2017-07-10T23:34:34Z")},
				{day: southernSummer, sunrise: p("2017-12-29T11:20:40Z"), sunset: p("2017-12-29T15:39:00Z")},
				{day: midSeason, sunrise: p("2017-10-15T08:19:47Z"), sunset: p("2017-10-15T18:04:34Z")},
			},
		},
	}
)

func p(s string) time.Time {
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		panic(err)
	}
	return t
}

func TestSunrise(t *testing.T) {
	for n, place := range places {
		for _, d := range place.times {
			name := fmt.Sprintf("%s on %v", n, d.day)
			t.Run(name, func(t *testing.T) {
				got := Sunrise(d.day, place.lat, place.lon)
				if got != d.sunrise {
					t.Errorf("got sunrise %s, want %s", got, d.sunrise)
				}
			})
		}
	}
}
func TestSunset(t *testing.T) {
	for n, place := range places {
		for _, d := range place.times {
			name := fmt.Sprintf("%s on %v", n, d.day)
			t.Run(name, func(t *testing.T) {
				got := Sunset(d.day, place.lat, place.lon)
				if got != d.sunset {
					t.Errorf("got sunrise %s, want %s", got, d.sunset)
				}
			})
		}
	}
}
