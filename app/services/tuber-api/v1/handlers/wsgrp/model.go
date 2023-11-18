package wsgrp

// ChatMessage represents a chat message with a user ID.
type ChatMessage struct {
	UserID  string `json:"user_id"`
	Message string `json:"message"`
}
