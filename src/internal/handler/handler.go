package handler

import (
	"fmt"

	"github.com/NeilElvirsson/Experience-points-IRL/internal/userrepository"
)

func Test(uR userrepository.UserRepository) {

	currentUser, err := uR.LoginUser("Mats", "stam")
	if err != nil {
		if err == userrepository.ErrUserNotFound {
			fmt.Println("user not found")
			return
		}
		fmt.Println("error fetsching user", err)

	}

	fmt.Printf("user: %s\n", currentUser.UserName)
	/*
		err := uR.AddUser(domain.User{
			UserName: "Mats",
			Password: "stam",
		})
		if err != nil {
			fmt.Printf("failed to add user %s\n", err.Error())
		}
	*/
}
