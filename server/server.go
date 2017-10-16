package server

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/AlKoFDC/gofarm/calltreeid"
)

type Farmer interface {
	List() []string
	Add(context.Context) string
	Kill(int) error
}

func Handler(f Farmer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		var response string
		var status int
		switch r.Method {
		case http.MethodGet:
			status, response = handleGet(ctx, f, r.RequestURI)
		case http.MethodPost:
			status = http.StatusMethodNotAllowed
			response = fmt.Sprintf("%s: %s", methodNotAllowed, r.Method)
		default:
			status = http.StatusNotImplemented
			response = fmt.Sprintf("%s: %s", methodNotImplemented, r.Method)
		}
		w.WriteHeader(status)
		w.Write([]byte(response))
	}
}

func handleGet(ctx context.Context, f Farmer, uri string) (int, string) {
	uri = strings.Trim(uri, uriSeparator)
	first, rest := shiftURI(uri)
	switch first {
	case "gophers":
		return handleGophers(ctx, f, rest)
	default:
		return http.StatusBadRequest, helpMessage
	}
}

func handleGophers(ctx context.Context, f Farmer, uri string) (int, string) {
	var (
		status   int
		response string
	)
	first, rest := shiftURI(uri)
	switch first {
	case "":
		list := f.List()
		status = http.StatusOK
		response = fmt.Sprintf("%s\n\n%s", emptyGophersList, gophersHelpMessage)
		if len(list) > 0 {
			response = fmt.Sprintf("%s\n\n  %s\n\n%s", fullGophersList, strings.Join(list, "\n  "), gophersHelpMessage)
		}
	case "add":
		status, response = handleGophersAdd(ctx, f, rest)
	case "kill":
		status, response = handleGophersKill(f, rest)
	default:
		status, response = http.StatusBadRequest, gophersHelpMessage
	}
	return status, response
}

func handleGophersAdd(ctx context.Context, f Farmer, uri string) (int, string) {
	if uri != "" {
		return http.StatusBadRequest, fmt.Sprintf("%s: %s\n\n%s", addGopherInvalidParams, uri, gophersHelpMessage)
	}
	return http.StatusOK, fmt.Sprintf("%s: %s!\n\n%s", addGopherSuccess, f.Add(ctx), gophersHelpMessage)
}

func handleGophersKill(f Farmer, uri string) (int, string) {
	idString, rest := shiftURI(uri)
	id, err := strconv.Atoi(idString)
	if err != nil || rest != "" {
		return http.StatusBadRequest, fmt.Sprintf("%s: %s\n\n%s", killGopherInvalidParams, uri, gophersHelpMessage)
	}
	killErr := f.Kill(id)
	if killErr != nil {
		return http.StatusInternalServerError, fmt.Sprintf("%s %d, error: %s.\n\n%s", killGopherFailure, id, killErr, gophersHelpMessage)
	}
	return http.StatusOK, fmt.Sprintf("%s %d.\n\n%s", killGopherSuccess, id, gophersHelpMessage)
}

func shiftURI(uri string) (string, string) {
	splitIndex := strings.Index(uri, uriSeparator)
	if splitIndex == -1 {
		return uri, ""
	}
	return uri[:splitIndex], uri[splitIndex+1:]
}

const (
	uriSeparator = "/"
	helpMessage  = `Available commands:
 - /gophers - returns list of current gophers`
	emptyGophersList   = `No gophers found in the farm.`
	fullGophersList    = `List of gophers found in the farm:`
	gophersHelpMessage = `Available commands:
 - /gophers/add  - spawns a new gopher
 - /gophers/kill/<gopher_id> - slaughters the gopher with <gopher_id>`
	addGopherSuccess        = `Successfully spawned gopher`
	killGopherSuccess       = `Successfully slaughtered gopher`
	killGopherFailure       = `Could not slaughter gopher`
	killGopherInvalidParams = "Parameters are not set properly for request"
	addGopherInvalidParams  = "Did not expect parameters for request"
)

const (
	methodNotAllowed     = "method not allowed"
	methodNotImplemented = "method not implemented"
)
