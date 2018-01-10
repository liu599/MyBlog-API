package main

import (
	"net/http"
	"fmt"
	"github.com/gorilla/mux"
	"strconv"
	"database/sql"
)

func (a *App) getUserDetail (w http.ResponseWriter, r *http.Request) {
	// The following lines is for validating the header
	if err := validateUserInfo(r); err != "right" {
		fmt.Printf("Not a Valid Query!\n")
		respondWithError(w, http.StatusBadRequest, "Authorization Failed")
		return
	}

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["uid"])
	fmt.Printf("%d", id)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	ud := userDetail{UID: id}
	fmt.Printf("%s", ud)
	if err := ud.getUserDetail(a.DB); err != nil {
		switch err {
		case sql.ErrNoRows:
			respondWithError(w, http.StatusNotFound, "User Not Found")
		default:
			respondWithError(w, http.StatusInternalServerError, "Server Internal Error")
		}
		return
	}

	respondWithJSON(w, http.StatusOK, ud)

}
