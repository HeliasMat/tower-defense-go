package rendering

import (
	"fmt"
	"strings"
	"sync"
	"tower-defense/internal/core"
	"tower-defense/internal/entities"
)

const (
	gameWidth      = 100
	gameHeight     = 30
	borderChar     = '█'
	cornerChar     = '█'
	enemyChar      = 'E'
	towerChar      = 'T'
	projectileChar = '•'
	sidebarWidth   = 25
	hudHeight      = 3
)

type Renderer struct {
	mu     sync.Mutex
	buffer [][]string
}

func NewRenderer() *Renderer {
	buffer := make([][]string, gameHeight)
	for i := range buffer {
		buffer[i] = make([]string, gameWidth)
	}
	return &Renderer{buffer: buffer}
}

func (r *Renderer) Render(gs *core.GameState) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.clearBuffer()
	r.drawGameArea(gs)
	r.drawWindow()
	r.drawHUD(gs)
	r.drawSidebar(gs)
	r.display()
}

func (r *Renderer) drawWindow() {
	// Draw vertical borders
	for y := 0; y < gameHeight; y++ {
		r.buffer[y][0] = string(borderChar)
		r.buffer[y][gameWidth-1] = string(borderChar)
		r.buffer[y][gameWidth-sidebarWidth] = string(borderChar) // Sidebar separator
	}
	// Draw horizontal borders
	for x := 0; x < gameWidth; x++ {
		r.buffer[0][x] = string(borderChar)
		r.buffer[hudHeight-1][x] = string(borderChar)
		r.buffer[gameHeight-1][x] = string(borderChar)
	}

	// Draw corners
	r.buffer[0][0] = string(cornerChar)
	r.buffer[0][gameWidth-1] = string(cornerChar)
	r.buffer[gameHeight-1][0] = string(cornerChar)
	r.buffer[gameHeight-1][gameWidth-1] = string(cornerChar)

	// Draw title
	title := " Tower Defense "
	titleStart := (gameWidth - len(title)) / 2
	for i, ch := range title {
		r.buffer[1][titleStart+i] = string(ch)
	}
}

func (r *Renderer) drawGameArea(gs *core.GameState) {
	r.drawPath(gs.GetEnemyPath())
	r.drawTowers(gs.GetTowers())
	r.drawEnemies(gs.GetEnemies())
}

func (r *Renderer) clearBuffer() {
	for y := range r.buffer {
		for x := range r.buffer[y] {
			r.buffer[y][x] = " "
		}
	}
}

func (r *Renderer) drawPath(enemyPath []entities.BaseEntity) {
	if len(enemyPath) < 2 {
		return // Need at least two points to draw a path
	}

	for i := 0; i < len(enemyPath)-1; i++ {
		start := enemyPath[i]
		end := enemyPath[i+1]

		// Convert world coordinates to screen coordinates
		startX, startY := r.worldToScreen(start.X, start.Y)
		endX, endY := r.worldToScreen(end.X, end.Y)

		// Draw line between start and end points
		r.drawLine(startX, startY, endX, endY)
	}
}

func (r *Renderer) drawLine(x1, y1, x2, y2 int) {
	dx := abs(x2 - x1)
	dy := abs(y2 - y1)
	sx, sy := 1, 1
	if x1 >= x2 {
		sx = -1
	}
	if y1 >= y2 {
		sy = -1
	}
	err := dx - dy

	for {
		if r.isInBounds(x1, y1) {
			r.buffer[y1][x1] = "."
		}
		if x1 == x2 && y1 == y2 {
			break
		}
		e2 := 2 * err
		if e2 > -dy {
			err -= dy
			x1 += sx
		}
		if e2 < dx {
			err += dx
			y1 += sy
		}
	}
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func (r *Renderer) drawTowers(towers []*entities.Tower) {
	for _, tower := range towers {
		x, y := tower.GetPosition()
		screenX, screenY := r.worldToScreen(x, y)
		if r.isInBounds(screenX, screenY) {
			r.buffer[screenY][screenX] = string(towerChar)
		}
	}
}

func (r *Renderer) drawEnemies(enemies []*entities.Enemy) {
	for _, enemy := range enemies {
		x, y := enemy.GetPosition()
		screenX, screenY := r.worldToScreen(x, y)
		if r.isInBounds(screenX, screenY) {
			r.buffer[screenY][screenX] = string(enemyChar)
		}
	}
}

func (r *Renderer) drawHUD(gs *core.GameState) {
	hudInfo := fmt.Sprintf("Wave: %d | Lives: %d | Money: %d", gs.GetWave(), gs.GetLives(), gs.GetMoney())
	r.drawText(gameHeight-1, 1, hudInfo)
}

func (r *Renderer) drawSidebar(gs *core.GameState) {
	sidebarX := gameWidth - sidebarWidth + 1
	r.drawText(3, sidebarX, "Tower Types:")
	r.drawText(4, sidebarX, "1. Basic Tower  $50")
	r.drawText(5, sidebarX, "2. Sniper Tower $100")
	r.drawText(6, sidebarX, "3. AOE Tower    $150")

	r.drawText(8, sidebarX, "Controls:")
	r.drawText(9, sidebarX, "B: Build mode")
	r.drawText(10, sidebarX, "U: Upgrade tower")
	r.drawText(11, sidebarX, "S: Sell tower")
	r.drawText(12, sidebarX, "P: Pause game")

	r.drawText(14, sidebarX, "Stats:")
	r.drawText(16, sidebarX, fmt.Sprintf("Towers Built: %d", len(gs.GetTowers())))
}

func (r *Renderer) drawText(y, x int, text string) {
	for i, ch := range text {
		if x+i < gameWidth-1 {
			r.buffer[y][x+i] = string(ch)
		}
	}
}

func (r *Renderer) display() {
	fmt.Print("\033[H\033[2J") // Clear the console

	var sb strings.Builder
	sb.Grow(gameWidth * gameHeight) // Pre-allocate buffer
	for _, row := range r.buffer {
		sb.WriteString(strings.Join(row, ""))
		sb.WriteRune('\n')
	}
	fmt.Print(sb.String())
}

func (r *Renderer) worldToScreen(x, y float64) (int, int) {
	screenX := int(x * float64(gameWidth-sidebarWidth-2) / 800)
	screenY := int(y * float64(gameHeight-hudHeight-3) / 600)

	// Ensure we're not writing to the border
	screenX = max(1, min(screenX, gameWidth-sidebarWidth-2))
	screenY = max(hudHeight, min(screenY, gameHeight-2))

	return screenX, screenY
}

func (r *Renderer) isInBounds(x, y int) bool {
	return x > 0 && x < gameWidth-sidebarWidth-1 && y >= hudHeight && y < gameHeight-1
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
