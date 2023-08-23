package main

import (
	"fmt"
	"math"
	"math/rand"
	"syscall/js"
	"time"
)

type Direction string
type GameStatus string

const (
	CanvasWidth  = 800
	CanvasHeight = 800
	RectSize     = 20

	DirectionUp    Direction = "up"
	DirectionDown  Direction = "down"
	DirectionLeft  Direction = "left"
	DirectionRight Direction = "right"

	GameStatusStopped GameStatus = "stopped"
	GameStatusRunning GameStatus = "running"
	GameStatusOver    GameStatus = "over"
)

type SnakeBody struct {
	PosX int
	PosY int
}

var (
	ctx            js.Value
	document       js.Value
	body           js.Value
	canvas         js.Value
	snakeBody      = []SnakeBody{}
	foodPosX       = 0
	foodPosY       = 0
	snakeDirection Direction
	gameStatus     GameStatus
)

func Game(this js.Value, args []js.Value) interface{} {
	document = js.Global().Get("document")
	body = document.Get("body")
	canvas = document.Call("getElementById", "canvas")
	ctx = canvas.Call("getContext", "2d")

	canvas.Set("width", js.ValueOf(CanvasWidth))
	canvas.Set("height", js.ValueOf(CanvasHeight))
	canvas.Get("style").Set("border", "1px solid black")

	button := document.Call("createElement", "button")
	button.Set("innerText", "Start Game!")
	button.Call("addEventListener", "click", js.FuncOf(Run))
	body.Call("appendChild", button)

	snakeHeadPosX := rand.Intn(CanvasWidth)
	snakeHeadPosY := rand.Intn(CanvasWidth)
	snakeBody = append(snakeBody, SnakeBody{PosX: snakeHeadPosX, PosY: snakeHeadPosY})
	snakeDirection = DirectionRight

	foodPosX = rand.Intn(CanvasWidth)
	foodPosY = rand.Intn(CanvasWidth)
	gameStatus = GameStatusStopped
	fmt.Println("GameStatus: ", gameStatus)

	registerCallbacks()

	go func() {
		for {

			// clear canvas
			ctx.Call("clearRect", 0, 0, CanvasWidth, CanvasHeight)

			// render the snake
			for i := 0; i < len(snakeBody); i++ {
				if i == 0 {
					drawSquare(snakeBody[i].PosX, snakeBody[i].PosY, "blue")
				} else {
					drawSquare(snakeBody[i].PosX, snakeBody[i].PosY, "purple")
				}
			}

			if gameStatus == GameStatusRunning {
				for i := len(snakeBody) - 1; i >= 0; i-- {
					if i == 0 {
						switch snakeDirection {
						case DirectionUp:
							snakeBody[i].PosY -= RectSize
						case DirectionDown:
							snakeBody[i].PosY += RectSize
						case DirectionLeft:
							snakeBody[i].PosX -= RectSize
						case DirectionRight:
							snakeBody[i].PosX += RectSize
						}
					} else {
						snakeBody[i].PosX = snakeBody[i-1].PosX
						snakeBody[i].PosY = snakeBody[i-1].PosY
					}
				}
			}

			// draw food
			drawCircle(foodPosX, foodPosY)

			// check if snake is out of canvas
			if snakeBody[0].PosX > CanvasWidth+RectSize || snakeBody[0].PosX < 0-RectSize || snakeBody[0].PosY > CanvasHeight+RectSize || snakeBody[0].PosY < 0-RectSize {
				fmt.Println("Game Over!")
				gameStatus = GameStatusOver
				break
			}

			// check if snake eat itself
			for i := 1; i < len(snakeBody); i++ {
				if snakeBody[0].PosX == snakeBody[i].PosX && snakeBody[0].PosY == snakeBody[i].PosY {
					fmt.Println("Game Over!")
					gameStatus = GameStatusOver
					break
				}
			}

			// check position of snake head and food with tolerance
			if (snakeBody[0].PosX >= foodPosX-20 && snakeBody[0].PosX <= foodPosX+20) && (snakeBody[0].PosY >= foodPosY-20 && snakeBody[0].PosY <= foodPosY+20) {
				fmt.Println("Snake eat food!")

				snakeBody = append(snakeBody, SnakeBody{PosX: snakeBody[len(snakeBody)-1].PosX - RectSize, PosY: snakeBody[len(snakeBody)-1].PosY - RectSize})
				foodPosX = rand.Intn(CanvasWidth)
				foodPosY = rand.Intn(CanvasWidth)
			}

			time.Sleep(100 * time.Millisecond)
		}

		js.Global().Call("alert", "Game Over")

		// remove button button
		body.Call("removeChild", button)
	}()

	return nil
}

func drawSquare(x, y int, color string) {
	ctx.Set("fillStyle", color)
	ctx.Call("fillRect", x, y, RectSize, RectSize)
}

func drawCircle(x, y int) {
	ctx.Set("fillStyle", "green")
	ctx.Call("beginPath")
	ctx.Call("arc", x, y, 10, 0, 2*math.Pi)
	ctx.Call("fill")
}

func registerCallbacks() {
	document.Call("addEventListener", "keydown", js.FuncOf(handleKeydown))
}

func handleKeydown(this js.Value, args []js.Value) interface{} {
	event := args[0]
	key := event.Get("key").String()

	switch key {
	case "ArrowUp":
		if snakeDirection == DirectionDown {
			return nil
		}
		snakeDirection = DirectionUp
	case "ArrowDown":
		if snakeDirection == DirectionUp {
			return nil
		}
		snakeDirection = DirectionDown
	case "ArrowLeft":
		if snakeDirection == DirectionRight {
			return nil
		}
		snakeDirection = DirectionLeft
	case "ArrowRight":
		if snakeDirection == DirectionLeft {
			return nil
		}
		snakeDirection = DirectionRight
	}

	return nil
}

func Run(this js.Value, args []js.Value) interface{} {
	this.Set("disabled", true)
	gameStatus = GameStatusRunning
	return nil
}

func HelloWorld(this js.Value, args []js.Value) interface{} {
	fmt.Println("Hello Golang WebAssembly! from HelloWorld func")
	return "Hello Golang WebAssembly! from return"
}

func main() {
	done := make(chan struct{})
	fmt.Println("Hallo from main func Golang WebAssembly!")

	funHelloWorld := js.FuncOf(HelloWorld)
	js.Global().Set("HelloWorld", funHelloWorld)
	defer funHelloWorld.Release()

	funGame := js.FuncOf(Game)
	js.Global().Set("Game", funGame)
	defer funGame.Release()

	<-done
}
