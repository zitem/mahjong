package mahjong

import (
	"strconv"
	"strings"
)

type NotTurn struct{}

func (err NotTurn) Error() string {
	return "Not Player's Turn"
}

type NoTiles struct {
	Tiles []Tile
}

func (err NoTiles) Error() string {
	msg := "Player has not all the tiles: "
	list := make([]string, 0)
	for _, tile := range err.Tiles {
		list = append(list, TilesName[tile.TileType]+"("+strconv.Itoa(int(tile.Id))+")")
	}
	msg += strings.Join(list, ", ")
	return msg
}

type WrongPhase struct {
	True  []Phase
	Wrong Phase
}

func (err WrongPhase) Error() string {
	str := make([]string, 0)
	for _, t := range err.True {
		str = append(str, Phase.Name(t))
	}
	true := strings.Join(str, "/")
	return "Player phase should be " + true + ", but is " + Phase.Name(err.Wrong) + " now."
}
