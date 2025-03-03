package handlers

import "net/http"

func redirect(w http.ResponseWriter, r *http.Request, redirectUrl string) {
	if r.Header.Get("Hx-Request") != "" {
		w.Header().Set("HX-Redirect", redirectUrl)
		http.Redirect(w, r, redirectUrl, http.StatusOK)
	} else {
		http.Redirect(w, r, redirectUrl, http.StatusSeeOther)
	}
}
