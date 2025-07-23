package models

import "sync"

type GameStateStruct struct {
	MU   sync.Mutex
	Data string
}

var GameState = &GameStateStruct{}
