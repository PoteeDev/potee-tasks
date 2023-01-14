package handlers

import "net/http"

type IndexData struct {
	SessionUser SessionUser
}

func (h *handler) ServeIndex(w http.ResponseWriter, r *http.Request) {
	var data IndexData
	name, role := h.GetUserName(r)
	data.SessionUser = SessionUser{name, role}
	page := "index.html"
	if h.t.Lookup(page) != nil {
		w.WriteHeader(200)
		h.t.ExecuteTemplate(w, page, &data)
		return
	}
	w.WriteHeader(404)
	w.Write([]byte("not found"))
}
