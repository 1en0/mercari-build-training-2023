package main

import (
	"crypto/sha256"
	"encoding/hex"
	"io/ioutil"
	"mime/multipart"
	"path"
)

func saveImage(img *multipart.FileHeader, filename string) error {
	imageFile, _ := img.Open()
	imgPath := path.Join(ImgDir, filename)
	data, _ := ioutil.ReadAll(imageFile)
	err := ioutil.WriteFile(imgPath, data, 0666)
	if err != nil {
		return err
	}
	return nil
}

func getImageHash(img *multipart.FileHeader) string {
	imageFile, _ := img.Open()
	data, _ := ioutil.ReadAll(imageFile)
	//print(img)
	hasher := sha256.New()
	hasher.Write(data)
	imgHash := hex.EncodeToString(hasher.Sum(nil))
	//img, err := c.FormFile("image")
	return imgHash
}
