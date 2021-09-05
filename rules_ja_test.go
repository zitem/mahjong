package mahjong

import (
	"fmt"
	"reflect"
	"testing"
)

func TestJapaneseBaseRule_Tiles(t *testing.T) {
	tiles := new(JapaneseBaseRule).Tiles()
	if tiles[0].TileType != Dots1 || tiles[0].Id != 0 {
		t.Error()
	}
	if tiles[135].TileType != Red || tiles[135].Id != 3 {
		t.Error()
	}
}

func args() ([]*WinningHandNormal, WinningHandBase) {
	//tts := SortTileTypes([]TileType{
	//	Characters2, Characters2, Characters3,
	//	Characters4, Characters5, Dots3,
	//	Dots4, Dots5, Bamboo3,
	//	Bamboo4, Bamboo5, West, West, Characters2})
	tts := []TileType{
		Dots1, Dots1, Dots1, Bamboo2, Bamboo2, Characters1, Characters2, Characters3, Green, Green, Green, Red, Red,
		Red,
	}
	tiles := toSampleTiles(tts)
	whb := WinningHandBase{
		Jun:      2,
		LastTile: Tile{South, 0},
		Player: &Player{
			FieldWind: SouthField,
			Score:     25000,
			Phase:     RemoveTile,
			Riichi:    2,
			Tiles:     tiles,
			XXXs:      nil,
			XYZs:      nil,
			XXXXs:     nil,
			Discards:  nil,
			LastDraw:  Tile{TileType: Red},
		},
		Round: Round{
			FieldWind:    SouthField,
			MaxFieldWind: 0,
			Number:       0,
			MaxNumber:    0,
			RealNumber:   0,
		},
		SortedTileTypes:       tts,
		RemainderTilesCanDraw: 13,
	}
	hands := whb.normalWin()

	return hands, whb
}

func TestFindYaku(t *testing.T) {
	hands, whb := args()
	for _, hand := range hands {
		hand.WinningHandBase = whb
		result := findYaku(hand, YakuTachi, make([]Yaku, 0))
		result = RealYaku(result, true)
		fmt.Println(result)
	}
}

func BenchmarkFindYaku(b *testing.B) {
	hands, whb := args()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, hand := range hands {
			hand.WinningHandBase = whb
			findYaku(hand, YakuTachi, make([]Yaku, 0))
			//fmt.Println(result)
		}
	}
}

func TestJapaneseBaseRule_Tsumo0(t *testing.T) {
	maj := Init(&JapaneseHanChanRule{})

	p0 := maj.Players.Now()
	p0.Phase.Change(RemoveTile)
	p0.Tiles = toSampleTiles(
		[]TileType{
			Dots1, Dots1, Dots1, Dots2, Dots2, Dots2, Dots3, Dots3, Dots3, Dots4, Dots4, Dots4, Dots5, Dots5,
		},
	)

	type fields struct {
		BaseRule BaseRule
	}
	type args struct {
		player *Player
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			"",
			fields{BaseRule: BaseRule{maj}},
			args{p0},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				rule := JapaneseBaseRule{
					BaseRule: tt.fields.BaseRule,
				}
				if err := rule.Tsumo(tt.args.player); (err != nil) != tt.wantErr {
					t.Errorf("Tsumo() error = %v, wantErr %v", err, tt.wantErr)
				}
			},
		)
	}
}

func TestJapaneseBaseRule_Tsumo1(t *testing.T) {
	maj := Init(&JapaneseHanChanRule{})
	p1 := maj.Players.Now()
	p1.Phase.Change(RemoveTile)
	p1.Tiles = toSampleTiles(
		[]TileType{
			Dots1, Dots1, Dots1, Dots2, Dots2, Dots2, Dots4, Dots4, Dots4, Dots5, Dots5,
		},
	)
	p1xxx := toSampleTiles([]TileType{Dots3, Dots3, Dots3})
	p1xyz := toSampleTiles([]TileType{Characters1, Characters2, Characters3})
	var p1arr [3]Tile
	copy(p1arr[:], p1xxx[:3])
	p1.XXXs = []Triplet{{p1arr, true}}
	var p2arr [3]Tile
	copy(p2arr[:], p1xyz[:3])
	p1.XYZs = []Sequential{{p2arr, true}}
	type fields struct {
		BaseRule BaseRule
	}
	type args struct {
		player *Player
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			"",
			fields{BaseRule: BaseRule{maj}},
			args{p1},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				rule := JapaneseBaseRule{
					BaseRule: tt.fields.BaseRule,
				}
				if err := rule.Tsumo(tt.args.player); (err != nil) != tt.wantErr {
					t.Errorf("Tsumo() error = %v, wantErr %v", err, tt.wantErr)
				}
			},
		)
	}
}

func TestJapaneseBaseRule_Tsumo2(t *testing.T) {
	maj := Init(&JapaneseHanChanRule{})
	p1 := maj.Players.Now()
	p1.Phase.Change(RemoveTile)
	p1.Tiles = toSampleTiles(
		[]TileType{
			Dots1, Dots2, Dots3, Characters1, Characters2, Characters3, Bamboo2, Bamboo3, Bamboo4, Bamboo5, Bamboo6,
			Bamboo7, East, East,
		},
	)
	type fields struct {
		BaseRule BaseRule
	}
	type args struct {
		player *Player
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			"",
			fields{BaseRule: BaseRule{maj}},
			args{p1},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				rule := JapaneseBaseRule{
					BaseRule: tt.fields.BaseRule,
				}
				if err := rule.Tsumo(tt.args.player); (err != nil) != tt.wantErr {
					t.Errorf("Tsumo() error = %v, wantErr %v", err, tt.wantErr)
				}
			},
		)
	}
}

func TestJapaneseBaseRule_CanRiichi(t *testing.T) {
	type fields struct {
		BaseRule BaseRule
	}
	type args struct {
		player *Player
	}
	rule := JapaneseHanChanRule{}
	maj := Mahjong{Rule: &rule}
	rule.Maj = &maj
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []Tile
		wantErr bool
	}{
		{
			name:   "",
			fields: fields{rule.BaseRule},
			args: args{
				player: &Player{
					Tiles: toSampleTiles(
						[]TileType{
							Characters5, Characters6, Characters7, Characters7, Characters7,
							Dots3, Dots4, Dots5, Dots8, Dots8, Bamboo5, Bamboo5, Bamboo5, East,
						},
					),
				},
			},
			want:    nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				rule := JapaneseBaseRule{
					BaseRule: tt.fields.BaseRule,
				}
				got, err := rule.CanRiichi(tt.args.player)
				if (err != nil) != tt.wantErr {
					t.Errorf("CanRiichi() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if reflect.DeepEqual(got, tt.want) {
					t.Errorf("CanRiichi() got = %v, want %v", got, tt.want)
				}
				for _, tile := range got {
					fmt.Println(TilesName[tile.TileType])
				}
			},
		)
	}
}

func TestJapaneseBaseRule_CanNineYaochus(t *testing.T) {
	type fields struct {
		BaseRule BaseRule
	}
	type args struct {
		player *Player
	}
	rule := JapaneseHanChanRule{}
	maj := Mahjong{Rule: &rule}
	rule.Maj = &maj
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			"",
			fields{rule.BaseRule},
			args{
				player: &Player{
					Tiles: toSampleTiles(
						[]TileType{
							Dots1, Dots3, Dots4, Dots5, Dots6, Dots8, Dots9, Dots9, Dots9, Bamboo1, Bamboo9, Characters1, Characters9, East,
						},
					),
				},
			},
			true,
		},
		{
			"",
			fields{rule.BaseRule},
			args{
				player: &Player{
					Tiles: toSampleTiles(
						[]TileType{
							Dots1, Dots3, Dots4, Dots5, Dots6, Dots8, Dots9, Bamboo1, Bamboo9, Characters1, Characters9, East, West, Green,
						},
					),
				},
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				rule := JapaneseBaseRule{
					BaseRule: tt.fields.BaseRule,
				}
				if err := rule.CanNineYaochus(tt.args.player); (err != nil) != tt.wantErr {
					t.Errorf("CanNineYaochus() error = %v, wantErr %v", err, tt.wantErr)
				}
			},
		)
	}
}
