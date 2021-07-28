package mahjong

import (
	"fmt"
	"testing"
)

func TestIsXYZ(t *testing.T) {
	if !IsXYZ(Bamboo1, Bamboo2, Bamboo3) {
		t.Error()
	}
	if !IsXYZ(Bamboo1, Bamboo3, Bamboo2) {
		t.Error()
	}
	if !IsXYZ(Bamboo2, Bamboo1, Bamboo3) {
		t.Error()
	}
	if !IsXYZ(Bamboo2, Bamboo3, Bamboo1) {
		t.Error()
	}
	if !IsXYZ(Bamboo3, Bamboo2, Bamboo1) {
		t.Error()
	}
	if !IsXYZ(Bamboo3, Bamboo1, Bamboo2) {
		t.Error()
	}
	if IsXYZ(Bamboo1, Bamboo1, Bamboo1) {
		t.Error()
	}
	if IsXYZ(Bamboo1, Bamboo2, Dots9) {
		t.Error()
	}
	if IsXYZ(Bamboo1, Dots8, Dots9) {
		t.Error()
	}
	if IsXYZ(East, South, West) {
		t.Error()
	}
	if IsXYZ(None, Dots1, Dots2) {
		t.Error()
	}
	if IsXYZ(Bamboo2, Bamboo3, Characters1) {
		t.Error()
	}
	if IsXYZ(Bamboo2, Bamboo1, Bamboo9) {
		t.Error()
	}
}

func TestWinningHandNormal_nineGates(t *testing.T) {
	tt := []TileType{
		Characters1, Characters1, Characters1, Characters2, Characters3, Characters4, Characters5,
		Characters6, Characters7, Characters7, Characters9, Characters9, Characters9, Characters9,
	}
	hand := WinningHandNormal{
		WinningHandBase: WinningHandBase{SortedTileTypes: tt, Player: &Player{Tiles: toSampleTiles(tt)}},
	}
	if !hand.nineGates() {
		t.Error()
	}
}

func BenchmarkWinningHandNormal_nineGates(b *testing.B) {
	tt := []TileType{
		Characters1, Characters1, Characters1, Characters2, Characters3, Characters4, Characters5,
		Characters6, Characters7, Characters7, Characters9, Characters9, Characters9, Characters9,
	}
	p := &Player{Tiles: toSampleTiles(tt)}
	hand := WinningHandNormal{WinningHandBase: WinningHandBase{SortedTileTypes: tt, Player: p}}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		hand.nineGates()
	}
}

func TestWinningHandNormal_doubleRun(t *testing.T) {
	type fields struct {
		WinningHandBase WinningHandBase
		XX              [2]TileType
		XXXs            [][3]TileType
		XYZs            [][3]TileType
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name: "case 2 true",
			fields: fields{
				XYZs: [][3]TileType{
					{Dots1, Dots2, Dots3},
					{Dots1, Dots2, Dots3},
				},
			},
			want: true,
		},
		{
			name: "case 2 false",
			fields: fields{
				XYZs: [][3]TileType{
					{Dots1, Dots2, Dots3},
					{Dots2, Dots3, Dots4},
				},
			},
			want: false,
		},
		{
			name: "case 3 true",
			fields: fields{
				XYZs: [][3]TileType{
					{Dots1, Dots2, Dots3},
					{Dots1, Dots2, Dots3},
					{Dots1, Dots2, Dots3},
				},
			},
			want: true,
		},
		{
			name: "case 3 false",
			fields: fields{
				XYZs: [][3]TileType{
					{Dots2, Dots3, Dots4},
					{Dots1, Dots2, Dots3},
					{Characters1, Characters2, Characters3},
				},
			},
			want: false,
		},
		{
			name: "case 4 true",
			fields: fields{
				XYZs: [][3]TileType{
					{Dots1, Dots2, Dots3},
					{Dots1, Dots2, Dots3},
					{Dots1, Dots2, Dots3},
					{Dots1, Dots2, Dots3},
				},
			},
			want: true,
		},
		{
			name: "case 4 false",
			fields: fields{
				XYZs: [][3]TileType{
					{Characters1, Characters2, Characters3},
					{Dots1, Dots2, Dots3},
					{Dots2, Dots3, Dots4},
					{Dots3, Dots4, Dots5},
				},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				hand := &WinningHandNormal{
					WinningHandBase: tt.fields.WinningHandBase,
					XX:              tt.fields.XX,
					XXXs:            tt.fields.XXXs,
					XYZs:            tt.fields.XYZs,
				}
				if got := hand.doubleRun(); got != tt.want {
					t.Errorf("doubleRun() = %v, notWant %v", got, tt.want)
				}
			},
		)
	}
}

func Benchmark_thirteenOrphansWin(b *testing.B) {
	tt := append(Yaochu, East)
	tt = SortTileTypes(tt)
	base := WinningHandBase{SortedTileTypes: tt}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		base.thirteenOrphansWin()
	}
}

func TestWinningHandNormal_twoDoubleRuns(t *testing.T) {
	type fields struct {
		WinningHandBase WinningHandBase
		XX              [2]TileType
		XXXs            [][3]TileType
		XYZs            [][3]TileType
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name: "case 3 false",
			fields: fields{
				XYZs: [][3]TileType{
					{Dots2, Dots3, Dots4},
					{Dots2, Dots3, Dots4},
					{Dots2, Dots3, Dots4},
				},
			},
			want: false,
		},
		{
			name: "case 4 true1",
			fields: fields{
				XYZs: [][3]TileType{
					{Dots1, Dots2, Dots3},
					{Dots1, Dots2, Dots3},
					{Dots1, Dots2, Dots3},
					{Dots1, Dots2, Dots3},
				},
			},
			want: true,
		},
		{
			name: "case 4 true2",
			fields: fields{
				XYZs: [][3]TileType{
					{Dots1, Dots2, Dots3},
					{Dots1, Dots2, Dots3},
					{Characters1, Characters2, Characters3},
					{Characters1, Characters2, Characters3},
				},
			},
			want: true,
		},
		{
			name: "case 4 true3",
			fields: fields{
				XYZs: [][3]TileType{
					{Characters1, Characters2, Characters3},
					{Dots1, Dots2, Dots3},
					{Dots1, Dots2, Dots3},
					{Characters1, Characters2, Characters3},
				},
			},
			want: true,
		},
		{
			name: "case 4 true4",
			fields: fields{
				XYZs: [][3]TileType{
					{Dots1, Dots2, Dots3},
					{Characters1, Characters2, Characters3},
					{Dots1, Dots2, Dots3},
					{Characters1, Characters2, Characters3},
				},
			},
			want: true,
		},
		{
			name: "case 4 true5",
			fields: fields{
				XYZs: [][3]TileType{
					{Characters1, Characters2, Characters3},
					{Dots1, Dots2, Dots3},
					{Characters1, Characters2, Characters3},
					{Dots1, Dots2, Dots3},
				},
			},
			want: true,
		},
		{
			name: "case 4 true6",
			fields: fields{
				XYZs: [][3]TileType{
					{Characters1, Characters2, Characters3},
					{Characters1, Characters2, Characters3},
					{Dots1, Dots2, Dots3},
					{Dots1, Dots2, Dots3},
				},
			},
			want: true,
		},
		{
			name: "case 4 true7",
			fields: fields{
				XYZs: [][3]TileType{
					{Dots1, Dots2, Dots3},
					{Characters1, Characters2, Characters3},
					{Characters1, Characters2, Characters3},
					{Dots1, Dots2, Dots3},
				},
			},
			want: true,
		},
		{
			name: "case 4 false1",
			fields: fields{
				XYZs: [][3]TileType{
					{Characters1, Characters2, Characters3},
					{Dots1, Dots2, Dots3},
					{Dots1, Dots2, Dots3},
					{Dots1, Dots2, Dots3},
				},
			},
			want: false,
		},
		{
			name: "case 4 false2",
			fields: fields{
				XYZs: [][3]TileType{
					{Dots1, Dots2, Dots3},
					{Characters1, Characters2, Characters3},
					{Dots1, Dots2, Dots3},
					{Dots1, Dots2, Dots3},
				},
			},
			want: false,
		},
		{
			name: "case 4 false3",
			fields: fields{
				XYZs: [][3]TileType{
					{Dots1, Dots2, Dots3},
					{Dots1, Dots2, Dots3},
					{Characters1, Characters2, Characters3},
					{Dots1, Dots2, Dots3},
				},
			},
			want: false,
		},
		{
			name: "case 4 false4",
			fields: fields{
				XYZs: [][3]TileType{
					{Dots1, Dots2, Dots3},
					{Dots1, Dots2, Dots3},
					{Dots1, Dots2, Dots3},
					{Characters1, Characters2, Characters3},
				},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				hand := &WinningHandNormal{
					WinningHandBase: tt.fields.WinningHandBase,
					XX:              tt.fields.XX,
					XXXs:            tt.fields.XXXs,
					XYZs:            tt.fields.XYZs,
				}
				if got := hand.twoDoubleRuns(); got != tt.want {
					t.Errorf("twoDoubleRuns() = %v, want %v", got, tt.want)
				}
			},
		)
	}
}

func TestWinningHandNormal_fullStraight(t *testing.T) {
	type fields struct {
		WinningHandBase WinningHandBase
		XX              [2]TileType
		XXXs            [][3]TileType
		XYZs            [][3]TileType
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name: "",
			fields: fields{
				XYZs: [][3]TileType{
					{Dots1, Dots2, Dots3},
					{Dots7, Dots8, Dots9},
					{Dots4, Dots5, Dots6},
				},
			},
			want: true,
		},
		{
			name: "",
			fields: fields{
				XYZs: [][3]TileType{
					{Dots1, Dots2, Dots3},
					{Dots7, Dots8, Dots9},
					{Dots4, Dots5, Dots6},
					{Dots4, Dots5, Dots6},
				},
			},
			want: true,
		},
		{
			name: "",
			fields: fields{
				XYZs: [][3]TileType{
					{Characters1, Characters2, Characters3},
					{Dots7, Dots8, Dots9},
					{Dots4, Dots5, Dots6},
				},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				hand := &WinningHandNormal{
					WinningHandBase: tt.fields.WinningHandBase,
					XX:              tt.fields.XX,
					XXXs:            tt.fields.XXXs,
					XYZs:            tt.fields.XYZs,
				}
				if got := hand.fullStraight(); got != tt.want {
					t.Errorf("fullStraight() = %v, want %v", got, tt.want)
				}
			},
		)
	}
}

func TestScoreSrc_ChildTsumo(t *testing.T) {
	tests := []struct {
		name  string
		score ScoreSrc
		want  int
		want1 int
	}{
		{
			"1",
			NewScore(40, 2),
			700,
			1300,
		},
		{
			"2",
			NewScore(80, 1),
			700,
			1300,
		},
		{
			"3",
			NewScore(50, 3),
			1600,
			3200,
		},
		{
			"4",
			NewScore(60, 3),
			2000,
			4000,
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				got, got1 := tt.score.ChildTsumo()
				if got != tt.want {
					t.Errorf("ChildTsumo() got = %v, want %v", got, tt.want)
				}
				if got1 != tt.want1 {
					t.Errorf("ChildTsumo() got1 = %v, want %v", got1, tt.want1)
				}
			},
		)
	}
}

func TestFindSequential(t *testing.T) {
	tts := []TileType{Dots1, Dots1, Dots1, Dots2, Dots2, Dots2, Dots3, Dots3, Dots3, Dots4, Dots4, Dots4, Dots5, Dots5}
	tiles := toSampleTiles(tts)
	//tiles := toSampleTiles([]TileType{Dots1, Dots2, Dots3, Dots3, Dots4, Dots5})
	r, _ := FindSequential(tiles)

	fmt.Println(r)
}

func TestFindTriplet(t *testing.T) {
	tiles := toSampleTiles([]TileType{Dots3, Dots3, Dots3, Dots4, Dots4, Dots4, Dots5, Dots5, Dots5})
	tri, _ := FindTriplet(tiles)
	for _, triplets := range tri {
		fmt.Println(triplets)
	}
	print(len(tri))
}

func TestWinningHandBase_normalWin(t *testing.T) {
	//tts := SortTileTypes([]TileType{
	//	Characters2, Characters2, Characters3,
	//	Characters4, Characters5, Dots3,
	//	Dots4, Dots5, Bamboo3,
	//	Bamboo4, Bamboo5, West, West, Characters2})
	tts := []TileType{Dots1, Dots1, Dots1, Dots2, Dots2, Dots2, Dots3, Dots3, Dots3, Dots4, Dots4, Dots4, Dots5, Dots5}
	tiles := toSampleTiles(tts)
	base := WinningHandBase{Player: &Player{Tiles: tiles}}

	r := base.normalWin()
	for _, normal := range r {
		for _, triplet := range normal.Triplets {
			for _, t := range triplet.TilesXXX {
				fmt.Print(TilesName[t.TileType], " ")
			}
			fmt.Print(", ")
		}
		for _, s := range normal.Sequential {
			for _, t := range s.TilesXYZ {
				fmt.Print(TilesName[t.TileType], " ")
			}
			fmt.Print(", ")
		}
		for _, s := range normal.Head {
			fmt.Print(TilesName[s.TileType], " ")
		}
		fmt.Println()
	}
}

func BenchmarkWinningHandBase_normalWin(b *testing.B) {
	tts := []TileType{Dots1, Dots1, Dots1, Dots2, Dots2, Dots2, Dots3, Dots3, Dots3, Dots4, Dots4, Dots4, Dots5, Dots5}
	tiles := toSampleTiles(tts)
	base := WinningHandBase{Player: &Player{Tiles: tiles}}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		base.normalWin()
	}
}
