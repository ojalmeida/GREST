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

/*
	CreateUser creates a new user from the passed profile.
	Username and Passkey are needed.
*/
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
		profile.Token = GenerateToken(64)
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
	GetPermission returns the permission from the wanted user.
	Token is needed.
*/
func GetPermission(profile Profile) (int, error) {
	row := LocalConn.QueryRow("SELECT perm FROM profile WHERE token = ?;", profile.Token)
	var perm int
	err := row.Scan(&perm)
	if err != nil {
		return 0, err
	} else {
		return perm, nil
	}
}

/*
	SetPermission is a func that has the ability to set up
	and change the permissions from a user.
*/
func SetPermission(profile Profile) error {
	transaction, err := LocalConn.Begin()
	if err != nil {
		return err
	}
	_, err = transaction.Exec("UPDATE profile SET perm = ? WHERE username = ?;", profile.Perm, profile.Username)
	if err != nil {
		return err
	}
	return nil
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

// TODO: Use hash function for authentication!!!!

/*
	GetToken seems to be a "login" function.
	It receives the credentials and returns the token.
*/
func GetToken(profile Profile) (string, error) {
	row := LocalConn.QueryRow("SELECT token FROM profile WHERE username = ? AND passkey = ?;", profile.Username, profile.Password)
	var token string
	err := row.Scan(&token)
	if err != nil {
		return "", err
	} else {
		return token, nil
	}
}

func RecreateToken(profile Profile) (string, error) {
	transaction, err := LocalConn.Begin()
	if err != nil {
		return "", err
	}
	_, err = transaction.Exec("UPDATE profile SET token = ? WHERE token = ?;", GenerateToken(64), profile.Perm)
	if err != nil {
		return "", err
	}
	return "OK", nil
}

/*
	GenerateToken fuction creates a new API token with random values
*/
func GenerateToken(length int) string {
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[r1.Intn(len(charset))]
	}
	return string(b)
}

// TODO: [ ] Use hash function for authentication!!!!
// DONE: [X] uniform function names
// DONE: [X] delegate just only thing to each function (maybe SetPermission and GetPermission instead of switch case ?) -> Clean Code chapter 3, by Robert C. Martin (https://cloudflare-ipfs.com/ipfs/bafykbzacedgzwm4qwdxqkq5oy6fxgwebffgxcncjaqkhxhhuawvbrmcdxe2u2?filename=Robert%20C.%20Martin%20-%20Clean%20Code_%20A%20Handbook%20of%20Agile%20Software%20Craftsmanship-Prentice%20Hall%20%282008%29.pdf)
