package main

import "fmt"
import "golang.org/x/crypto/bcrypt"

func generatePasswordHash(password string) {
	bytes, _ := bcrypt.GenerateFromPassword([]byte(password), 14)
	fmt.Println(string(bytes))
}
func main() {
	generatePasswordHash("pass1")
	generatePasswordHash("pass2")
	generatePasswordHash("pass3")
	generatePasswordHash("pass4")
	generatePasswordHash("pass5")
	generatePasswordHash("pass6")
	generatePasswordHash("pass7")
	generatePasswordHash("pass8")
}
