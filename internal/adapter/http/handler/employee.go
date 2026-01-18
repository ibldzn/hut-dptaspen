package handler

import "net/http"

func (h *Handler) GetPresentEmployees(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "/Users/haytsam/Projects/hut/spinner-v2/employees.json")
}
