package main

import (
	arg "github.com/yellowsink/termvid/args"
	"github.com/yellowsink/termvid/processing"
	"io/fs"
	"os"
	"path"
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
		framerate = processing.PreProcess(args.InputPath, tempDir)
		audioPath = path.Join(tempDir, "audio.wav")
	}

	// oh my god why does this lang not allow unused vars
	_, _, _ = frames, audioPath, framerate

	err = os.RemoveAll(tempDir)
	if err != nil {
		return
	}
}
