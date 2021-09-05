package mahjong

import (
	"container/ring"
)

type Player struct {
	FieldWind
	Score int
	Phase
	jun      Jun
	Riichi   Jun
	Tiles    []Tile
	XXXs     []Triplet
	XYZs     []Sequential
	XXXXs    []Quad
	Discards []DiscardTile
	LastDraw Tile
	Through  bool
}

func (player *Player) HasDiscarded(tile Tile) bool {
	for _, discard := range player.Discards {
		if discard.Tile == tile {
			return true
		}
	}
	return false
}

func (player *Player) GetTilesIndexes(tiles ...Tile) ([]int, error) {
	indexes := make([]int, 0)
	for i, t := range player.Tiles {
		for _, tile := range tiles {
			if t == tile {
				indexes = append(indexes, i)
				//tiles = append(tiles[:j], tiles[j+1:]...)
			}
		}
	}
	if len(indexes) != len(tiles) {
		return nil, NoTiles{Tiles: tiles}
	}
	return indexes, nil
}

func (player *Player) GetTileTypeIndexes(tileType TileType) []int {
	indexes := make([]int, 0)
	for i, tile := range player.Tiles {
		if tile.TileType == tileType {
			indexes = append(indexes, i)
		}
	}
	return indexes
}

func (player *Player) Wind(round Round) FieldWind {
	return (player.FieldWind + FieldWind(round.Number)) % 4
}

func (player *Player) IsParent(round Round) bool {
	return player.Wind(round) == round.FieldWind
}

func (player *Player) Concealed() bool {
	return len(player.Tiles) >= 13
}

type Phase int8

func (phase Phase) Name() string {
	switch phase {
	case Idle:
		return "Idle"
	case AddTile:
		return "AddTile"
	case RemoveTile:
		return "RemoveTile"
	case AddTileKan:
		return "AddTileKan"
	default:
		return "Undefined Phase"
	}
}

func (phase Phase) Check(true ...Phase) error {
	for _, t := range true {
		if phase == t {
			return nil
		}
	}
	return WrongPhase{Wrong: phase, True: true}
}

func (phase *Phase) Change(now Phase) {
	*phase = now
}

type Players struct {
	ring *ring.Ring
}

func NewPlayers(n int) Players {
	r := ring.New(n)
	for i := 0; i < 4; i++ {
		r.Value = new(Player)
		r = r.Next()
	}
	return Players{r}
}

func (players *Players) Now() *Player {
	return players.ring.Value.(*Player)
}

func (players *Players) Next() *Player {
	return players.ring.Next().Value.(*Player)
}

func (players *Players) Move(dice int) *Player {
	return players.ring.Move(dice).Value.(*Player)
}

func (players *Players) ToNext() *Player {
	players.ring = players.ring.Next()
	return players.ring.Value.(*Player)
}

func (players *Players) FindField(field FieldWind) *Player {
	var p *Player
	players.ring.Do(
		func(v interface{}) {
			if v.(*Player).FieldWind == field {
				p = v.(*Player)
			}
		},
	)
	return p
}

func (players *Players) FindWind(wind FieldWind, round Round) *Player {
	var p *Player
	players.ring.Do(
		func(v interface{}) {
			if v.(*Player).Wind(round) == wind {
				p = v.(*Player)
			}
		},
	)
	return p
}

func (players *Players) Parent(round Round) *Player {
	return players.FindWind(EastField, round)
}

func (players *Players) Start() *Player {
	return players.FindField(EastField)
}

func (players *Players) Left(player *Player) *Player {
	return players.find(player).Prev().Value.(*Player)
}

func (players *Players) Right(player *Player) *Player {
	return players.find(player).Next().Value.(*Player)
}

func (players *Players) Toimen(player *Player) *Player {
	return players.find(player).Move(2).Value.(*Player)
}

func (players *Players) Set(player *Player) {
	players.ring = players.find(player)
}

func (players *Players) Do(fn func(players *Player)) {
	players.ring.Do(
		func(v interface{}) {
			fn(v.(*Player))
		},
	)
}

func (players *Players) find(player *Player) *ring.Ring {
	r := players.ring
	for i := 0; i < players.ring.Len(); i++ {
		if r.Value == player {
			return r
		}
		r = r.Next()
	}
	return nil
}

const (
	//Pon or Kan
	Idle Phase = iota

	//Chii or Draw
	AddTile

	//AnKan(keep phase) or Dahai
	RemoveTile

	//Draw(kan)
	AddTileKan
)
