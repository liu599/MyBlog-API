package main

import (
	"net/http"
	"fmt"
	"github.com/gorilla/mux"
	"strconv"
	"database/sql"
	"encoding/json"
	"html"
	"time"
	"strings"
)

type postTimeIndex struct {
	Index string `json:"ptIndex"`
	Year string `json:"year"`
	Mo string `json:"month"`
	List []int `json:"list"`
}

type postInfo struct {
	Posts []post
	Pager pager
	Indexes []postTimeIndex
}

func (a *App) getPost (w http.ResponseWriter, r *http.Request) {

	// The following lines are for validating the header
	if err := validateUserInfo(r); err != "right" {
		fmt.Printf("Not a Valid Query!\n")
		respondWithError(w, http.StatusBadRequest, "Authorization Failed")
		return
	}

	type articleList struct {
		List []int `json:"list"`
	}
	var tem articleList
	// The following is used to add a Person
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&tem); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid Request Payload")
		return
	}

	defer r.Body.Close()

	fmt.Printf("%v", tem)
	var pc []post

	for _, pin := range tem.List {
		p := post{PID: pin}
		if err := p.getPost(a.DB); err != nil {
			switch err {
			case sql.ErrNoRows:
				respondWithError(w, http.StatusNotFound, "User Not Found")
			default:
				respondWithError(w, http.StatusInternalServerError, "Server Internal Error")
			}
			return
		}
		p.Body = html.UnescapeString(p.Body)
		pc = append(pc, p)
	}

	respondWithJSON(w, http.StatusOK, pc)

}

func (a *App) getPosts(w http.ResponseWriter, r *http.Request) {
	// The following lines is for validating the header
	if err := validateUserInfo(r); err != "right" {
		fmt.Printf("Not a Valid Query!\n")
		respondWithError(w, http.StatusBadRequest, "Authorization Failed")
		return
	}
	/*
	type pager struct {
		Start int `json:"start"`
		Count int `json:"count"`
	}
	var page pager
	if r.Body == nil {
		respondWithError(w, http.StatusBadRequest, "Illegal Request")
		return
	}
	err := json.NewDecoder(r.Body).Decode(&page)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Illegal Request")
		return
	}
	fmt.Println(page.Count, page.Start)
	*/

	/*以下是以form格式请求*/
	count, _ := strconv.Atoi(r.FormValue("count"))
	start, _ := strconv.Atoi(r.FormValue("start"))
	fmt.Println(start, count)

	if count > 20 || count < 1 {
		count = 20
	}

	if start <= 0 {
		start = 0
	}

	posts, err := getPosts(a.DB, start, count)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Server Internal Error")
		return
	}

	totalNumber, err1 := getPostsCountNumber(a.DB)

	if err1 != nil {
		respondWithError(w, http.StatusInternalServerError, "Server Internal Error")
		return
	}
	/* 创建Pager用于翻页 */
	var pager pager
	/* 创建时间索引 */
	var postInfoInstance postInfo
	// var postIndexInstance postTimeIndex
	pager.TotalNumber = totalNumber
	pager.Page = (start + 1) / count + 1
	pager.Count = count
	postInfoInstance.Pager = pager
	postInfoInstance.Posts = posts

	// TODO: INTO FUNCTION
	dic := make(map[string]bool)
	var result []postTimeIndex

	allPostTimes, err3 := getPostsCreateTime(a.DB)

	if err3 != nil {
		respondWithError(w, http.StatusInternalServerError, "Server Internal Error")
		return
	}

	for _, ps := range allPostTimes {
		ts := ps.Created
		timeString := time.Time.String(time.Unix(ts, 0))
		mk := strings.Split(timeString, " ")[0]
		yyy := strings.Split(mk, "-")[0]
		mmm := strings.Split(mk, "-")[1]
		kkk := yyy + mmm
		//fmt.Println(kkk)
		//fmt.Println(dic[kkk])
		if !dic[kkk] {
			dic[kkk] = true
			var temp postTimeIndex
			temp.Index = kkk
			temp.Mo = mmm
			temp.Year = yyy
			temp.List = append(temp.List, ps.PID)
			result = append(result, temp)
		} else {
			for i, re := range result {
				if re.Index == kkk {
					//fmt.Println("12332323")
					re.List = append(re.List, ps.PID)
					result[i].List = re.List
					//fmt.Println(result)
				}
			}
		}
	}
	postInfoInstance.Indexes = result

	respondWithJSON(w, http.StatusOK, postInfoInstance)
}

func (a *App) getPostsByCategory(w http.ResponseWriter, r *http.Request) {

	// The following lines is for validating the header
	if err := validateUserInfo(r); err != "right" {
		fmt.Printf("Not a Valid Query!\n")
		respondWithError(w, http.StatusBadRequest, "Authorization Failed")
		return
	}
	/*以下是以form格式请求*/
	count, _ := strconv.Atoi(r.FormValue("count"))
	start, _ := strconv.Atoi(r.FormValue("start"))
	vars := mux.Vars(r)
	cid, err2 := strconv.Atoi(vars["cid"])
	if err2 != nil {
		// Do search the category database for cid
		var cc category
		cc.CName = vars["cid"]
		err3 := cc.getCategoryByName(a.DB)
		if err3 != nil {
			respondWithError(w, http.StatusNotFound, "Fail to get the content")
		}
		cid = cc.CID
		fmt.Println(cid)
	}

	if count > 20 || count < 1 {
		count = 20
	}

	if start <= 0 {
		start = 0
	}

	posts, err := getPostsByCategory(a.DB, cid, start, count)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Server Internal Error")
		return
	}

	totalNumber, err1 := getPostsCountNumberByCategory(a.DB, cid)

	if err1 != nil {
		respondWithError(w, http.StatusInternalServerError, "Server Internal Error")
		return
	}

	/*创建Pager用于翻页*/
	var pager pager
	var postInfoInstance postInfo
	pager.TotalNumber = totalNumber
	pager.Page = (start + 1) / count + 1
	pager.Count = count
	postInfoInstance.Pager = pager
	postInfoInstance.Posts = posts

	// TODO: INTO FUNCTION
	dic := make(map[string]bool)
	var result []postTimeIndex

	allPostTimes, err3 := getPostsCreateTime(a.DB)

	if err3 != nil {
		respondWithError(w, http.StatusInternalServerError, "Server Internal Error")
		return
	}

	for _, ps := range allPostTimes {
		ts := ps.Created
		timeString := time.Time.String(time.Unix(ts, 0))
		mk := strings.Split(timeString, " ")[0]
		yyy := strings.Split(mk, "-")[0]
		mmm := strings.Split(mk, "-")[1]
		kkk := yyy + mmm
		//fmt.Println(kkk)
		//fmt.Println(dic[kkk])
		if !dic[kkk] {
			dic[kkk] = true
			var temp postTimeIndex
			temp.Index = kkk
			temp.Mo = mmm
			temp.Year = yyy
			temp.List = append(temp.List, ps.PID)
			result = append(result, temp)
		} else {
			for i, re := range result {
				if re.Index == kkk {
					//fmt.Println("12332323")
					re.List = append(re.List, ps.PID)
					result[i].List = re.List
					//fmt.Println(result)
				}
			}
		}
	}
	postInfoInstance.Indexes = result

	respondWithJSON(w, http.StatusOK, postInfoInstance)

}

func (a *App) createPost (w http.ResponseWriter, r *http.Request) {

	var p post

	// The following is used to add a Person
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&p); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid Request Payload")
		return
	}

	defer r.Body.Close()

	if err := p.createPost(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, p)

}

func (a *App) updatePost (w http.ResponseWriter, r *http.Request) {

	var p post
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&p); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid Request Payload")
		return
	}

	defer r.Body.Close()


	if err := p.updatePost(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, p)
}

func (a *App) deletePost(w http.ResponseWriter, r *http.Request) {

	var p post
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&p); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid Request Payload")
		return
	}

	if err := p.deletePost(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"result":"success"})

}