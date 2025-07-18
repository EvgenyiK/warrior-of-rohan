package game

import (
	"github.com/hajimehoshi/ebiten/v2"
	"image"
	_ "image/png"
	"log"
	"os"
)

const (
	screenWidth  = 320
	screenHeight = 240
	frameOX      = 0
	frameOY      = 32
	frameWidth   = 32
	frameHeight  = 32
	frameCount   = 8
)

var (
	runnerImage *ebiten.Image
)

type Game struct {
	count int
}

func (g *Game) Update() error {
	g.count++
	return nil
}

func loadImage() {
	file, err := os.Open("assets/img/thorfinn.png")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		log.Fatal(err)
	}

	runnerImage = ebiten.NewImageFromImage(img)
}

func (g *Game) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(-float64(frameWidth)/2, -float64(frameHeight)/2)
	op.GeoM.Translate(screenWidth/2, screenHeight/2)
	i := (g.count / 5) % frameCount
	sx, sy := frameOX+i*frameWidth, frameOY
	screen.DrawImage(runnerImage.SubImage(image.Rect(sx, sy, sx+frameWidth, sy+frameHeight)).(*ebiten.Image), op)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func RunEbiten() {
	//img, _, err := image.Decode(bytes.NewReader(images.Runner_png))
	//if err != nil {
	//log.Fatal(err)
	//}
	//runnerImage = ebiten.NewImageFromImage(img)

	loadImage()

	ebiten.SetWindowSize(screenWidth*2, screenHeight*2)
	ebiten.SetWindowTitle("RohanWarrior")
	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}
