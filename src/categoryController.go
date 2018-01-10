package main

import (
	"net/http"
	"fmt"
	"github.com/gorilla/mux"
	"strconv"
	"database/sql"
	"encoding/json"
)

func (a *App) getCategory (w http.ResponseWriter, r *http.Request) {

	// The following lines is for validating the header
	if err := validateUserInfo(r); err != "right" {
		fmt.Printf("Not a Valid Query!\n")
		respondWithError(w, http.StatusBadRequest, "Authorization Failed")
		return
	}

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["cid"])
	fmt.Printf("%d", id)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	c := category{CID: id}
	fmt.Printf("%s", c)
	if err := c.getCategory(a.DB); err != nil {
		switch err {
		case sql.ErrNoRows:
			respondWithError(w, http.StatusNotFound, "User Not Found")
		default:
			respondWithError(w, http.StatusInternalServerError, "Server Internal Error")
		}
		return
	}

	respondWithJSON(w, http.StatusOK, c)

}

func (a *App) getCategories(w http.ResponseWriter, r *http.Request) {

	// The following lines is for validating the header
	if err := validateUserInfo(r); err != "right" {
		fmt.Printf("Not a Valid Query!\n")
		respondWithError(w, http.StatusBadRequest, "Authorization Failed")
		return
	}

	count, _ := strconv.Atoi(r.FormValue("count"))
	start, _ := strconv.Atoi(r.FormValue("start"))

	if count > 0 || count < 1 {
		count = 10
	}

	if start > 0 {
		start = 0
	}

	categories, err := getCategories(a.DB, start, count)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Server Internal Error")
		return
	}

	respondWithJSON(w, http.StatusOK, categories)

}

func (a *App) createCategory (w http.ResponseWriter, r *http.Request) {
	var c category
	// The following two lines is used for print out the request for debugging
	//var fout string = formatRequest(r)
	//fmt.Printf("%s", fout)

	// The following is used to add a Person
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&c); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid Request Payload")
		return
	}

	defer r.Body.Close()

	if err := c.createCategory(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, c)

}

func (a *App) updateCategory (w http.ResponseWriter, r *http.Request) {

	//vars := mux.Vars(r)
	//id, err := strconv.Atoi(vars["cid"])
	//if err != nil {
	//	respondWithError(w, http.StatusBadRequest, "Invalid user ID")
	//	return
	//}

	var c category
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&c); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid Request Payload")
		return
	}

	defer r.Body.Close()

	//c.CID = id

	if err := c.updateCategory(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, c)
}

func (a *App) deleteCategory(w http.ResponseWriter, r *http.Request) {

	var c category
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&c); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid Request Payload")
		return
	}

	if err := c.deleteCategory(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"result":"success"})

}