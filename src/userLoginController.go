package main

import (
	"net/http"
	"fmt"
	"github.com/gorilla/mux"
	"strconv"
	"database/sql"
	"encoding/json"
)

type userinfo struct {
	UserLogin user `json:"login"`
	UserDetail userDetail `json:"detail"`
}

type userInfoArray struct {
	UserLogin []user `json:"login"`
	UserDetail []userDetail `json:"detail"`
}

func (a *App) getUser (w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["uid"])
	fmt.Printf("%d", id)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	u := user{UID: id}
	fmt.Printf("%s", u)
	if err := u.getUser(a.DB); err != nil {
		switch err {
		case sql.ErrNoRows:
			respondWithError(w, http.StatusNotFound, "User Not Found")
		default:
			respondWithError(w, http.StatusInternalServerError, "Server Internal Error")
		}
		return
	}

	// GetUserDetail

	uD := userDetail{UID: id}

	if err = uD.getUserDetail(a.DB); err != nil {
		switch err {
		case sql.ErrNoRows:
			uD.Intro = ""
			uD.Groupid = 0
			uD.Avatar = 1
			uD.Nick = ""
			uD.Url = ""
		default:
			respondWithError(w, http.StatusInternalServerError, "The detail table has error, throw it!")
		}
	}

	var uinfo userinfo
	uinfo.UserLogin = u
	uinfo.UserDetail = uD

	respondWithJSON(w, http.StatusOK, uinfo)

}

func (a *App) getUsers(w http.ResponseWriter, r *http.Request) {

	count, _ := strconv.Atoi(r.FormValue("count"))
	start, _ := strconv.Atoi(r.FormValue("start"))

	if count > 20 || count < 1 {
		count = 20
	}

	if start <= 0 {
		start = 0
	}

	users, err := getUsers(a.DB, start, count)

	usersDetail, err1 := getUserDetails(a.DB, start, count)

	if err != nil || err1 != nil {
		respondWithError(w, http.StatusInternalServerError, "Server Internal Error")
		return
	}

	var userArr userInfoArray
	userArr.UserLogin = users
	userArr.UserDetail = usersDetail

	respondWithJSON(w, http.StatusOK, userArr)

}

func (a *App) createUser (w http.ResponseWriter, r *http.Request) {
	var u user
	// The following two lines is used for print out the request for debugging
	var fout string = formatRequest(r)
	fmt.Printf("%s", fout)

	// The following is used to add a Person
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&u); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid Request Payload")
		return
	}

	defer r.Body.Close()

	if err := u.createUser(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, u)

}

func (a *App) updateUser (w http.ResponseWriter, r *http.Request) {

	if r.Body == nil {
		respondWithError(w, http.StatusBadRequest, "Illegal Request")
		return
	}

	var bodyMessage map[string]map[string]interface{}
	err := json.NewDecoder(r.Body).Decode(&bodyMessage)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "The request data cannot be read")
		return
	}
	fmt.Println(bodyMessage["data"])
	dataRec := bodyMessage["data"]

	var uif userinfo
	//TODO: 验证用户输入
	type_cast(dataRec["uid"])
	type_cast(dataRec["avatar"])
	type_cast(dataRec["introduction"])
	type_cast(dataRec["nick"])
	type_cast(dataRec["url"])

	uif.UserDetail.UID = int(dataRec["uid"].(float64))
	uif.UserDetail.Avatar = int(dataRec["avatar"].(float64))
	uif.UserDetail.Intro = dataRec["introduction"].(string)
	uif.UserDetail.Nick = dataRec["nick"].(string)
	uif.UserDetail.Url = dataRec["url"].(string)
	//
	uif.UserLogin.UID = int(dataRec["uid"].(float64))
	uif.UserLogin.Mail = dataRec["mail"].(string)
	uif.UserLogin.Name = dataRec["username"].(string)

	fmt.Println(uif)

	//return
	//vars := mux.Vars(r)
	//id, err := strconv.Atoi(vars["uid"])
	//if err != nil {
	//	respondWithError(w, http.StatusBadRequest, "Invalid user ID")
	//	return
	//}
	//
	//var u user
	//decoder := json.NewDecoder(r.Body)
	//if err := decoder.Decode(&u); err != nil {
	//	respondWithError(w, http.StatusBadRequest, "Invalid Request Payload")
	//	return
	//}

	defer r.Body.Close()


	if err := uif.updateUser(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, uif)

	//u.UID = id
	//
	//if err := u.updateUser(a.DB); err != nil {
	//	respondWithError(w, http.StatusInternalServerError, err.Error())
	//	return
	//}
	//
	//respondWithJSON(w, http.StatusOK, u)
}

func (a *App) deleteUser(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["uid"])

	fmt.Printf("/n%d/n",id)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid User ID")
		return
	}

	u := user{UID: id}

	if err := u.deleteUser(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"result":"success"})

}
