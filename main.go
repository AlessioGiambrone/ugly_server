package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
	"strings"
)

func handleMin(cConstraint Constraint, cName string, queryValues *url.Values) {
	if cConstraint.Min != nil {
		f, err := strconv.ParseFloat(queryValues.Get(cName), 32)
		if err == nil && f < *cConstraint.Min {
			queryValues.Set(cName, strconv.FormatFloat(*cConstraint.Min, 'G', -1, 32))
		}
	}
}

func handleMax(cConstraint Constraint, cName string, queryValues *url.Values) {
	if cConstraint.Max != nil {
		f, err := strconv.ParseFloat(queryValues.Get(cName), 32)
		if err == nil && f > *cConstraint.Max {
			queryValues.Set(cName, strconv.FormatFloat(*cConstraint.Max, 'G', -1, 32))
		}
	}
}

func handleRound(cConstraint Constraint, cName string, queryValues *url.Values) {
	if cConstraint.Round != nil {
		f, err := strconv.ParseFloat(queryValues.Get(cName), 32)
		if err == nil {
			queryValues.Set(cName, fmt.Sprintf("%."+strconv.Itoa(*cConstraint.Round)+"f", f))
		}
	}
}

func applyConstraints(u url.URL, conf map[string]Constraint) string {
	var newValues url.Values = u.Query()
	for cName, cConstraint := range conf {
		if _, ok := u.Query()[cName]; ok == true {
			handleMax(cConstraint, cName, &newValues)
			handleMin(cConstraint, cName, &newValues)
			handleRound(cConstraint, cName, &newValues)
		}
	}
	return newValues.Encode()
}

func singleJoiningSlash(a, b string) string {
	aslash := strings.HasSuffix(a, "/")
	bslash := strings.HasPrefix(b, "/")
	switch {
	case aslash && bslash:
		return a + b[1:]
	case !aslash && !bslash:
		return a + "/" + b
	}
	return a + b
}

// UglyReverseProxy is an ugly copy of httputil's NewSingleHostReverseProxy
// that adds the functionality of filtering and editing http request parameters
// by using user defined constraints.
func UglyReverseProxy(target *url.URL, constraints map[string]Constraint) *httputil.ReverseProxy {
	targetQuery := target.RawQuery
	director := func(req *http.Request) {
		req.URL.Scheme = target.Scheme
		req.URL.Host = target.Host
		req.URL.Path = singleJoiningSlash(target.Path, req.URL.Path)
		if targetQuery == "" || req.URL.RawQuery == "" {
			req.URL.RawQuery = applyConstraints(*req.URL, constraints)
		} else {
			req.URL.RawQuery = targetQuery + "&" + req.URL.RawQuery
		}
		if _, ok := req.Header["User-Agent"]; !ok {
			// explicitly disable User-Agent so it's not set to default value
			req.Header.Set("User-Agent", "")
		}
	}
	return &httputil.ReverseProxy{Director: director}
}

func main() {
	loadConfig()
	ourPort := fmt.Sprintf(":%d", Config.Port)
	rpURL, err := url.Parse(fmt.Sprintf(Config.ProxiedService))
	if err != nil {
		log.Panic("invalid proxy url")
	}
	log.Println("Serving", rpURL, "at port", ourPort)
	reverseProxy := UglyReverseProxy(rpURL, Config.Constraints)
	if err := http.ListenAndServe(ourPort, reverseProxy); err != nil {
		log.Fatal(err)
	}
}
