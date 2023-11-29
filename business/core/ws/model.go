package ws

import "github.com/google/uuid"

type ChatMessage struct {
	UserID   uuid.UUID
	Username string
	ImageURL string
	Message  string
}
