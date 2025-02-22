package e

const (
	SUCCESS               = 200
	UpdatePasswordSuccess = 201
	NotExistInentifier    = 202
	ERROR                 = 500
	InvalidParams         = 400
	ErrorDatabase         = 40001

	WebsocketLinkSuccess    = 50000
	WebsocketSuccessMessage = 50001
	WebsocketSuccess        = 50002
	WebsocketEnd            = 50003
	WebsocketOnlineReply    = 50004
	WebsocketOfflineReply   = 50005
	WebsocketLimit          = 50006
)

const (
	ChatSystemMsg      = 0
	ChatUserCommonMsg  = 1
	ChatMessageHistory = 2
	ChatMessageNew     = 3
	ChatMessageFile    = 4
)
