package main

import (
	"fmt"
	"syscall/js"
	"time"
)

const (
	colorRed     = "\033[31m"
	colorGreen   = "\033[32m"
	colorYellow  = "\033[33m"
	colorBlue    = "\033[34m"
	colorMagenta = "\033[35m"
	colorCyan    = "\033[36m"
	colorReset   = "\033[0m"
)

type Logger struct {
	enabled bool
}

func NewLogger() *Logger {
	return &Logger{enabled: true}
}

func (l *Logger) Info(format string, args ...interface{}) {
	if l.enabled {
		timestamp := time.Now().Format("15:04:05")
		fmt.Printf(colorBlue+"[INFO] %s: "+format+colorReset+"\n", append([]interface{}{timestamp}, args...)...)
	}
}

func (l *Logger) Success(format string, args ...interface{}) {
	if l.enabled {
		timestamp := time.Now().Format("15:04:05")
		fmt.Printf(colorGreen+"[SUCCESS] %s: "+format+colorReset+"\n", append([]interface{}{timestamp}, args...)...)
	}
}

func (l *Logger) Warning(format string, args ...interface{}) {
	if l.enabled {
		timestamp := time.Now().Format("15:04:05")
		fmt.Printf(colorYellow+"[WARNING] %s: "+format+colorReset+"\n", append([]interface{}{timestamp}, args...)...)
	}
}

func (l *Logger) Error(format string, args ...interface{}) {
	if l.enabled {
		timestamp := time.Now().Format("15:04:05")
		fmt.Printf(colorRed+"[ERROR] %s: "+format+colorReset+"\n", append([]interface{}{timestamp}, args...)...)
	}
}

type Point struct {
	X, Y int
}

type Snake struct {
	segments  []Point
	direction Point
	growing   bool
}

type Game struct {
	snake    Snake
	food     Point
	score    int
	gameOver bool
	isPaused bool
	gridSize int
	width    int
	height   int
	ctx      js.Value
	document js.Value
	logger   *Logger
}

func NewGame() *Game {
	logger := NewLogger()
	logger.Info("Initializing new game...")

	doc := js.Global().Get("document")
	canvas := doc.Call("getElementById", "game-canvas")
	ctx := canvas.Call("getContext", "2d")

	width := canvas.Get("width").Int()
	height := canvas.Get("height").Int()
	gridSize := 20

	logger.Success("Canvas initialized: %dx%d pixels, grid size: %d", width, height, gridSize)

	game := &Game{
		gridSize: gridSize,
		width:    width,
		height:   height,
		ctx:      ctx,
		document: doc,
		logger:   logger,
	}

	game.initializeGame()
	return game
}

func (g *Game) initializeGame() {
	g.logger.Info("Starting new game...")

	startX := (g.width / g.gridSize / 2) * g.gridSize
	startY := (g.height / g.gridSize / 2) * g.gridSize

	g.snake = Snake{
		segments: []Point{{startX, startY}},
		direction: Point{
			X: g.gridSize,
			Y: 0,
		},
	}

	g.score = 0
	g.gameOver = false
	g.isPaused = false
	g.updateScore()
	g.spawnFood()
	g.updatePauseButton()

	g.logger.Success("Game initialized: Snake positioned at (%d, %d)", startX, startY)
}

func (g *Game) updatePauseButton() {
	pauseBtn := g.document.Call("getElementById", "pause-btn")
	if g.isPaused {
		pauseBtn.Set("innerText", "Resume")
		pauseBtn.Get("classList").Call("add", "paused")
		g.logger.Info("Game paused")
	} else {
		pauseBtn.Set("innerText", "Pause")
		pauseBtn.Get("classList").Call("remove", "paused")
		g.logger.Info("Game resumed")
	}
}

func (g *Game) spawnFood() {
	maxX := g.width / g.gridSize
	maxY := g.height / g.gridSize

	x := js.Global().Get("Math").Call("floor", js.Global().Get("Math").Call("random").Float()*float64(maxX)).Int() * g.gridSize
	y := js.Global().Get("Math").Call("floor", js.Global().Get("Math").Call("random").Float()*float64(maxY)).Int() * g.gridSize

	g.food = Point{x, y}
	g.logger.Info("Food spawned at (%d, %d)", x, y)
}

func (g *Game) update() {
	if g.gameOver || g.isPaused {
		return
	}

	newHead := Point{
		X: g.snake.segments[0].X + g.snake.direction.X,
		Y: g.snake.segments[0].Y + g.snake.direction.Y,
	}

	if newHead.X < 0 || newHead.X >= g.width || newHead.Y < 0 || newHead.Y >= g.height {
		g.gameOver = true
		g.logger.Error("Game Over: Wall collision at (%d, %d)", newHead.X, newHead.Y)
		return
	}

	for _, segment := range g.snake.segments {
		if newHead.X == segment.X && newHead.Y == segment.Y {
			g.gameOver = true
			g.logger.Error("Game Over: Self collision at (%d, %d)", newHead.X, newHead.Y)
			return
		}
	}

	g.snake.segments = append([]Point{newHead}, g.snake.segments...)
	if !g.snake.growing {
		g.snake.segments = g.snake.segments[:len(g.snake.segments)-1]
	}
	g.snake.growing = false

	if newHead.X == g.food.X && newHead.Y == g.food.Y {
		g.score++
		g.snake.growing = true
		g.logger.Success("Food eaten! Score: %d, Snake length: %d", g.score, len(g.snake.segments))
		g.spawnFood()
		g.updateScore()
	}
}

func (g *Game) draw() {
	g.ctx.Set("fillStyle", "#f0f0f0")
	g.ctx.Call("fillRect", 0, 0, g.width, g.height)

	g.ctx.Set("fillStyle", "#4CAF50")
	for _, segment := range g.snake.segments {
		g.ctx.Call("fillRect", segment.X, segment.Y, g.gridSize-1, g.gridSize-1)
	}

	g.ctx.Set("fillStyle", "#FF5722")
	g.ctx.Call("fillRect", g.food.X, g.food.Y, g.gridSize-1, g.gridSize-1)

	if g.gameOver {
		g.ctx.Set("fillStyle", "rgba(0, 0, 0, 0.5)")
		g.ctx.Call("fillRect", 0, 0, g.width, g.height)
		g.ctx.Set("fillStyle", "white")
		g.ctx.Set("font", "30px Arial")
		g.ctx.Call("fillText", "Game Over!", g.width/2-70, g.height/2)
	} else if g.isPaused {
		g.ctx.Set("fillStyle", "rgba(0, 0, 0, 0.5)")
		g.ctx.Call("fillRect", 0, 0, g.width, g.height)
		g.ctx.Set("fillStyle", "white")
		g.ctx.Set("font", "30px Arial")
		g.ctx.Call("fillText", "Paused", g.width/2-50, g.height/2)
	}
}

func (g *Game) updateScore() {
	g.document.Call("getElementById", "score-value").Set("innerText", g.score)
	g.logger.Info("Score updated: %d", g.score)
}

func (g *Game) handleKeydown(this js.Value, args []js.Value) interface{} {
	if g.gameOver {
		return nil
	}

	event := args[0]
	key := event.Get("key").String()

	if key == "p" || key == "P" {
		g.togglePause()
		return nil
	}

	if g.isPaused {
		return nil
	}

	currentX := g.snake.direction.X
	currentY := g.snake.direction.Y

	switch key {
	case "ArrowUp":
		if currentY == 0 {
			g.snake.direction = Point{0, -g.gridSize}
			g.logger.Info("Direction changed: UP")
		}
	case "ArrowDown":
		if currentY == 0 {
			g.snake.direction = Point{0, g.gridSize}
			g.logger.Info("Direction changed: DOWN")
		}
	case "ArrowLeft":
		if currentX == 0 {
			g.snake.direction = Point{-g.gridSize, 0}
			g.logger.Info("Direction changed: LEFT")
		}
	case "ArrowRight":
		if currentX == 0 {
			g.snake.direction = Point{g.gridSize, 0}
			g.logger.Info("Direction changed: RIGHT")
		}
	}

	return nil
}

func (g *Game) togglePause() {
	g.isPaused = !g.isPaused
	g.updatePauseButton()
}

func (g *Game) handlePauseClick(this js.Value, args []js.Value) interface{} {
	if !g.gameOver {
		g.togglePause()
	}
	return nil
}

func (g *Game) handleNewGameClick(this js.Value, args []js.Value) interface{} {
	g.logger.Info("Starting new game...")
	g.initializeGame()
	return nil
}

func (g *Game) gameLoop(this js.Value, args []js.Value) interface{} {
	g.update()
	g.draw()
	js.Global().Call("setTimeout", js.FuncOf(g.gameLoop), 100)
	return nil
}

func main() {
	logger := NewLogger()

	fmt.Printf("\n%s╔════════════════════════════════════╗%s\n", colorCyan, colorReset)
	fmt.Printf("%s║      WASM Snake Game Started      ║%s\n", colorCyan, colorReset)
	fmt.Printf("%s║ Open http://localhost:8080 to play ║%s\n", colorCyan, colorReset)
	fmt.Printf("%s╚════════════════════════════════════╝%s\n\n", colorCyan, colorReset)

	logger.Info("Starting WASM Snake Game...")
	logger.Info("Controls:")
	logger.Info("- Arrow keys to move")
	logger.Info("- 'P' to pause/resume")
	logger.Info("- Use the New Game button to restart")

	game := NewGame()

	js.Global().Get("document").Call("addEventListener", "keydown", js.FuncOf(game.handleKeydown))
	game.document.Call("getElementById", "pause-btn").Call("addEventListener", "click", js.FuncOf(game.handlePauseClick))
	game.document.Call("getElementById", "new-game-btn").Call("addEventListener", "click", js.FuncOf(game.handleNewGameClick))

	logger.Success("Game initialized and ready to play!")

	js.Global().Call("setTimeout", js.FuncOf(game.gameLoop), 100)

	select {}
}
