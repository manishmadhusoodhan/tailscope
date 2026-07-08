package api

import (
	"net/http"
)

func (s *Server) deployments(
	w http.ResponseWriter,
	r *http.Request,
) {

	namespace :=
		r.URL.Query().
			Get("namespace")

	if namespace == "" {

		namespace = "dev09"
	}

	items, err :=
		s.kube.ListDeployments(
			r.Context(),
			namespace,
		)

	if err != nil {

		http.Error(
			w,
			err.Error(),
			http.StatusInternalServerError,
		)

		return
	}

	jsonResponse(
		w,
		items,
	)
}
