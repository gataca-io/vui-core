package setup

import (
	"expvar"
	"fmt"
	"io/ioutil"
	"net/http"
)

var version []byte

func metricsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	first := true
	report := func(key string, value interface{}) {
		if !first {
			fmt.Fprintf(w, ",\n")
		}
		first = false
		if str, ok := value.(string); ok {
			fmt.Fprintf(w, "%q: %q", key, str)
		} else {
			fmt.Fprintf(w, "%q: %v", key, value)
		}
	}

	fmt.Fprintf(w, "{\n")
	expvar.Do(func(kv expvar.KeyValue) {
		report(kv.Key, kv.Value)
	})
	fmt.Fprintf(w, "\n}\n")
}

func versionHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if len(version) == 0 {
		version, _ = ioutil.ReadFile("version")
	}
	fmt.Fprintf(w, "%s", version)
	w.WriteHeader(http.StatusOK)
}

func SetupMetrics(serverAddress string) {
	mux := http.NewServeMux()
	mux.HandleFunc("/debug/vars", metricsHandler)
	mux.HandleFunc("/version", versionHandler)
	_ = http.ListenAndServe(serverAddress, mux)
}
