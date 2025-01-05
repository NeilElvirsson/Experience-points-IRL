package server

import (
	"crypto/sha512"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/NeilElvirsson/Experience-points-IRL/internal/models"
	"github.com/NeilElvirsson/Experience-points-IRL/internal/userrepository"
)

//struct server exporteras, 2 vriabler 1 hoststring en port int
//new dunc tar in host och port, sätter ny instans på severn och returnerar

type Server struct {
	host           string
	port           int
	router         *http.ServeMux
	userRepository userrepository.UserRepository
}

func New(host string, port int, user userrepository.UserRepository) Server {

	s := Server{
		host:           host,
		port:           port,
		router:         http.NewServeMux(),
		userRepository: user,
	}

	s.router.HandleFunc("GET /health", s.health)
	s.router.HandleFunc("POST /user/add", s.addUser)

	return s
}

func (this Server) health(w http.ResponseWriter, req *http.Request) {

	w.WriteHeader(http.StatusNoContent)
}

func (this Server) Start() {

	addr := fmt.Sprintf("%s:%d", this.host, this.port)
	fmt.Printf("Starting server on %s\n", addr)

	err := http.ListenAndServe(addr, this.router)

	if err != nil {
		panic(err)
	}

}

func (this Server) addUser(w http.ResponseWriter, req *http.Request) {

	bytes, err := io.ReadAll(req.Body)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var body addUserRequestBody
	err = json.Unmarshal(bytes, &body)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	hash := sha512.New()
	hashedPassword := hash.Sum([]byte(body.Password))
	formatedString := fmt.Sprintf("%x", hashedPassword)
	fmt.Println(formatedString)

	err = this.userRepository.AddUser(models.User{
		UserName: body.UserName,
		Password: formatedString,
	})
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	w.WriteHeader(http.StatusAccepted)

}
