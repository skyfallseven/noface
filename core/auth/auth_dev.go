/*
NoFace Development Build
User Authentication and Stuff
*/
package auth

import (
	"golang.org/x/crypto/bcrypt"
	"crypto/sha512"
	"encoding/base64"
	"encoding/hex"
	"log"
	"io"
	"time"
	"math/rand"
	"strings"
)

/*
Name:	newID
Return:	int of user ID
~Server Side~
*/
func NewID() int {
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)
	return r1.Intn(999999999)
}

/*
Name:	hashPass
Param: 	pass - plaintext of user password
Return:	hashed password as string
~Client Side~
*/
func HashPass(pass string) string {
	bytes, _ := bcrypt.GenerateFromPassword([]byte(pass), 14)
	return string(bytes)
}

/*
Name:	passMatch
Param:	hash - hash of user's password (stored in DB)
		pass - user supplied plaintext password
Return: true if password matches hash
TODO: Transmit pass securely??? Client or server side???
*/
func PassMatch(hash, pass string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(pass))
	return err == nil
}

/*
Name:	makeToken
Param:	id - user id
Pre:	user must complete both authentication stages
Return:	string representing session token for user
~Server Side~
*/
func MakeToken(id int) string {
	var tok []string
	tok = append(tok, string(id), time.Now().String())
	//Build Secret
	secret := make([]byte, 12)
	_, err := rand.Read(secret)

	if err != nil {
		log.Fatal(err)
	}
	s := base64.StdEncoding.EncodeToString(secret)
	tok = append(tok, s)
	raw := strings.Join(tok, "")

	//SHA-512 Hashing
	h := sha512.New()
	io.WriteString(h, raw)
	return hex.EncodeToString(h.Sum(nil))
}
