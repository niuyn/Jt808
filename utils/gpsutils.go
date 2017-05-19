package utils

import "math"

const (
	MAX_LATITUDE = 90.0
	MIN_LATITUDE = -90.0

	MAX_LONGITUDE = 180.0
	MIN_LONGITUDE = -180.0

	pi   = math.Pi
	a    = 6378245.0
	ee   = 0.00669342162296594323
	x_pi = pi * 3000.0 / 180.0
)

func BDtoWGS(lat, lon float64) (float64, float64) {
	lat, lon = BDtoGCJ(lat, lon)
	lat, lon = GCJtoWGS(lat, lon)
	return lat, lon
}

func WGStoBD(lat, lon float64) (float64, float64) {
	lat, lon = WGStoGCJ(lat, lon)
	lat, lon = GCJtoBD(lat, lon)
	return lat, lon
}

func GCJtoBD(lat, lon float64) (bd_lat, bd_lon float64) {
	x := lon
	y := lat
	z := math.Sqrt(x*x+y*y) + 0.00002*math.Sin(y*x_pi)
	theta := math.Atan2(y, x) + 0.000003*math.Cos(x*x_pi)
	bd_lon = z*math.Cos(theta) + 0.0065
	bd_lat = z*math.Sin(theta) + 0.006
	return
}

func BDtoGCJ(lat, lon float64) (gg_lat, gg_lon float64) {
	x := lon - 0.0065
	y := lat - 0.006
	z := math.Sqrt(x*x+y*y) - 0.00002*math.Sin(y*x_pi)
	theta := math.Atan2(y, x) - 0.000003*math.Cos(x*x_pi)
	gg_lon = z * math.Cos(theta)
	gg_lat = z * math.Sin(theta)
	return
}

func WGStoGCJ(lat, lon float64) (mgLat, mgLon float64) {
	dLat := transformLat(lon-105.0, lat-35.0)
	dLon := transformLon(lon-105.0, lat-35.0)
	radLat := lat / 180.0 * pi
	magic := math.Sin(radLat)
	magic = 1 - ee*magic*magic
	sqrtMagic := math.Sqrt(magic)
	dLat = (dLat * 180.0) / ((a * (1 - ee)) / (magic * sqrtMagic) * pi)
	dLon = (dLon * 180.0) / (a / sqrtMagic * math.Cos(radLat) * pi)
	mgLat = lat + dLat
	mgLon = lon + dLon
	return
}

func GCJtoWGS(lat, lon float64) (latitude, lontitude float64) {
	_lat, _lon := transform_gps(lat, lon)
	latitude = lat*2 - _lat
	lontitude = lon*2 - _lon
	return
}

func transformLat(lat, lon float64) float64 {
	ret := -100.0 + 2.0*lat + 3.0*lon + 0.2*lon*lon + 0.1*lat*lon + 0.2*math.Sqrt(math.Abs(lat))
	ret += (20.0*math.Sin(6.0*lat*pi) + 20.0*math.Sin(2.0*lat*pi)) * 2.0 / 3.0
	ret += (20.0*math.Sin(lon*pi) + 40.0*math.Sin(lon/3.0*pi)) * 2.0 / 3.0
	ret += (160.0*math.Sin(lon/12.0*pi) + 320*math.Sin(lon*pi/30.0)) * 2.0 / 3.0
	return ret
}

func transformLon(lat, lon float64) float64 {
	ret := 300.0 + lat + 2.0*lon + 0.1*lat*lat + 0.1*lat*lon + 0.1*math.Sqrt(math.Abs(lat))
	ret += (20.0*math.Sin(6.0*lat*pi) + 20.0*math.Sin(2.0*lat*pi)) * 2.0 / 3.0
	ret += (20.0*math.Sin(lat*pi) + 40.0*math.Sin(lat/3.0*pi)) * 2.0 / 3.0
	ret += (150.0*math.Sin(lat/12.0*pi) + 300.0*math.Sin(lat/30.0*pi)) * 2.0 / 3.0
	return ret
}

func transform_gps(lat, lon float64) (mgLat, mgLon float64) {
	dLat := transformLat(lon-105.0, lat-35.0)
	dLon := transformLon(lon-105.0, lat-35.0)
	radLat := lat / 180.0 * pi
	magic := math.Sin(radLat)
	magic = 1 - ee*magic*magic
	sqrtMagic := math.Sqrt(magic)
	dLat = (dLat * 180.0) / ((a * (1 - ee)) / (magic * sqrtMagic) * pi)
	dLon = (dLon * 180.0) / (a / sqrtMagic * math.Cos(radLat) * pi)
	mgLat = lat + dLat
	mgLon = lon + dLon
	return
}
