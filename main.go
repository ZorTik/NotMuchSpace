package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
)

const (
	Player       EntityType = "player"
	CommonEntity EntityType = "common_entity"
)

const (
	Up    Direction = "w"
	Left  Direction = "a"
	Down  Direction = "s"
	Right Direction = "d"
)

type Game struct {
	field  []Entity
	bounds Bounds
}

type Bounds struct {
	width, height int
}

type Direction string
type EntityType string

type Entity struct {
	id         string
	entityType EntityType
	x, y       int
}

func (game *Game) forOf(cons func(entity *Entity, index int)) {
	for i, entity := range game.field {
		cons(&entity, i)
	}
}

func (game *Game) addEntity(entity *Entity) {
	fmt.Println(game.field)
	game.field = append(game.field, *entity)
	fmt.Println(game.field)
}

func (game *Game) getEntityAt(x, y int) (*Entity, int) {
	var result *Entity = nil
	var index = -1

	game.forOf(func(entity *Entity, _index int) {
		if entity.x == x && entity.y == y {
			result = entity
			index = _index
		}
	})

	return result, index
}

func (game *Game) generateEntity(_type EntityType) {
	var x int
	var y int
	for {
		_x := rand.Intn(game.bounds.width)
		_y := rand.Intn(game.bounds.height)

		if entity, _ := game.getEntityAt(_x, _y); entity == nil {
			x = _x
			y = _y
			break
		}
	}

	game.addEntity(&Entity{
		id:         string(rune(rand.Int())),
		entityType: _type,
		x:          x,
		y:          y,
	})
}

func (game *Game) getPlayer() (*Entity, int) {
	var result *Entity = nil
	var index = -1

	game.forOf(func(entity *Entity, _index int) {
		if entity.entityType == Player {
			result = entity
			index = _index
		}
	})

	return result, index
}

func (game *Game) move(eIndex int, direction Direction) bool {
	entity := game.field[eIndex]

	dx, dy := DirectionToXY(direction)
	if nbh, _ := game.getEntityAt(entity.x+dx, entity.y+dy); nbh != nil {
		// There is already an entity at the target position
		return false
	}

	switch direction {
	case Up:
		entity.y--
		break
	case Left:
		entity.x--
		break
	case Down:
		entity.y++
		break
	case Right:
		entity.x++
		break
	}
	game.field[eIndex] = entity

	return true
}

func (game *Game) moveAI(index int) {
	entity := game.field[index]
	x := entity.x
	y := entity.y

	moves := []Direction{
		Up, Left, Down, Right,
	}

	moves = *Shuffle(moves)

	for _index, move := range moves {
		_xC, _yC := DirectionToXY(move)
		e, _ := game.getEntityAt(x+_xC, y+_yC)
		if e == nil {
			game.move(index, move)
		} else if e != nil && _index == len(moves)-1 && entity.entityType == Player {
			fmt.Println("No more moves. Game over.")
			os.Exit(0)
		}
	}
}

func (game *Game) handleInput(reader *bufio.Reader) {
	for {
		input, _, _ := reader.ReadLine()

		_, index := game.getPlayer()

		if game.move(index, Direction(input)) {
			break
		}

		fmt.Println("You can't move here.")
	}
	game.forOf(func(entity *Entity, index int) {
		if entity.entityType != Player {
			game.moveAI(index)
		}
	})
	game.generateEntity(CommonEntity)
}

func (game *Game) render() {
	tempField := map[[2]int]*Entity{}

	game.forOf(func(e *Entity, _ int) {
		coords := [2]int{e.x, e.y}
		tempField[coords] = e
	})

	for y := 0; y < game.bounds.height; y++ {
		for x := 0; x < game.bounds.width; x++ {
			coords := [2]int{x, y}

			mark := "#"
			if entity, has := tempField[coords]; has {
				switch entity.entityType {
				case Player:
					mark = "O"
					break
				case CommonEntity:
					mark = "X"
					break
				}
				fmt.Println(entity.entityType)
			}

			print(mark + " ")
		}
		print("\n")
	}
}

func loop(game *Game) {
	reader := bufio.NewReader(os.Stdin)
	for {
		game.render()
		game.handleInput(reader)

		if player, _ := game.getPlayer(); OutOfBounds(player, &game.bounds) {
			fmt.Println("Out of bounds. Game over.")
			break
		}
	}
}

func prepare(bounds Bounds) Game {
	game := Game{bounds: bounds}
	game.generateEntity(Player)

	return game
}

func main() {
	game := prepare(Bounds{
		width:  10,
		height: 10,
	})
	loop(&game)
}
