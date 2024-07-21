package main

import (
	"fmt"
	"math/rand"
	"time"
	"tower-defense/internal/core"
	"tower-defense/internal/rendering"
)

const (
	gameWidth     = 800
	gameHeight    = 600
	frameDuration = time.Second / 60 // 60 FPS
)

func main() {
	rand.Seed(time.Now().UnixNano())
	gameState := core.NewGameState()
	renderer := rendering.NewRenderer()

	// Set up initial game elements
	setupGame(gameState)

	// Game loop
	ticker := time.NewTicker(frameDuration)
	defer ticker.Stop()

	for !gameState.IsGameOver() {
		select {
		case <-ticker.C:
			handleInput(gameState)
			gameState.Update()
			renderer.Render(gameState)
		}
	}

	fmt.Printf("Game Over! You survived %d waves and earned %d money.\n", gameState.GetWave(), gameState.GetMoney())
}

func setupGame(gs *core.GameState) {
	// Add some initial towers
	gs.AddTower(core.BasicTower, 210, 300)
	gs.AddTower(core.SniperTower, 300, 300)
	// Start the first wave
	gs.NextWave()
}

func handleInput(gs *core.GameState) {
	// This is a placeholder for handling user input
	// In a real game, you'd handle keyboard/mouse events here
	// For now, we'll just randomly upgrade or sell a tower occasionally
	if rand.Float32() < 0.01 { // 1% chance each frame
		if len(gs.GetTowers()) > 0 {
			towerIndex := rand.Intn(len(gs.GetTowers()))
			if rand.Float32() < 0.5 {
				gs.UpgradeTower(towerIndex)
			} else {
				gs.SellTower(towerIndex)
			}
		}
	}

	// Randomly toggle pause
	if rand.Float32() < 0.001 { // 0.1% chance each frame
		gs.TogglePause()
	}
}
