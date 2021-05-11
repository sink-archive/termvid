package processing

import (
	"fmt"
	"image"
	"image/draw"
	_ "image/draw"
	_ "image/jpeg"
	"math"
	"os"
	"strconv"
	"time"
)

func imgToAscii(filePath string, width int, height int) string {
	file, err := os.Open(filePath)
	if err != nil {
		return ""
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {

		}
	}(file)
	img, _, err := image.Decode(file)
	if err != nil {
		return ""
	}

	frame := ""

	rect := img.Bounds()
	rgba := image.NewRGBA(rect)
	draw.Draw(rgba, rect, img, rect.Min, draw.Src)

	pixels := rgba.Pix
	stride := rgba.Stride
	min := rect.Min

	for y := 0; y < height; y += 2 {
		for x := 0; x < width; x++ {
			topStart := (y-min.Y)*stride + (x-min.X)*4
			btmStart := (y+1-min.Y)*stride + (x-min.X)*4
			topR, topG, topB := pixels[topStart], pixels[topStart+1], pixels[topStart+2]
			btmR, btmG, btmB := pixels[btmStart], pixels[btmStart+1], pixels[btmStart+2]

			frame += "\u001b[38;2;" +
				strconv.Itoa(int(topR)) +
				";" +
				strconv.Itoa(int(topG)) +
				";" +
				strconv.Itoa(int(topB)) +
				";48;2;" +
				strconv.Itoa(int(btmR)) +
				";" +
				strconv.Itoa(int(btmG)) +
				";" +
				strconv.Itoa(int(btmB)) +
				"m▀"
		}

		frame += "\n"
	}

	return frame
}

func BatchToAscii(filePaths []string, width int, height int) []string {
	startTime := time.Now()

	padAmount := len(strconv.Itoa(len(filePaths)))

	print("Converting to ASCII art      ")
	var working []string
	for i, file := range filePaths {
		working = append(working, imgToAscii(file, width, height))

		if i%10 == 0 {
			// progress
			fmt.Printf("%0"+strconv.Itoa(padAmount)+"d / %d [%3d%%]",
				i,
				len(filePaths),
				100*i/len(filePaths))
			// move back
			moveBack := 10 + padAmount*2
			//print("\033[0;0H")
			print("\033[" + strconv.Itoa(moveBack) + "D")
		}
	}

	timeTaken := time.Since(startTime)
	println("Done in " +
		strconv.Itoa(int(math.Floor(timeTaken.Minutes()))) +
		"m " +
		strconv.FormatFloat(timeTaken.Seconds(), 'f', 2, 64) +
		"s")
	return working
}
