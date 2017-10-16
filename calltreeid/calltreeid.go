package calltreeid

import (
	"context"
	"net/http"

	"github.com/AlKoFDC/gofarm/rand"
)

type contextKey string

const requestKey = "CallTreeID"
const callTreeID contextKey = requestKey

func CreateRequestContext(r *http.Request) context.Context {
	cti := r.Header.Get(requestKey)
	if cti == "" {
		cti = rand.String(30)
	}
	return InNewContext(r.Context(),  cti)
}

func InNewContext(ctx context.Context, cti string) context.Context {
	return context.WithValue(ctx, callTreeID, cti)
}

func FromContext(ctx context.Context) string {
	cti, _ := ctx.Value(callTreeID).(string)
	return cti
}
