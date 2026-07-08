package api

import (
	"encoding/json"
	"net/http"
)

type StartLogsRequest struct {
	Namespace string   `json:"namespace"`
	Pods      []string `json:"pods"`
}

func (s *Server) startLogs(
	w http.ResponseWriter,
	r *http.Request,
) {

	var req StartLogsRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if req.Namespace == "" {
		http.Error(w, "namespace required", http.StatusBadRequest)
		return
	}

	if len(req.Pods) == 0 {
		http.Error(w, "at least one pod required", http.StatusBadRequest)
		return
	}

	s.app.StartFollowing(
		req.Namespace,
		req.Pods,
	)

	w.WriteHeader(http.StatusAccepted)
}

func (s *Server) stopLogs(
	w http.ResponseWriter,
	r *http.Request,
) {

	s.app.StopFollowing()

	w.WriteHeader(
		http.StatusOK,
	)

}
