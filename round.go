package mahjong

type Round struct {
	//場
	FieldWind, MaxFieldWind FieldWind

	//局
	Number, MaxNumber int8

	//累計局数
	RealNumber int8

	//本場
	Honba int8
}

func NewRound(maxField FieldWind, maxNumber int8) *Round {
	return &Round{MaxFieldWind: maxField, MaxNumber: maxNumber}
}

func (round *Round) ToNext(normal bool) {
	if normal {
		round.Number++
		round.Honba = 0
		if round.Number > round.MaxNumber {
			round.FieldWind++
		}
	} else {
		round.Honba++
	}
	round.RealNumber++
}

type FieldWind int8

const (
	EastField FieldWind = iota
	SouthField
	WestField
	NorthField
)

//巡
type Jun int8

func (n Jun) First() bool {
	return n == 0
}
