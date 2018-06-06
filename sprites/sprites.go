package sprites

import (
	"image"
	"os"

	_ "image/png"

	"github.com/faiface/pixel"
)

func LoadPicture(path string) pixel.Picture {
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}

	defer file.Close()
	img, _, err := image.Decode(file)
	if err != nil {
		panic(err)
	}

	return pixel.PictureDataFromImage(img)
}

func LoadSprite(path string) (sprite *pixel.Sprite, bounds pixel.Rect) {
	pic := LoadPicture(path)
	bounds = pic.Bounds()

	return pixel.NewSprite(pic, bounds), bounds
}
