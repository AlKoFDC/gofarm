package server

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"context"
)

type mockedFarm struct {
	addFn  func(context.Context) string
	listFn func() []string
	killFn func(int) error
}

func (mf mockedFarm) List() []string {
	if mf.listFn == nil {
		return []string{}
	}
	return mf.listFn()
}

func (mf mockedFarm) Kill(index int) error {
	if mf.killFn == nil {
		return nil
	}
	return mf.killFn(index)
}

func (mf mockedFarm) Add(ctx context.Context) string {
	if mf.addFn == nil {
		return ""
	}
	return mf.addFn(ctx)
}

var _ Farmer = mockedFarm{}

func TestGetHandlers(t *testing.T) {
	const (
		partialAddGopherSuccess = "Successfully spawned gopher"
		gophersTemplate         = "/gophers"
		addGopherTemplate       = "/gophers/add"
		killGopherTemplate      = "/gophers/kill"
		invalidGopherRequest    = "/gophers/haha"
		invalidGopherAddRequest = "/gophers/add/haha"

		testGopherString     = "test gopher info"
		killTestErrorMessage = "kill gopher id is not"
		existingGopherID     = 1
		nonExistingGopherID  = 2
	)

	mockListFarm := mockedFarm{
		listFn: func() []string {
			return []string{testGopherString}
		},
	}

	mockKillFarm := mockedFarm{
		killFn: func(i int) error {
			if i == existingGopherID {
				return nil
			}
			return fmt.Errorf("%s %d", killTestErrorMessage, existingGopherID)
		},
	}

	for _, test := range []struct {
		farm   mockedFarm
		url    string
		expect string
	}{
		{url: "/", expect: helpMessage},
		{url: gophersTemplate, expect: emptyGophersList},
		{url: gophersTemplate, expect: gophersHelpMessage},
		{url: addGopherTemplate, expect: partialAddGopherSuccess},
		{url: invalidGopherAddRequest, expect: addGopherInvalidParams},
		{url: invalidGopherAddRequest, expect: gophersHelpMessage},
		{url: invalidGopherRequest, expect: gophersHelpMessage},
		{url: gophersTemplate, expect: fullGophersList, farm: mockListFarm},
		{url: gophersTemplate, expect: testGopherString, farm: mockListFarm},
		{url: gophersTemplate, expect: gophersHelpMessage, farm: mockListFarm},
		{url: killGopherTemplate, expect: killGopherInvalidParams, farm: mockKillFarm},
		{url: killGopherTemplate, expect: gophersHelpMessage, farm: mockKillFarm},
		{url: fmt.Sprint(killGopherTemplate, uriSeparator, existingGopherID), expect: killGopherSuccess, farm: mockKillFarm},
		{url: fmt.Sprint(killGopherTemplate, uriSeparator, existingGopherID), expect: gophersHelpMessage, farm: mockKillFarm},
		{url: fmt.Sprint(killGopherTemplate, uriSeparator, nonExistingGopherID), expect: killGopherFailure, farm: mockKillFarm},
		{url: fmt.Sprint(killGopherTemplate, uriSeparator, nonExistingGopherID), expect: gophersHelpMessage, farm: mockKillFarm},
	} {
		t.Run(test.url, func(t *testing.T) {
			testServer := httptest.NewServer(Handler(test.farm))
			defer testServer.Close()

			res, err := http.Get(testServer.URL + test.url)
			if err != nil {
				t.Fatal(err)
			}

			testBody, err := ioutil.ReadAll(res.Body)
			res.Body.Close()
			if err != nil {
				t.Fatal(err)
			}

			if !strings.Contains(string(testBody), test.expect) {
				t.Errorf("did not get expected content '%s' in response: '%s'", test.expect, testBody)
			}
		})
	}
}

func TestOtherMethods(t *testing.T) {
	testServer := httptest.NewServer(Handler(mockedFarm{}))
	defer testServer.Close()

	var requestBody io.Reader
	t.Run("POST method", func(t *testing.T) {
		var contentType string
		res, err := http.Post(testServer.URL, contentType, requestBody)
		if err != nil {
			t.Fatal(err)
		}

		testBody, err := ioutil.ReadAll(res.Body)
		res.Body.Close()
		if err != nil {
			t.Fatal(err)
		}

		if !strings.Contains(string(testBody), methodNotAllowed) {
			t.Errorf("did not get expected content '%s' in response: '%s'", methodNotAllowed, testBody)
		}
	})

	t.Run("PUT method", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodPut, testServer.URL, requestBody)
		if err != nil {
			t.Fatal(err)
		}

		res, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatal(err)
		}

		testBody, err := ioutil.ReadAll(res.Body)
		res.Body.Close()
		if err != nil {
			t.Fatal(err)
		}

		if !strings.Contains(string(testBody), methodNotImplemented) {
			t.Errorf("did not get expected content '%s' in response: '%s'", methodNotImplemented, testBody)
		}
	})
}

func TestShiftURI(t *testing.T) {
	for _, test := range []struct {
		uri, first, rest string
	}{
		{},
		{uri: "single", first: "single"},
		{uri: "one/two", first: "one", rest: "two"},
		{uri: "one/two/three", first: "one", rest: "two/three"},
	} {
		t.Run(test.uri, func(t *testing.T) {
			if first, rest := shiftURI(test.uri); first != test.first || rest != test.rest {
				t.Errorf("for URI: %q\nexpected to get first: %q and rest: %q,\nbut instead got first: %q and rest: %q", test.uri, test.first, test.rest, first, rest)
			}
		})
	}
}
