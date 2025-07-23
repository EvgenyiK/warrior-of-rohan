package game

import (
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"image/color"
	_ "image/jpeg" // для декодирования JPEG
	_ "image/png"  // для PNG
	"io/ioutil"
	"log"
	"math/rand"
	"path/filepath"
	"strconv"
	"time"
	"warrior-of-rohan/internal/models"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text"
)

type Game struct{}

const (
	screenWidth  = 800
	screenHeight = 600
)

var (
	playerImage            *ebiten.Image
	backgroundImage        *ebiten.Image
	fontFace               font.Face
	backgrounds            []*ebiten.Image
	currentText            string
	playerX                float64 = screenWidth / 2
	playerY                float64 = screenHeight / 2
	speed                  float64 = 4
	currentBackgroundIndex         = 0
)

func init() {
	var err error

	rand.Seed(time.Now().UnixNano())
	backgroundImage = ebiten.NewImage(screenWidth, screenHeight)
	//fillBackgroundWithColor(backgroundImage, color.RGBA{R: 100, G: 150, B: 200, A: 255})

	// Используем ebitenutil.NewImageFromFile для загрузки изображения
	playerImage, _, err = ebitenutil.NewImageFromFile("assets/img/thorfinn.png")
	if err != nil {
		log.Fatal("Ошибка загрузки изображения:", err)
	}

	fontBytes, err := ioutil.ReadFile("assets/fonts/Norse-Regular.ttf")
	if err != nil {
		log.Fatal(err)
	}

	fontParsed, err := opentype.Parse(fontBytes)
	if err != nil {
		log.Fatal(err)
	}
	fontFace, err = opentype.NewFace(fontParsed, &opentype.FaceOptions{
		Size: 24,
		DPI:  72,
	})
	if err != nil {
		log.Fatal(err)
	}

}

func fillBackgroundWithColor(img *ebiten.Image, c color.Color) {
	img.Fill(c)
}

func loadBackgrounds() {
	// Путь к папке с изображениями
	folderPath := "assets/backgrounds"

	// Получаем список файлов в папке
	files, err := ioutil.ReadDir(folderPath)
	if err != nil {
		log.Println("Ошибка чтения папки:", err)
		return
	}

	for _, file := range files {
		// Проверяем, что это файл и что он имеет расширение png или jpg
		if file.IsDir() {
			continue
		}
		ext := filepath.Ext(file.Name())
		if ext != ".png" && ext != ".jpg" && ext != ".jpeg" {
			continue
		}

		// Полный путь к файлу
		path := filepath.Join(folderPath, file.Name())

		// Загружаем изображение из файла
		img, _, err := ebitenutil.NewImageFromFile(path)
		if err != nil {
			log.Println("Ошибка загрузки файла", path, ":", err)
			continue
		}

		backgrounds = append(backgrounds, img)
	}
}

func (g *Game) Update() error {
	moved := false

	if ebiten.IsKeyPressed(ebiten.KeyArrowUp) {
		playerY -= speed
		moved = true
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowDown) {
		playerY += speed
		moved = true
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) {
		playerX -= speed
		moved = true
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowRight) {
		playerX += speed
		moved = true
	}

	if moved {
		/*newColor := color.RGBA{
			R: uint8(rand.Intn(256)),
			G: uint8(rand.Intn(256)),
			B: uint8(rand.Intn(256)),
			A: 255,
		}*/
		//fillBackgroundWithColor(backgroundImage, newColor)
		loadBackgrounds()

		// Обновляем текст (например, показываем координаты)
		models.GameState.MU.Lock()
		models.GameState.Data = "X:" + strconv.Itoa(int(playerX)) + " Y:" + strconv.Itoa(int(playerY))
		currentData := models.GameState.Data
		models.GameState.MU.Unlock()

		currentText = currentData
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

	if ebiten.IsKeyPressed(ebiten.Key1) {
		handleChoice(0)
	} else if ebiten.IsKeyPressed(ebiten.Key2) {
		handleChoice(1)
	} else if ebiten.IsKeyPressed(ebiten.Key3) {
		handleChoice(2)
	} else if ebiten.IsKeyPressed(ebiten.Key4) && len(backgrounds) > 3 {
		handleChoice(3)
	}

	return nil
}

func handleChoice(index int) {
	switch index {
	case 0:
		currentText = "Вы выбрали путь через Морию."
		currentBackgroundIndex = 0
	case 1:
		currentText = "Вы выбрали путешествие через Горы."
		currentBackgroundIndex = 1
	case 2:
		currentText = "Вы выбрали лес Лотлориен."
		currentBackgroundIndex = 2
	case 3:
		currentText = "Вы выбрали Ривенделл."
		currentBackgroundIndex = 3
	}
}

func (g *Game) Draw(screen *ebiten.Image) {
	// Рисуем фон первым слоем
	screen.DrawImage(backgroundImage, nil)

	// Рисуем персонажа
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(playerX, playerY)
	screen.DrawImage(playerImage, op)

	// Рисуем текст поверх всего
	text.Draw(screen, currentText, fontFace, 10, 50, color.White)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func RunEbiten(dataChan <-chan string) {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Simple Ebiten Game")

	go func() {
		for data := range dataChan {
			models.GameState.MU.Lock()
			models.GameState.Data = data // обновляем состояние игры данными из канала
			models.GameState.MU.Unlock()
		}
	}()

	loadBackgrounds()

	currentBackgroundIndex = 0
	currentText = "Добро пожаловать в мир Властелина Колец! Сделайте выбор."

	ebiten.SetWindowSize(800, 600)
	ebiten.SetWindowTitle("Властелин Колец: Выбор пути")

	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}
