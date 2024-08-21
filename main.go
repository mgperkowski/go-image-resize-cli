package main

import (
	"fmt"
	"log"
	"path/filepath"
	"sync"

	"errors"
	"image"
	"os"
	"strconv"

	"github.com/disintegration/imaging"

	"github.com/mgperkowski/goasyncawait/async"
)

type Image struct {
	path      string
	name      string
	img       image.Image
	value     int
	dimension string
}

func processArgs(args []string) (string, string, int, error) {
	if len(args) != 4 {
		log.Fatal("Usage: go-resize ./path/to/image.jpg -h 300 or go-resize ./path/to/dir -w 500")
		return "", "", 0, errors.New("Invalid number of arguments")
	}

	path := args[1]
	flag := args[2]
	value := args[3]

	valueInt, err := strconv.Atoi(value)

	if err != nil {
		log.Fatal("Usage: go-resize ./path/to/image.jpg -h 300 or go-resize ./path/to/dir -w 500")
		return "", "", 0, errors.New("Invalid value -- must be an integer")
	}

	if flag != "-h" && flag != "-w" {
		log.Fatal("Usage: go-resize ./path/to/image.jpg -h 300 or go-resize ./path/to/dir -w 500")
		return "", "", 0, errors.New("Invalid flag -- must be -h or -w")
	}

	return path, flag, valueInt, nil
}

func resizeImages(path string, flag string, value int) ([]*Image, error) {
	var promises []*async.Promise
	var resizedImages []*Image
	var mutex sync.Mutex

	isDir, _ := isDirectory(path)

	if isDir {
		files, err := os.ReadDir(path)

		if err != nil {
			return nil, err
		}

		for _, file := range files {

			isImg, _ := isImage(filepath.Join(path, file.Name()))

			if isImg {
				p := async.NewPromise(func(resolve func(interface{}), reject func(error)) {
					image, err := imaging.Open(filepath.Join(path, file.Name()))

					if err != nil {
						reject(err)
					}

					fmt.Println("Resizing: ", file.Name())

					if flag == "-h" {
						resized := imaging.Resize(image, 0, value, imaging.Lanczos)
						imageStruct := Image{path: path, name: file.Name(), img: resized, value: value, dimension: "h"}
						mutex.Lock()
						resizedImages = append(resizedImages, &imageStruct)
						mutex.Unlock()
					} else {
						resized := imaging.Resize(image, value, 0, imaging.Lanczos)
						imageStruct := Image{path: path, name: file.Name(), img: resized, value: value, dimension: "w"}
						mutex.Lock()
						resizedImages = append(resizedImages, &imageStruct)
						mutex.Unlock()
					}

					resolve(nil)
				})
				promises = append(promises, p)
			}
		}

		_, err = async.AwaitAll(promises)

		if err != nil {
			return nil, err
		} else {
			return resizedImages, nil
		}
	} else {

		isImg, _ := isImage(path)

		if isImg {

			fileName := getFileName(path)

			directoryPath := path[:len(path)-len(fileName)]

			image, err := imaging.Open(path)

			if err != nil {
				return nil, err
			}

			fmt.Println("Resizing: ", fileName)

			if flag == "-h" {
				resized := imaging.Resize(image, 0, value, imaging.Lanczos)
				imageStruct := Image{path: directoryPath, name: fileName, img: resized, value: value, dimension: "h"}
				resizedImages = append(resizedImages, &imageStruct)
			} else {
				resized := imaging.Resize(image, value, 0, imaging.Lanczos)
				imageStruct := Image{path: directoryPath, name: fileName, img: resized, value: value, dimension: "w"}
				resizedImages = append(resizedImages, &imageStruct)
			}

			return resizedImages, nil
		} else {
			return nil, errors.New("invalid image file")
		}
	}
}

func getFileName(path string) string {
	file, err := os.Open(path)

	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	fileInfo, err := file.Stat()

	if err != nil {
		log.Fatal(err)
	}

	return fileInfo.Name()
}

func isImage(path string) (bool, error) {
	file, err := os.Open(path)

	if err != nil {
		return false, err
	}

	defer file.Close()

	_, _, err = image.Decode(file)

	if err != nil {
		return false, err
	} else {
		return true, nil
	}
}

func isDirectory(path string) (bool, error) {
	fileInfo, err := os.Stat(path)

	if err != nil {
		return false, err
	}

	return fileInfo.IsDir(), nil
}

func directoryExists(path string) (bool, error) {
	_, err := os.Stat(path)

	if err != nil {
		return false, err
	}

	return true, nil
}

func saveImages(images []*Image) {

	resizedDir := filepath.Join(images[0].path, "Resized_Images")

	dirExists, _ := directoryExists(resizedDir)

	if !dirExists {

		err := os.Mkdir(resizedDir, 0755)

		if err != nil {
			log.Fatal(err)
		}
	}

	for _, img := range images {

		valueAsString := strconv.Itoa(img.value)

		path := filepath.Join(resizedDir, img.dimension+valueAsString+"-"+img.name)
		err := imaging.Save(img.img, path)

		if err != nil {
			log.Fatal(err)
		} else {
			fmt.Println("Image saved to: ", path)
		}
	}
}

func main() {

	path, flag, value, err := processArgs(os.Args)

	if err != nil {
		log.Fatal(err)
		return
	}

	resizedImages, err := resizeImages(path, flag, value)

	if err != nil {
		log.Fatal(err)
		return
	}

	saveImages(resizedImages)

}
