package main

import (
	"image/color"
	"math"
)

var wallsCreated uint = 0

// The structure of walls
type Wall struct {
	id     uint
	Start  Vect2   `json:"start"`
	End    Vect2   `json:"end"`
	Radius float64 `json:"radius"`
}

func NewWall(x1, x2, y1, y2, radius float64) *Wall {
	wallsCreated++
	return &Wall{
		id:     wallsCreated,
		Start:  Vect2{x1, y1},
		End:    Vect2{x2, y2},
		Radius: radius,
	}
}

type WallList [](*Wall)

func NewEmptyWallList() WallList {
	return make(WallList, 0)
}

func (wall *Wall) Length() float64 {
	return norm(substract(wall.End, wall.Start))
}

func (wall *Wall) Direction() Vect2 {
	return scalar_times(substract(wall.End, wall.Start), 1/wall.Length())
}

func (wall *Wall) Angle() float64 {
	var angle float64
	if wall.End.X-wall.Start.X != 0 {
		angle = math.Atan((wall.End.Y - wall.Start.Y) / (wall.End.X - wall.Start.X))
	} else {
		if wall.End.X > wall.Start.X {
			angle = -math.Pi / 2
		} else {
			angle = math.Pi / 2
		}
	}
	return angle
}

func instanciateWalls() WallList {
	list := make(WallList, 0)
	list = append(list, NewWall(7.5, 7.5, 0, 14, 0.2))
	list = append(list, NewWall(7.5, 7.5, 16, 30, 0.2))

	// // Digas
	// list = append(list, NewWall(0, 10, 0, 10, 0.2))
	// list = append(list, NewWall(30, 20, 0, 10, 0.2))
	// list = append(list, NewWall(0, 10, 30, 20, 0.2))
	// list = append(list, NewWall(30, 20, 30, 20, 0.2))
	// // Small doors
	// list = append(list, NewWall(10, 14, 10, 10, 0.3))
	// list = append(list, NewWall(20, 16, 10, 10, 0.3))

	// list = append(list, NewWall(10, 14, 20, 20, 0.3))
	// list = append(list, NewWall(20, 16, 20, 20, 0.3))

	// list = append(list, NewWall(10, 10, 10, 14, 0.3))
	// list = append(list, NewWall(10, 10, 20, 16, 0.3))

	// list = append(list, NewWall(20, 20, 10, 14, 0.3))
	// list = append(list, NewWall(20, 20, 20, 16, 0.3))
	return list
}

func (wall *Wall) Draw(cv *Canvas) {
	gc := cv.gc
	ratio := cv.getRatio()
	theta := wall.Angle()
	R := wall.Radius * ratio
	x1 := wall.Start.X * ratio
	y1 := wall.Start.Y * ratio
	x2 := wall.End.X * ratio
	y2 := wall.End.Y * ratio

	// Params
	gc.SetFillColor(color.RGBA{0x00, 0x00, 0x00, 0xff})
	gc.SetLineWidth(0)

	// Draw the wall
	gc.MoveTo(x1+R*math.Sin(theta), y1-R*math.Cos(theta))
	gc.LineTo(x2+R*math.Sin(theta), y2-R*math.Cos(theta))
	gc.ArcTo(x2, y2, R, R, theta-math.Pi/2, math.Pi)
	gc.LineTo(x1-R*math.Sin(theta), y1+R*math.Cos(theta))
	gc.ArcTo(x1, y1, R, R, theta+math.Pi/2, math.Pi)
	gc.Close()
	gc.FillStroke()
}
