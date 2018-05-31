/*
NoFace Development Build
Cryptography Core
Reading comms keys and pushing chinese chars
*/
package crypto //Change to crypto on build

import (
	"fmt"
	"log"
	"os/exec"
)

/*
Name:	crypt
Param:	method - encrypt or decrypt
		text - plain/cipher text
		pass - symmetric key 
Return:	plain or ciphertext 
*/
func Crypt(method, text, pass string) string {
	cmd := exec.Command("python3", "crypt.py", method, text, pass)
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("Output failed.")
		log.Fatal(err)
	}
	return string(output)
}

/*
func main() {
	fmt.Println(Crypt("dec", "♈ㆩ윢ꜙᆧ丞投ꝙ㙍곇倿圦캎똱", "password!"))
}*/
