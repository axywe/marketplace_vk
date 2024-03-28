package util

import (
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"net/http"
)

func IsValidImage(url string) bool {
	response, err := http.Get(url)
	if err != nil {
		fmt.Println("Ошибка при загрузке изображения:", err)
		return false
	}
	defer response.Body.Close()

	img, format, err := image.Decode(response.Body)
	if err != nil {
		fmt.Println("Ошибка при декодировании изображения:", err)
		return false
	}

	isValidFormat := false
	switch format {
	case "jpeg", "png", "gif":
		isValidFormat = true
	}
	if !isValidFormat {
		fmt.Println("Недопустимый формат изображения")
		return false
	}

	bounds := img.Bounds()
	width := bounds.Max.X
	height := bounds.Max.Y

	if width > 1024 || height > 1024 {
		fmt.Println("Размер изображения превышает 1024x1024 пикселей")
		return false
	}

	return true
}
