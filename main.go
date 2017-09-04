package main

import (
	"net/http"
	"net/url"
	"context"
	"fmt"
	"time"
	"strings"
	"os"
)

var globalCtx, globalCancel = context.WithCancel(context.Background())
var myFarm farm

func init() {
	myFarm = NewGopherFarm(globalCtx)
	// Add some initial gophers.
	myFarm.Add("jay")
	myFarm.Add("alex")
}

func mainHandler(r *http.Request) (string, error) {
	ctx, cancel := context.WithTimeout(r.Context(), 5 * time.Second)
	defer cancel()
	return handleMethod(ctx, r)
}

func handleMethod(ctx context.Context, r *http.Request) (string, error) {
	switch r.Method {
	case http.MethodGet:
		return handleGet(ctx, r.URL)
	}
	return "", fmt.Errorf("method not supported: %s", r.Method)
}

var getHelp = fmt.Sprintf(`Possible requests:
/%s	Returns a list of Gophers.
/%s	Returns this help.`,
	gopher, help,
)

func handleGet(ctx context.Context, request *url.URL) (string, error) {
	path := strings.Split(request.Path, string(os.PathSeparator))
	if len(path) <= 0 || path[0] == "" {
		return "", fmt.Errorf("path empty\n%s", getHelp)
	}
	switch path[0] {
		case help:
			return getHelp, nil
		case gopher:
			return myFarm.GopherList(), nil
	}
	return "", fmt.Errorf("cannot parse '%s'\n%s", request.Path, getHelp)
}
