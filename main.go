package main

import (
	"fmt"
	arg "github.com/yellowsink/termvid/args"
	"github.com/yellowsink/termvid/processing"
	"io/fs"
	"os"
	"path"
	"strconv"
	"strings"
	"time"
)

func main() {
	args := arg.ProcessArgs()

	tempDir := args.TempFolderPath
	if len(tempDir) == 0 {
		tempDir = path.Join(os.TempDir(), "termvid")
	}

	err := os.RemoveAll(tempDir)
	if err != nil {
		return
	}

	err = os.Mkdir(tempDir, os.ModeDir|fs.ModePerm)
	if err != nil {
		return
	}

	var frames []string
	var audioPath string
	var framerate float64

	if !args.UseSavedFrames {
		framerate = preProcess(args.InputPath, tempDir)
		audioPath = path.Join(tempDir, "audio.wav")
	}

	// oh my god why does this lang not allow unused vars
	_, _, _ = frames, audioPath, framerate
}

func preProcess(inputPath string, tempDir string) float64 {
	startTime := time.Now()

	print("Reading metadata             ")

	// get metadata
	meta := processing.GetMeta(inputPath)
	if meta == nil {
		return 0
	}

	// get framerate
	frStr := processing.FirstVideoStream(meta).GetAvgFrameRate()
	frFract := strings.Split(frStr, "/")
	frNum, err := strconv.ParseFloat(frFract[0], 64)
	frDen, err := strconv.ParseFloat(frFract[1], 64)
	if err != nil {
		return 0
	}
	framerate := frNum / frDen
	if framerate == 0 {
		return 0
	}

	metaTime := time.Since(startTime).Milliseconds()
	fmt.Printf("Done in %sms\n", strconv.Itoa(int(metaTime)))

	print("Extracting Audio             ")

	audioPath := path.Join(tempDir, "audio.wav")
	err = processing.ExtractAudio(inputPath, audioPath)
	if err != nil {
		return 0
	}

	audtime := float64(time.Since(startTime).Milliseconds()) / 1000
	fmt.Printf("Done in %ss\n", strconv.FormatFloat(audtime, 'f', 2, 64))

	print("Splitting into images        ")

	imgDir := path.Join(tempDir, "rawframes")
	err = processing.ExtractImages(inputPath, imgDir)
	if err != nil {
		return 0
	}

	imgtime := float64(time.Since(startTime).Milliseconds()) / 1000
	fmt.Printf("Done in %ss", strconv.FormatFloat(imgtime, 'f', 2, 64))

	return framerate
}
