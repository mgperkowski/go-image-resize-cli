package main

import (
	"log"
	"path/filepath"

	"errors"
	"image"
	"os"
	"strconv"

	"github.com/disintegration/imaging"

	"github.com/mgperkowski/goasyncawait/async"

	"time"
)

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

func resizeImages(path string, flag string, value int) error {

	startTime := time.Now()
	resizedCount := 0

	var promises []*async.Promise

	isDir, _ := isDirectory(path)

	if isDir {

		resizedDir := filepath.Join(path, "Resized_Images")

		saveDirExists, _ := directoryExists(resizedDir)

		files, err := os.ReadDir(path)

		if err != nil {
			return err
		}

		for _, file := range files {

			isImg, _ := isImage(filepath.Join(path, file.Name()))

			if isImg {

				if !saveDirExists {

					err := os.Mkdir(resizedDir, 0755)

					if err != nil {
						log.Fatal(err)
					}

					saveDirExists = true

				}

				p := async.NewPromise(func(resolve func(interface{}), reject func(error)) {
					image, err := imaging.Open(filepath.Join(path, file.Name()))

					if err != nil {
						reject(err)
					}

					if flag == "-h" {
						resized := imaging.Resize(image, 0, value, imaging.Lanczos)
						err := imaging.Save(resized, filepath.Join(resizedDir, "h"+strconv.Itoa(value)+"-"+file.Name()))

						if err != nil {
							reject(err)
						} else {
							resizedCount++
							log.Println("Image resized and saved to: ", filepath.Join(resizedDir, "h"+strconv.Itoa(value)+"-"+file.Name()))
						}
					} else {
						resized := imaging.Resize(image, value, 0, imaging.Lanczos)
						err := imaging.Save(resized, filepath.Join(resizedDir, "w"+strconv.Itoa(value)+"-"+file.Name()))

						if err != nil {
							reject(err)
						} else {
							resizedCount++
							log.Println("Image resized and saved to: ", filepath.Join(resizedDir, "w"+strconv.Itoa(value)+"-"+file.Name()))
						}
					}

					resolve(nil)
				})
				promises = append(promises, p)
			}
		}

		_, err = async.AwaitAll(promises)

		elapsedTime := time.Since(startTime)

		log.Println("Resized ", resizedCount, " image(s) in ", elapsedTime)

		if err != nil {
			return err
		} else {
			return nil
		}
	} else {

		isImg, _ := isImage(path)

		if isImg {

			fileName := getFileName(path)

			directoryPath := path[:len(path)-len(fileName)]

			image, err := imaging.Open(path)

			if err != nil {
				return err
			}

			resizedDir := filepath.Join(directoryPath, "Resized_Images")

			saveDirExists, _ := directoryExists(resizedDir)

			if !saveDirExists {

				err := os.Mkdir(resizedDir, 0755)

				if err != nil {
					log.Fatal(err)
				}

				saveDirExists = true
			}

			if flag == "-h" {
				resized := imaging.Resize(image, 0, value, imaging.Lanczos)
				err := imaging.Save(resized, filepath.Join(resizedDir, "h"+strconv.Itoa(value)+"-"+fileName))

				if err != nil {
					return err
				} else {
					resizedCount++
					log.Println("Image resized and saved to: ", filepath.Join(resizedDir, "h"+strconv.Itoa(value)+"-"+fileName))
				}
			} else {
				resized := imaging.Resize(image, value, 0, imaging.Lanczos)
				err := imaging.Save(resized, filepath.Join(resizedDir, "h"+strconv.Itoa(value)+"-"+fileName))

				if err != nil {
					return err
				} else {
					resizedCount++
					log.Println("Image resized and saved to: ", filepath.Join(resizedDir, "h"+strconv.Itoa(value)+"-"+fileName))
				}
			}

			elapsedTime := time.Since(startTime)

			log.Println("Resized 1 image in ", elapsedTime)

			return nil
		} else {
			return errors.New("invalid image file")
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

func main() {

	path, flag, value, err := processArgs(os.Args)

	if err != nil {
		log.Fatal(err)
		return
	}

	err = resizeImages(path, flag, value)

	if err != nil {
		log.Fatal(err)
		return
	}

}
