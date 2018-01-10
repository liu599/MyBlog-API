package main

import (
	"net/http"
	"fmt"
	"database/sql"
	"regexp"
	"os"
	"github.com/dgrijalva/jwt-go"
	"time"
	"errors"
	"github.com/dgrijalva/jwt-go/request"
	"strings"
	"encoding/json"
)

type administrator struct {
	Name string `json:"username"`
	Password string `json:"password"`
}


func (a *App) GenerateTokenMiddleware(w http.ResponseWriter, r *http.Request) {
	//body := []byte(`
	//	<h1>This page is used to generate token</h1>
	//`)
	//w.Header().Add("Content-Type", "text/html; charset=utf-8")
	//w.Header().Add("Content-Length", fmt.Sprintf("%d", len(body)))
	//w.Write(body)
	if r.Body != nil {
		// json格式
		a.validateJSONLoginMsg(w, r)
	} else {
		// Form格式
		a.validateFormLoginMsg(w, r)
	}

}

func validateUserInput(w http.ResponseWriter, user string, pwd string) error {
	// 用户登录时需要做字符验证
	if len(user)==0 || len(pwd)==0 {
		respondWithError(w, http.StatusBadRequest, "用户名和密码其中之一不能为空")
		return errors.New("用户名和密码其中之一不能为空")
	}

	if m, _ := regexp.MatchString("^\\p{Han}+$", user); m {
		respondWithError(w, http.StatusBadRequest, "用户名不能含有中文字符")
		return errors.New("用户名不能含有中文字符")
	}

	return nil
}


func (a *App) validateJSONLoginMsg(w http.ResponseWriter, r *http.Request) {
	var u administrator
	err := json.NewDecoder(r.Body).Decode(&u)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, "用户JSON格式不合法")
	}

	var user string = u.Name
	var pwd string = u.Password

	// 验证用户输入
	err1 := validateUserInput(w, user, pwd)
	if err1 != nil {
		return
	}

	// 去服务器验证用户名密码
	token, err := validateUsernameAndPassword(a.DB, user, pwd)

	if err != nil {
		respondWithError(w, http.StatusForbidden, err.Error())
		return
	}

	fmt.Printf("%v, %v, %v", user, pwd, token)

	respondWithJSON(w, http.StatusOK, map[string]string{"token": token})
}

func (a *App) validateFormLoginMsg(w http.ResponseWriter, r *http.Request) {
	// 支持以Form格式登陆
	var user string= r.FormValue("user")
	var pwd  string= r.FormValue("pwd")
	fmt.Println(user)
	fmt.Println(pwd)

	// 验证用户输入
	err1 := validateUserInput(w, user, pwd)
	if err1 != nil {
		return
	}
	// 去服务器验证用户名密码
	token, err := validateUsernameAndPassword(a.DB, user, pwd)

	if err != nil {
		respondWithError(w, http.StatusForbidden, err.Error())
		return
	}

	fmt.Printf("%v, %v, %v", user, pwd, token)

	respondWithJSON(w, http.StatusOK, map[string]string{"token": token})
}

func validateUsernameAndPassword(db *sql.DB, user string, pwd string) (string, error) {
	var password string
	var uid int
	var err error = errors.New("authorization failed, check the user name and password")
	// 有可能是通过邮箱验证
	if m, _ := regexp.MatchString(`^([\w\.\_]{2,10})@(\w{1,}).([a-z]{2,4})$`, user); m {
		fmt.Println("email")
		// 查询数据库密码
		statement := fmt.Sprintf("SELECT password, uid FROM user_login_nh WHERE mail='%s'", user)
		db.QueryRow(statement).Scan(&password, &uid)

		fmt.Printf("\n%v, %v\n", password, uid)

		if password == pwd {
			statement = fmt.Sprintf("UPDATE user_login_nh SET logged=%d WHERE uid=%d'", time.Now().Unix(), uid)
			_, err3 := db.Exec(statement)
			if err3 != nil {
				return "cannot update loggin time for the current user, may be a server failure!", err3
			}
			return generateToken(uid, os.Getenv("NEKOHAND_AUTHORIZATION"))
		}

	}
	// 如果不是
	fmt.Println("not an email")

	// 查询数据库密码
	statement := fmt.Sprintf("SELECT password, uid FROM user_login_nh WHERE name='%s'", user)
	db.QueryRow(statement).Scan(&password, &uid)

	fmt.Printf("\n%v\n", password)

	fmt.Printf("\n%v, %v\n", password, uid)

	if password == pwd {
		statement = fmt.Sprintf("UPDATE user_login_nh SET logged='%d' WHERE uid='%d'", time.Now().Unix(), uid)
		_, err3 := db.Exec(statement)
		if err3 != nil {
			return "cannot update loggin time for the current user, may be a server failure!", err3
		}
		return generateToken(uid, os.Getenv("NEKOHAND_AUTHORIZATION"))
	}

	return "Authorization failed", err
}

func generateToken(uid int, environmentVariable string) (string, error) {

	token := jwt.New(jwt.SigningMethodHS256)
	claims := make(jwt.MapClaims)
	claims["exp"] = time.Now().Add(time.Hour * time.Duration(1)).Unix()
	claims["iat"] = time.Now().Unix()
	token.Claims = claims
	fmt.Printf("\n%v\n", uid)
	// str1 := strconv.Itoa(uid)
	tokenString, err := token.SignedString([]byte(environmentVariable))

	fmt.Printf("\n%v\n", tokenString)

	if err != nil {
		fmt.Println("error in convert tokenString")
		return "Cannot convert secret string", err
	}

	return tokenString, nil
}

func ValidateTokenMiddleware(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	fmt.Printf("\n%s\n", r.RequestURI)
	// 如果是登陆请求不需要判断token
	if r.RequestURI == "/token.get" {
		next(w, r)
		return
	}
	// 使用auth作为后台标志, 如果没有auth的话不做验证
	if !strings.Contains(r.RequestURI, "/auth/") {
		next(w, r)
		return
	}
	// 验证Token
	tokenRequest, err := request.ParseFromRequest(r, request.AuthorizationHeaderExtractor,
		func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("NEKOHAND_AUTHORIZATION")), nil
		})

	if err == nil {
		if tokenRequest.Valid {
			next(w, r)
		} else {
			respondWithError(w, http.StatusUnauthorized, "Token is expired, please log in again!")
			// fmt.Fprint(w, "Token is not valid")
		}
	} else {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized access to this resource")
		// fmt.Fprint(w, "Unauthorized access to this resource")
	}

}
