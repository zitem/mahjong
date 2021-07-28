package mahjong

import "testing"

func Test(t *testing.T) {
	players := NewPlayers(4)
	for i := 0; i < 4; i++ {
		*players.Now() = Player{FieldWind: FieldWind(i), Score: 25000, Tiles: make([]Tile, 0)}
		players.ToNext()
	}
	
	for i := 0; i < 4; i++ {
		now := players.Now()
		if now.FieldWind != FieldWind(i) {
			t.Error()
		}
		players.ToNext()
	}
}

func TestPlayers_find(t *testing.T) {
	players := NewPlayers(4)
	for i := 0; i < 4; i++ {
		*players.Now() = Player{FieldWind: FieldWind(i), Score: 25000, Tiles: make([]Tile, 0)}
		players.ToNext()
	}
	
	now := players.Now()
	find := players.find(now)
	if find.Value != now {
		t.Error()
	}
	left := players.Left(now)
	find = players.find(left)
	if find.Value != left {
		t.Error()
	}
	right := players.Right(now)
	find = players.find(right)
	if find.Value != right {
		t.Error()
	}
	if left.FieldWind != NorthField || right.FieldWind != SouthField {
		t.Error()
	}
}
