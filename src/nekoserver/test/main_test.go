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

var Pass_gen = "000000"

var Neko_token = "4a2c4b"

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

	os.Setenv("PASS_GEN", Pass_gen)

	os.Setenv("NEKO_TOKEN", Neko_token)

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
const categoryTableCreationQuery = `CREATE TABLE IF NOT EXISTS category
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
	router.AssignFrontendRouter(engine)
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
	req, _ := http.NewRequest("GET", "/v2/frontend/status", strings.NewReader(form.Encode()))
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
	//fmt.Println("123456")
	t.Skipped()
	req, _ := http.NewRequest("GET", "/v2/frontend/categories", nil)
	response := executeRequest(req)
	fmt.Println(response.Body)
}

func TestFetchPosts(t *testing.T) {
	//db, _ := _func.MySqlGetDB("nekohand")
	//insertPost(db)
	t.Skipped()
	form := url.Values{}
	fmt.Println(bson.NewObjectId().Hex(), "ID")
	form.Add("pageNumber",  "1")
	form.Add("pageSize", "10")
	req, _ := http.NewRequest("POST", "/v2/frontend/posts", strings.NewReader(form.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	response := executeRequest(req)
	fmt.Println(response.Body)
}

func TestFetchOnePost(t *testing.T) {
	t.Skipped()
	form := url.Values{}
	req, _ := http.NewRequest("POST", "/v2/frontend/post/5b6c42b25c964c10a4c68d19", strings.NewReader(form.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	response := executeRequest(req)
	fmt.Println(response.Body)
}

func TestFetchPostsByCategory(t *testing.T) {
	t.Skipped()
	form := url.Values{}
	form.Add("pageNumber",  "1")
	form.Add("pageSize", "10")
	req, _ := http.NewRequest("POST", "/v2/frontend/posts/5b6c42b25c964c10a4c68d1a", strings.NewReader(form.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	response := executeRequest(req)
	fmt.Println(response.Body)
}

func TestFetchPostsChronology(t *testing.T) {
	t.Skipped()
	form := url.Values{}
	req, _ := http.NewRequest("GET", "/v2/frontend/posts-chronology", strings.NewReader(form.Encode()))
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
	t.Skipped()
	form := url.Values{}
	form.Add("token", "0003020")
	req, _ := http.NewRequest("POST", "/v2/frontend/comments/5b6c5ced5c964c2770f66e8a", strings.NewReader(form.Encode()))
	response := executeRequest(req)
	fmt.Println(response.Body)
}

func TestCreateComment(t *testing.T) {
	t.Skipped()
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
	req, _ := http.NewRequest("POST", "/v2/frontend/c2a5cc3b070", bytes.NewBuffer(jsonStr))
	req.Header.Add("Content-Type", "application/json")
	response := executeRequest(req)
	fmt.Println(response.Body)
}

var usrr = "eddie32"

var pwdd = "safdafasd"

func insertUser(db *sqlx.DB) {
	fmt.Println(pwdd)
	fmt.Println(os.Getenv("PASS_GEN"))
	dk, _ := scrypt.Key([]byte(pwdd), []byte(os.Getenv("PASS_GEN")), 16384, 8, 1, 32)
	fmt.Println(dk, base64.StdEncoding.EncodeToString(dk))
	statement := fmt.Sprintf("INSERT INTO user (userid, name, password, mail, createdAt, loggedAt) VALUES('%s', '%s', '%s','%s', '%d', '%d')", bson.NewObjectId().Hex(), usrr, base64.StdEncoding.EncodeToString(dk), "xxxs@qq.com", time.Now().Unix(), time.Now().Unix())
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
	form.Add("username", usrr)
	form.Add("password", pwdd)
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


	if body := fmt.Sprintf("%v", responseBody.API_TOKEN); body != "" {
		//var p data.Post
		//p.Id = "5bc76af45c964c2e0c24109b"
		//p.Body = "xin"
		//p.Password = "wwwwwwwwwww"
		//p.Author = usrr
		//p.Slug = "abcdefg"
		//p.PTitle = "abasdfaasfasdfasfg"
		//p.Category = "5b6c42a95c964c0eb4896fe9"
		//p.Status = "Public"
		//mp, _ := json.Marshal(p)
		//req2, _ := http.NewRequest("POST", "/v2/backend/auth/post.edit", bytes.NewBuffer(mp))
		//req2.Header.Set("Authorization", body)
		//req2.Header.Set("User", usrr)
		//req2.Header.Add("Content-Type", "application/json")
		//response2 := executeRequest(req2)
		//fmt.Println(response2.Body)

		//form := url.Values{}
		//form.Add("cid", "5bc77a445c964c284c07adcd")
		//req3, _ := http.NewRequest("POST", "/v2/backend/auth/category.delete", strings.NewReader(form.Encode()))
		//req3.Header.Set("Authorization", body)
		//req3.Header.Set("User", usrr)
		//req3.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		//response3 := executeRequest(req3)
		//fmt.Println(response3.Body)

		form := url.Values{}
		form.Add("pid", "5b6c46f25c964c0784be5c22")
		req3, _ := http.NewRequest("POST", "/v2/backend/auth/post.delete", strings.NewReader(form.Encode()))
		//req3.Header.Set("Authorization", body)
		//req3.Header.Set("User", usrr)
		req3.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		response3 := executeRequest(req3)
		fmt.Println(response3.Body)

		//var ccc data.Category
		//ccc.Id = "5bc77b7f5c964c2fb0ccdc27"
		//ccc.CName = "asdfsadfafEDDSAFDFDASFsdfdfedda"
		//ccc.CLink = "aaDASFs"
		//ccc.CInfo = "adfaDAFSADFsf"
		//mp5, _ := json.Marshal(ccc)
		//req4, _ := http.NewRequest("POST", "/v2/backend/auth/category.edit", bytes.NewBuffer(mp5))
		//req4.Header.Set("Authorization", body)
		//req4.Header.Set("User", usrr)
		//req4.Header.Add("Content-Type", "application/json")
		//response4 := executeRequest(req4)
		//fmt.Println(response4.Body)
	} else {
		t.Errorf("Error Generate Token")
	}
}

func TestFetchPostByTime(t *testing.T) {
	form := url.Values{}
	form.Add("t", "1533821000")
	req4, _ := http.NewRequest("POST", "/v2/frontend/po/t", strings.NewReader(form.Encode()))
	req4.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	response4 := executeRequest(req4)
	fmt.Println(response4.Body)
}