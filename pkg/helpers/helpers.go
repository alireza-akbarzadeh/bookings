package helpers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	jsonpretty "github.com/ansidev/json-pretty"
)

func PrettyDuration(d time.Duration) string {
	if d < time.Microsecond {
		return fmt.Sprintf("%dns", d.Nanoseconds())
	}
	if d < time.Millisecond {
		return fmt.Sprintf("%.2fµs", float64(d)/1000)
	}
	if d < time.Second {
		return fmt.Sprintf("%.2fms", float64(d)/1e6)
	}
	return fmt.Sprintf("%.2fs", d.Seconds())
}

// ResponseWriter writes a pretty-printed JSON response
func ResponseWriter(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	// Marshal data to JSON
	raw, err := json.Marshal(data)
	if err != nil {
		http.Error(w, `{"error":"failed to marshal json"}`, http.StatusInternalServerError)
		return
	}

	// Pretty-print JSON
	pretty := jsonpretty.Pretty(raw)

	// Write to client
	w.Write([]byte(pretty))
}
