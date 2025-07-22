package game

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	screenWidth  = 800
	screenHeight = 600
)

var (
	playerImage *ebiten.Image
	playerX     float64 = screenWidth / 2
	playerY     float64 = screenHeight / 2
	speed       float64 = 4
)

func init() {
	var err error
	// Используем ebitenutil.NewImageFromFile для загрузки изображения
	playerImage, _, err = ebitenutil.NewImageFromFile("assets/img/thorfinn.png")
	if err != nil {
		log.Fatal("Ошибка загрузки изображения:", err)
	}
}

type Game struct{}

func (g *Game) Update() error {
	if ebiten.IsKeyPressed(ebiten.KeyArrowUp) {
		playerY -= speed
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowDown) {
		playerY += speed
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) {
		playerX -= speed
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowRight) {
		playerX += speed
	}

	// Ограничения по границам экрана
	if playerX < 0 {
		playerX = 0
	}
	if playerY < 0 {
		playerY = 0
	}
	maxX := float64(screenWidth - playerImage.Bounds().Dx())
	maxY := float64(screenHeight - playerImage.Bounds().Dy())
	if playerX > maxX {
		playerX = maxX
	}
	if playerY > maxY {
		playerY = maxY
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(playerX, playerY)
	screen.DrawImage(playerImage, op)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func RunEbiten() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Simple Ebiten Game")

	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}
