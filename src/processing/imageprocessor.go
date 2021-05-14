package processing

import (
	"fmt"
	"gopkg.in/gographics/imagick.v3/imagick"
	"math"
	"sort"
	"strconv"
	"sync"
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

			topR *= 255
			topG *= 255
			topB *= 255
			btmR *= 255
			btmG *= 255
			btmB *= 255

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
	//var working []string
	var working []struct {
		int
		string
	}

	// create a goroutine for every 100 images
	for i := 0; i < len(filePaths); i += 100 {

		waitingGroup := sync.WaitGroup{}
		waitingGroup.Add(int(math.Min(100, float64(len(filePaths)-i))))

		end := int(math.Min(float64(i+100), float64(len(filePaths))))

		for j, file := range filePaths[i:end] {
			go func(k int, f string) {
				ascii := imgToAscii(f, width, height)
				working = append(working, struct {
					int
					string
				}{k, ascii})
				// progress
				fmt.Printf("%0"+strconv.Itoa(padAmount)+"d / %d [%3d%%]",
					len(working),
					len(filePaths),
					100*len(working)/len(filePaths))
				// move back
				moveBack := 10 + padAmount*2
				//print("\033[0;0H")
				print("\033[" + strconv.Itoa(moveBack) + "D")

				waitingGroup.Done()
			}(i+j, file)
		}

		waitingGroup.Wait()

	}

	timeTaken := time.Since(startTime)
	println("Done in " +
		strconv.Itoa(int(math.Floor(timeTaken.Minutes()))) +
		"m " +
		strconv.FormatFloat(timeTaken.Seconds(), 'f', 2, 64) +
		"s")

	sort.SliceStable(working, func(i, j int) bool {
		return working[i].int < working[j].int
	})

	var framesSorted []string
	for _, frame := range working {
		framesSorted = append(framesSorted, frame.string)
	}

	return framesSorted
}

func MagickInit() { imagick.Initialize() }
func MagickEnd()  { imagick.Terminate() }
