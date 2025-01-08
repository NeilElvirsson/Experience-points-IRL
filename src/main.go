package main

import (
	"github.com/NeilElvirsson/Experience-points-IRL/internal/logrepository"
	"github.com/NeilElvirsson/Experience-points-IRL/internal/server"
	"github.com/NeilElvirsson/Experience-points-IRL/internal/sessionhandler"
	"github.com/NeilElvirsson/Experience-points-IRL/internal/userrepository"
)

func main() {

	//handler.Test(userrepository.New("../database.db"))
	userRepository := userrepository.New("../database.db?_fk=true")
	sessionHandler := sessionhandler.New()
	logRepository := logrepository.New("../database.db?_fk=true")

	connectServer := server.New("localhost", 42069, userRepository, sessionHandler, logRepository)

	connectServer.Start()

}
