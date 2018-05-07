package test

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"
	"time"

	"nekoserver/middleware/data"
	"nekoserver/middleware/func"
	"nekoserver/router"

	gin "github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"gopkg.in/mgo.v2/bson"
)



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

var db *sqlx.DB

func TestMain(m *testing.M) {

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
    id         VARCHAR(50) UNIQUE NOT NULL,
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
    coid       INT(11) NOT NULL AUTO_INCREMENT PRIMARY KEY,
    id         VARCHAR(50) UNIQUE NOT NULL,
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
	engine.Use()

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
	form.Add("token", "0003020")
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
	checkResponseCode(t, http.StatusOK, http.StatusOK)
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