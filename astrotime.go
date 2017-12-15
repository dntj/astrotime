// Package astrotime implements NAAA.
// NAA - NOAA's Astronomical Algorithms
// (JavaScript web page
//  http://www.srrb.noaa.gov/highlights/sunrise/sunrise.html by
//  Chris Cornwall, Aaron Horiuchi and Chris Lehman)
// Ported to C++ by Pete Gray (petegray@ieee.org), July 2006
// Released as Open Source and can be used in any way, as long as the
// above description remains in place.
package astrotime

import (
	"math"
	"time"
)

const (
	radToDeg  = 180 / math.Pi
	degToRad  = math.Pi / 180
	radToGrad = 200 / math.Pi
	gradToDeg = math.Pi / 200

	oneDay = time.Hour * 24
)

// julianDate converts a Time to a Julian date.
func julianDate(t time.Time) float64 {
	y := t.Year()
	m := int(t.Month())
	d := t.Day()
	hh := t.Hour()
	mm := t.Minute()
	ss := t.Second()
	ms := t.Nanosecond() / 1e6

	// Calc integer part (days)
	jday := (1461*(y+4800+(m-14)/12))/4 + (367*(m-2-12*((m-14)/12)))/12 - (3*((y+4900+(m-14)/12)/100))/4 + d - 32075

	// Calc floating point part (fraction of a day)
	jdatetime := float64(jday) + (float64(hh)-12.0)/24.0 + (float64(mm) / 1440.0) + (float64(ss) / 86400.0) + (float64(ms) / 86400000.0)

	// Adjust to UT
	_, zoneOffset := t.Zone()

	return jdatetime + float64(zoneOffset)/86400
}

// julianCentury converts a Julian Day to centuries since J2000.0.
func julianCentury(t float64) float64 {
	return (t - 2451545) / 36525
}

// julianDateFromJulianCentury converts centuries since J2000.0 to Julian Day.
func julianDateFromJulianCentury(t float64) float64 {
	return t*36525.0 + 2451545.0
}

// solarGeoMeanLon calculates the Geometric Mean Longitude of the Sun.
func solarGeoMeanLon(t float64) float64 {
	lon := math.Mod(280.46646+t*(36000.76983+0.0003032*t), 360)
	if lon > 0.0 {
		return lon
	}

	return lon + 360
}

// eclipticMeanObliquity calculates the mean obliquity of the ecliptic.
func eclipticMeanObliquity(t float64) float64 {
	seconds := 21.448 - t*(46.8150+t*(0.00059-t*(0.001813)))
	return 23.0 + (26.0+(seconds/60.0))/60.0
}

// obliquityCorrection calculates the corrected obliquity of the ecliptic.
func obliquityCorrection(t float64) float64 {
	e0 := eclipticMeanObliquity(t)
	omega := 125.04 - 1934.136*t
	return e0 + 0.00256*math.Cos(omega*degToRad)
}

// earthOrbitEccentricity calculates the eccentricity of earth's orbit.
func earthOrbitEccentricity(t float64) float64 {
	return 0.016708634 - t*(0.000042037+0.0000001267*t)
}

// meanSolarAnomaly calculates the Geometric Mean Anomaly of the Sun.
func meanSolarAnomaly(t float64) float64 {
	return 357.52911 + t*(35999.05029-0.0001537*t)
}

// equationOfTime calculates the difference between true solar time and mean solar time.
func equationOfTime(t float64) float64 {
	epsilon := obliquityCorrection(t)
	l0 := solarGeoMeanLon(t)
	e := earthOrbitEccentricity(t)
	m := meanSolarAnomaly(t)

	y := math.Tan(degToRad * epsilon / 2.0)
	y *= y

	sin2l0 := math.Sin(2.0 * degToRad * l0)
	sinm := math.Sin(degToRad * m)
	cos2l0 := math.Cos(2.0 * degToRad * l0)
	sin4l0 := math.Sin(4.0 * degToRad * l0)
	sin2m := math.Sin(2.0 * degToRad * m)

	Etime := y*sin2l0 - 2.0*e*sinm + 4.0*e*y*sinm*cos2l0 - 0.5*y*y*sin4l0 - 1.25*e*e*sin2m

	return radToDeg * Etime * 4.0
}

// solarEqOfCenter calculates the equation of center for the sun.
func solarEqOfCenter(t float64) float64 {
	m := meanSolarAnomaly(t)
	mrad := degToRad * m
	sinm := math.Sin(mrad)
	sin2m := math.Sin(mrad + mrad)
	sin3m := math.Sin(mrad + mrad + mrad)
	return sinm*(1.914602-t*(0.004817+0.000014*t)) + sin2m*(0.019993-0.000101*t) + sin3m*0.000289
}

// solarTrueLon calculates the true longitude of the sun.
func solarTrueLon(t float64) float64 {
	l0 := solarGeoMeanLon(t)
	c := solarEqOfCenter(t)
	return l0 + c
}

// solarApparentLon calculates the apparent longitude of the sun.
func solarApparentLon(t float64) float64 {
	o := solarTrueLon(t)
	omega := 125.04 - 1934.136*t
	return o - 0.00569 - 0.00478*math.Sin(degToRad*omega)
}

// solarDeclination calculates the declination of the sun.
func solarDeclination(t float64) float64 {
	e := obliquityCorrection(t)
	lambda := solarApparentLon(t)
	sint := math.Sin(degToRad*e) * math.Sin(degToRad*lambda)
	return radToDeg * math.Asin(sint)
}

// hourAngleSunrise calculates the hour angle of the sun at sunrise for the latitude.
func hourAngleSunrise(lat, solarDec float64) float64 {
	latRad := degToRad * lat
	sdRad := degToRad * solarDec
	return -math.Acos(math.Cos(degToRad*90.833)/(math.Cos(latRad)*math.Cos(sdRad)) - math.Tan(latRad)*math.Tan(sdRad))
}

// solNoonUTC calculates the Universal Coordinated Time (UTC) of solar noon for the
// given day at the given location on earth.
func solNoonUTC(t, longitude float64) float64 {
	// First pass uses approximate solar noon to calculate eqtime
	tnoon := julianCentury(julianDateFromJulianCentury(t) - longitude/360.0)
	eqTime := equationOfTime(tnoon)
	solNoonUTC := 720 - (longitude * 4) - eqTime
	newt := julianCentury(julianDateFromJulianCentury(t) - 0.5 + solNoonUTC/1440.0)
	eqTime = equationOfTime(newt)
	return 720 - (longitude * 4) - eqTime
}

// sunriseUTC calculates the UTC sunrise for the given day at the given location.
func sunriseUTC(jd, latitude, longitude float64) float64 {
	t := julianCentury(jd)

	// *** Find the time of solar noon at the location, and use
	//     that declination. This is better than start of the
	//     Julian day

	noonmin := solNoonUTC(t, longitude)
	tnoon := julianCentury(jd + noonmin/1440.0)

	// *** First pass to approximate sunrise (using solar noon)

	eqTime := equationOfTime(tnoon)
	solarDec := solarDeclination(tnoon)
	hourAngle := hourAngleSunrise(latitude, solarDec)

	delta := radToDeg*hourAngle - longitude
	timeDiff := 4 * delta
	timeUTC := 720 + timeDiff - eqTime

	// *** Second pass includes fractional jday in gamma calc

	newt := julianCentury(julianDateFromJulianCentury(t) + timeUTC/1440.0)
	eqTime = equationOfTime(newt)
	solarDec = solarDeclination(newt)
	hourAngle = hourAngleSunrise(latitude, solarDec)
	delta = radToDeg*hourAngle - longitude
	timeDiff = 4 * delta
	timeUTC = 720 + timeDiff - eqTime
	return timeUTC
}

// Sunrise calculates the sunrise, in local time, on the day t at the
// location specified in longitude and latitude.
func Sunrise(t time.Time, latitude, longitude float64) time.Time {
	jd := julianDate(t)
	sr := time.Duration(math.Floor(sunriseUTC(jd, latitude, longitude)*60) * 1e9)
	loc, _ := time.LoadLocation("UTC")
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, loc).Add(sr).In(t.Location())
}

// hourAngleSunset calculates the hour angle of the sun at sunset for the latitude.
func hourAngleSunset(lat, solarDec float64) float64 {
	latRad := degToRad * lat
	sdRad := degToRad * solarDec

	HA := (math.Acos(math.Cos(degToRad*90.833)/(math.Cos(latRad)*math.Cos(sdRad)) - math.Tan(latRad)*math.Tan(sdRad)))

	return -HA // in radians
}

// sunsetUTC calculates the Universal Coordinated Time (UTC) of sunset
// for the given day at the given location on earth.
func sunsetUTC(jd, latitude, longitude float64) float64 {
	t := julianCentury(jd)

	// *** Find the time of solar noon at the location, and use
	//     that declination. This is better than start of the
	//     Julian day

	noonmin := solNoonUTC(t, longitude)
	tnoon := julianCentury(jd + noonmin/1440.0)

	// First calculates sunrise and approx length of day

	eqTime := equationOfTime(tnoon)
	solarDec := solarDeclination(tnoon)
	hourAngle := hourAngleSunset(latitude, solarDec)

	delta := -longitude - radToDeg*hourAngle
	timeDiff := 4 * delta
	timeUTC := 720 + timeDiff - eqTime

	// first pass used to include fractional day in gamma calc

	newt := julianCentury(julianDateFromJulianCentury(t) + timeUTC/1440.0)
	eqTime = equationOfTime(newt)
	solarDec = solarDeclination(newt)
	hourAngle = hourAngleSunset(latitude, solarDec)

	delta = -longitude - radToDeg*hourAngle
	timeDiff = 4 * delta
	return 720 + timeDiff - eqTime
}

// Sunset calculates the sunset, in local time, on the day t at the
// location specified in longitude and latitude.
func Sunset(t time.Time, latitude, longitude float64) time.Time {
	jd := julianDate(t)
	ss := time.Duration(math.Floor(sunsetUTC(jd, latitude, longitude)*60) * 1e9)
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC).Add(ss).In(t.Location())
}

// NextSunrise returns date/time of the next sunrise after after
func NextSunrise(after time.Time, latitude, longitude float64) time.Time {
	s := Sunrise(after, latitude, longitude)
	if after.Before(s) {
		return s
	}

	return Sunrise(after.Add(oneDay), latitude, longitude)
}

// NextSunset returns date/time of the next sunset after after
func NextSunset(after time.Time, latitude, longitude float64) time.Time {
	s := Sunset(after, latitude, longitude)
	if after.Before(s) {
		return s
	}

	return Sunset(after.Add(oneDay), latitude, longitude)
}
