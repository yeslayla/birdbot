package common

type GameModule interface {
	SendMessage(user string, message string)
	RecieveMessage(user User, message string)
}
