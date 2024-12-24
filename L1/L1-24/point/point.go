package point

import (
	"math"
)

type Point struct {
	x, y float64
}

func NewPoint(x, y float64) *Point {
	return &Point{x: x, y: y}
}

func (a *Point) GetDistance(b *Point) float64 {
	return math.Sqrt(math.Pow((b.GetX()-a.GetX()), 2) + math.Pow((b.GetY()-a.GetY()), 2))
}

func (p *Point) GetX() float64 {
	return p.x
}

func (p *Point) GetY() float64 {
	return p.y
}
