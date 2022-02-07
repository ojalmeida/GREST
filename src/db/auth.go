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

// maybe ValidateToken() ?
func TokenValidation(r *http.Request) Profile {
	var profile Profile
	Token := r.Header.Get("GToken")
	// if len(Token) > 64 {return "NULL",false}

	if Token != "" {
		row := LocalConn.QueryRow("SELECT username FROM profile WHERE token = ?;", profile.Token)
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
			_, err := LocalConn.Exec("UPDATE profile SET lstUse = datetime(strftime('%s','now'),'unixepoch', 'localtime') WHERE username = ?;", User)
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

	stt, err := transaction.Prepare("INSERT INTO profile (username, passkey, perm, token, created, lstUse) VALUES (?,?,?,?,datetime(strftime('%s','now'),'unixepoch', 'localtime'),datetime(strftime('%s','now'),'unixepoch', 'localtime'));")

	if err != nil {
		return err
	}

	if profile.Token == "" {
		profile.Token = Randomstring(64)
	}

	_, err = stt.Exec(profile.Username, profile.Password, 0, profile.Token)
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
func Permission(profile Profile, action string) (int, error) {
	switch action {
	case "set":
		transaction, err := LocalConn.Begin()
		if err != nil {
			return 0, err
		}
		_, err = transaction.Exec("UPDATE profile SET perm = ? WHERE username = ?;", profile.Perm, profile.Username)
		if err != nil {
			return 0, err
		}
	case "get":
		row := LocalConn.QueryRow("SELECT perm FROM profile WHERE username = ?;", profile.Username)
		var perm int
		err := row.Scan(&perm)
		if err != nil {
			return 0, err
		} else {
			return perm, nil
		}
	default:

	}
	return 0, nil
}

/*
	DropUser has the function to delete user.
*/
func DropUser(p Profile) error {
	_, err := LocalConn.Exec("DELETE FROM profile WHERE username = ?;", p.Username)
	if err != nil {
		return err
	}
	return nil
}

// maybe GenerateToken() ?
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

// TODO: uniform function names
// TODO: delegate just only thing to each function (maybe SetPermission and GetPermission instead of switch case ?) -> Clean Code chapter 3, by Robert C. Martin (https://cloudflare-ipfs.com/ipfs/bafykbzacedgzwm4qwdxqkq5oy6fxgwebffgxcncjaqkhxhhuawvbrmcdxe2u2?filename=Robert%20C.%20Martin%20-%20Clean%20Code_%20A%20Handbook%20of%20Agile%20Software%20Craftsmanship-Prentice%20Hall%20%282008%29.pdf)
