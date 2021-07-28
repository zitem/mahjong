package mahjong

import (
	"errors"
	"fmt"
	"math/rand"
	"time"
)

type Mahjong struct {
	Round
	Rule
	Players
	Result
	Seed           *rand.Rand
	NextTile       uint8
	Tiles          []Tile
	LastTile       Tile
	LastTilePlayer *Player
	KanCount       uint8
}

func Init(rule Rule) *Mahjong {
	seed := time.Now().UnixNano()
	return InitWithSeed(rule, seed)
}

func InitWithSeed(rule Rule, seed int64) *Mahjong {
	maj := &Mahjong{
		Rule:    rule,
		Seed:    rand.New(rand.NewSource(seed)),
		Players: rule.PlayersSitDown(),
	}
	maj.Rule.Init(maj)
	return maj
}

func (maj *Mahjong) Start() error {
	maj.Shuffle()
	err := maj.Stack()
	if err != nil {
		return err
	}
	maj.Haipai()
	return nil
}

func (maj *Mahjong) Shuffle() {
	tiles := maj.Rule.Tiles()
	maj.Seed.Shuffle(len(tiles), func(i, j int) { tiles[i], tiles[j] = tiles[j], tiles[i] })
	maj.Tiles = tiles
}

func (maj *Mahjong) Stack() error {
	player := maj.Players.Now()
	if player.Wind(maj.Round) != EastField {
		return errors.New("Is not east Player's turn. ")
	}
	dice := maj.Dice() + maj.Dice()
	kaimenPlayer := maj.Players.Move(dice)
	maj.Output("プレイヤー", kaimenPlayer.FieldWind, "開門、サイコロ", dice)
	indexes := JapaneseBaseRule{}.DoraHints(1, false)
	maj.Output("ドラ", maj.Tiles[indexes[0]].TileType)

	return nil
}

func (maj *Mahjong) Haipai() {
	for i := 0; i < 12; i++ {
		for j := 0; j < 4; j++ {
			maj.draw()
		}
		maj.Players.ToNext()
	}
	maj.draw()
	maj.draw()
	maj.Players.ToNext()
	maj.draw()
	maj.Players.ToNext()
	maj.draw()
	maj.Players.ToNext()
	maj.draw()
	firstPlayer := maj.Players.ToNext()
	firstPlayer.Phase = RemoveTile
	firstPlayer.Jun++
}

func (maj *Mahjong) draw() Tile {
	tile := maj.Tiles[maj.NextTile]
	player := maj.Players.Now()
	player.LastDraw = tile
	player.Tiles = append(player.Tiles, tile)
	maj.NextTile++
	return tile
}

func (maj *Mahjong) DrawKan(player *Player) (Tile, error) {
	err := player.Phase.Check(AddTileKan)
	if err != nil {
		return Tile{}, err
	}
	defer player.Phase.Change(RemoveTile)

	if maj.KanCount > 4 {
		return Tile{}, errors.New("Kanned 4 times. ")
	}
	tile := maj.Tiles[uint8(len(maj.Tiles))-maj.KanCount]
	player.LastDraw = tile
	player.Tiles = append(player.Tiles, tile)
	return tile, nil
}

func (maj *Mahjong) Draw(player *Player) (Tile, error) {
	err := player.Phase.Check(AddTile)
	if err != nil {
		return Tile{}, err
	}
	defer player.Phase.Change(RemoveTile)

	return maj.draw(), nil
}

func (maj *Mahjong) CanDahai(player *Player, tile Tile) (int, error) {
	err := player.Phase.Check(RemoveTile)
	if err != nil {
		return 0, err
	}
	if player.Riichi.First() && tile != player.LastDraw {
		return 0, errors.New("Riichi Player can only discard the last tile. ")
	}
	for i, t := range player.Tiles {
		if t == tile {
			return i, nil
		}
	}
	return 0, errors.New("Has not the tile. ")
}

func (maj *Mahjong) dahai(player *Player, i int) {
	tile := player.Tiles[i]
	player.Tiles = append(player.Tiles[:i], player.Tiles[i+1:]...)
	maj.LastTile = tile
	maj.LastTilePlayer = player
	player.Discards = append(player.Discards, tile)
	maj.toNextPlayer()
}

func (maj *Mahjong) Dahai(player *Player, tile Tile) error {
	i, err := maj.CanDahai(player, tile)
	if err != nil {
		return err
	}
	maj.dahai(player, i)
	player.Phase.Change(Idle)
	return nil
}

func (maj *Mahjong) CanChii(player *Player) ([]TilesXY, error) {
	err := player.Phase.Check(AddTile)
	if err != nil {
		return nil, err
	}

	result := make([]TilesXY, 0)
	sortedTiles := SortTiles(player.Tiles)
	the3Suits := split3Suits(sortedTiles)
	for _, suit := range the3Suits {
		xys := for2Tile(suit)
		for _, tilesXY := range xys {
			if IsXYZ(tilesXY[0].TileType, tilesXY[1].TileType, maj.LastTile.TileType) {
				result = append(result, tilesXY)
			}
		}
	}
	if len(result) != 0 {
		return result, nil
	}
	return nil, errors.New("Can not chii the tile. ")
}

func (maj *Mahjong) Chii(player *Player, tileA, tileB Tile) error {
	if !IsXYZ(tileA.TileType, tileB.TileType, maj.LastTile.TileType) {
		return errors.New("Can not pon the tile. ")
	}
	indexes, err := player.GetTilesIndexes(tileA, tileB)
	if err != nil {
		return err
	}
	err = player.Phase.Check(AddTile)
	if err != nil {
		return err
	}
	defer player.Phase.Change(RemoveTile)

	xyz := [3]Tile{}
	for i, index := range indexes {
		index -= i
		xyz[i+1] = player.Tiles[index]
		player.Tiles = append(player.Tiles[:index], player.Tiles[index+1:]...)
	}
	xyz[0] = maj.LastTile
	player.XYZs = append(player.XYZs, Sequential{xyz, false})
	maj.LastTile = Tile{}
	return nil
}

func (maj *Mahjong) CanPon(player *Player) ([]TilesXX, error) {
	if player.HasDiscarded(maj.LastTile) {
		return nil, errors.New("Discarded. ")
	}
	indexes := player.GetTileTypeIndexes(maj.LastTile.TileType)
	switch len(indexes) {
	case 2:
		return []TilesXX{
			{player.Tiles[indexes[0]], player.Tiles[indexes[1]]},
		}, nil
	case 3:
		return []TilesXX{
			{player.Tiles[indexes[0]], player.Tiles[indexes[1]]},
			{player.Tiles[indexes[1]], player.Tiles[indexes[2]]},
			{player.Tiles[indexes[0]], player.Tiles[indexes[2]]},
		}, nil
	default:
		return nil, errors.New("Can not pon the tile. ")
	}
}

//todo: rotate & move tiles for fuuro
func (maj *Mahjong) Pon(player *Player, tileA, tileB Tile) error {
	err := player.Phase.Check(Idle)
	if err != nil {
		return err
	}
	if player.HasDiscarded(maj.LastTile) {
		return errors.New("Discarded. ")
	}
	if !IsXXX(tileA.TileType, tileB.TileType, maj.LastTile.TileType) {
		return errors.New("Can not pon the tile. ")
	}
	indexes, err := player.GetTilesIndexes(tileA, tileB)
	if err != nil {
		return err
	}
	maj.playerTakeTurn(player)
	player.Phase.Change(RemoveTile)

	xxx := [3]Tile{}
	for i, index := range indexes {
		index -= i
		xxx[i] = player.Tiles[index]
		player.Tiles = append(player.Tiles[:index], player.Tiles[index+1:]...)
	}
	xxx[2] = maj.LastTile
	player.XXXs = append(player.XXXs, Triplet{xxx, false})
	maj.LastTile = Tile{}
	return nil
}

func (maj *Mahjong) CanKan(player *Player) (TilesXXX, error) {
	if player.HasDiscarded(maj.LastTile) {
		return TilesXXX{}, errors.New("Discarded. ")
	}
	indexes := player.GetTileTypeIndexes(maj.LastTile.TileType)
	if len(indexes) != 3 {
		return TilesXXX{}, errors.New("Can not kan the tile. ")
	}
	return TilesXXX{player.Tiles[indexes[0]], player.Tiles[indexes[1]], player.Tiles[indexes[2]]}, nil
}

func (maj *Mahjong) Kan(player *Player) error {
	err := player.Phase.Check(Idle, AddTile)
	if err != nil {
		return err
	}

	indexes := player.GetTileTypeIndexes(maj.LastTile.TileType)
	if len(indexes) < 3 {
		return errors.New("Can not kan. ")
	}
	maj.playerTakeTurn(player)
	player.Phase.Change(AddTileKan)

	xxxx := Quad{Concealed: false, Jun: maj.Jun()}
	for i, index := range indexes {
		index -= i
		player.Tiles = append(player.Tiles[:index], player.Tiles[index+1:]...)
	}
	xxxx.TilesXXXX = maj.LastTile.TileType
	player.XXXXs = append(player.XXXXs, xxxx)
	maj.LastTile = Tile{}
	return nil
}

func (maj *Mahjong) CanAnKan(player *Player) ([]TileType, error) {
	sorted := SortTiles(player.Tiles)
	xxxxs := make([]TileType, 0)
	for i := 0; i < len(sorted)-3; i++ {
		if IsXXXX(sorted[i].TileType, sorted[i+1].TileType, sorted[i+2].TileType, sorted[i+3].TileType) {
			xxxxs = append(xxxxs, sorted[i].TileType)
		}
	}
	if len(xxxxs) == 0 {
		return nil, errors.New("Can not ankan. ")
	}
	return xxxxs, nil
}

func (maj *Mahjong) AnKan(player *Player, tileType TileType) error {
	err := player.Phase.Check(RemoveTile)
	if err != nil {
		return err
	}

	indexes := player.GetTileTypeIndexes(tileType)

	if len(indexes) != 4 {
		return errors.New("Can not kan the tileType. ")
	}
	xxxx := Quad{Concealed: true, Jun: maj.Jun()}
	for i, index := range indexes {
		index -= i
		player.Tiles = append(player.Tiles[:index], player.Tiles[index+1:]...)
	}
	xxxx.TilesXXXX = tileType
	player.XXXXs = append(player.XXXXs, xxxx)
	return nil
}

func (maj *Mahjong) CanKaKan(player *Player) ([]Tile, error) {
	options := make([]Tile, 0)
	for _, xxx := range player.XXXs {
		for _, tile := range player.Tiles {
			if xxx.TilesXXX[0].TileType == tile.TileType {
				options = append(options, tile)
			}
		}
	}
	if len(options) == 0 {
		return nil, errors.New("Can not kakan. ")
	}
	return options, nil
}

func (maj *Mahjong) KaKan(player *Player, tile Tile) error {
	for i, xxx := range player.XXXs {
		if xxx.TilesXXX[0].TileType == tile.TileType {
			player.XXXs = append(player.XXXs[:i], player.XXXs[i+1:]...)
			player.XXXXs = append(player.XXXXs, Quad{tile.TileType, false, maj.Jun()})
			return nil
		}
	}
	return errors.New("Can not kakan. ")
}

func (maj *Mahjong) CanRiichi(player *Player) ([]Tile, error) {
	err := player.Phase.Check(RemoveTile)
	if err != nil {
		return nil, err
	}
	return maj.Rule.CanRiichi(player)
}

func (maj *Mahjong) Riichi(player *Player, tile Tile) error {
	err := player.Phase.Check(RemoveTile)
	if err != nil {
		return err
	}
	err = maj.Rule.Riichi(player, tile)
	if err == nil {
		player.Phase.Change(Idle)
	}
	return err
}

func (maj *Mahjong) CanRon(player *Player) ([]Agari, error) {
	err := player.Phase.Check(Idle)
	if err != nil {
		return nil, err
	}
	return maj.Rule.CanRon(player)
}

func (maj *Mahjong) Ron(player *Player) error {
	err := player.Phase.Check(Idle)
	if err != nil {
		return err
	}
	return maj.Rule.Ron(player)
}

func (maj *Mahjong) CanTsumo(player *Player) ([]Agari, error) {
	err := player.Phase.Check(RemoveTile)
	if err != nil {
		return nil, err
	}
	return maj.Rule.CanTsumo(player)
}

func (maj *Mahjong) Tsumo(player *Player) error {
	err := player.Phase.Check(RemoveTile)
	if err != nil {
		return err
	}
	return maj.Rule.Tsumo(player)
}

func (maj *Mahjong) CanNineYaochus(player *Player) error {
	if !maj.IsTurn(player) {
		return NotTurn{}
	}
	return maj.Rule.CanNineYaochus(player)
}

func (maj *Mahjong) NineYaochus(player *Player) error {
	if !maj.IsTurn(player) {
		return NotTurn{}
	}
	return maj.Rule.NineYaochus(player)
}

func (maj *Mahjong) CanRestart() error {
	if maj.Result.Done() {
		return errors.New("gaming")
	}
	return nil
}

func (maj *Mahjong) Restart() error {
	err := maj.CanRestart()
	if err != nil {
		return err
	}
	err = maj.Start()
	if err != nil {
		return err
	}
	maj.Result.Init()
	return nil
}

func (maj *Mahjong) Dice() int {
	dice := maj.Seed.Intn(6) + 1
	maj.Output("Dice:", dice)
	return dice
}

func (maj *Mahjong) IsTurn(player *Player) bool {
	return player == maj.Players.Now()
}

func (maj *Mahjong) Jun() Jun {
	var max Jun
	maj.Players.Do(
		func(player *Player) {
			if player.Jun > max {
				max = player.Jun
			}
		},
	)
	return max
}

func (maj *Mahjong) Dora() uint8 {
	return maj.KanCount + 1
}

func (maj *Mahjong) RemainderTilesAll() uint8 {
	return uint8(len(maj.Tiles)) - maj.NextTile
}

func (maj *Mahjong) RemainderTilesCanDraw() uint8 {
	return maj.RemainderTilesAll() - maj.Rule.WallTilesCannotDraw()
}

func (maj *Mahjong) playerTakeTurn(player *Player) {
	maj.Players.Now().Phase.Change(Idle)
	maj.Players.Set(player)
	maj.Players.Right(maj.LastTilePlayer).Phase.Change(Idle)
	player.Jun++
}

func (maj *Mahjong) toNextPlayer() *Player {
	player := maj.Players.ToNext()
	player.Phase.Change(AddTile)
	player.Jun++
	return player
}

type PlayerActions struct {
	Draw         bool
	Dahai        bool
	Chii         bool
	ChiiOption   []TilesXY
	Pon          bool
	PonOption    []TilesXX
	Kan          bool
	KanOption    TilesXXX
	AnKan        bool
	AnKanOption  []TileType
	KaKan        bool
	KaKanOption  []Tile
	Riichi       bool
	RiichiOption []Tile
	Ron          bool
	Tsumo        bool
	NineYaochus  bool
}

func (maj *Mahjong) PlayerCan(player *Player) *PlayerActions {
	if player == nil {
		panic("player == nil")
	}

	pa := new(PlayerActions)

	if player.Phase.Check(AddTile) == nil {
		pa.Draw = true
	}
	if player.Phase.Check(RemoveTile) == nil {
		pa.Dahai = true
	}
	if option, err := maj.CanChii(player); err == nil {
		pa.Chii = true
		pa.ChiiOption = option
	}
	if option, err := maj.CanPon(player); err == nil {
		pa.Pon = true
		pa.PonOption = option
	}
	if option, err := maj.CanKan(player); err == nil {
		pa.Kan = true
		pa.KanOption = option
	}
	if option, err := maj.CanAnKan(player); err == nil {
		pa.AnKan = true
		pa.AnKanOption = option
	}
	if option, err := maj.CanKaKan(player); err == nil {
		pa.KaKan = true
		pa.KaKanOption = option
	}
	if option, err := maj.CanRiichi(player); err == nil {
		pa.Riichi = true
		pa.RiichiOption = option
	}
	if _, err := maj.CanRon(player); err == nil {
		pa.Ron = true
	}
	if _, err := maj.CanTsumo(player); err == nil {
		pa.Tsumo = true
	}
	if maj.CanNineYaochus(player) == nil {
		pa.NineYaochus = true
	}

	return pa
}

func (maj *Mahjong) Output(v ...interface{}) {
	fmt.Println(v...)
}
