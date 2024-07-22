package handlers

import (
	"fmt"
	"time"
)

// Froamter for message sent by clients to other clients
func FormatMessage(c Client, m string) string {
	timeStamp := time.Now().Format("2006-01-02 15:04:05")
	msg := fmt.Sprintf(Yellow+"[%s][%s]: %s"+ResetColor, timeStamp, c.Name, m)
	return msg
}

// Formater for prompt msg
func FormatPrompt(c Client) string {
	timeStamp := time.Now().Format("2006-01-02 15:04:05")
	msg := fmt.Sprintf("[%s][%s]: ", timeStamp, c.Name)
	return msg
}
