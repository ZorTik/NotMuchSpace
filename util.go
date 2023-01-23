package main

import "math/rand"

func Shuffle[n any](arr []n) *[]n {
	for i := len(arr) - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		arr[i], arr[j] = arr[j], arr[i]
	}

	return &arr
}

func OutOfBounds(entity *Entity, bounds *Bounds) bool {
	return entity.x < 0 || entity.x >= bounds.width || entity.y < 0 || entity.y >= bounds.height
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
