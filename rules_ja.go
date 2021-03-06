package mahjong

import "errors"

type JapaneseBaseRule struct {
	BaseRule
}

func (JapaneseBaseRule) TileAmount() uint8 {
	return 136
}

func (rule JapaneseBaseRule) WallTilesCannotDraw() uint8 {
	return 14 + rule.Maj.KanCount
}

func (rule JapaneseBaseRule) WallLastTile() Tile {
	return rule.Maj.Tiles[rule.WallTilesCannotDraw()]
}

func (JapaneseBaseRule) RyanShanAmount() uint8 {
	return 4
}

func (rule JapaneseBaseRule) CanRiichi(player *Player) ([]Tile, error) {
	tiles := make([]Tile, 0)
	for tileType := Dots1; tileType <= Red; tileType++ {
		if discard, err := rule.canRiichi(player, tileType); err == nil {
			tiles = append(tiles, discard...)
		}
	}
	if len(tiles) == 0 {
		return nil, errors.New("Player can not riichi. ")
	}
	keys := make(map[Tile]bool)
	var uniqueTiles []Tile
	for _, tile := range tiles {
		if _, value := keys[tile]; !value {
			keys[tile] = true
			uniqueTiles = append(uniqueTiles, tile)
		}
	}
	return uniqueTiles, nil
}

func (rule JapaneseBaseRule) canRiichi(player *Player, add TileType) ([]Tile, error) {
	tiles := make([]Tile, 0)
	for i, t := range player.Tiles {
		removed := make([]Tile, len(player.Tiles))
		copy(removed, player.Tiles)
		removed[i] = Tile{TileType: add}
		removed = SortTiles(removed)
		base := rule.Maj.NewWinningHandBase(player, nil, removed)
		baseWin := base.normalWin()
		_7Win := base.is7PairsWin()
		_13Win := base.thirteenOrphansWin()
		if baseWin != nil || _7Win != nil || _13Win != nil {
			tiles = append(tiles, t)
		}
	}
	if len(tiles) != 0 {
		return tiles, nil
	}
	return nil, errors.New("Player can not riichi. ")
}

func (rule JapaneseBaseRule) Riichi(player *Player, tile Tile) error {
	index, err := rule.Maj.CanDahai(player, tile)
	if err != nil {
		return err
	}
	ok := false
	for tileType := Dots1; tileType <= Red; tileType++ {
		removed := make([]Tile, len(player.Tiles))
		copy(removed, player.Tiles)
		removed[index] = Tile{TileType: tileType}
		removed = SortTiles(removed)
		base := rule.Maj.NewWinningHandBase(player, nil, removed)
		baseWin := base.normalWin()
		_7Win := base.is7PairsWin()
		_13Win := base.thirteenOrphansWin()
		if baseWin != nil || _7Win != nil || _13Win != nil {
			ok = true
			break
		}
	}
	if !ok {
		return errors.New("cannot riichi")
	}
	rule.Maj.dahai(player, index)
	player.Riichi = rule.Maj.Jun()
	rule.Maj.Output("Riichi")
	return nil
}

func (rule JapaneseBaseRule) CanRon(player *Player) ([]Agari, error) {
	if rule.FuriTen(player) {
		return nil, errors.New("FuriTen")
	}
	if rule.Maj.LastTilePlayer == player {
		return nil, errors.New("")
	}
	return rule.CanAgari(player, rule.Maj.LastTile)
}

func (rule JapaneseBaseRule) Ron(player *Player) error {
	agaris, err := rule.CanRon(player)
	if err != nil {
		return err
	}
	menZen := player.Concealed()
	maxScoreSrc := ScoreSrc(0)
	var result Agari
	for _, agari := range agaris {
		if src := NewScore(agari.Fu, CountFan(agari.YakuTachi, menZen)); src > maxScoreSrc {
			maxScoreSrc = src
			result = agari
		}
	}
	rule.Maj.Output("Ron!", "Base score:", maxScoreSrc)
	for _, yaku := range result.YakuTachi {
		rule.Maj.Output(yaku.Name)
	}
	if player.IsParent(rule.Maj.Round) {
		s := maxScoreSrc.ParentRon()
		player.Score += s
		rule.Maj.LastTilePlayer.Score -= s
	} else {
		s := maxScoreSrc.ChildRon()
		player.Score += s
		rule.Maj.LastTilePlayer.Score -= s
	}
	return nil
}

func (rule JapaneseBaseRule) CanTsumo(player *Player) ([]Agari, error) {
	return rule.CanAgari(player, Tile{})
}

func (rule JapaneseBaseRule) Tsumo(player *Player) error {
	agaris, err := rule.CanTsumo(player)
	if err != nil {
		return err
	}
	menZen := player.Concealed()
	maxScoreSrc := ScoreSrc(0)
	for _, agari := range agaris {
		if src := NewScore(agari.Fu, CountFan(agari.YakuTachi, menZen)); src > maxScoreSrc {
			maxScoreSrc = src
		}
	}
	rule.Maj.Output("Tsumo!", "Base score:", maxScoreSrc)
	round := rule.Maj.Round
	if player.IsParent(round) {
		s := maxScoreSrc.ParentTsumo()
		player.Score += s * 3
		rule.Maj.Players.Do(
			func(p *Player) {
				if p != player {
					p.Score -= s
				}
			},
		)
	} else {
		child, parent := maxScoreSrc.ChildTsumo()
		player.Score += child*2 + parent
		rule.Maj.Players.Do(
			func(p *Player) {
				if p == player {
					return
				}
				if p.IsParent(round) {
					p.Score -= parent
				} else {
					p.Score -= child
				}
			},
		)
	}
	return nil
}

//????????????
func (rule JapaneseBaseRule) CanNineYaochus(player *Player) error {
	if !rule.Maj.Jun().First() {
		return errors.New("not first jun")
	}
	m := make(map[TileType]bool)
	for _, tile := range player.Tiles {
		if !tile.IsYaochu() {
			continue
		}
		if _, t := m[tile.TileType]; !t {
			m[tile.TileType] = true
		}
	}
	if len(m) < 9 {
		return errors.New("yaochus not enough")
	}
	return nil
}

func (rule JapaneseBaseRule) NineYaochus(player *Player) error {
	if err := rule.CanNineYaochus(player); err != nil {
		return err
	}
	rule.Maj.Output("????????????")
	return nil
}

func (rule JapaneseBaseRule) FuriTen(player *Player) bool {
	for _, discard := range player.Discards {
		if _, err := rule.CanAgari(player, discard.Tile); err != nil {
			return true
		}
	}
	if player.Through {
		return true
	}
	return false
}

func (rule JapaneseBaseRule) DoraHints(count uint8, ura bool) []uint8 {
	indexes := make([]uint8, 0)
	for i := uint8(0); i < count; i += 2 {
		indexes = append(indexes, rule.TileAmount()-rule.RyanShanAmount()-i)
		if ura {
			indexes = append(indexes, rule.TileAmount()-rule.RyanShanAmount()-i-1)
		}
	}
	return indexes
}

func (JapaneseBaseRule) PlayersSitDown() Players {
	players := NewPlayers(4)
	for i := 0; i < 4; i++ {
		*players.Now() = Player{FieldWind: FieldWind(i), Score: 25000, Tiles: make([]Tile, 0)}
		players.ToNext()
	}
	return players
}

func (rule JapaneseBaseRule) CanAgari(player *Player, last Tile) ([]Agari, error) {
	agaris := rule.Agaris(player, last)
	if len(agaris) > 0 {
		return agaris, nil
	}
	return nil, errors.New("Player cannot agari. ")
}

type Agari struct {
	YakuTachi []Yaku
	Fu        int
}

func (rule JapaneseBaseRule) Agaris(player *Player, last Tile) []Agari {
	menZen := player.Concealed()
	var agaris []Agari
	tiles := make([]Tile, len(player.Tiles))
	copy(tiles, player.Tiles)
	if last.TileType != None {
		tiles = append(tiles, last)
	}
	base := rule.Maj.NewWinningHandBase(player, rule.Maj.Players.Now(), SortTiles(tiles))
	if handsBase := base.normalWin(); handsBase != nil {
		for _, hand := range handsBase {
			yakuTachi := RealYaku(FindYaku(hand), menZen)
			if len(yakuTachi) > 0 {
				agaris = append(agaris, Agari{YakuTachi: yakuTachi, Fu: hand.CountFu(menZen)})
			}
		}
	}
	if hand7 := base.is7PairsWin(); hand7 != nil {
		yakuTachi := RealYaku(FindYaku(hand7), player.Concealed())
		agaris = append(agaris, Agari{YakuTachi: yakuTachi, Fu: hand7.CountFu(menZen)})
	}
	if hand13 := base.thirteenOrphansWin(); hand13 != nil {
		agaris = append(agaris, Agari{[]Yaku{{Name: "????????????", FanFR: ??????, FanMZ: ??????}}, 0})
	}

	return agaris
}

type ResultType int8

const (
	NoResult    ResultType = 0
	AgariResult ResultType = 1
	DrawResult  ResultType = 2
)

type ResultData struct {
	Tsumo bool
	*Player
	Agari
}

type Result struct {
	ResultType
	data []ResultData
}

func (result *Result) AddDraw() {
	result.ResultType = DrawResult
}

func (result *Result) AddAgari(tsumo bool, player *Player, agari Agari) {
	result.ResultType = AgariResult
	resultData := ResultData{tsumo, player, agari}
	result.data = append(result.data, resultData)
}
func (result *Result) Init() {
	result.ResultType = NoResult
	result.data = nil
}

func (result *Result) Done() bool {
	if result.ResultType == NoResult {
		return false
	}
	return true
}

func (rule JapaneseBaseRule) Tiles() []Tile {
	tiles := make([]Tile, rule.TileAmount())
	for i := range tiles {
		tiles[i] = Tile{TileType: TileType(i/4 + 1), Id: int8(i % 4)}
	}
	return tiles
}

type JapaneseTonPuuRule struct {
	JapaneseBaseRule
}

func (JapaneseTonPuuRule) MaxRound() *Round {
	return NewRound(TonBa, 4)
}

type JapaneseHanChanRule struct {
	JapaneseBaseRule
}

func (JapaneseHanChanRule) MaxRound() *Round {
	return NewRound(NanBa, 4)
}

const (
	TonBa = EastField
	NanBa = SouthField
	shaBa = WestField
	PeiBa = NorthField
)

type Fan int8

const (
	?????? Fan = 0
	?????? Fan = 1
	?????? Fan = 2
	?????? Fan = 3
	?????? Fan = 5
	?????? Fan = 6
	?????? Fan = 13
)

type Yaku struct {
	Name    string
	FanFR   Fan
	FanMZ   Fan
	Upgrade []Yaku
	Check   func(WinningHand) bool
}

func FindYaku(hand WinningHand) []Yaku {
	return findYaku(hand, YakuTachi, make([]Yaku, 0))
}
func findYaku(hand WinningHand, progress []Yaku, result []Yaku) []Yaku {
	for _, yaku := range progress {
		if yaku.Check(hand) {
			l := len(result)
			newResult := findYaku(hand, yaku.Upgrade, result)
			if len(newResult) != l {
				result = newResult
			} else {
				result = append(result, yaku)
			}
		}
	}
	return result
}

func RealYaku(yakuTachi []Yaku, menZen bool) []Yaku {
	nonYakuMan := make([]Yaku, 0)
	yakuMan := make([]Yaku, 0)
	for _, yaku := range yakuTachi {
		//judge concealed
		if !menZen && yaku.FanFR == ?????? {
			continue
		}
		yaku.Upgrade = nil

		//judge yakuMan
		if yaku.FanMZ == ?????? {
			yakuMan = append(yakuMan, yaku)
		} else {
			nonYakuMan = append(nonYakuMan, yaku)
		}
	}

	if len(yakuMan) == 0 {
		return nonYakuMan
	} else {
		return yakuMan
	}
}

func CountFan(yakuTachi []Yaku, menZen bool) Fan {
	var fan Fan
	for _, yaku := range yakuTachi {
		if menZen {
			fan += yaku.FanMZ
		} else {
			fan += yaku.FanFR
		}
	}
	return fan
}

var YakuTachi = []Yaku{
	{
		Name:  "??????",
		FanFR: ??????,
		FanMZ: ??????,
		Check: WinningHand.readyHand,
		Upgrade: []Yaku{
			{
				Name:  "???????????????",
				FanFR: ??????,
				FanMZ: ??????,
				Check: WinningHand.doubleReady,
			},
		},
	},
	{
		Name:  "??????",
		FanFR: ??????,
		FanMZ: ??????,
		Check: WinningHand.oneShot,
	},
	{
		Name:  "??????????????????",
		FanFR: ??????,
		FanMZ: ??????,
		Check: WinningHand.selfPick,
	},
	{
		Name:  "?????????",
		FanFR: ??????,
		FanMZ: ??????,
		Check: WinningHand.allSimples,
	},
	{
		Name:  "??????",
		FanFR: ??????,
		FanMZ: ??????,
		Check: WinningHand.allRuns,
	},
	{
		Name:  "?????????",
		FanFR: ??????,
		FanMZ: ??????,
		Check: WinningHand.doubleRun,
		Upgrade: []Yaku{
			{
				Name:  "?????????",
				FanFR: ??????,
				FanMZ: ??????,
				Check: WinningHand.twoDoubleRuns,
			},
		},
	},
	{
		Name:  "??????(???)",
		FanFR: ??????,
		FanMZ: ??????,
		Check: WinningHand.honorTiles,
		Upgrade: []Yaku{
			{
				Name:  "?????????",
				FanFR: ??????,
				FanMZ: ??????,
				Check: WinningHand.prevailingWind,
			},
			{
				Name:  "?????????",
				FanFR: ??????,
				FanMZ: ??????,
				Check: WinningHand.playerWind,
			},
			{
				Name:  "?????????",
				FanFR: ??????,
				FanMZ: ??????,
				Check: WinningHand.whiteDragon,
			},
			{
				Name:  "?????????",
				FanFR: ??????,
				FanMZ: ??????,
				Check: WinningHand.greenDragon,
			},
			{
				Name:  "?????????",
				FanFR: ??????,
				FanMZ: ??????,
				Check: WinningHand.redDragon,
			},
		},
	},
	{
		Name:  "????????????",
		FanFR: ??????,
		FanMZ: ??????,
		Check: WinningHand.kingsTileDraw,
	},
	{
		Name:  "??????",
		FanFR: ??????,
		FanMZ: ??????,
		Check: WinningHand.addAQuad,
	},
	{
		Name:  "????????????",
		FanFR: ??????,
		FanMZ: ??????,
		Check: WinningHand.finalTurnWinSeaMoon,
	},
	{
		Name:  "????????????",
		FanFR: ??????,
		FanMZ: ??????,
		Check: WinningHand.finalTurnWinRiverFish,
	},
	{
		Name:  "????????????",
		FanFR: ??????,
		FanMZ: ??????,
		Check: WinningHand.threeColorRuns,
	},
	{
		Name:  "????????????",
		FanFR: ??????,
		FanMZ: ??????,
		Check: WinningHand.fullStraight,
	},
	{
		Name:  "???????????????",
		FanFR: ??????,
		FanMZ: ??????,
		Check: WinningHand.terminalOrHonorInEachSet,
		Upgrade: []Yaku{
			{
				Name:  "?????????",
				FanFR: ??????,
				FanMZ: ??????,
				Check: WinningHand.allTerminalsAndHonors,
			},
			{
				Name:  "???????????????",
				FanFR: ??????,
				FanMZ: ??????,
				Check: WinningHand.terminalInEachSet,
			},
		},
	},
	{
		Name:  "?????????",
		FanFR: ??????,
		FanMZ: ??????,
		Check: WinningHand.sevenPairs,
	},
	{
		Name:  "?????????",
		FanFR: ??????,
		FanMZ: ??????,
		Check: WinningHand.allTripletHand,
	},
	{
		Name:  "?????????",
		FanFR: ??????,
		FanMZ: ??????,
		Check: WinningHand.threeClosedTriples,
	},
	{
		Name:  "????????????",
		FanFR: ??????,
		FanMZ: ??????,
		Check: WinningHand.threeColourTriplets,
	},
	{
		Name:  "?????????",
		FanFR: ??????,
		FanMZ: ??????,
		Check: WinningHand.threeKans,
	},
	{
		Name:  "?????????",
		FanFR: ??????,
		FanMZ: ??????,
		Check: WinningHand.littleThreeDragons,
	},
	{
		Name:  "?????????",
		FanFR: ??????,
		FanMZ: ??????,
		Check: WinningHand.halfFlush,
		Upgrade: []Yaku{
			{
				Name:  "?????????",
				FanFR: ??????,
				FanMZ: ??????,
				Check: WinningHand.fullFlush,
			},
		},
	},
	{
		Name:  "????????????",
		FanFR: ??????,
		FanMZ: ??????,
		Check: WinningHand.thirteenOrphans,
	},
	{
		Name:  "?????????",
		FanFR: ??????,
		FanMZ: ??????,
		Check: WinningHand.fourConcealedTriples,
	},
	{
		Name:  "?????????",
		FanFR: ??????,
		FanMZ: ??????,
		Check: WinningHand.bigThreeDragons,
	},
	{
		Name:  "?????????",
		FanFR: ??????,
		FanMZ: ??????,
		Check: WinningHand.allHonors,
	},
	{
		Name:  "?????????",
		FanFR: ??????,
		FanMZ: ??????,
		Check: WinningHand.littleFourWinds,
	},
	{
		Name:  "?????????",
		FanFR: ??????,
		FanMZ: ??????,
		Check: WinningHand.bigFourWinds,
	},
	{
		Name:  "?????????",
		FanFR: ??????,
		FanMZ: ??????,
		Check: WinningHand.allGreen,
	},
	{
		Name:  "?????????",
		FanFR: ??????,
		FanMZ: ??????,
		Check: WinningHand.allTerminals,
	},
	{
		Name:  "?????????",
		FanFR: ??????,
		FanMZ: ??????,
		Check: WinningHand.fourKans,
	},
	{
		Name:  "????????????",
		FanFR: ??????,
		FanMZ: ??????,
		Check: WinningHand.nineGates,
	},
	{
		Name:  "??????",
		FanFR: ??????,
		FanMZ: ??????,
		Check: WinningHand.heavenlyHand,
	},
	{
		Name:  "??????",
		FanFR: ??????,
		FanMZ: ??????,
		Check: WinningHand.handOfEarth,
	},
}
