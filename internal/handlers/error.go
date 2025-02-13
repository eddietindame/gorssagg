package handlers

import "net/http"

func Err(w http.ResponseWriter, r *http.Request) {
	respondWithError(w, 400, "Something went wrong")
}
