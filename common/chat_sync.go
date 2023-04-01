package common

type ChatSyncModule interface {
	SendMessage(user string, message string)
	RecieveMessage(user User, message string)
}
