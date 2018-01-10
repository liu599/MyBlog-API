package main

import (
	"net/http"
	"html/template"
	// "os"
	"path/filepath"
	"fmt"
	// "bytes"
	"github.com/gorilla/mux"
	"strconv"
	"database/sql"
	"gopkg.in/russross/blackfriday.v2"
)

type Purchase struct {
	Amount int
}



func (a *App) postWriter (w http.ResponseWriter, r *http.Request) {

	// https://studygolang.com/articles/3174

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["pid"])
	fmt.Printf("%d\n", id)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	p := post{PID: id}

	if err := p.getPost(a.DB); err != nil {
		switch err {
		case sql.ErrNoRows:
			respondWithError(w, http.StatusNotFound, "Post Not Found")
		default:
			respondWithError(w, http.StatusInternalServerError, "Server Internal Error")
		}
		return
	}


	output := blackfriday.Run([]byte(p.Body))

	//cwd, _ := os.Getwd()
	// https://stackoverflow.com/questions/18436543/golang-new-template-not-working
	t := template.New("post.gohtml") //创建一个模板
	abs, _ := filepath.Abs("./tmpl/post.gohtml")
	tmpp, _ := t.ParseFiles(abs)  //解析模板文件
	tmpp.Execute(w, " ")
	w.Write(output)

	p.Body = string(output)
	p.updatePost(a.DB)

}
