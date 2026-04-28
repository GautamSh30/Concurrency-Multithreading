package main

import "fmt"

type User struct {
	Name  string
	Email string
}

func (u User) String() string {
	return fmt.Sprintf("%s <%s>", u.Name, u.Email)
}

func main() {
	u := User{Name: "Gaut", Email: "gaut@example.com"}
	fmt.Println(u) // Calls String() automatically
}
