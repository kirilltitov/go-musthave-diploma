package utils

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kinbiko/jsonassert"
)

func testHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusBadRequest)
	io.WriteString(w, "hello world")
}

func TestWithLogging(t *testing.T) {
	ja := jsonassert.New(t)

	buf := bytes.Buffer{}
	Log.SetOutput(&buf)

	r := httptest.NewRequest(http.MethodGet, "/abc", nil)
	w := httptest.NewRecorder()

	WithLogging(testHandler)(w, r)

	ja.Assertf(
		buf.String(),
		`{
			"duration_μs": "<<PRESENCE>>",
			"level": "info",
			"method": "GET",
			"msg": "Served HTTP request",
			"size": 11,
			"status": 400,
			"time": "<<PRESENCE>>",
			"uri": "/abc"
		}`,
	)
}
