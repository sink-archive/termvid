package main

import (
	"github.com/faiface/beep"
	"github.com/faiface/beep/speaker"
	"github.com/faiface/beep/wav"
	arg "github.com/yellowsink/termvid/args"
	"github.com/yellowsink/termvid/frameio"
	"github.com/yellowsink/termvid/player"
	"github.com/yellowsink/termvid/processing"
	"io/fs"
	"os"
	"path"
	"sort"
	"strconv"
	"strings"
	"time"
)

func main() {
	processing.MagickInit()
	defer processing.MagickEnd()

	args := arg.ProcessArgs()

	tempDir, err := prepareTempDir(args)
	if err != nil {
		return
	}

	var frames []string
	var audioPath string
	var framerate float64

	if args.UseSavedFrames {
		playSaved(args, audioPath, tempDir)
		return
	}

	framerate = processing.PreProcess(args.InputPath, tempDir, args.Width, args.Height)
	audioPath = path.Join(tempDir, "audio.wav")

	dir, err := os.Open(path.Join(tempDir, "rawframes"))
	if err != nil {
		return
	}
	files, err := dir.ReadDir(0)
	var filePaths []string
	for _, file := range files {
		filePaths = append(filePaths, path.Join(path.Join(tempDir, "rawframes"), file.Name()))
	}

	sort.SliceStable(filePaths, func(i, j int) bool {
		split1 := strings.Split(filePaths[i], "/")
		split2 := strings.Split(filePaths[j], "/")
		path1 := split1[len(split1)-1]
		path2 := split2[len(split2)-1]
		parsed1, _ := strconv.Atoi(path1[:len(path1)-4])
		parsed2, _ := strconv.Atoi(path2[:len(path2)-4])
		return parsed1 < parsed2
	})

	frames = processing.BatchToAscii(filePaths, args.Width, args.Height)

	if len(args.AsciiSavePath) != 0 {
		save(args, audioPath, frames, framerate)
	} else {
		play(frames, framerate, audioPath, tempDir)
	}
}

func playSaved(args arg.Args, audioPath string, tempDir string) {
	savedFrames, err := frameio.LoadFrames(args.InputPath)
	if err != nil {
		return
	}
	frames, audio, framerate := savedFrames.Frames, savedFrames.AudioWav, savedFrames.Framerate

	audioPath = path.Join(tempDir, "audio.wav")
	os.WriteFile(audioPath, audio, 0644)

	play(frames, framerate, audioPath, tempDir)
}

func save(args arg.Args, audioPath string, frames []string, framerate float64) {
	audio, err := os.ReadFile(audioPath)
	if err != nil {
		return
	}
	savedFrames := frameio.SavedFrames{
		Frames:    frames,
		Framerate: framerate,
		AudioWav:  audio,
	}
	frameio.SaveFrames(savedFrames, args.AsciiSavePath)
}

func prepareTempDir(args arg.Args) (string, error) {
	tempDir := args.TempFolderPath
	if len(tempDir) == 0 {
		tempDir = path.Join(os.TempDir(), "termvid")
	}

	err := os.RemoveAll(tempDir)
	if err != nil {
		return "", err
	}

	err = os.Mkdir(tempDir, os.ModeDir|fs.ModePerm)
	if err != nil {
		return "", err
	}

	return tempDir, nil
}

func playAudio(streamer beep.Streamer, format beep.Format) (chan bool, error) {
	err := speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
	if err != nil {
		return nil, err
	}

	done := make(chan bool)
	speaker.Play(beep.Seq(streamer, beep.Callback(func() {
		done <- true
	})))

	return done, nil
}

func getStreamer(audioPath string) (beep.StreamSeekCloser, beep.Format, error) {
	audioF, err := os.Open(audioPath)
	if err != nil {
		return nil, beep.Format{}, err
	}
	streamer, format, err := wav.Decode(audioF)
	if err != nil {
		return nil, beep.Format{}, err
	}

	return streamer, format, nil
}

func play(frames []string, framerate float64, audioPath string, tempDir string) {
	streamer, format, err := getStreamer(audioPath)
	if err != nil {
		return
	}
	audioDoneChan, err := playAudio(streamer, format)
	if err != nil {
		return
	}

	player.PlayAscii(frames, framerate)

	err = os.RemoveAll(tempDir)
	if err != nil {
		return
	}

	<-audioDoneChan        // wait for audio to finish
	err = streamer.Close() // close the streamer now were done with it
	if err != nil {
		return
	}
}
