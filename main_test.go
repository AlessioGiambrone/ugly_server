package main

import (
	"net/url"
	"strconv"
	"testing"
)

func parseValue(q url.Values, key string) float64 {
 v, _ :=strconv.ParseFloat(q[key][0], 32)
 return v
}

func TestUrlRound(t *testing.T) {
	u, _ := url.Parse("http://go.on.bike/?lat=57&lng=58.9&test3=1.05")
	var round = 1
	var max float64 = 4
	var min float64 = 3
	conf := map[string]Constraint{
		"lat":   {nil, &max, &min},
		"lng":   {&round, nil, &min},
		"test3": {&round, nil, &min},
	}
	newQueryString := applyConstraints(*u, conf)
	newQuery, _ := url.ParseQuery(newQueryString)

	if v := parseValue(newQuery, "lat"); v > max {
		t.Errorf("max value not respected %v", v)
	}

	if v := parseValue(newQuery, "lng"); v == 58.9 {
		t.Errorf("incorrect round value %v", v)
	}

	if v := parseValue(newQuery, "test3"); v == 1.1 {
		t.Errorf("incorrect round value %v", v)
	}

}

func TestUrlMin(t *testing.T) {
	u, _ := url.Parse("http://go.on.bike/?lat=-57.999")
  var round = 2
	var min float64 = -3
	conf := map[string]Constraint{
		"lat":   {&round, nil, &min},
	}
	newQuery, _ := url.ParseQuery(applyConstraints(*u, conf))

	if v := parseValue(newQuery, "lat") ; v < min {
		t.Errorf("min value not respected %v", v)
	}

}
