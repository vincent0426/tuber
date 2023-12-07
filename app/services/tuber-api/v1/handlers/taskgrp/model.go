package taskgrp

// AppSendTaskToMQ contains information needed to send a task to MQ.
type AppPublishSendEmailTask struct {
	Email     string `json:"email"`
	DelayTime int    `json:"delayTime"`
}
