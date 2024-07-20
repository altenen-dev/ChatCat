package handlers

import (
	"fmt"
	"time"
)

func FormatMessage(c Client, m string) string {
	timeStamp := time.Now().Format("2006-01-02 15:04:05")
	msg := fmt.Sprintf("[%s][%s]: %s", timeStamp, c.Name, m)
	return msg
}

func FormatPrompt(c Client) string {
	timeStamp := time.Now().Format("2006-01-02 15:04:05")
	msg := fmt.Sprintf("[%s][%s]: ", timeStamp, c.Name)
	return msg
}
