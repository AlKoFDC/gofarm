package main

import (
	"testing"
	"net/http"
	"net/url"
)

func TestShouldParseURL(t *testing.T) {
	for desc, testIO := range map[string]struct{
		rawURL string
	}{
		//"empty": {rawURL: ""},
		//"empty path": {rawURL: "http://tamagophers:1800/"},
		"gopher": {rawURL: "gopher"},
		"help": {rawURL: "help"},
	}{
		t.Run(desc, func(t *testing.T) {
			gophurl, err := url.Parse(testIO.rawURL)
			if err != nil {
				t.Fatalf("cannot parse URL '%s': %s", testIO.rawURL, err)
			}
			r := http.Request{
				Method: http.MethodGet,
				URL: gophurl, 
			}
			if response, err := mainHandler(&r); err != nil {
				t.Error(err)
			} else {
				t.Log(response)
			}
		})
	}
}
