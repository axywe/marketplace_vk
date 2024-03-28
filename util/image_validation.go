package util

import (
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"net/http"
)

func IsValidImage(url string) int {
	response, err := http.Get(url)
	if err != nil {
		return -1
	}
	defer response.Body.Close()

	img, format, err := image.Decode(response.Body)
	if err != nil {
		return -2
	}

	isValidFormat := false
	switch format {
	case "jpeg", "png", "gif":
		isValidFormat = true
	}
	if !isValidFormat {
		return -3
	}

	bounds := img.Bounds()
	width := bounds.Max.X
	height := bounds.Max.Y

	if width > 1024 || height > 1024 {
		return -4
	}

	return 0
}
