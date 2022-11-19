package main

import (
	"errors"
)

type Users struct {
	Users []User `toml:"users"`
}

type User struct {
	Name        string
	Uuid        string
	OfflineUuid string
}

// func getUserByUuid(uuid string) (User, error) {

// 	for _, user := range users.Users {
// 		if user.Uuid == uuid || user.OfflineUuid == uuid {
// 			return user, nil
// 		}
// 	}

// 	return User{}, errors.New("Uuid no encontrado en el sistema")

// }

func getUserByName(name string) (User, error) {

	for _, user := range users.Users {
		if user.Name == name {
			return user, nil
		}
	}

	return User{}, errors.New("Nombre no encontrado en el sistema")

}

func getUsersNames() []string {

	var names []string

	for _, user := range users.Users {
		names = append(names, user.Name)
	}

	return names

}
