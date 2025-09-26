package httpx

import (
	"net/http"

	"github.com/gorilla/schema"
)

func DecodeURLQuery(r *http.Request, v any) error {
	decoder := schema.NewDecoder()
	decoder.IgnoreUnknownKeys(true)

	return decoder.Decode(v, r.URL.Query())
}
