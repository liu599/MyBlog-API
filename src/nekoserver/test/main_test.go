package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"
	"time"

	"golang.org/x/crypto/scrypt"
	"nekoserver/middleware/data"
	"nekoserver/middleware/func"
	"nekoserver/router"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"gopkg.in/mgo.v2/bson"
)

var db *sqlx.DB

func ensureTableExists(db *sqlx.DB) {
	if _, err := db.Exec(userTableCreationQuery); err != nil {
		log.Fatal(err)
	}
	if _, err := db.Exec(postTableCreationQuery); err != nil {
		log.Fatal(err)
	}
	if _, err := db.Exec(categoryTableCreationQuery); err != nil {
		log.Fatal(err)
	}
	if _, err := db.Exec(commentTableCreationQuery); err != nil {
		log.Fatal(err)
	}
}

func TestMain(m *testing.M) {

	os.Setenv("PASS_GEN", "asdfasd")

	os.Setenv("NEKO_TOKEN", "c131c35a24d")

	database := data.Database{
		Driver: "mysql",
		MaxIdle: 2,
		MaxOpen: 15,
		Name: "nekohand",
		Source: "root:86275198@tcp(127.0.0.1:3306)/nekohand?charset=utf8",
	}

	var Apps = make(map[string]data.Database)

	Apps["nekohand"] = database

	_func.AssignAppDataBaseList(Apps)

	_func.AssignDatabaseFromList([]string{"nekohand"})

	db, err := _func.MySqlGetDB("nekohand")
	if err != nil {
		fmt.Println("Error Database Connection")
		return
	}
	ensureTableExists(db)
	//ensureCategoryTableExists(db)
	//ensureRelationshipTableExists(db)
	code := m.Run()
	//clearTable(db)
	//clearCategoryTable(db)
	//clearRelationshipTable(db)
	os.Exit(code)
}

func clearTable(db *sqlx.DB) {
	db.Exec("DELETE FROM user")
	db.Exec("ALTER TABLE user AUTO_INCREMENT = 1")
}
const userTableCreationQuery = `
CREATE TABLE IF NOT EXISTS user
(
    uid        INT(11) NOT NULL AUTO_INCREMENT PRIMARY KEY,
    userid     VARCHAR(50) UNIQUE NOT NULL,
	name       VARCHAR(32)  NOT NULL,
	password   VARCHAR(100)  NOT NULL,
	mail       VARCHAR(200)  NOT NULL,
	createdAt  INT(64)  NOT NULL,
	loggedAt   INT(64) NOT NULL
) character set = utf8`

const postTableCreationQuery = `
CREATE TABLE IF NOT EXISTS post
(
    pid        INT(11) NOT NULL AUTO_INCREMENT PRIMARY KEY,
    id         VARCHAR(50) UNIQUE NOT NULL,
	author	   VARCHAR(50) NOT NULL,
	category   VARCHAR(50) NOT NULL,
	body	   TEXT(1000) NOT NULL,
	ptitle     VARCHAR(32)  NOT NULL,
	slug       VARCHAR(32)  NOT NULL,
	password   VARCHAR(32)  NOT NULL,
	createdAt  INT(64)  NOT NULL,
	modifiedAt   INT(64) NOT NULL
) character set = utf8`

const categoryTableCreationQuery = `
CREATE TABLE IF NOT EXISTS category
(
    cid        INT(11) NOT NULL AUTO_INCREMENT PRIMARY KEY,
    id         VARCHAR(50) UNIQUE NOT NULL,
	cname	   VARCHAR(50) UNIQUE NOT NULL,
	cinfo     VARCHAR(32) NULL
) character set = utf8`

const commentTableCreationQuery = `
CREATE TABLE IF NOT EXISTS comment
(
    comid      INT(11) NOT NULL AUTO_INCREMENT PRIMARY KEY,
    commentid  VARCHAR(50) UNIQUE NOT NULL,
    pid        VARCHAR(50) NOT NULL,
	author	   VARCHAR(50) NOT NULL,
	mail	   VARCHAR(50) NOT NULL,
	url	       VARCHAR(200) NOT NULL,
	ip         VARCHAR(80) NOT NULL,
	prid       VARCHAR(50) NOT NULL,
	body	   TEXT(1000) NOT NULL,
	createdAt  INT(64)  NOT NULL,
	modifiedAt INT(64) NOT NULL
) character set = utf8`


// 发送请求
func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	gin.SetMode(gin.ReleaseMode)
	engine := gin.New()
	//engine.Use(auth.TokenAuthMiddleware())

	// Router
	router.AssignBackendRouter(engine)
	engine.ServeHTTP(rr, req)

	return rr
}
// 检查Response
func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
}

func TestEmptyTable(t *testing.T) {
	// clearTable(db)
	form := url.Values{}
	req, _ := http.NewRequest("GET", "/v2/backend/status", strings.NewReader(form.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	response := executeRequest(req)
	fmt.Println(response.Body)
	//var r common.ResponseBody
	//decoder := json.NewDecoder(response.Body)
	//if err := decoder.Decode(&r); err != nil {
	//	fmt.Println("JSON Illegal")
	//	return
	//}
	//checkResponseCode(t, http.StatusOK, http.StatusOK)
	//if body := fmt.Sprintf("%v", r.Data); body != "[]" {
	//	t.Errorf("Expected an empty array. Got %s", body)
	//}
}

var NewName = "Test"
var NewDate = time.Now().Unix()

func insertOneData(id string, db *sqlx.DB) {
	//fmt.Println(id, "Inserted")
	statement := fmt.Sprintf("INSERT INTO user (id, name, password, mail, createAt, loggedAt) VALUES('%s', '%s', '%s', '%s', '%d', '%d')",
		id, NewName, bson.NewObjectId().Hex(), "1234567890@qq.com", NewDate, NewDate)
	_, err := db.Exec(statement)

	if err != nil {
		fmt.Println("Database error")
	}
}

func insertCategories(db *sqlx.DB) {
	statement := fmt.Sprintf("INSERT INTO category (id, cname, cinfo) VALUES('%s', '%s', '%s')", bson.NewObjectId().Hex(), "plsss", "1231232")
	_, err := db.Exec(statement)

	if err != nil {
		fmt.Println("Database error")
	}
}

func insertPost(db *sqlx.DB) {
	statement := fmt.Sprintf("INSERT INTO post (id, author, category, body, ptitle, slug, password, createdAt, modifiedAt) VALUES('%s', '%s', '%s','%s', '%s', '%s','%s', '%d', '%d')", bson.NewObjectId().Hex(), "eddie32", "5b6c42b25c964c10a4c68d1a", "This is post test", "Osafdafn", "Sec", "ADFSdfasdafdsDF", time.Now().Unix(), time.Now().Unix())
	_, err := db.Exec(statement)

	if err != nil {
		fmt.Println("Database error")
	}
}



func TestFetchCategories(t *testing.T) {
	//db, _ := _func.MySqlGetDB("nekohand")
	//insertCategories(db)
	fmt.Println("123456")
	req, _ := http.NewRequest("GET", "/v2/backend/categories", nil)
	response := executeRequest(req)
	fmt.Println(response.Body)
}

func TestFetchPosts(t *testing.T) {
	//db, _ := _func.MySqlGetDB("nekohand")
	//insertPost(db)
	form := url.Values{}
	fmt.Println(bson.NewObjectId().Hex(), "ID")
	form.Add("token", "0003020")
	form.Add("pageNumber",  "1")
	form.Add("pageSize", "10")
	req, _ := http.NewRequest("POST", "/v2/backend/posts", strings.NewReader(form.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	response := executeRequest(req)
	fmt.Println(response.Body)
}

func TestFetchOnePost(t *testing.T) {
	form := url.Values{}
	req, _ := http.NewRequest("POST", "/v2/backend/post/5b6c42b25c964c10a4c68d19", strings.NewReader(form.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	response := executeRequest(req)
	fmt.Println(response.Body)
}

func TestFetchPostsByCategory(t *testing.T) {
	form := url.Values{}
	form.Add("token", "0003020")
	form.Add("pageNumber",  "1")
	form.Add("pageSize", "10")
	form.Add("category", "5b6c42b25c964c10a4c68d1a")
	req, _ := http.NewRequest("POST", "/v2/backend/posts/5b6c42b25c964c10a4c68d1a", strings.NewReader(form.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	response := executeRequest(req)
	fmt.Println(response.Body)
}

func TestFetchPostsChronology(t *testing.T) {
	form := url.Values{}
	req, _ := http.NewRequest("GET", "/v2/backend/posts-chronology", strings.NewReader(form.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	response := executeRequest(req)
	fmt.Println(response.Body)
}

func insertComment(db *sqlx.DB) {
	statement := fmt.Sprintf("INSERT INTO comment (commentid, pid, author, mail, url, ip, prid, body, createdAt, modifiedAt) VALUES('%s', '%s', '%s','%s', '%s', '%s', '%s', '%s', '%d', '%d')", bson.NewObjectId().Hex(), "5b6c5ced5c964c2770f66e8a", "eddie32", "1185414132@qq.com", "https://bbs.nekohand.moe", "102.22.22.22", "5bb8436c5c964c04e0367836", "This is a reply 2", time.Now().Unix(), time.Now().Unix())
	_, err := db.Exec(statement)

	if err != nil {
		fmt.Println("Database error")
	}
}

func TestFetchComments(t *testing.T) {
	//db, _ := _func.MySqlGetDB("nekohand")
	//insertComment(db)
	form := url.Values{}
	form.Add("token", "0003020")
	req, _ := http.NewRequest("POST", "/v2/backend/comments/5b6c5ced5c964c2770f66e8a", strings.NewReader(form.Encode()))
	response := executeRequest(req)
	fmt.Println(response.Body)
}

func TestCreateComment(t *testing.T) {
	comet := &data.Comment{}
	comet.COMID = bson.NewObjectId().Hex()
	comet.PID = "5b72f09a5c964c32f078402c"
	comet.Author = "eddjkladfja"
	comet.Mail = "460512944@qq.com"
	comet.Url = "https://bbs.sss.com"
	comet.Ip = "102.22.8.225"
	comet.Prid = "0"
	comet.Body = "nice blog"
	comet.CreatedAt = time.Now().Unix()
	comet.ModifiedAt = time.Now().Unix()
	jsonStr, err := json.Marshal(comet)
	if err != nil {
		panic(err)
	}
	req, _ := http.NewRequest("POST", "/v2/backend/c2a5cc3b070", bytes.NewBuffer(jsonStr))
	req.Header.Add("Content-Type", "application/json")
	response := executeRequest(req)
	fmt.Println(response.Body)
}


func insertUser(db *sqlx.DB) {
	dk, _ := scrypt.Key([]byte("w23456789"), []byte(os.Getenv("PASS_GEN")), 16384, 8, 1, 32)
	statement := fmt.Sprintf("INSERT INTO user (userid, name, password, mail, createdAt, loggedAt) VALUES('%s', '%s', '%s','%v', '%d', '%d')", bson.NewObjectId().Hex(), "tokeiwwww", dk, "xxxs@qq.com", time.Now().Unix(), time.Now().Unix())
	_, err := db.Exec(statement)

	if err != nil {
		panic(err)
		fmt.Println("Database error")
	}
}
func TestAuth(t *testing.T) {
	//db, _ := _func.MySqlGetDB("nekohand")
	//insertUser(db)
	form := url.Values{}
	form.Add("username", "tokeiwwww")
	form.Add("password", "w23456789")
	req, _ := http.NewRequest("POST", "/v2/backend/token.get", strings.NewReader(form.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	response := executeRequest(req)
	fmt.Println(response.Body)

	type RB struct {
		API_TOKEN string
	}

	var responseBody RB

	decoder := json.NewDecoder(response.Body)
	if err := decoder.Decode(&responseBody); err != nil {
		fmt.Println("JSON Illegal")
		return
	}

	if body := fmt.Sprintf("%v", responseBody.API_TOKEN); body != "[]" {
		req2, _ := http.NewRequest("POST", "/v2/backend/auth/post.create", nil)
		req2.Header.Set("Authorization", body)
		req2.Header.Set("User", "tokeiwwww")
		req2.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		response2 := executeRequest(req2)
		fmt.Println(response2.Body)
		//t.Errorf("Expected an empty array. Got %s", body)
	}

}
