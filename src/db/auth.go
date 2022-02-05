package db

import (
	"log"
	"math/rand"
	"net/http"
	"time"
)


/*
	TokenValidation validates the token if authentication required.
	The function returns a string with the Username and bool value
	which references to whether the Token is from valid user or not.
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
			// Update last interaction with this profile.
			_,err := LocalConn.Exec("UPDATE Profile SET LstUse = datetime(strftime('%s','now'),'unixepoch', 'localtime') WHERE Username = ?;",User)
			if err != nil {
				log.Println(err.Error())
			}
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

	_, err = stt.Exec(profile.Username,profile.Password,0,profile.Token)
	if err != nil {
		return err
	}

	err = transaction.Commit()
	if err != nil {
		return err
	}
	return nil
}

/*
	Permission handle user permissions read/write/edit
	Parameters: Profile and action. Return error.
	Actions:
	set - Set/change permission.
	get - Return permission.
 */
func Permission(profile Profile, action string) (string,error) {
	switch action {
	case "set":
		transaction, err := LocalConn.Begin()
		if err != nil {
			return "",err
		}
		_,err = transaction.Exec("UPDATE Profile SET Perm = ? WHERE Username = ?;",profile.Perm,profile.Username)
		if err != nil {
			return "OK",err
		}
	case "get":
		row := LocalConn.QueryRow("SELECT Perm FROM Profile WHERE Username = ?;",profile.Username)
		var perm string
		err := row.Scan(&perm)
		if err != nil {
			return "",err
		} else {
			return perm,nil
		}
	default:

	}
	return "",nil
}


/*
	DropUser has the function to delete user.
 */
func DropUser(p Profile) error {
	_,err := LocalConn.Exec("DELETE FROM Profile WHERE Username = ?;",p.Username)
	if err != nil {
		return err
	}
	return nil
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