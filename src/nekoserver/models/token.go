package models

import (
	"bytes"

	"nekoserver/middleware/data"
)

func TokenCheckUser(usr data.User) (error, bool) {
	err, uk := UserFetch(usr.Name)
	if err != nil {
		return err, false
	}
	if !bytes.Equal(uk.Password, usr.Password) {
		return err, false
	}
	return nil, true
}