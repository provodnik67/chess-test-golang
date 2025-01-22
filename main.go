package main

import (
	"chess/pkg/models"
	"errors"
	"fmt"
	"math"

	"github.com/fogleman/gg"
)

const (
	White  = "w"
	Black  = "b"
	Pawn   = "pawn"
	Rook   = "rook"
	Knight = "knight"
	Bishop = "bishop"
	Queen  = "queen"
	King   = "king"
)

func main() {
	field := models.Field{Size: [2]int{8, 8}}
	for i := 0; i <= 7; i++ {
		field.Pieces = append(field.Pieces, models.Piece{Name: Pawn, CurrentPosition: [2]int{i, 1}, Color: White})
		field.Pieces = append(field.Pieces, models.Piece{Name: Pawn, CurrentPosition: [2]int{i, 6}, Color: Black})
	}
	field.Pieces = append(field.Pieces, models.Piece{Name: Rook, CurrentPosition: [2]int{0, 0}, Color: White})
	field.Pieces = append(field.Pieces, models.Piece{Name: Rook, CurrentPosition: [2]int{7, 0}, Color: White})
	field.Pieces = append(field.Pieces, models.Piece{Name: Rook, CurrentPosition: [2]int{0, 7}, Color: Black})
	field.Pieces = append(field.Pieces, models.Piece{Name: Rook, CurrentPosition: [2]int{7, 7}, Color: Black})

	field.Pieces = append(field.Pieces, models.Piece{Name: Knight, CurrentPosition: [2]int{1, 0}, Color: White})
	field.Pieces = append(field.Pieces, models.Piece{Name: Knight, CurrentPosition: [2]int{6, 0}, Color: White})
	field.Pieces = append(field.Pieces, models.Piece{Name: Knight, CurrentPosition: [2]int{1, 7}, Color: Black})
	field.Pieces = append(field.Pieces, models.Piece{Name: Knight, CurrentPosition: [2]int{6, 7}, Color: Black})

	field.Pieces = append(field.Pieces, models.Piece{Name: Bishop, CurrentPosition: [2]int{2, 0}, Color: White})
	field.Pieces = append(field.Pieces, models.Piece{Name: Bishop, CurrentPosition: [2]int{5, 0}, Color: White})
	field.Pieces = append(field.Pieces, models.Piece{Name: Bishop, CurrentPosition: [2]int{2, 7}, Color: Black})
	field.Pieces = append(field.Pieces, models.Piece{Name: Bishop, CurrentPosition: [2]int{5, 7}, Color: Black})

	field.Pieces = append(field.Pieces, models.Piece{Name: Queen, CurrentPosition: [2]int{3, 0}, Color: White})
	field.Pieces = append(field.Pieces, models.Piece{Name: Queen, CurrentPosition: [2]int{3, 7}, Color: Black})

	field.Pieces = append(field.Pieces, models.Piece{Name: King, CurrentPosition: [2]int{4, 0}, Color: White})
	field.Pieces = append(field.Pieces, models.Piece{Name: King, CurrentPosition: [2]int{4, 7}, Color: Black})

	_, err := MakeAMove(&field, [2]int{3, 1}, [2]int{3, 3})
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	_, err = MakeAMove(&field, [2]int{3, 6}, [2]int{3, 4})
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	_, err = MakeAMove(&field, [2]int{4, 1}, [2]int{4, 3})
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	//черные атакуют
	_, err = TryToAttack(&field, [2]int{3, 4}, [2]int{4, 3})
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	const S = 1024
	dc := gg.NewContext(S, S)
	dc.SetRGB(1, 1, 1)
	dc.Clear()
	dc.SetRGB(0, 0, 0)
	dc.SetLineWidth(2)
	for i := 0; i < 8; i++ {
		dc.DrawLine(float64(i), (S/8)*(float64(i)+1), 1024, (S/8)*(float64(i)+1))
		dc.DrawLine((S/8)*(float64(i)+1), S, (S/8)*(float64(i)+1), 0)
	}
	dc.Stroke()
	for _, p := range field.Pieces {
		dc.DrawStringAnchored(fmt.Sprintf("%s%s", p.Name, p.Color), (S/8)*(float64(p.CurrentPosition[0])+0.5), (S/8)*(math.Abs(float64(p.CurrentPosition[1])-7)+0.5), 0.5, 0.5)
	}
	dc.SavePNG("out.png")
}

// TODO на последней линии противника пешка становится любой другой фигурой на выбор
// TODO добавить проверку выхода за игровое поле
func MakeAMove(field *models.Field, from [2]int, to [2]int) (bool, error) {
	pieceToMove := len(field.Pieces)
	for i, p := range field.Pieces {
		if p.CurrentPosition[0] == from[0] && p.CurrentPosition[1] == from[1] {
			pieceToMove = i
		}
		if p.CurrentPosition[0] == to[0] && p.CurrentPosition[1] == to[1] {
			return false, fmt.Errorf("%d - %d the place is occupied", to[0], to[1])
		}
	}
	if pieceToMove == len(field.Pieces) {
		return false, fmt.Errorf("%d - %d the piece doesn't exist", from[0], from[1])
	}
	field.Pieces[pieceToMove].PeaceMove(to[0], to[1])
	return true, nil
}

func TryToAttack(field *models.Field, from [2]int, to [2]int) (bool, error) {
	attackingPiece := len(field.Pieces)
	attackedPiece := attackingPiece
	for i, p := range field.Pieces {
		if p.CurrentPosition[0] == from[0] && p.CurrentPosition[1] == from[1] {
			attackingPiece = i
		}
		if p.CurrentPosition[0] == to[0] && p.CurrentPosition[1] == to[1] {
			attackedPiece = i
		}
	}
	if attackingPiece == len(field.Pieces) {
		return false, errors.New("attacking piece doesn't exist")
	}
	if attackedPiece == len(field.Pieces) {
		return false, errors.New("attacked piece doesn't exist")
	}
	var OldPrevious = field.Pieces[attackingPiece].PreviousPosition
	field.Pieces[attackingPiece].PreviousPosition = field.Pieces[attackingPiece].CurrentPosition
	field.Pieces[attackingPiece].CurrentPosition = to
	if !field.Pieces[attackingPiece].IsOurPathIsRight(true) {
		field.Pieces[attackingPiece].CurrentPosition = field.Pieces[attackingPiece].PreviousPosition
		field.Pieces[attackingPiece].PreviousPosition = OldPrevious
		return false, fmt.Errorf("%s%s %d%d wrong move", field.Pieces[attackingPiece].Name, field.Pieces[attackingPiece].Color, to[0], to[1])
	}
	field.Pieces = append(field.Pieces[:attackedPiece], field.Pieces[attackedPiece+1:]...)
	return true, nil
}
