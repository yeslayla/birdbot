package common

type ExternalChatManager interface {
	SendMessage(user string, message string)
}

type ExternalChatModule interface {
	Initialize(ExternalChatManager)

	RecieveMessage(user User, message string)
}
