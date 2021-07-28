package main

import (
	"bufio"
	"fmt"
	"github.com/zitem/mahjong"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

var players *mahjong.Players
var maj *mahjong.Mahjong
var TileInterface = mahjong.TilesName
var self *mahjong.Player

func main() {
	Game()
}

//1607773938258683200
//1607774238180340101
func Tenho() {
	i := int64(0)
	for {
		i++
		seed := time.Now().UnixNano() + i
		fmt.Println("seed:", seed)
		maj = mahjong.InitWithSeed(&mahjong.JapaneseHanChanRule{}, seed)
		players = &maj.Players
		self = players.Now()
		err := maj.Start()
		if err != nil {
			log.Fatalln(err)
		}
		PrintPlayerStatus(self)
		if _, err := maj.CanTsumo(self); err == nil {
			break
		} else {
		}
	}
	fmt.Println("TRY:", i)
}

func Game() {
	//seed := int64(1607774238180340101) //Tenho
	seed := time.Now().UnixNano()
	fmt.Println("seed:", seed)
	maj = mahjong.InitWithSeed(&mahjong.JapaneseHanChanRule{}, seed)
	players = &maj.Players
	self = players.Now()

	err := maj.Start()
	if err != nil {
		log.Fatalln(err)
	}
	//Tenho()
	players.Do(
		func(player *mahjong.Player) {
			PrintPlayerStatus(player)
			fmt.Println()
		},
	)
	PrintActions(self)
	Cli()
}

func PrintPlayerStatus(player *mahjong.Player) {
	line1 := strings.Join(
		[]string{
			"プレイヤー" + WindInterface[player.FieldWind],
			strconv.Itoa(int(player.Score)),
			WindInterface[player.Wind(maj.Round)],
			player.Phase.Name(),
			"Draw: " + TileInterface[player.LastDraw.TileType],
		}, " | ",
	)
	if player == maj.Players.Now() {
		line1 += " <-"
	}
	fmt.Println(line1)
	names := make([]string, 0)
	for _, tile := range mahjong.SortTiles(player.Tiles) {
		names = append(names, TileInterface[tile.TileType])
	}
	s := strings.Join(names, " ")
	for _, xyz := range player.XYZs {
		s += " | " + TileInterface[xyz.TilesXYZ[0].TileType] + " " + TileInterface[xyz.TilesXYZ[1].TileType] + " " + TileInterface[xyz.TilesXYZ[2].TileType]
	}
	for _, xxx := range player.XXXs {
		s += " | " + TileInterface[xxx.TilesXXX[0].TileType] + " " + TileInterface[xxx.TilesXXX[1].TileType] + " " + TileInterface[xxx.TilesXXX[2].TileType]
	}
	for _, xxxx := range player.XXXXs {
		s += " | " + TileInterface[xxxx.TilesXXXX] + " " + TileInterface[xxxx.TilesXXXX] +
			" " + TileInterface[xxxx.TilesXXXX] + " " + TileInterface[xxxx.TilesXXXX]
	}
	fmt.Println(s)
}

func OpponentsTsumoGiri() {
	if players.Now() != self {

	}
}

func PrintActions(player *mahjong.Player) {
	actions := maj.PlayerCan(player)
	if actions.Chii {
		fmt.Println(WindInterface[player.FieldWind], "Chii?")
		fmt.Println(TssXY(actions.ChiiOption))
	}
	if actions.Pon {
		fmt.Println(WindInterface[player.FieldWind], "Pon?")
		fmt.Println(TssXX(actions.PonOption))
	}
	if actions.Kan {
		fmt.Println(WindInterface[player.FieldWind], "Kan?")
		fmt.Println(TssXXX(actions.KanOption))
	}
	if actions.AnKan {
		fmt.Println(WindInterface[player.FieldWind], "AnKan?")
		fmt.Println(TssType(actions.AnKanOption))
	}
	if actions.KaKan {
		fmt.Println(WindInterface[player.FieldWind], "KaKan?")
		fmt.Println(TssTile(actions.KaKanOption))
	}
	if actions.Riichi {
		fmt.Println(WindInterface[player.FieldWind], "Riichi?")
		fmt.Println(TssTile(actions.RiichiOption))
	}
	if actions.Ron {
		fmt.Println(WindInterface[player.FieldWind], "Ron?")
	}
	if actions.Tsumo {
		fmt.Println(WindInterface[player.FieldWind], "Tsumo?")
	}
	if actions.NineYaochus {
		fmt.Println(WindInterface[player.FieldWind], "NineYaochus?")
	}

}

func Cli() {
	scanner := bufio.NewScanner(os.Stdin)
	var cmd []string
	count := 0
	for scanner.Scan() {
		cmd = strings.Fields(scanner.Text())
		count = len(cmd)
		if count == 0 {
			continue
		}
		switch cmd[0] {
		case "draw":
			if count != 2 {
				break
			}
			p, err := strconv.Atoi(cmd[1])
			if err != nil {
				break
			}
			_, err = maj.Draw(players.FindField(mahjong.FieldWind(p)))
			if err != nil {
				fmt.Println(err)
				continue
			}
			PrintPlayerStatus(players.FindField(mahjong.FieldWind(p)))
			continue
		case "dahai":
			if count != 3 {
				break
			}
			p, err := strconv.Atoi(cmd[1])
			if err != nil {
				break
			}
			err = maj.Dahai(
				players.FindField(mahjong.FieldWind(p)),
				firstTile(players.FindField(mahjong.FieldWind(p)), InterfaceTile[cmd[2]]),
			)
			if err != nil {
				fmt.Println(err)
				continue
			}
			PrintPlayerStatus(players.FindField(mahjong.FieldWind(p)))
			players.Do(PrintActions)
			continue
		case "riichi":
			if count != 3 {
				break
			}
			p, err := strconv.Atoi(cmd[1])
			if err != nil {
				break
			}
			err = maj.Riichi(
				players.FindField(mahjong.FieldWind(p)),
				firstTile(players.FindField(mahjong.FieldWind(p)), InterfaceTile[cmd[2]]),
			)
			if err != nil {
				fmt.Println(err)
				continue
			}
			PrintPlayerStatus(players.FindField(mahjong.FieldWind(p)))
			players.Do(PrintActions)
			continue
		case "chii": //chii 0 2m 1 4m 2
			if count != 6 {
				break
			}
			p, err := strconv.Atoi(cmd[1])
			if err != nil {
				break
			}
			t1 := InterfaceTile[cmd[2]]
			n1, err := strconv.Atoi(cmd[3])
			if err != nil {
				break
			}
			t2 := InterfaceTile[cmd[4]]
			n2, err := strconv.Atoi(cmd[5])
			if err != nil {
				break
			}
			err = maj.Chii(
				players.FindField(mahjong.FieldWind(p)),
				mahjong.Tile{TileType: t1, Id: int8(n1)},
				mahjong.Tile{TileType: t2, Id: int8(n2)},
			)
			if err != nil {
				fmt.Println(err)
				continue
			}
			PrintPlayerStatus(players.FindField(mahjong.FieldWind(p)))
			continue
		case "pon": //pon 0 2m 1 4m 2
			if count != 6 {
				break
			}
			p, err := strconv.Atoi(cmd[1])
			if err != nil {
				break
			}
			t1 := InterfaceTile[cmd[2]]
			n1, err := strconv.Atoi(cmd[3])
			if err != nil {
				break
			}
			t2 := InterfaceTile[cmd[4]]
			n2, err := strconv.Atoi(cmd[5])
			if err != nil {
				break
			}
			err = maj.Pon(
				players.FindField(mahjong.FieldWind(p)),
				mahjong.Tile{TileType: t1, Id: int8(n1)},
				mahjong.Tile{TileType: t2, Id: int8(n2)},
			)
			if err != nil {
				fmt.Println(err)
				continue
			}
			PrintPlayerStatus(players.FindField(mahjong.FieldWind(p)))
			continue
		case "kan":
			if count != 2 {
				break
			}
			p, err := strconv.Atoi(cmd[1])
			if err != nil {
				break
			}
			err = maj.Kan(players.FindField(mahjong.FieldWind(p)))
			if err != nil {
				fmt.Println(err)
				continue
			}
			PrintPlayerStatus(players.FindField(mahjong.FieldWind(p)))
			continue
		case "ron":
			if count != 2 {
				break
			}
			p, err := strconv.Atoi(cmd[1])
			if err != nil {
				break
			}
			err = maj.Ron(players.FindField(mahjong.FieldWind(p)))
			if err != nil {
				fmt.Println(err)
				continue
			}
			PrintPlayerStatus(players.FindField(mahjong.FieldWind(p)))
			continue
		case "tsumo":
			if count != 2 {
				break
			}
			p, err := strconv.Atoi(cmd[1])
			if err != nil {
				break
			}
			err = maj.Tsumo(players.FindField(mahjong.FieldWind(p)))
			if err != nil {
				fmt.Println(err)
				continue
			}
			PrintPlayerStatus(players.FindField(mahjong.FieldWind(p)))
			continue
		case "show":
			switch count {
			case 1:
				players.Do(
					func(player *mahjong.Player) {
						PrintPlayerStatus(player)
						fmt.Println()
					},
				)
				continue
			case 2:
				p, err := strconv.Atoi(cmd[1])
				if err != nil {
					break
				}
				PrintPlayerStatus(players.FindField(mahjong.FieldWind(p)))
				continue
			}
		case "s":
			if count != 1 {
				break
			}
			p := players.Now()
			_, err := maj.Draw(p)
			if err != nil {
				fmt.Println(err)
			}
			err = maj.Dahai(p, p.LastDraw)
			if err != nil {
				fmt.Println(err)
			}
			PrintPlayerStatus(p)
			continue
		case "d":
			if count != 1 {
				break
			}
			p := players.Now()
			_, err := maj.Draw(p)
			if err != nil {
				fmt.Println(err)
			}
			PrintPlayerStatus(p)
			PrintActions(p)
			continue
		}
	}
	println("undefined input")
}

func TssType(tss []mahjong.TileType) string {
	var s []string
	for _, tileType := range tss {
		str := TileInterface[tileType]
		s = append(s, str)
	}
	return strings.Join(s, " | ")
}

func TssTile(tss []mahjong.Tile) string {
	var s []string
	for _, tile := range tss {
		str := TileInterface[tile.TileType]
		s = append(s, str)
	}
	return strings.Join(s, " | ")
}

func TssXX(tss []mahjong.TilesXX) string {
	var s []string
	for _, tiles := range tss {
		str := TileInterface[tiles[0].TileType]
		str += "(" + strconv.Itoa(int(tiles[0].Id)) + ") "
		str += TileInterface[tiles[1].TileType]
		str += "(" + strconv.Itoa(int(tiles[1].Id)) + ")"
		s = append(s, str)
	}
	return strings.Join(s, " | ")
}

func TssXXX(tss ...mahjong.TilesXXX) string {
	var s []string
	for _, tiles := range tss {
		str := TileInterface[tiles[0].TileType]
		str += "(" + strconv.Itoa(int(tiles[0].Id)) + ") "
		str += TileInterface[tiles[1].TileType]
		str += "(" + strconv.Itoa(int(tiles[1].Id)) + ") "
		str += TileInterface[tiles[2].TileType]
		str += "(" + strconv.Itoa(int(tiles[2].Id)) + ")"
		s = append(s, str)
	}
	return strings.Join(s, " | ")
}

func TssXY(tss []mahjong.TilesXY) string {
	var s []string
	for _, tiles := range tss {
		str := TileInterface[tiles[0].TileType]
		str += "(" + strconv.Itoa(int(tiles[0].Id)) + ") "
		str += TileInterface[tiles[1].TileType]
		str += "(" + strconv.Itoa(int(tiles[1].Id)) + ")"
		s = append(s, str)
	}
	return strings.Join(s, " | ")
}

func firstTile(player *mahjong.Player, tileType mahjong.TileType) mahjong.Tile {
	for _, tile := range player.Tiles {
		if tile.TileType == tileType {
			return tile
		}
	}
	return mahjong.Tile{}
}

var WindInterface = map[mahjong.FieldWind]string{
	mahjong.EastField:  "東",
	mahjong.SouthField: "南",
	mahjong.WestField:  "西",
	mahjong.NorthField: "北",
}

var InterfaceTile = map[string]mahjong.TileType{
	"1p": mahjong.Dots1,
	"2p": mahjong.Dots2,
	"3p": mahjong.Dots3,
	"4p": mahjong.Dots4,
	"5p": mahjong.Dots5,
	"6p": mahjong.Dots6,
	"7p": mahjong.Dots7,
	"8p": mahjong.Dots8,
	"9p": mahjong.Dots9,
	"1s": mahjong.Bamboo1,
	"2s": mahjong.Bamboo2,
	"3s": mahjong.Bamboo3,
	"4s": mahjong.Bamboo4,
	"5s": mahjong.Bamboo5,
	"6s": mahjong.Bamboo6,
	"7s": mahjong.Bamboo7,
	"8s": mahjong.Bamboo8,
	"9s": mahjong.Bamboo9,
	"1m": mahjong.Characters1,
	"2m": mahjong.Characters2,
	"3m": mahjong.Characters3,
	"4m": mahjong.Characters4,
	"5m": mahjong.Characters5,
	"6m": mahjong.Characters6,
	"7m": mahjong.Characters7,
	"8m": mahjong.Characters8,
	"9m": mahjong.Characters9,
	"1z": mahjong.East,
	"2z": mahjong.South,
	"3z": mahjong.West,
	"4z": mahjong.North,
	"5z": mahjong.White,
	"6z": mahjong.Green,
	"7z": mahjong.Red,
}
