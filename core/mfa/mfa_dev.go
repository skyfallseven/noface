/*
NoFace Development Build
Multi-Factor Authentication
Using Google Authenticator
*/
package mfa

import (
	"crypto/rand"
	"encoding/base32"
	"fmt"
	"io/ioutil"
	"net/url"

	dgoogauth "github.com/dgryski/dgoogauth"
	qr "github.com/rsc/qr"
)

/*
Generate a new secret for Google Auth
Returns: string
*/
func NewSecret() string {
	secret := make([]byte, 6)
	_, err := rand.Read(secret)

	if err != nil {
		panic(err)
	}

	return base32.StdEncoding.EncodeToString(secret)

}

/*
Compile the otpauth URL for the QR code to be scanned into gAuth
Params:
	secret: string, base32 encoded token to generate QR code
	id: string, user id cast to string
	issuer: string, name of NoFace Server
	qrFile: string, path to write the QR code to on local machine
Returns:
	true if QR code was written
*/
func NewQR(secret, id, issuer, qrFile string) bool {
	URL, err := url.Parse("otpauth://totp")
	if err != nil {
		fmt.Printf("URL could not be parsed.\n")
		return false
	}

	URL.Path += "/" + url.PathEscape(issuer) + ":" + url.PathEscape(id)

	params := url.Values{}
	params.Add("secret", secret)
	params.Add("issuer", issuer)
	URL.RawQuery = params.Encode()

	code, err := qr.Encode(URL.String(), qr.Q)
	if err != nil {
		fmt.Printf("QR code could not be encoded\n")
		return false
	}

	b := code.PNG()
	err = ioutil.WriteFile(qrFile, b, 0600)
	if err != nil {
		fmt.Printf("QR code could not be written to file\n")
		return false
	}
	return true
}

/*
Creates new otp object for a user to auth with
Params:
	secret: string, base32 encoded string tied to user
Returns:
	OTPConfig
*/
func NewOTP(secret string) *dgoogauth.OTPConfig {
	otpc := &dgoogauth.OTPConfig{
		Secret:		secret,
		WindowSize:	3,
		HotpCounter:0,
	}
	return otpc
}

/*
Authenticates user given a token and an otp object
Params:
	otpc: OTPConfig, object used to authenticate user
	token: string, code from Google Authenticator app
Returns:
	true if authenticated
*/
func AuthUser(otpc dgoogauth.OTPConfig, token string) bool{
	val, err := otpc.Authenticate(token)
	if err != nil {
		fmt.Println("Could not authenticate token.")
		return false
	}

	if !val {
		return false
	}

	return true
}


