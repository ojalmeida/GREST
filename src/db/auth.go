package db

import (
	"math/rand"
	"net/http"
	"time"
)


/*
	TokenValidation validates the token if authentication required.
	The function returns a string with the Username a bool value
	which references to if the Token is from a valid user or not.
 */
func TokenValidation(r *http.Request) Profile {
	var profile Profile
	Token := r.Header.Get("GToken")
	// if len(Token) > 64 {return "NULL",false}

	if Token != "" {
		row := LocalConn.QueryRow("SELECT Username FROM Profile WHERE Token = ?;",profile.Token)
		var User string
		err := row.Scan(&User)
		if err != nil {
			profile.Authenticated = false
		} else if User == "" {
			profile.Authenticated = false
		} else {
			profile.Username = User
			profile.Authenticated = true
		}
	} else {
		profile.Authenticated = false
	}
	return profile
}

func CreateUser(profile Profile) error {
	transaction, err := LocalConn.Begin()
	if err != nil {
		return err
	}

	stt,err := transaction.Prepare("INSERT INTO Profile (Username, Passkey, Perm, Token, Created, LstUse) VALUES (?,?,?,?,datetime(strftime('%s','now'),'unixepoch', 'localtime'),datetime(strftime('%s','now'),'unixepoch', 'localtime'));")

	if err != nil {
		return err
	}
	
	if profile.Token == "" {
		profile.Token = Randomstring(64)
	}

	_, err = stt.Exec(profile.Username,profile.Password,7,profile.Token)
	if err != nil {
		return err
	}

	err = transaction.Commit()
	if err != nil {
		return err
	}

	return nil
}

func Permission(profile Profile, action string) error {
	switch action {
	case "set":
		
	}
}


func Randomstring(length int) string {
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[r1.Intn(len(charset))]
	}
	return string(b)
}