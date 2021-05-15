package processing

import (
	"fmt"
	"github.com/floostack/transcoder"
	"github.com/floostack/transcoder/ffmpeg"
	_ "github.com/floostack/transcoder/ffmpeg"
	"io/fs"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"
	"time"
)

func getMeta(filePath string) transcoder.Metadata {
	ffmpegPath, ffprobePath := getBinPaths()
	tc := ffmpeg.
		New(&ffmpeg.Config{
			FfprobeBinPath: ffprobePath,
			FfmpegBinPath:  ffmpegPath,
		}).
		Input(filePath)

	meta, err := tc.GetMetadata()
	if err != nil {
		return nil
	}
	return meta
}

func firstVideoStream(meta transcoder.Metadata) transcoder.Streams {
	for _, stream := range meta.GetStreams() {
		if stream.GetCodecType() == "video" {
			return stream
		}
	}

	return nil
}

func extractAudio(inPath string, outPath string) error {
	audioFormat := "wav"
	skipVideo := true
	opts := ffmpeg.Options{
		OutputFormat: &audioFormat,
		SkipVideo:    &skipVideo,
	}

	_, err := getTranscoder(inPath).
		WithOptions(opts).
		Output(outPath).
		Start(opts)
	if err != nil {
		return err
	}

	return nil
}

func extractImages(inPath string, outDir string, width int, height int) error {
	/*imgFormat := "bmp"
	opts := ffmpeg.Options{
		OutputFormat: &imgFormat,
	}*/

	err := os.Mkdir(outDir, os.ModeDir|fs.ModePerm)
	if err != nil {
		return err
	}

	/*_, err = getTranscoder(inPath).
		WithOptions(opts).
		Output(path.Join(outDir, "%6d.bmp")).
		Start(opts)
	if err != nil {
		return err
	}*/

	// haha ffmpeg go brr until this damn library gets fixed
	cmd := exec.Command(
		"ffmpeg",
		"-i",
		inPath,
		"-s",
		strconv.Itoa(width)+"x"+strconv.Itoa(height),
		/*"-q:v",
		"1",*/
		path.Join(outDir, "%6d.png"))

	err = cmd.Run()
	if err != nil {
		return err
	}

	return nil
}

func getTranscoder(inputPath string) transcoder.Transcoder {
	ffmpegPath, ffprobePath := getBinPaths()

	return ffmpeg.New(&ffmpeg.Config{
		FfprobeBinPath: ffprobePath,
		FfmpegBinPath:  ffmpegPath,
	}).
		Input(inputPath)
}

func getBinPaths() (string, string) {
	ffprobePath := ""
	if _, err := os.Stat("/usr/bin/ffprobe"); !os.IsNotExist(err) {
		ffprobePath = "/usr/bin/ffprobe"
	} else {
		if _, err := os.Stat("/Program Files/ffmpeg/ffprobe.exe"); !os.IsNotExist(err) {
			ffprobePath = "/Program Files/ffmpeg/ffprobe.exe"
		}
	}

	ffmpegPath := ""
	if _, err := os.Stat("/usr/bin/ffmpeg"); !os.IsNotExist(err) {
		ffmpegPath = "/usr/bin/ffmpeg"
	} else {
		if _, err := os.Stat("/Program Files/ffmpeg/ffmpeg.exe"); !os.IsNotExist(err) {
			ffmpegPath = "/Program Files/ffmpeg/ffmpeg.exe"
		}
	}

	return ffmpegPath, ffprobePath
}

func PreProcess(inputPath string, tempDir string, width int, height int) float64 {
	startTime := time.Now()

	print("Reading metadata             ")

	// get metadata
	meta := getMeta(inputPath)
	if meta == nil {
		return 0
	}

	// get framerate
	frStr := firstVideoStream(meta).GetAvgFrameRate()
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
	err = extractAudio(inputPath, audioPath)
	if err != nil {
		return 0
	}

	audtime := float64(time.Since(startTime).Milliseconds()) / 1000
	fmt.Printf("Done in %ss\n", strconv.FormatFloat(audtime, 'f', 2, 64))

	print("Splitting into images        ")

	imgDir := path.Join(tempDir, "rawframes")
	err = extractImages(inputPath, imgDir, width, height)
	if err != nil {
		return 0
	}

	imgtime := float64(time.Since(startTime).Milliseconds()) / 1000
	fmt.Printf("Done in %ss\n", strconv.FormatFloat(imgtime, 'f', 2, 64))

	return framerate
}
