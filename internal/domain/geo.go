package domain

import "math"

type Coordinate struct{ Latitude, Longitude float64 }

func (c Coordinate) Valid() bool {
	return c.Latitude >= -90 && c.Latitude <= 90 && c.Longitude >= -180 && c.Longitude <= 180 && !math.IsNaN(c.Latitude) && !math.IsNaN(c.Longitude)
}
func (c Coordinate) Distance(other Coordinate) float64 {
	lat := c.Latitude - other.Latitude
	lon := c.Longitude - other.Longitude
	return math.Sqrt(lat*lat + lon*lon)
}
func (c Coordinate) Within(other Coordinate, radius float64) bool { return c.Distance(other) <= radius }
func AverageCoordinate(points []Coordinate) Coordinate {
	if len(points) == 0 {
		return Coordinate{}
	}
	lat, lon := 0.0, 0.0
	for _, point := range points {
		lat += point.Latitude
		lon += point.Longitude
	}
	return Coordinate{Latitude: lat / float64(len(points)), Longitude: lon / float64(len(points))}
}
