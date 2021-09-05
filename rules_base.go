package mahjong

import (
	"errors"
	"math"
	"unsafe"
)

type Rule interface {
	Init(maj *Mahjong)
	PlayersSitDown() Players
	Tiles() []Tile
	MaxRound() *Round
	WallTilesCannotDraw() uint8
	CanRiichi(player *Player) ([]Tile, error)
	Riichi(player *Player, tile Tile) error
	CanRon(player *Player) ([]Agari, error)
	Ron(player *Player) error
	CanTsumo(player *Player) ([]Agari, error)
	Tsumo(player *Player) error
	CanNineYaochus(player *Player) error
	NineYaochus(player *Player) error
	CanAgari(player *Player, last Tile) ([]Agari, error)
}

type BaseRule struct {
	Maj *Mahjong
}

func (rule *BaseRule) Init(maj *Mahjong) {
	rule.Maj = maj
}

func Ceil(src int, degree int) int {
	tail := src % degree
	if tail == 0 {
		return src
	} else {
		return (src/degree + 1) * degree
	}
}

func IsXYZ(ttA, ttB, ttC TileType) bool {
	if !ttA.SameSuit(ttB) || !ttB.SameSuit(ttC) {
		return false
	}
	a := ttA.Number()
	b := ttB.Number()
	c := ttC.Number()
	sequence := (a-b)*(b-c)*2 + 1
	if sequence != 3 && sequence != -3 {
		return false
	}
	return true
}

func IsXX(ttA, ttB TileType) bool {
	if ttA != ttB {
		return false
	}
	return true
}

func IsXXX(ttA, ttB, ttC TileType) bool {
	if ttA != ttB || ttB != ttC {
		return false
	}
	return true
}

func IsXXXX(ttA, ttB, ttC, ttD TileType) bool {
	if ttA != ttB || ttB != ttC || ttC != ttD {
		return false
	}
	return true
}

type WinningHand interface {
	readyHand() bool
	doubleReady() bool
	oneShot() bool
	selfPick() bool
	allSimples() bool
	allRuns() bool
	doubleRun() bool
	twoDoubleRuns() bool
	honorTiles() bool
	prevailingWind() bool
	playerWind() bool
	whiteDragon() bool
	greenDragon() bool
	redDragon() bool
	kingsTileDraw() bool
	addAQuad() bool
	finalTurnWinSeaMoon() bool
	finalTurnWinRiverFish() bool
	threeColorRuns() bool
	fullStraight() bool
	terminalOrHonorInEachSet() bool
	allTerminalsAndHonors() bool
	terminalInEachSet() bool
	sevenPairs() bool
	allTripletHand() bool
	threeClosedTriples() bool
	threeColourTriplets() bool
	threeKans() bool
	littleThreeDragons() bool
	halfFlush() bool
	fullFlush() bool
	thirteenOrphans() bool
	fourConcealedTriples() bool
	allHonors() bool
	bigThreeDragons() bool
	littleFourWinds() bool
	bigFourWinds() bool
	allGreen() bool
	allTerminals() bool
	fourKans() bool
	nineGates() bool
	heavenlyHand() bool
	handOfEarth() bool

	CountFu(bool) int
}

//和了形
type WinningHandBase struct {
	Rule
	Jun
	LastTile              Tile
	Player                *Player
	Atm                   *Player
	Round                 Round
	SortedTileTypes       []TileType
	SortedHandTiles       []Tile
	RemainderTilesCanDraw uint8
}

func (maj *Mahjong) NewWinningHandBase(player, atm *Player, sortedHandTiles []Tile) *WinningHandBase {

	return &WinningHandBase{
		maj.Rule,
		maj.Jun(),
		maj.LastTile,
		player,
		atm,
		maj.Round,
		SortTileTypes(toTileTypes(sortedHandTiles)),
		sortedHandTiles,
		maj.RemainderTilesCanDraw(),
	}
}

func (base *WinningHandBase) newWinningHandNormal(
	triplets []Triplet, sequential []Sequential, quad []Quad, head TilesXX,
) *WinningHandNormal {
	hand := &WinningHandNormal{
		WinningHandBase: *base, Triplets: triplets, Sequential: sequential, Quad: quad, Head: head,
	}

	hand.XX[0] = head[0].TileType
	hand.XX[1] = head[1].TileType
	hand.XXXs = make([][3]TileType, len(triplets))
	hand.XYZs = make([][3]TileType, len(sequential))

	for i, triplet := range triplets {
		t := triplet.TilesXXX[0].TileType
		hand.XXXs[i] = [3]TileType{t, t, t}
	}
	for i, sequential := range sequential {
		for j, tile := range sequential.TilesXYZ {
			hand.XYZs[i][j] = tile.TileType
		}
	}
	return hand
}

func (base *WinningHandBase) normalWin() []*WinningHandNormal {
	player := base.Player
	tiles := base.SortedHandTiles
	hands := make([]*WinningHandNormal, 0)
	for i := 0; i < len(tiles)-1; i++ {
		if !IsXX(tiles[i].TileType, tiles[i+1].TileType) {
			continue
		}
		head := TilesXX{tiles[i], tiles[i+1]}
		newTT := make([]Tile, len(tiles))
		copy(newTT, tiles)
		others := append(newTT[:i], newTT[i+2:]...)
		i++

		tripletsCombines, another := FindTriplet(others)
		for j, combine := range tripletsCombines {
			o := tripletsCombines[len(tripletsCombines)-j-1]
			oo := SortTiles(tripletsSlice2Tiles(o, another))
			seq, err := FindSequential(oo)
			if err != nil {
				continue
			}
			combine = append(combine, player.XXXs...)
			seq = append(seq, player.XYZs...)
			hand := base.newWinningHandNormal(combine, seq, player.XXXXs, head)
			hands = append(hands, hand)
		}
	}
	if len(hands) == 0 {
		return nil
	}
	return hands
}

func (base *WinningHandBase) is7PairsWin() *WinningHand7 {
	hand := base.is7PairsWith4SameWin()
	if hand == nil {
		return nil
	}
	sortedTileTypes := base.SortedTileTypes
	for i := 0; i < len(sortedTileTypes)-2; i++ {
		if IsXX(sortedTileTypes[i], sortedTileTypes[i+2]) {
			return nil
		}
	}
	return hand
}

func (base *WinningHandBase) is7PairsWith4SameWin() *WinningHand7 {
	sortedTileTypes := base.SortedTileTypes
	hand := WinningHand7{WinningHandBase: *base}
	for i := 0; i < 14; i += 2 {
		if !IsXX(sortedTileTypes[i], sortedTileTypes[i+1]) {
			return nil
		}
		hand.XX[i/2] = [2]TileType{sortedTileTypes[i], sortedTileTypes[i]}
	}
	return &hand
}

func (base *WinningHandBase) thirteenOrphansWin() *WinningHand13 {
	sortedTileTypes := base.SortedTileTypes
	kokusi := new(WinningHand13)
	kokusi.WinningHandBase = *base
	keys := make(map[TileType]bool)
	var tileTypes13 []TileType
	for _, tileType := range sortedTileTypes {
		if !tileType.IsYaochu() {
			return nil
		}
		if _, value := keys[tileType]; !value {
			keys[tileType] = true
			tileTypes13 = append(tileTypes13, tileType)
		} else {
			kokusi.Yaochu = tileType
		}

	}
	if len(tileTypes13) != 13 {
		return nil
	}

	s1 := *(*string)(unsafe.Pointer(&tileTypes13))
	s2 := *(*string)(unsafe.Pointer(&Yaochu))
	if s1 != s2 {
		return nil
	}

	return kokusi
}

func (base *WinningHandBase) readyHand() bool {
	return base.Player.Riichi.First()
}

func (base *WinningHandBase) doubleReady() bool {
	return base.Player.Riichi == 1
}

func (base *WinningHandBase) oneShot() bool {
	return base.Player.Riichi == base.Jun
}

func (base *WinningHandBase) selfPick() bool {
	return base.Player.LastDraw != base.LastTile
}

func (base *WinningHandBase) allSimples() bool {
	for _, tileType := range base.SortedTileTypes {
		if tileType.IsYaochu() {
			return false
		}
	}
	return true
}

func (base *WinningHandBase) finalTurnWinSeaMoon() bool {
	return base.selfPick() && base.RemainderTilesCanDraw == 0
}

func (base *WinningHandBase) finalTurnWinRiverFish() bool {
	return !base.selfPick() && base.RemainderTilesCanDraw == 0
}

func (base *WinningHandBase) halfFlush() bool {
	tileTypes := base.SortedTileTypes
	if !tileTypes[0].IsSuit() {
		return false
	}
	suit := tileTypes[0].Suit()
	for i := 1; i < len(tileTypes); i++ {
		if !tileTypes[i].IsHonor() && tileTypes[i].Suit() != suit {
			return false
		}
	}
	return true
}

func (base *WinningHandBase) fullFlush() bool {
	tileTypes := base.SortedTileTypes
	if !tileTypes[len(tileTypes)-1].IsSuit() {
		return false
	}

	return true
}

func (base *WinningHandBase) allHonors() bool {
	for _, tileType := range base.SortedTileTypes {
		if !tileType.IsHonor() {
			return false
		}
	}
	return true
}

func (base *WinningHandBase) heavenlyHand() bool {
	return base.Jun == 1 && base.Player.IsParent(base.Round) && base.selfPick()
}

func (base *WinningHandBase) handOfEarth() bool {
	return base.Jun == 1 && !base.Player.IsParent(base.Round) && base.selfPick()
}

//一般的な和了形
type WinningHandNormal struct {
	WinningHandBase
	XX   [2]TileType
	XXXs [][3]TileType
	XYZs [][3]TileType

	Triplets   []Triplet
	Sequential []Sequential
	Quad       []Quad
	Head       TilesXX
}

func (hand *WinningHandNormal) allRuns() bool {
	player := hand.Player
	if len(hand.XYZs) != 4 || hand.XX[0].IsDragon() {
		return false
	}
	if hand.XX[0].IsActiveWind(player.FieldWind) || hand.XX[0].IsActiveWind(hand.Round.FieldWind) {
		return false
	}
	has := false
	for _, xyz := range hand.XYZs {
		for i, tileType := range xyz {
			if tileType == player.LastDraw.TileType && i != 1 && !(tileType.Number() == 3 && i == 2) && !(tileType.Number() == 7 && i == 0) {
				has = true
			}
		}
	}
	return has
}

func (hand *WinningHandNormal) doubleRun() bool {
	l := len(hand.XYZs)
	switch l {
	case 2:
		return hand.XYZs[0] == hand.XYZs[1]
	case 3:
		return hand.XYZs[0] == hand.XYZs[1] || hand.XYZs[1] == hand.XYZs[2] || hand.XYZs[0] == hand.XYZs[2]
	case 4:
		return hand.XYZs[0] == hand.XYZs[1] || hand.XYZs[1] == hand.XYZs[2] || hand.XYZs[0] == hand.XYZs[2] ||
			hand.XYZs[0] == hand.XYZs[3] || hand.XYZs[1] == hand.XYZs[3] || hand.XYZs[2] == hand.XYZs[3]
	default:
		return false
	}
}

func (hand *WinningHandNormal) twoDoubleRuns() bool {
	l := len(hand.XYZs)
	if l != 4 {
		return false
	}
	xyz0 := hand.XYZs[0]
	theSameIndex := 0
	for i := 0; i < l; i++ {
		if hand.XYZs[i] == xyz0 {
			theSameIndex = i
		}
	}
	if theSameIndex == 0 {
		return false
	}
	var sample [3]TileType
	for i := 1; i < l; i++ {
		if i == theSameIndex {
			continue
		}
		if sample == hand.XYZs[i] {
			return true
		} else {
			sample = hand.XYZs[i]
		}
	}
	return false
}

func (hand *WinningHandNormal) honorTiles() bool {
	var i8 int8 = 0
	for _, xxx := range hand.XXXs {
		if xxx[0].IsActiveWind(hand.Player.Wind(hand.Round)) {
			i8++
		}
		if xxx[0].IsActiveWind(hand.Round.FieldWind) {
			i8++
		}
		if xxx[0].IsDragon() {
			i8++
		}
	}
	return i8 > 0
}

func (hand *WinningHandNormal) prevailingWind() bool {
	wind := hand.Round.FieldWind
	for _, xxx := range hand.XXXs {
		if xxx[0].IsActiveWind(wind) {
			return true
		}
	}
	return false
}

func (hand *WinningHandNormal) playerWind() bool {
	wind := hand.Player.Wind(hand.Round)
	for _, xxx := range hand.XXXs {
		if xxx[0].IsActiveWind(wind) {
			return true
		}
	}
	return false
}

func (hand *WinningHandNormal) whiteDragon() bool {
	for _, xxx := range hand.XXXs {
		if White == xxx[0] {
			return true
		}
	}
	return false
}

func (hand *WinningHandNormal) greenDragon() bool {
	for _, xxx := range hand.XXXs {
		if Green == xxx[0] {
			return true
		}
	}
	return false
}

func (hand *WinningHandNormal) redDragon() bool {
	for _, xxx := range hand.XXXs {
		if Red == xxx[0] {
			return true
		}
	}
	return false
}

func (hand *WinningHandNormal) kingsTileDraw() bool {
	for _, quad := range hand.Player.XXXXs {
		if quad.Jun == hand.Jun {
			return true
		}
	}
	return false
}

func (hand *WinningHandNormal) addAQuad() bool {
	if hand.Atm == hand.Player || hand.Atm == nil {
		return false
	}
	for _, quad := range hand.Atm.XXXXs {
		if quad.Jun == hand.Jun {
			return true
		}
	}
	return false
}

func (hand *WinningHandNormal) threeColorRuns() bool {
	xyz := hand.XYZs
	count := 0
	if len(xyz) < 3 {
		return false
	}
	number := xyz[0][0].Number()
	for i := 1; i < len(xyz); i++ {
		if xyz[i][0].Number() == number {
			count++
		}
	}
	return count >= 2
}

func (hand *WinningHandNormal) fullStraight() bool {
	l := len(hand.XYZs)
	if l < 3 {
		return false
	}

	tts := make([]TileType, 0)
	for _, xyz := range hand.XYZs {
		tts = append(tts, xyz[0])
	}
	split3Suits := split3SuitsType(tts)
	var b1, b4, b7 bool
	for _, suit := range split3Suits {
		if len(suit) >= 3 {
			for _, tileType := range suit {
				switch tileType.Number() {
				case 1:
					b1 = true
				case 4:
					b4 = true
				case 7:
					b7 = true
				}
			}
		}
	}
	return b1 && b4 && b7
}

func (hand *WinningHandNormal) terminalOrHonorInEachSet() bool {
	if !hand.XX[0].IsYaochu() {
		return false
	}
	for _, xxx := range hand.XXXs {
		if !xxx[0].IsYaochu() {
			return false
		}
	}
	for _, xyz := range hand.XYZs {
		if !xyz[0].IsYaochu() && !xyz[2].IsYaochu() {
			return false
		}
	}
	return true
}

func (hand *WinningHandNormal) allTerminalsAndHonors() bool {
	if len(hand.XXXs) != 4 {
		return false
	}
	for _, xxx := range hand.XXXs {
		if !xxx[0].IsYaochu() {
			return false
		}
	}
	return true
}

func (hand *WinningHandNormal) terminalInEachSet() bool {
	if !hand.XX[0].IsTerminals() {
		return false
	}
	for _, xxx := range hand.XXXs {
		if !xxx[0].IsTerminals() {
			return false
		}
	}
	for _, xyz := range hand.XYZs {
		if !xyz[0].IsTerminals() && !xyz[2].IsTerminals() {
			return false
		}
	}
	return true
}

func (hand *WinningHandNormal) sevenPairs() bool {
	return false
}

func (hand *WinningHandNormal) allTripletHand() bool {
	return len(hand.XXXs) == 4
}

func (hand *WinningHandNormal) threeClosedTriples() bool {
	count := len(hand.XXXs)
	count -= len(hand.Player.XXXs)
	for _, xxxx := range hand.Player.XXXXs {
		if !xxxx.Concealed {
			count--
		}
	}
	return count >= 3
}

func (hand *WinningHandNormal) threeColourTriplets() bool {
	if len(hand.XXXs) < 3 {
		return false
	}

	tts := make([]TileType, 0)
	for _, xyz := range hand.XYZs {
		tts = append(tts, xyz[0])
	}
	split3Suits := split3SuitsType(tts)
	for _, suit := range split3Suits {
		if len(suit) == 3 {
			return true
		}
	}
	return false
}

func (hand *WinningHandNormal) threeKans() bool {
	return len(hand.Player.XXXXs) == 3
}

func (hand *WinningHandNormal) littleThreeDragons() bool {
	if !hand.XX[0].IsDragon() {
		return false
	}
	n := 1
	for _, xxx := range hand.XXXs {
		if xxx[0].IsDragon() {
			n++
		}
	}
	return n == 3
}

func (hand *WinningHandNormal) thirteenOrphans() bool {
	return false
}

func (hand *WinningHandNormal) fourConcealedTriples() bool {
	if !hand.allTriples() {
		return false
	}
	return hand.Player.Concealed()
}

func (hand *WinningHandNormal) bigThreeDragons() bool {
	count := 0
	for _, xxx := range hand.XXXs {
		if xxx[0].IsDragon() {
			count++
		}
	}
	return count >= 3
}

func (hand *WinningHandNormal) littleFourWinds() bool {
	count := 0
	for _, xxx := range hand.XXXs {
		if xxx[0].IsWind() {
			count++
		}
	}
	return count == 3 && hand.XX[0].IsWind()
}

func (hand *WinningHandNormal) bigFourWinds() bool {
	count := 0
	for _, xxx := range hand.XXXs {
		if xxx[0].IsWind() {
			count++
		}
	}
	return count == 4
}

func (hand *WinningHandNormal) allGreen() bool {
	for _, tileType := range hand.SortedTileTypes {
		if !tileType.IsGreen() {
			return false
		}
	}
	return true
}

func (hand *WinningHandNormal) allTerminals() bool {
	for _, tileType := range hand.SortedTileTypes {
		if !tileType.IsTerminals() {
			return false
		}
	}
	return true
}

func (hand *WinningHandNormal) fourKans() bool {
	return len(hand.Player.XXXXs) == 4
}

func (hand *WinningHandNormal) nineGates() bool {
	sortedTileTypes := hand.SortedTileTypes
	if !hand.Player.Concealed() || !hand.fullFlush() {
		return false
	}
	return sortedTileTypes[0] == sortedTileTypes[1] && sortedTileTypes[2] == sortedTileTypes[1] &&
		sortedTileTypes[13] == sortedTileTypes[12] && sortedTileTypes[12] == sortedTileTypes[11]
}

func (hand *WinningHandNormal) allTriples() bool {
	if len(hand.XXXs) != 4 {
		return false
	}
	return true
}

func (hand *WinningHandNormal) CountFu(menZen bool) int {
	fu := 20
	tsumo := hand.selfPick()

	if hand.allRuns() {
		if !menZen {
			//食い平和
			return 30
		}

		//ツモ平和
		return fu
	}

	//刻子
	for i, xxx := range hand.XXXs {
		yaochu := xxx[0].IsYaochu()
		switch hand.isAnKoAnKan(i) {
		case MinKo:
			if yaochu {
				fu += 4
			} else {
				fu += 2
			}
		case AnKo:
			if yaochu {
				fu += 8
			} else {
				fu += 4
			}
		case MinKan:
			if yaochu {
				fu += 16
			} else {
				fu += 8
			}
		case AnKan:
			if yaochu {
				fu += 32
			} else {
				fu += 16
			}

		}
	}

	//雀頭
	if hand.XX[0].IsDragon() {
		//三元牌
		fu += 2
	} else {
		//場風
		if hand.XX[0].IsActiveWind(hand.Round.FieldWind) {
			fu += 2
		}
		//自風
		if hand.XX[0].IsActiveWind(hand.Player.Wind(hand.Round)) {
			fu += 2
		}
	}

	//待ち
	max := 0
	for _, xyz := range hand.XYZs {
		index, tileType := findTileType(hand.LastTile.TileType, xyz)
		if tileType == None {
			continue
		}
		if index == 1 {
			//嵌張待ち
			max = 2
		} else {
			//辺張待ち
			if tileType.Number() == 7 && index == 0 {
				max = 2
			} else if tileType.Number() == 3 && index == 2 {
				max = 2
			}
		}
	}
	fu += max

	if hand.XX[0] == hand.LastTile.TileType {
		//単騎待ち
		fu += 2
	}

	if tsumo {
		//ツモ符
		fu += 2
	} else if menZen {
		//門前加符
		fu += 10
	}

	return fu
}

//0:明刻 1:暗刻 2:明槓 3:暗槓
func (hand *WinningHandNormal) isAnKoAnKan(i int) FuuroType {
	ko := hand.XXXs[i][0]
	for _, px := range hand.Player.XXXs {
		if px.TilesXXX[0].TileType == ko {
			return MinKo
		}
	}
	for _, px := range hand.Player.XXXXs {
		if px.TilesXXXX == ko {
			if px.Concealed {
				return AnKan
			} else {
				return MinKan
			}
		}
	}
	return AnKo
}

//七対子の和了形
type WinningHand7 struct {
	WinningHandBase
	XX [7][2]TileType
}

func (hand *WinningHand7) allRuns() bool {
	return false
}

func (hand *WinningHand7) doubleRun() bool {
	return false
}

func (hand *WinningHand7) twoDoubleRuns() bool {
	return false
}

func (hand *WinningHand7) honorTiles() bool {
	return false
}

func (hand *WinningHand7) prevailingWind() bool {
	return false
}

func (hand *WinningHand7) playerWind() bool {
	return false
}

func (hand *WinningHand7) whiteDragon() bool {
	return false
}

func (hand *WinningHand7) greenDragon() bool {
	return false
}

func (hand *WinningHand7) redDragon() bool {
	return false
}

func (hand *WinningHand7) kingsTileDraw() bool {
	return false
}

func (hand *WinningHand7) addAQuad() bool {
	return false
}

func (hand *WinningHand7) threeColorRuns() bool {
	return false
}

func (hand *WinningHand7) fullStraight() bool {
	return false
}

//judge allTerminalsAndHonors
func (hand *WinningHand7) terminalOrHonorInEachSet() bool {
	for _, xx := range hand.XX {
		if !xx[0].IsYaochu() {
			return false
		}
	}
	return true
}

//auto upgrade to allTerminalsAndHonors
func (hand *WinningHand7) allTerminalsAndHonors() bool {
	return true
}

func (hand *WinningHand7) terminalInEachSet() bool {
	return false
}

func (hand *WinningHand7) sevenPairs() bool {
	return true
}

func (hand *WinningHand7) allTripletHand() bool {
	return false
}

func (hand *WinningHand7) threeClosedTriples() bool {
	return false
}

func (hand *WinningHand7) threeColourTriplets() bool {
	return false
}

func (hand *WinningHand7) threeKans() bool {
	return false
}

func (hand *WinningHand7) littleThreeDragons() bool {
	return false
}

func (hand *WinningHand7) thirteenOrphans() bool {
	return false
}

func (hand *WinningHand7) fourConcealedTriples() bool {
	return false
}

func (hand *WinningHand7) bigThreeDragons() bool {
	return false
}

func (hand *WinningHand7) littleFourWinds() bool {
	return false
}

func (hand *WinningHand7) bigFourWinds() bool {
	return false
}

func (hand *WinningHand7) allGreen() bool {
	return false
}

func (hand *WinningHand7) allTerminals() bool {
	return false
}

func (hand *WinningHand7) fourKans() bool {
	return false
}

func (hand *WinningHand7) nineGates() bool {
	return false
}

func (hand *WinningHand7) CountFu(bool) int {
	return 25
}

//国士無双の和了形
type WinningHand13 struct {
	WinningHandBase
	Yaochu TileType
}

type ScoreSrc int

const (
	満貫   ScoreSrc = 2000
	跳満   ScoreSrc = 3000
	倍満   ScoreSrc = 4000
	三倍満  ScoreSrc = 6000
	数え役満 ScoreSrc = 8000
)

//切り上げ満貫
func NewScore(fu int, fan Fan) ScoreSrc {
	if fan < 5 {
		if (fan == 3 && fu >= 60) || (fan == 4 && fu >= 30) {
			return 満貫
		} else {
			for i := 0; i < 2+int(fan); i++ {
				fu *= 2
			}
			return ScoreSrc(fu)
		}
	} else if fan == 5 {
		return 満貫
	} else if fan == 6 || fan == 7 {
		return 跳満
	} else if fan >= 8 && fan <= 10 {
		return 倍満
	} else if fan == 11 || fan == 12 {
		return 三倍満
	} else {
		return 数え役満
	}
}

func (score ScoreSrc) Ceil() int {
	return Ceil(int(score), 100)
}

func (score ScoreSrc) ChildRon() int {
	return (score * 4).Ceil()
}

func (score ScoreSrc) ParentRon() int {
	return (score * 6).Ceil()
}

func (score ScoreSrc) ChildTsumo() (int, int) {
	return score.Ceil(), (score * 2).Ceil()
}

func (score ScoreSrc) ParentTsumo() int {
	return (score * 2).Ceil()
}

func FindTriplet(tiles []Tile) ([][]Triplet, []Tile) {
	triplets := make([]Triplet, 0)
	l := len(tiles) - 2
	for i := 0; i < l; i++ {
		if !IsXXX(tiles[i].TileType, tiles[i+1].TileType, tiles[i+2].TileType) {
			continue
		}
		l -= 3
		triplets = append(triplets, Triplet{TilesXXX{tiles[i], tiles[i+1], tiles[i+2]}, true})
		tiles = append(tiles[:i], tiles[i+3:]...)
		i--
	}

	length := int(math.Pow(2, float64(len(triplets))))
	result := make([][]Triplet, length)
	index := 1

	for _, n := range triplets {
		max := index
		for i := 0; i < max; i++ {
			result[index] = findTriplet(result[i], n)
			index++
		}
	}

	return result, tiles
}

func findTriplet(triplets []Triplet, triplet Triplet) []Triplet {
	dst := make([]Triplet, len(triplets)+1)
	copy(dst, triplets)
	dst[len(triplets)] = triplet
	return dst
}

func FindSequential(tiles []Tile) ([]Sequential, error) {
	l := len(tiles)
	_3suits := split3Suits(tiles)
	seq := make([]Sequential, 0)
	for _, suit := range _3suits {
		seq = findSequential(suit, 1, seq)
	}
	if len(seq)*3 == l {
		return seq, nil
	}
	return nil, errors.New("")
}

func findSequential(tiles []Tile, n int8, result []Sequential) []Sequential {
	if n == 8 {
		return result
	}
	if indexes, err := GetIndexes(tiles, n, n+1, n+2); err == nil {
		result = append(
			result,
			Sequential{TilesXYZ: TilesXYZ{tiles[indexes[0]], tiles[indexes[1]], tiles[indexes[2]]}, Concealed: true},
		)
		RemoveIndexes(&tiles, indexes...)
		return findSequential(tiles, n, result)
	} else {
		return findSequential(tiles, n+1, result)
	}
}

func GetIndexes(tiles []Tile, n ...int8) ([]int, error) {
	indexes := make([]int, 0)
	for _, nn := range n {
		for i, t := range tiles {
			if t.TileType.IsSuit() && t.TileType.Number() == nn {
				indexes = append(indexes, i)
				break
			}
		}
	}
	if len(indexes) != len(n) {
		return nil, errors.New("")
	}
	return indexes, nil
}

func RemoveIndexes(tiles *[]Tile, n ...int) {
	l := len(n)
	for i := 0; i < l; i++ {
		n[i] -= i
		*tiles = append((*tiles)[:n[i]], (*tiles)[n[i]+1:]...)
	}
}

func tripletsSlice2Tiles(t []Triplet, others []Tile) []Tile {
	r := make([]Tile, 0)
	for _, triplet := range t {
		for _, xxx := range triplet.TilesXXX {
			r = append(r, xxx)
		}
	}
	r = append(r, others...)
	return r
}
