package sessionhandler

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// Package sessionhandler manages user session lifecycle, including creation, validation,
// expiration, and invalidation of sessions. It stores active sessions in memory and provides
// methods to interact with session data during authentication and authorization.

var ErrSessionNotFound = errors.New("Session not found")

type Session struct {
	UserName   string
	UserId     string
	SessionId  uuid.UUID
	Expiration time.Time
}

type SessionHandler struct {
	activeSessions map[string]Session
}

func New() SessionHandler {

	return SessionHandler{
		activeSessions: make(map[string]Session),
	}
}

func (this SessionHandler) StartSession(userName string, userId string) (string, error) {

	sessionId, err := uuid.NewRandom()
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	expiration := time.Now().Add(15 * time.Minute)

	session := Session{
		UserName:   userName,
		UserId:     userId,
		SessionId:  sessionId,
		Expiration: expiration,
	}
	this.activeSessions[sessionId.String()] = session

	return sessionId.String(), nil

}

func (this SessionHandler) GetSession(sessionId string) (Session, error) {

	activeSession, found := this.activeSessions[sessionId]

	if !found {

		return Session{}, ErrSessionNotFound
	}
	return activeSession, nil

}

func (this SessionHandler) InValidateSession(SessionId string) {

	delete(this.activeSessions, SessionId)

}
