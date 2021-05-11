package processing

import (
	"fmt"
	"gopkg.in/gographics/imagick.v3/imagick"
	"math"
	"strconv"
	"time"
)

func imgToAscii(filePath string, width int, height int) string {
	mw := imagick.NewMagickWand()
	defer mw.Destroy()
	err := mw.ReadImage(filePath)
	if err != nil {
		return ""
	}
	filter := imagick.FILTER_POINT
	if mw.GetImageWidth() != uint(width) || mw.GetImageHeight() != uint(height) {
		err := mw.ResizeImage(uint(width), uint(height), filter)
		if err != nil {
			return ""
		}
	}

	iterator := mw.NewPixelIterator()

	frame := ""

	for y := 0; y < height; y += 2 {
		topRow := iterator.GetNextIteratorRow()
		btmRow := iterator.GetNextIteratorRow()

		for x := 0; x < width; x++ {
			topR, topG, topB := topRow[x].GetRed(), topRow[x].GetGreen(), topRow[x].GetBlue()
			btmR, btmG, btmB := btmRow[x].GetRed(), btmRow[x].GetGreen(), btmRow[x].GetBlue()

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

func MagickInit() { imagick.Initialize() }
func MagickEnd()  { imagick.Terminate() }
