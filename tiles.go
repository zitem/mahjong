package mahjong

import sort2 "sort"

type Tile struct {
	TileType
	Id int8
}

func (tile Tile) IsRed() bool {
	return tile.Id == 0 && (tile.TileType == Dots5 || tile.TileType == Bamboo5 || tile.TileType == Characters5)
}

func (tileType TileType) IsSuit() bool {
	if tileType > Characters9 || tileType < Dots1 {
		return false
	}
	return true
}

func (tileType TileType) IsDots() bool {
	if tileType > Dots9 || tileType < Dots1 {
		return false
	}
	return true
}

func (tileType TileType) IsBamboo() bool {
	if tileType > Bamboo9 || tileType < Bamboo1 {
		return false
	}
	return true
}

func (tileType TileType) IsCharacter() bool {
	if tileType > Characters9 || tileType < Characters1 {
		return false
	}
	return true
}

func (tileType TileType) Suit() Suit {
	return Suit(tileType-1) / 9
}

func (tileType TileType) SameSuit(other TileType) bool {
	if !tileType.IsSuit() || tileType.Suit() != other.Suit() {
		return false
	}
	return true
}

func (tileType TileType) IsHonor() bool {
	if tileType > Red || tileType < East {
		return false
	}
	return true
}

func (tileType TileType) IsWind() bool {
	if tileType > North || tileType < East {
		return false
	}
	return true
}

func (tileType TileType) IsActiveWind(w1 FieldWind) bool {
	if !tileType.IsWind() {
		return false
	}
	wind := tileType.Wind()
	return wind == w1
}

func (tileType TileType) Wind() FieldWind {
	return FieldWind(tileType - East)
}

func (tileType TileType) IsDragon() bool {
	if tileType > Red || tileType < White {
		return false
	}
	return true
}

func (tileType TileType) IsTerminals() bool {
	if tileType != Bamboo1 && tileType != Bamboo9 && tileType != Dots1 &&
		tileType != Dots9 && tileType != Characters1 && tileType != Characters9 {
		return false
	}
	return true
}

func (tileType TileType) IsGreen() bool {
	if tileType != Bamboo2 && tileType != Bamboo3 && tileType != Bamboo4 &&
		tileType != Bamboo6 && tileType != Bamboo8 && tileType != Green {
		return false
	}
	return true
}

func (tileType TileType) Number() int8 {
	n := int8(tileType % 9)
	if n == 0 {
		return 9
	}
	return n
}

func (tileType TileType) IsYaochu() bool {
	return !((tileType > Dots1 && tileType < Dots9) ||
		(tileType > Bamboo1 && tileType < Bamboo9) ||
		(tileType > Characters1 && tileType < Characters9))
}

func SortTileTypes(tileTypes []TileType) []TileType {
	sort2.Slice(tileTypes, func(i, j int) bool { return tileTypes[i] < tileTypes[j] })
	return tileTypes
}

func SortTiles(tiles []Tile) []Tile {
	sort2.Slice(tiles, func(i, j int) bool {
		return uint8(tiles[i].Id)+uint8(tiles[i].TileType)*4 < uint8(tiles[j].Id)+uint8(tiles[j].TileType)*4
	})
	return tiles
}

func split3Suits(tiles []Tile) [3][]Tile {
	dots := make([]Tile, 0)
	bamboo := make([]Tile, 0)
	character := make([]Tile, 0)

	for _, tile := range tiles {
		if tile.IsDots() {
			dots = append(dots, tile)
		} else if tile.IsBamboo() {
			bamboo = append(bamboo, tile)
		} else if tile.IsCharacter() {
			character = append(character, tile)
		}
	}

	return [3][]Tile{dots, bamboo, character}
}

func split3SuitsType(tileTypes []TileType) [3][]TileType {
	dots := make([]TileType, 0)
	bamboo := make([]TileType, 0)
	character := make([]TileType, 0)

	for _, tileType := range tileTypes {
		if tileType.IsDots() {
			dots = append(dots, tileType)
		} else if tileType.IsBamboo() {
			bamboo = append(bamboo, tileType)
		} else if tileType.IsCharacter() {
			character = append(character, tileType)
		}
	}

	return [3][]TileType{dots, bamboo, character}
}

func for2Tile(tiles []Tile) []TilesXY {
	var a Tile
	var b Tile
	ab := make([]TilesXY, 0)
	for i := 0; i < len(tiles); i++ {
		a = tiles[i]
		for j := i + 1; j < len(tiles); j++ {
			b = tiles[j]
			ab = append(ab, TilesXY{a, b})
		}
	}
	return ab
}

func toTileTypes(tiles []Tile) []TileType {
	tts := make([]TileType, 0)
	for _, tile := range tiles {
		tts = append(tts, tile.TileType)
	}
	return tts
}

func toSampleTiles(tileTypes []TileType) []Tile {
	tiles := make([]Tile, 0)
	for _, tileType := range tileTypes {
		tiles = append(tiles, Tile{TileType: tileType})
	}
	return tiles
}

func ToSampleTiles(tileTypes []TileType) []Tile {
	return toSampleTiles(tileTypes)
}

func findTileType(tileType TileType, xyz [3]TileType) (int, TileType) {
	for i, t := range xyz {
		if t == tileType {
			return i, t
		}
	}
	return -1, None
}

type TilesXX [2]Tile
type TilesXXX [3]Tile
type TilesXY [2]Tile
type TilesXYZ [3]Tile

type Sequential struct {
	TilesXYZ
	Concealed bool
}

func (sequential *Sequential) ToTileType() []TileType {
	tt := make([]TileType, len(sequential.TilesXYZ))
	for i, tile := range sequential.TilesXYZ {
		tt[i] = tile.TileType
	}
	return tt
}

type Triplet struct {
	TilesXXX
	Concealed bool
}

func (triplet *Triplet) ToTileType() []TileType {
	tt := make([]TileType, len(triplet.TilesXXX))
	for i, tile := range triplet.TilesXXX {
		tt[i] = tile.TileType
	}
	return tt
}

type Quad struct {
	TilesXXXX TileType
	Concealed bool
	Jun
}

type DiscardTile struct {
	Tile
	Jun
	TsumoGiri bool
}

type FuuroType int8

const (
	MinKo FuuroType = iota
	AnKo
	MinKan
	AnKan
	ShunTsu
)

type Suit int8

const (
	DotsSuit Suit = iota
	BambooSuit
	CharacterSuit
)

type TileType int8

var Yaochu = []TileType{Dots1, Dots9, Bamboo1, Bamboo9, Characters1, Characters9, East, South, West, North, White, Green, Red}

const (
	None TileType = iota
	Dots1
	Dots2
	Dots3
	Dots4
	Dots5
	Dots6
	Dots7
	Dots8
	Dots9
	Bamboo1
	Bamboo2
	Bamboo3
	Bamboo4
	Bamboo5
	Bamboo6
	Bamboo7
	Bamboo8
	Bamboo9
	Characters1
	Characters2
	Characters3
	Characters4
	Characters5
	Characters6
	Characters7
	Characters8
	Characters9

	East
	South
	West
	North
	White
	Green
	Red

	PlumBlossom
	Orchid
	Chrysanthemum
	Bamboo
	Spring
	Summer
	Autumn
	Winter
)

var TilesName = map[TileType]string{
	Dots1:       "1p", //🀙
	Dots2:       "2p", //🀚
	Dots3:       "3p", //🀛
	Dots4:       "4p", //🀜
	Dots5:       "5p", //🀝
	Dots6:       "6p", //🀞
	Dots7:       "7p", //🀟
	Dots8:       "8p", //🀠
	Dots9:       "9p", //🀡
	Bamboo1:     "1s", //🀐
	Bamboo2:     "2s", //🀑
	Bamboo3:     "3s", //🀒
	Bamboo4:     "4s", //🀓
	Bamboo5:     "5s", //🀔
	Bamboo6:     "6s", //🀕
	Bamboo7:     "7s", //🀖
	Bamboo8:     "8s", //🀗
	Bamboo9:     "9s", //🀘
	Characters1: "1m", //🀇
	Characters2: "2m", //🀈
	Characters3: "3m", //🀉
	Characters4: "4m", //🀊
	Characters5: "5m", //🀋
	Characters6: "6m", //🀌
	Characters7: "7m", //🀍
	Characters8: "8m", //🀎
	Characters9: "9m", //🀏
	East:        "1z", //🀀
	South:       "2z", //🀁
	West:        "3z", //🀂
	North:       "4z", //🀃
	White:       "5z", //🀆
	Green:       "6z", //🀅
	Red:         "7z", //🀄
}
