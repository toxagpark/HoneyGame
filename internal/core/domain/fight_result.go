package domain

// FightResult — исход боя: кто победил, кто проиграл и ставка.
type FightResult struct {
	WinnerUserID int
	LoserUserID  int
	Amount       int64
}
