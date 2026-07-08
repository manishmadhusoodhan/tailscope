package api

import (
	"net/http"
)

func (s *Server) pods(
	w http.ResponseWriter,
	r *http.Request,
) {

	namespace := r.URL.Query().Get("namespace")

	if namespace == "" {
		http.Error(w, "namespace required", http.StatusBadRequest)
		return
	}

	pods, err := s.kube.PodsForNamespace(
		r.Context(),
		namespace,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jsonResponse(
		w,
		pods,
	)
}
