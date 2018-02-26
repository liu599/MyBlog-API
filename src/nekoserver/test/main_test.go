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

	"nekoserver/middleware/func"
	"nekoserver/router"

	"nekoserver/middleware/data"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"

)



func ensureTableExists(db *sqlx.DB) {
	if _, err := db.Exec(tableCreationQuery); err != nil {
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
		Source: "root:rrrrr@tcp(127.0.0.1:3306)/nekohand?charset=utf8",
	}

	var Apps = make(map[string]data.Database)

	Apps["nekohand"] = database

	_func.AssignAppDataBaseList(Apps)

	_func.AssignDatabaseFromList([]string{"nekohand"})

	dbc, err := _func.MySqlGetDB("nekohand")
	if err != nil {
		fmt.Println("Error Database Connection")
		return
	}
	db = dbc
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
const tableCreationQuery = `
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
	clearTable(db)
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
