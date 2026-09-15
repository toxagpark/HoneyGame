package domain

type User struct {
	ID       int
	TgUserID int64
	UserName string
	Honey    int64
}

func NewUser(
	id int,
	tgUserID int64,
	userName string,
	honey int64,
) User {
	return User{
		ID:       id,
		TgUserID: tgUserID,
		UserName: userName,
		Honey:    honey,
	}
}

func NewUserUninitialized(
	tgUserID int64,
	userName string,
) User {
	return User{
		ID:       uninitializedID,
		TgUserID: tgUserID,
		UserName: userName,
	}
}
