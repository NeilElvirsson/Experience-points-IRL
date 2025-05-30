package server

import (
	"context"
	"crypto/sha512"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/NeilElvirsson/Experience-points-IRL/internal/logrepository"
	"github.com/NeilElvirsson/Experience-points-IRL/internal/models"
	"github.com/NeilElvirsson/Experience-points-IRL/internal/sessionhandler"
	"github.com/NeilElvirsson/Experience-points-IRL/internal/taskrepository"
	"github.com/NeilElvirsson/Experience-points-IRL/internal/userrepository"
)

// Package server sets up and starts the HTTP server, handles routing, and manages API endpoints.
// It includes middleware for authentication and coordinates request handling by interacting with
// user, session, task, and log repositories.

type Server struct {
	host           string
	port           int
	router         *http.ServeMux
	userRepository userrepository.UserRepository
	sessionHandler sessionhandler.SessionHandler
	logRepository  logrepository.LogRepository
	taskRepository taskrepository.Taskrepository
}

func New(host string,
	port int,
	user userrepository.UserRepository,
	sessionHandler sessionhandler.SessionHandler,
	logRepository logrepository.LogRepository,
	taskRepository taskrepository.Taskrepository) Server {

	s := Server{
		host:           host,
		port:           port,
		router:         http.NewServeMux(),
		userRepository: user,
		sessionHandler: sessionHandler,
		logRepository:  logRepository,
		taskRepository: taskRepository,
	}

	s.router.Handle("GET /health", s.authMiddleware(s.health()))
	s.router.HandleFunc("POST /user/add", s.addUser)
	s.router.HandleFunc("POST /user/login", s.loginUser)
	s.router.Handle("GET /user/validate", s.authMiddleware(s.validateUser()))
	s.router.Handle("POST /log", s.authMiddleware(s.addLogEntry()))
	s.router.Handle("POST /task/add", s.authMiddleware(s.addTask()))
	s.router.Handle("GET /log", s.authMiddleware(s.getLogs()))
	s.router.Handle("GET /log/xp", s.authMiddleware(s.getXpLevel()))
	s.router.Handle("POST /user/logout", s.authMiddleware(s.logoutUser()))

	return s
}

func (this Server) authMiddleware(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		sessionId := req.Header.Get("x-session")

		session, err := this.sessionHandler.GetSession(sessionId)

		if err != nil {
			if errors.Is(err, sessionhandler.ErrSessionNotFound) {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			fmt.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if session.Expiration.Before(time.Now()) {
			this.sessionHandler.InValidateSession(sessionId)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		fmt.Printf("Authenticated user %s\n", session.UserName)
		ctx := req.Context()
		ctx = context.WithValue(ctx, "session", session)

		req = req.Clone(ctx)

		next.ServeHTTP(w, req)
	})
}

func (this Server) health() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {

		w.WriteHeader(http.StatusNoContent)
	})

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

func (this Server) loginUser(w http.ResponseWriter, req *http.Request) {

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

	user, err := this.userRepository.LoginUser(body.UserName, formatedString)
	if err != nil {
		if err == userrepository.ErrUserNotFound {
			fmt.Println("Unauthorized")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	sessionId, err := this.sessionHandler.StartSession(user.UserName, user.UserId)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Add("x-session", sessionId)

	fmt.Println("userId: ", user.UserId)
	fmt.Println("sessionId: ", sessionId)

	w.WriteHeader(http.StatusOK)

}

func (this Server) validateUser() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {

		session, ok := req.Context().Value("session").(sessionhandler.Session)

		if !ok {
			fmt.Println("Could not cast session")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Write([]byte(session.UserName))

	})
}

func (this Server) addLogEntry() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {

		session, ok := req.Context().Value("session").(sessionhandler.Session)
		if !ok {
			fmt.Println("Could not cast session")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		bytes, err := io.ReadAll(req.Body)
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		var body addLogTaskRequestBody
		err = json.Unmarshal(bytes, &body)
		if err != nil {

			fmt.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		err = this.logRepository.AddLogEntry(session.UserId, body.TaskId)
		if err != nil {
			fmt.Println("Failed to add log", err)
			w.WriteHeader(http.StatusInternalServerError)
			return

		}
		fmt.Printf("Add task %s to user %s", body.TaskId, session.UserId)
		w.WriteHeader(http.StatusOK)

	})
}

func (this Server) addTask() http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {

		_, ok := req.Context().Value("session").(sessionhandler.Session)
		if !ok {
			fmt.Println("Could not cast session")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		bytes, err := io.ReadAll(req.Body)
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		var body addTaskRequestBody
		err = json.Unmarshal(bytes, &body)
		if err != nil {

			fmt.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		err = this.taskRepository.AddTask(body.TaskName, body.XpValue)
		if err != nil {
			fmt.Println("Failed to add task", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		fmt.Printf("Added %s with value %v", body.TaskName, body.XpValue)

		w.WriteHeader(http.StatusOK)

	})

}

func (this Server) getLogs() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {

		session, ok := req.Context().Value("session").(sessionhandler.Session)
		if !ok {
			fmt.Println("could not cast session")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		logs, err := this.logRepository.GetLogs(session.UserId)
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		var responsBody []getResponseBody

		responsBody = []getResponseBody{}

		for _, log := range logs {

			responsBody = append(responsBody, getResponseBody{
				TaskId:    log.TaskId,
				Timestamp: log.Timestamp,
				TaskName:  log.TaskName,
				XpValue:   log.XpValue,
			})

		}

		bytes, err := json.Marshal(responsBody)
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		_, err = w.Write(bytes)
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

	})
}

func (this Server) getXpLevel() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {

		session, ok := req.Context().Value("session").(sessionhandler.Session)
		if !ok {
			fmt.Println("Could not cast session")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		xpSummary, err := this.logRepository.GetXpLevel(session.UserId)
		if err != nil {
			fmt.Println("Could not get xp level")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		responsBody := getXpSummaryBody{
			TotalXp:  xpSummary.TotalXp,
			Level:    xpSummary.Level,
			Progress: xpSummary.Progress,
		}

		bytes, err := json.Marshal(responsBody)
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		_, err = w.Write(bytes)
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return

		}

	})

}

func (this Server) logoutUser() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {

		session, ok := req.Context().Value("session").(sessionhandler.Session)
		if !ok {
			fmt.Println("Could not cast session")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		this.sessionHandler.InValidateSession(session.SessionId.String())
		w.WriteHeader(http.StatusOK)

	})
}
