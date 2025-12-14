package middleware

import (
	"bytes"
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"

	jsonpretty "github.com/ansidev/json-pretty"
)

// ANSI color codes
const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[34m"
	colorCyan   = "\033[36m"
	colorGray   = "\033[90m"
	colorWhite  = "\033[97m"
	colorBold   = "\033[1m"
)

// responseRecorder wraps http.ResponseWriter to capture status and body
type responseRecorder struct {
	http.ResponseWriter
	status int
	body   *bytes.Buffer
}

func (rw *responseRecorder) WriteHeader(statusCode int) {
	rw.status = statusCode
	rw.ResponseWriter.WriteHeader(statusCode)
}

func (rw *responseRecorder) Write(b []byte) (int, error) {
	rw.body.Write(b)
	return rw.ResponseWriter.Write(b)
}

// Logger returns a Chi-compatible middleware
func Logger() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			// Wrap ResponseWriter
			rec := &responseRecorder{
				ResponseWriter: w,
				status:         200,
				body:           bytes.NewBuffer(nil),
			}

			// Process request
			next.ServeHTTP(rec, r)

			// Calculate latency
			latency := time.Since(start)
			latencyStr := formatLatency(latency)

			// Status & method color
			statusColor := getStatusColor(rec.status)
			methodColor := getMethodColor(r.Method)

			// Get client IP
			clientIP := getClientIP(r)

			// Print log header
			fmt.Printf("%s %s%-7s%s %s%s%s %s%3d%s | %s | %s\n",
				colorGray+time.Now().Format("15:04:05")+colorReset,
				methodColor, r.Method, colorReset,
				colorWhite, r.RequestURI, colorReset,
				statusColor, rec.status, colorReset,
				latencyStr,
				clientIP,
			)

			// Pretty-print JSON if response body is JSON
			contentType := rec.Header().Get("Content-Type")
			if strings.Contains(contentType, "application/json") && rec.body.Len() > 0 {
				prettyJSON := jsonpretty.Pretty(rec.body.Bytes())
				if len(prettyJSON) > 0 {
					fmt.Printf("%s📦 Response:%s\n%s\n", colorCyan, colorReset, string(prettyJSON))
				}
			}
		})
	}
}

// Helper: format latency
func formatLatency(d time.Duration) string {
	switch {
	case d < time.Microsecond:
		return fmt.Sprintf("%dns", d.Nanoseconds())
	case d < time.Millisecond:
		return fmt.Sprintf("%.2fµs", float64(d.Nanoseconds())/1000)
	case d < time.Second:
		return fmt.Sprintf("%.2fms", float64(d.Nanoseconds())/1e6)
	default:
		return fmt.Sprintf("%.2fs", d.Seconds())
	}
}

// Helper: get status color
func getStatusColor(code int) string {
	switch {
	case code >= 200 && code < 300:
		return colorGreen + colorBold
	case code >= 300 && code < 400:
		return colorCyan + colorBold
	case code >= 400 && code < 500:
		return colorYellow + colorBold
	default:
		return colorRed + colorBold
	}
}

// Helper: get method color
func getMethodColor(method string) string {
	switch method {
	case "GET":
		return colorBlue + colorBold
	case "POST":
		return colorGreen + colorBold
	case "PUT":
		return colorYellow + colorBold
	case "DELETE":
		return colorRed + colorBold
	case "PATCH":
		return colorCyan + colorBold
	default:
		return colorWhite + colorBold
	}
}

// Helper: get client IP
func getClientIP(r *http.Request) string {
	// Check X-Forwarded-For first
	xff := r.Header.Get("X-Forwarded-For")
	if xff != "" {
		parts := strings.Split(xff, ",")
		return strings.TrimSpace(parts[0])
	}

	// Fallback to RemoteAddr
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return host
}
