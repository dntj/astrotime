package astrotime_test

import (
	"fmt"
	"time"

	"github.com/dntj/astrotime"
)

func ExampleNextSunrise() {
	loc, _ := time.LoadLocation("US/Eastern")
	t := time.Date(2017, 12, 15, 10, 14, 0, 0, loc)
	sr := astrotime.NextSunrise(t, 38.8895, 77.0352)

	tzname, _ := sr.Zone()
	fmt.Printf("The next sunrise at the Washington Monument is %d:%02d %s on %d/%d/%d.\n", sr.Hour(), sr.Minute(), tzname, sr.Month(), sr.Day(), sr.Year())
	// Output: The next sunrise at the Washington Monument is 7:20 EST on 12/16/2017.
}
