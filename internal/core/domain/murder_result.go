package domain

// MurderResult — исход ограбления улья: куда ткнула лапа, где лежал мёд,
// сколько мёда выиграно (может быть отрицательным — проигранная ставка)
// и новый баланс игрока.
type MurderResult struct {
	PawHive   int
	HoneyHive int
	Honey     int64
	NewHoney  int64
}

func NewMurderResult(
	pawHive int,
	honeyHive int,
	honey int64,
	newHoney int64,
) MurderResult {
	return MurderResult{
		PawHive:   pawHive,
		HoneyHive: honeyHive,
		Honey:     honey,
		NewHoney:  newHoney,
	}
}
