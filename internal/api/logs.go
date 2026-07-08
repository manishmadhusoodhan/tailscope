package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"tailscope/internal/store"
)

func (s *Server) latestLogs(
	w http.ResponseWriter,
	r *http.Request,
) {

	limit := 100

	if value :=
		r.URL.Query().Get("limit"); value != "" {

		if parsed, err :=
			strconv.Atoi(value); err == nil {

			limit = parsed
		}
	}

	logs :=
		s.store.Latest(limit)

	jsonResponse(
		w,
		logs,
	)
}

func (s *Server) searchLogs(
	w http.ResponseWriter,
	r *http.Request,
) {

	query := store.Query{

		TraceID: r.URL.Query().Get("traceId"),

		CorrelationID: r.URL.Query().Get("correlationId"),

		Level: r.URL.Query().Get("level"),

		Pod: r.URL.Query().Get("pod"),

		Text: r.URL.Query().Get("text"),
	}

	logs :=
		s.store.Search(query)

	jsonResponse(
		w,
		logs,
	)
}

func jsonResponse(
	w http.ResponseWriter,
	data interface{},
) {

	w.Header().
		Set(
			"Content-Type",
			"application/json",
		)

	json.NewEncoder(w).
		Encode(data)
}
