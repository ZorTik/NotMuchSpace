package main

import "math/rand"

func Shuffle[n any](arr []n) *[]n {
	for i := len(arr) - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		arr[i], arr[j] = arr[j], arr[i]
	}

	return &arr
}

func EntityOutOfBounds(entity *Entity, bounds *Bounds) bool {
	return OutOfBounds(entity.x, entity.y, bounds)
}

func OutOfBounds(x, y int, bounds *Bounds) bool {
	return x < 0 || x >= bounds.width || y < 0 || y >= bounds.height
}

func DirectionToXY(direction Direction) (int, int) {
	switch direction {
	case Up:
		return 0, -1
	case Down:
		return 0, 1
	case Left:
		return -1, 0
	case Right:
		return 1, 0
	}

	return 0, 0
}

func XYToDirection(x, y int) Direction {
	switch {
	case x == 0 && y == -1:
		return Up
	case x == 0 && y == 1:
		return Down
	case x == -1 && y == 0:
		return Left
	case x == 1 && y == 0:
		return Right
	}

	return ""
}

func CoordsToIndex(x, y, width int) int {
	return x + y*width
}
