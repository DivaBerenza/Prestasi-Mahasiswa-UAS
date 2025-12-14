// package main

// import (
// 	"fmt"

// 	"golang.org/x/crypto/bcrypt"
// )

// func main() {
// 	pass := "admin123" // password yang ingin di-hash
// 	hashed, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
// 	if err != nil {
// 		panic(err)
// 	}
// 	fmt.Println(string(hashed))
// }
