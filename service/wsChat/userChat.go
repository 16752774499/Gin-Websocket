package wsChat

func SendMsg(userId int, message Message) {

	for id, conn := range chatUserConnManager.User {
		if userId == id {
			continue
		}
		err := conn.WriteJSON(message)
		if err != nil {
			return
		}
	}

}
