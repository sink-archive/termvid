package main

import (
	arg "github.com/yellowsink/termvid/args"
	"github.com/yellowsink/termvid/processing"
	"io/fs"
	"os"
	"path"
)

func main() {
	processing.MagickInit()
	defer processing.MagickEnd()

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
		framerate = processing.PreProcess(args.InputPath, tempDir, args.Width, args.Height)
		audioPath = path.Join(tempDir, "audio.wav")
	}

	dir, err := os.Open(path.Join(tempDir, "rawframes"))
	if err != nil {
		return
	}
	files, err := dir.ReadDir(0)
	var filePaths []string
	for _, file := range files {
		filePaths = append(filePaths, path.Join(path.Join(tempDir, "rawframes"), file.Name()))
	}

	processing.BatchToAscii(filePaths, args.Width, args.Height)

	// oh my god why does this lang not allow unused vars
	_, _, _ = frames, audioPath, framerate

	err = os.RemoveAll(tempDir)
	if err != nil {
		return
	}
}
