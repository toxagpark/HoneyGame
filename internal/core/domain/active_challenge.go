package domain

type ActiveChallenge struct {
	ID            int
	CreatorUserID int
	CreatorName   string
	Amount        int64
}

func NewActiveChallenge(
	id int,
	creatorUserID int,
	creatorName string,
	amount int64,
) ActiveChallenge {
	return ActiveChallenge{
		ID:            id,
		CreatorUserID: creatorUserID,
		CreatorName:   creatorName,
		Amount:        amount,
	}
}
