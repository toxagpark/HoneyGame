package domain

type User struct {
	ID       int
	TgChatID int64
	UserName string
	Honey    int64
}

func NewUser(
	id int,
	tgChatId int64,
	userName string,
	honey int64,
) User {
	return User{
		ID:       id,
		TgChatID: tgChatId,
		UserName: userName,
		Honey:    honey,
	}
}

func NewUserUninitialized(
	tgChatId int64,
	userName string,
) User {
	return User{
		ID:       uninitializedID,
		TgChatID: tgChatId,
		UserName: userName,
	}
}
