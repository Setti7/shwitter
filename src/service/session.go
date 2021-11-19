package service

import (
	"github.com/Setti7/shwitter/session"
	"sync"
)

var onceSession sync.Once

func initSession() {
	services.Session = session.New()
}

func Session() *session.Session {
	onceSession.Do(initSession)

	return services.Session
}
