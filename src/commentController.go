package main

import (
	"fmt"
	"net/http"
	"github.com/gorilla/mux"
	"strconv"
	"database/sql"
	"encoding/json"
)

func (a *App) getComment (w http.ResponseWriter, r *http.Request) {

	// The following lines is for validating the header
	if err := validateUserInfo(r); err != "right" {
		fmt.Printf("Not a Valid Query!\n")
		respondWithError(w, http.StatusBadRequest, "Authorization Failed")
		return
	}

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["coid"])
	fmt.Printf("%d", id)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	co := comment{COID: id}
	fmt.Printf("%s", co)
	if err := co.getComment(a.DB); err != nil {
		switch err {
		case sql.ErrNoRows:
			respondWithError(w, http.StatusNotFound, "User Not Found")
		default:
			respondWithError(w, http.StatusInternalServerError, "Server Internal Error")
		}
		return
	}

	respondWithJSON(w, http.StatusOK, co)

}


func (a *App) getComments (w http.ResponseWriter, r *http.Request) {

	// The following lines is for validating the header
	if err := validateUserInfo(r); err != "right" {
		fmt.Printf("Not a Valid Query!\n")
		respondWithError(w, http.StatusBadRequest, "Authorization Failed")
		return
	}

	count, _ := strconv.Atoi(r.FormValue("count"))
	start, _ := strconv.Atoi(r.FormValue("start"))
	vars := mux.Vars(r)
	pid, err := strconv.Atoi(vars["pid"])

	if count > 0 || count < 1 {
		count = 10
	}

	if start > 0 {
		start = 0
	}

	comments, err := getComments(a.DB, pid, start, count)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Server Internal Error")
		return
	}

	respondWithJSON(w, http.StatusOK, comments)

}

func (a *App) createComment (w http.ResponseWriter, r *http.Request) {
	var co comment
	// The following two lines is used for print out the request for debugging
	var fout string = formatRequest(r)
	fmt.Printf("%s", fout)
	// The following lines is for validating the header
	if err := validateUserInfo(r); err != "right" {
		fmt.Printf("Not a Valid Query!\n")
		respondWithError(w, http.StatusBadRequest, "Authorization Failed")
		return
	}
	// The following is used to add a Person
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&co); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid Request Payload")
		return
	}

	defer r.Body.Close()

	if err := co.createComment(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, co)

}


func (a *App) updateComment (w http.ResponseWriter, r *http.Request) {

	// The following lines is for validating the header
	if err := validateUserInfo(r); err != "right" {
		fmt.Printf("Not a Valid Query!\n")
		respondWithError(w, http.StatusBadRequest, "Authorization Failed")
		return
	}

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["coid"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	var co comment
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&co); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid Request Payload")
		return
	}

	defer r.Body.Close()

	co.COID = id

	if err := co.updateComment(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, co)
}


func (a *App) deleteComment(w http.ResponseWriter, r *http.Request) {

	// The following lines is for validating the header
	if err := validateUserInfo(r); err != "right" {
		fmt.Printf("Not a Valid Query!\n")
		respondWithError(w, http.StatusBadRequest, "Authorization Failed")
		return
	}


	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["coid"])

	fmt.Printf("/n%d/n",id)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid User ID")
		return
	}

	co := comment{COID: id}

	if err := co.deleteComment(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"result":"success"})

}