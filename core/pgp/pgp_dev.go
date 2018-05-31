/*
NoFace Development Build
PGP Core Dev
Thanks to jyap808 on GitHub
I probably have no idea what I'm doing
*/

package pgp //switch to pgp when done

import (
	"golang.org/x/crypto/openpgp"
	"encoding/base64"
	"fmt"
	"log"
	"bytes"
	"io/ioutil"
	"os"

)

// Location of server public key (read in from config)
const pubKey = "./nfAlphaKey"

/*
Name:	readPGPKey
Param:	path - file path to ASCII armored PGP key
Return:	EntityList of the armored key ring
*/
func ReadPGPKey(path string) openpgp.EntityList {
	//Read path into file
	keypath, err := os.Open(path)
	if err != nil {
		fmt.Println("Reading file failed.")
		log.Fatal(err)
	}

	//Read armored public key into Entity
	entityList, err := openpgp.ReadArmoredKeyRing(keypath)
	if err != nil {
		fmt.Println("Reading armored key failed.")
		log.Fatal(err)
	}
	return entityList
}

/*
Name:	encryptMessage
Param:	message - plaintext to encrypt
		entity - EntityList of ASCII armored public key
Return:	base64 encoded encrypted message
*/
func EncryptMsg(msg string, entity openpgp.EntityList) string {
	buf := new(bytes.Buffer)
	w, err := openpgp.Encrypt(buf, entity, nil, nil, nil)
	if err != nil {
		fmt.Println("Encryption failed.")
		log.Fatal(err)
	}
	_, err = w.Write([]byte(msg))
	if err != nil {
		fmt.Println("Byte conversion failed.")
		log.Fatal(err)
	}

	err = w.Close()
	if err != nil {
		fmt.Println("Writer close failed.")
		log.Fatal(err)
	}
	// now convert to base64
	bytes, err := ioutil.ReadAll(buf)
	return base64.StdEncoding.EncodeToString(bytes)
}

/*
Name:	DecryptPrivKey
Param:	entList - EntityList of armored PGP private key
		pass	- passhrase to decrypt private key with
Return:	decrypted private key as an EntityList
~Can be used for server decryption of plaintext~
*/
func DecryptPrivKey(entList openpgp.EntityList, pass string) {
	ent := entList[0]
	passphrase := []byte(pass)

	// Decrypt priv key
	if ent.PrivateKey != nil && ent.PrivateKey.Encrypted {
		err := ent.PrivateKey.Decrypt(passphrase)
		if err != nil {
			fmt.Println("Private key decryption failed.")
			log.Fatal(err)
		}
	}

	// Decrypt all subkeys
	for _, subkey := range ent.Subkeys {
		if subkey.PrivateKey != nil && subkey.PrivateKey.Encrypted {
			err := subkey.PrivateKey.Decrypt(passphrase)
			if err != nil {
				fmt.Println("Subkey decryption failed.")
				log.Fatal(err)
			}
		}
	}
}

/*
Name:	DecryptMsg
Param:	msg - base64 encoding string of message
		entList - entity list of decrypted private key
Pre:	msg must be encrypted with corresponding public key
Return:	string of cleartext message
*/
func DecryptMsg(msg string, entList openpgp.EntityList) string {
	plain, err := base64.StdEncoding.DecodeString(msg)
	if err != nil {
		fmt.Println("Base64 Decode Fail.")
		log.Fatal(err)
	}

	mdigest, err := openpgp.ReadMessage(bytes.NewBuffer(plain), entList, nil, nil)
	if err != nil {
		fmt.Println("Error reading message with key.")
		log.Fatal(err)
	}

	bytes, err := ioutil.ReadAll(mdigest.UnverifiedBody)
	return string(bytes)
}
/*
func main() {
	e := ReadPGPKey("./nfPubKey")
	md := EncryptMsg("Brennan Is Gay", e)
	d := ReadPGPKey("./nfPrivKey")
	DecryptPrivKey(d, "sUUUp3rN0f@c3!")
	fmt.Println(DecryptMsg(md, d))
}*/
