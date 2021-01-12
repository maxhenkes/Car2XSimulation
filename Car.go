package main

import (
	"fmt"
	"github.com/faiface/pixel"
	"math"
	"math/rand"
)

type Move struct {
	x, y int
}

type Car struct {
	x, y             int
	id               string
	direction        Move
	sprite           *pixel.Sprite
	status           string
	ReceivedMessages []Message
}

var CarSprites = []*pixel.Sprite{LoadAndSprite("assets/car1.png"), LoadAndSprite("assets/car2.png"), LoadAndSprite("assets/car3.png"), LoadAndSprite("assets/car4.png")}

func (car Car) RenderCar() {

	mat := pixel.IM
	mat = mat.Moved(pixel.V(streetMap.tiles[car.x][car.y].x, streetMap.tiles[car.x][car.y].y))
	mat = mat.Rotated(pixel.V(streetMap.tiles[car.x][car.y].x, streetMap.tiles[car.x][car.y].y), rotateDirection(car.direction))

	mat = moveToLane(mat, car)

	car.sprite.Draw(mainWindow, mat)
}

func moveToLane(mat pixel.Matrix, car Car) pixel.Matrix {

	distance := 8.0

	switch car.direction {
	case UP:
		return mat.Moved(pixel.V(distance, 0))
	case DOWN:
		return mat.Moved(pixel.V(-distance, 0))
	case LEFT:
		return mat.Moved(pixel.V(0, distance))
	case RIGHT:
		return mat.Moved(pixel.V(0, -distance))
	default:
		return mat
	}

}

func (car *Car) MoveCar() {

	var movePool []Move

	switch {
	case compareDir(UP, car.direction):
		for f := range dir[0] {
			if isInside(*car, f, dir[0]) {
				if streetMap.tiles[car.x+dir[0][f].x][car.y+dir[0][f].y].tileType > 0 {
					movePool = append(movePool, dir[0][f])
				}
			}
		}
	case compareDir(DOWN, car.direction):
		for f := range dir[1] {
			if isInside(*car, f, dir[1]) {
				if streetMap.tiles[car.x+dir[1][f].x][car.y+dir[1][f].y].tileType > 0 {
					movePool = append(movePool, dir[1][f])
				}
			}
		}
	case compareDir(LEFT, car.direction):
		for f := range dir[2] {
			if isInside(*car, f, dir[2]) {
				if streetMap.tiles[car.x+dir[2][f].x][car.y+dir[2][f].y].tileType > 0 {
					movePool = append(movePool, dir[2][f])
				}
			}
		}
	case compareDir(RIGHT, car.direction):
		for f := range dir[3] {
			if isInside(*car, f, dir[3]) {
				if streetMap.tiles[car.x+dir[3][f].x][car.y+dir[3][f].y].tileType > 0 {
					movePool = append(movePool, dir[3][f])
				}
			}
		}
	}
	if len(movePool) > 0 {
		i := rand.Intn(len(movePool))

		if !isOccupied(movePool[i], car) {
			car.addDir(movePool[i])
		}
	} else {
		switch {
		case compareDir(car.direction, UP):
			car.addDir(DOWN)
		case compareDir(car.direction, DOWN):
			car.addDir(UP)
		case compareDir(car.direction, LEFT):
			car.addDir(RIGHT)
		case compareDir(car.direction, RIGHT):
			car.addDir(LEFT)
		}
	}

}

func isOccupied(move Move, car *Car) bool {
	x := car.x + move.x
	y := car.y + move.y

	for _, c := range streetMap.cars {
		if c.x == x && c.y == y {
			if car.direction == c.direction {
				return true
			}
		}
	}
	return false
}

func (car *Car) addDir(m Move) {

	car.x = car.x + m.x
	car.y = car.y + m.y

	if streetMap.tiles[car.x][car.y].obstacle {
		message := NewActiveMessage(*car, *NewMessage(*car))
		AddMessage(*message)
		car.receiveMessage(message.message)
	}

	car.direction = m
}

func (car *Car) update() {
	car.checkMessages()
	car.MoveCar()

}

func (car *Car) receiveMessage(message Message) {

	fmt.Println("rec")

	for _, f := range car.ReceivedMessages {
		if message.isEqual(f) {
			fmt.Println("was equal")
			return
		}
	}
	car.ReceivedMessages = append(car.ReceivedMessages, message)
	fmt.Println(car.ReceivedMessages)
	car.checkMessages()
}

func (car *Car) checkMessages() {

	cleanedMessages := make([]Message, 0)

	for _, f := range car.ReceivedMessages {
		if f.timeCounter > 0 {
			f.timeCounter--
			cleanedMessages = append(cleanedMessages, f)
		}
	}

	for _, f := range cleanedMessages {

		if f.hopCounter > 0 {
			newMessage := f
			newMessage.hopCounter--
			var newActiveMessage = NewActiveMessage(*car, newMessage)
			fmt.Println("sent new mesg")
			futureInstance = append(futureInstance, newActiveMessage)
			fmt.Println(futureInstance)
		}
	}
	car.ReceivedMessages = cleanedMessages
}

func isInside(car Car, f int, mov []Move) bool {
	return car.x+mov[f].x <= streetMap.size-1 && car.y+mov[f].y <= streetMap.size-1 && car.x+mov[f].x >= 0 && car.y+mov[f].y >= 0
}

func rotateDirection(m Move) float64 {

	var deg float64

	switch {
	case compareDir(m, UP):
		deg = 0
	case compareDir(m, DOWN):
		deg = 180
	case compareDir(m, LEFT):
		deg = 90
	case compareDir(m, RIGHT):
		deg = 270
	default:
		deg = 0
	}

	return deg * (math.Pi / 180)
}
