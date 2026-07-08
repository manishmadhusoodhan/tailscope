package api

import (
	"encoding/json"
	"net/http"
)

func (s *Server) namespaces(
	w http.ResponseWriter,
	r *http.Request,
) {

	namespaces, err :=
		s.kube.Namespaces(
			r.Context(),
		)

	if err != nil {

		http.Error(
			w,
			err.Error(),
			http.StatusInternalServerError,
		)

		return

	}

	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	json.NewEncoder(w).Encode(
		namespaces,
	)

}
