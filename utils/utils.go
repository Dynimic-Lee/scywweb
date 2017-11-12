package utils

import (
	"fmt"
	"image/jpeg"
	"image/png"
	"io"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
)

func GetImageSize(path string) (string, error) {
	f, ferr := os.Open(path)
	defer f.Close()
	if ferr != nil {
		log.Printf("[utils] GetImageSize file open fail err : [%s]", ferr.Error())
		return "", ferr
	}

	switch strings.ToLower(filepath.Ext(f.Name())) {
	case ".jpg", ".jpeg":
		j, jerr := jpeg.DecodeConfig(f)
		if jerr == nil {
			return fmt.Sprintf("%dx%d", j.Width, j.Height), nil
		}
	case ".png":
		p, perr := png.DecodeConfig(f)
		if perr == nil {
			return fmt.Sprintf("%dx%d", p.Width, p.Height), nil
		}
	default:
		return "undefine format", nil
	}

	return "", nil
}

func SaveFile(imageData multipart.File, imageHeader *multipart.FileHeader) error {
	filename := imageHeader.Filename
	//fmt.Println(header.Filename)

	out, err := os.Create("./data/image/" + filename)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, imageData)
	if err != nil {
		return err
	}

	return nil
}
