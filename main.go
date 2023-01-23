// Author: ZorTik
// Date: 23/01 2023
// Elapsed: <1d
package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
)

const (
	Player       EntityType = "player"
	CommonEntity EntityType = "common_entity"
)

const (
	Up         Direction = "w"
	Left       Direction = "a"
	Down       Direction = "s"
	Right      Direction = "d"
	Directions           = "wasd"
)

type Game struct {
	field          []Entity
	bounds         Bounds
	entityGenSteps int
}

type Bounds struct {
	width, height int
}

type Direction string
type EntityType string
type EntityGenerationFunc func(*Entity) *Entity

var (
	// SetTransparent Sets the entity to be transparent after generated
	SetTransparent EntityGenerationFunc = func(entity *Entity) *Entity {
		entity.transparent = true
		return entity
	}
)

type Entity struct {
	id          string
	entityType  EntityType
	x, y        int
	transparent bool
	handler     EntityHandler
}

type OnCollidedFunc func(*Entity, *Entity)

type EntityHandler struct {
	// onCollided Only if Entity is transparent
	onCollided []OnCollidedFunc // Player, Target
}

var startTime = time.Now()
var moves = 0

func (game *Game) forOf(cons func(entity *Entity, index int)) {
	// Makes Game structure iterable
	for i, entity := range game.field {
		cons(&entity, i)
	}
}

func (game *Game) addEntity(entity *Entity) {
	game.field = append(game.field, *entity)
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

func (game *Game) generateBuff(_type EntityType, collidedFuncs ...OnCollidedFunc) {
	game.generateEntity(_type, SetTransparent, func(entity *Entity) *Entity {
		for _, collidedFunc := range collidedFuncs {
			// When player stands on the buff
			entity.handler.onCollided = append(entity.handler.onCollided, collidedFunc)
		}
		return entity
	})
}

func (game *Game) generateEntity(_type EntityType, chain ...EntityGenerationFunc) {
	var x int
	var y int
	for {
		// generates random position until empty one is found
		_x := rand.Intn(game.bounds.width)
		_y := rand.Intn(game.bounds.height)

		if entity, _ := game.getEntityAt(_x, _y); entity == nil {
			x = _x
			y = _y
			break
		}
	}

	generatedEntity := &Entity{
		id:          string(rune(rand.Int())),
		entityType:  _type,
		x:           x,
		y:           y,
		transparent: false,
	}

	for _, f := range chain {
		generatedEntity = f(generatedEntity)
	}

	game.addEntity(generatedEntity)
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

func (game *Game) checkCanMove(eIndex int) (dx, dy int, canMove bool) {
	entity := game.field[eIndex]

	moves := []Direction{Up, Left, Down, Right}

	moves = *Shuffle(moves)

	for _, move := range moves {
		if &entity == nil {
			continue
		}

		_dx, _dy := DirectionToXY(move)
		if game.canMove(&entity, _dx, _dy) {
			return _dx, _dy, true
		}
	}

	return 0, 0, false
}

func (game *Game) canMove(entity *Entity, dx, dy int) bool {
	if nbh, _ := game.getEntityAt(entity.x+dx, entity.y+dy); nbh != nil && !nbh.transparent {
		// There is already an entity at the target position
		return false
	} else if OutOfBounds(entity.x+dx, entity.y+dy, &game.bounds) {
		// The target position is out of bounds
		return false
	}
	return true
}

func (game *Game) move(eIndex int, direction Direction) bool {
	entity := game.field[eIndex]

	dx, dy := DirectionToXY(direction)
	if !game.canMove(&entity, dx, dy) {
		return false
	}

	entity.x += dx
	entity.y += dy

	if ent, _ := game.getEntityAt(entity.x, entity.y); ent != nil {
		defer (func() {
			// Call onCollided handlers after reassign
			for _, f := range entity.handler.onCollided {
				f(&entity, ent)
			}
		})()
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

	for _, move := range moves {
		_xC, _yC := DirectionToXY(move)
		e, _ := game.getEntityAt(x+_xC, y+_yC)
		if e == nil {
			game.move(index, move)
		}
	}
}

func (game *Game) handleInput(reader *bufio.Reader) {
	for {
		input, _, _ := reader.ReadLine()

		_, index := game.getPlayer()

		if game.move(index, Direction(input)) {
			// Move was successful, breaking attempt loop
			break
		} else if !strings.Contains(Directions, string(input)) {
			// Invalid direction
			fmt.Println("Directions: " + Directions)
			continue
		}

		fmt.Println("You can't move here.")
	}
	game.forOf(func(entity *Entity, index int) {
		if entity.entityType != Player {
			game.moveAI(index)
		}
	})
	for i := 0; i < game.entityGenSteps; i++ {
		game.generateEntity(CommonEntity)
	}
}

func (game *Game) render() {
	tempField := map[int]Entity{}

	game.forOf(func(e *Entity, _ int) {
		tempField[CoordsToIndex(e.x, e.y, game.bounds.width)] = *e
	})

	for y := 0; y < game.bounds.height; y++ {
		for x := 0; x < game.bounds.width; x++ {
			index := CoordsToIndex(x, y, game.bounds.width)

			mark := "#"
			if entity, has := tempField[index]; has {

				typ := entity.entityType

				switch typ {
				case Player:
					mark = "O"
					break
				case CommonEntity:
					mark = "X"
					break
				}
			}

			print(mark + " ")
		}
		print("\n")
	}
}

func loop(game *Game) {
	reader := bufio.NewReader(os.Stdin)
	for {
		moves++

		game.render()
		game.handleInput(reader)

		_, pIndex := game.getPlayer()

		if _, _, canMove := game.checkCanMove(pIndex); !canMove {
			exit("No more moves. Game over.")
			break
		}
	}
}

func prepare(bounds Bounds) Game {
	game := Game{bounds: bounds, entityGenSteps: 2}
	game.generateEntity(Player)

	return game
}

func exit(msg string) {
	fmt.Println(msg)
	fmt.Println("------------------------------------")
	fmt.Println("Moves:", moves)
	fmt.Println("Time:", time.Since(startTime))
	fmt.Println("------------------------------------")
	os.Exit(0)
}

func main() {
	game := prepare(Bounds{
		width:  10,
		height: 10,
	})
	loop(&game)
}
