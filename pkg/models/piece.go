package models

import (
	"fmt"
	"math"
)

type Piece struct {
	Name             string
	Color            string
	Health           int
	BaseDamage       int
	PreviousPosition [2]int
	CurrentPosition  [2]int
}

func (p Piece) IsOurPathIsRight() bool {
	if p.Name == "pawn" {
		return p.PreviousPosition[0] == p.CurrentPosition[0] && (math.Abs(float64(p.PreviousPosition[1]-p.CurrentPosition[1])) == 1 || (p.PreviousPosition[1] == 1 && math.Abs(float64(p.PreviousPosition[1]-p.CurrentPosition[1])) == 2))
	}
	if p.Name == "rook" {
		return p.PreviousPosition[0] == p.CurrentPosition[0] || p.PreviousPosition[1] == p.CurrentPosition[1]
	}
	if p.Name == "knight" {
		return (math.Abs(float64(p.CurrentPosition[0]-p.PreviousPosition[0])) == 2 && math.Abs(float64(p.CurrentPosition[1]-p.PreviousPosition[1])) == 1) || (math.Abs(float64(p.CurrentPosition[1]-p.PreviousPosition[1])) == 2 && math.Abs(float64(p.CurrentPosition[0]-p.PreviousPosition[0])) == 1)
	}
	if p.Name == "bishop" {
		return math.Abs(float64(p.CurrentPosition[0]-p.PreviousPosition[0])) == math.Abs(float64(p.CurrentPosition[1]-p.PreviousPosition[1]))
	}
	if p.Name == "queen" {
		return true
	}
	if p.Name == "king" {
		return (math.Abs(float64(p.CurrentPosition[0]-p.PreviousPosition[0])) == 1 && p.CurrentPosition[1] == p.PreviousPosition[1]) || (math.Abs(float64(p.CurrentPosition[1]-p.PreviousPosition[1])) == 1 && p.CurrentPosition[0] == p.PreviousPosition[0]) || (math.Abs(float64(p.CurrentPosition[1]-p.PreviousPosition[1])) == 1 && math.Abs(float64(p.CurrentPosition[0]-p.PreviousPosition[0])) == 1)
	}

	return false
}

func (p *Piece) Move(x int, y int) {
	var OldPrevious = p.PreviousPosition
	p.PreviousPosition = p.CurrentPosition
	p.CurrentPosition[0] = x
	p.CurrentPosition[1] = y
	if !p.IsOurPathIsRight() {
		p.CurrentPosition = p.PreviousPosition
		p.PreviousPosition = OldPrevious
		fmt.Println("Wrong move")
	}
}
