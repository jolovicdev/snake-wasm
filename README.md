# Snake Game WebAssembly

A classic Snake game implemented in Go and compiled to WebAssembly, featuring smooth gameplay.

## Prerequisites

- Go 1.22 or higher
- Web browser with WebAssembly support

## Quick Start

Clone and run:
```sh
git clone <repository-url>
cd snake-wasm
make run
```
Open [http://localhost:8080](http://localhost:8080) in your browser.

## Controls

- Arrow keys to move snake
- 'P' key or Pause button to pause/resume
- New Game button to restart
- Avoid hitting walls or yourself
- Eat food (red squares) to grow and score points

## Features

- Score tracking
- Pause functionality
- Smooth controls
- New Game option

## Commands

- `make run` - Build and start server
- `make build` - Only build WASM
- `make clean` - Clean generated files