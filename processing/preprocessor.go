package processing

import (
	"github.com/floostack/transcoder"
	"github.com/floostack/transcoder/ffmpeg"
	_ "github.com/floostack/transcoder/ffmpeg"
	"io/fs"
	"os"
	"os/exec"
	"path"
)

func GetMeta(filePath string) transcoder.Metadata {
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

func FirstVideoStream(meta transcoder.Metadata) transcoder.Streams {
	for _, stream := range meta.GetStreams() {
		if stream.GetCodecType() == "video" {
			return stream
		}
	}

	return nil
}

func ExtractAudio(inPath string, outPath string) error {
	audioFormat := "wav"
	skipVideo := true
	opts := ffmpeg.Options{
		OutputFormat: &audioFormat,
		SkipVideo:    &skipVideo,
	}

	_, err := GetTranscoder(inPath).
		WithOptions(opts).
		Output(outPath).
		Start(opts)
	if err != nil {
		return err
	}

	return nil
}

func ExtractImages(inPath string, outDir string) error {
	/*imgFormat := "bmp"
	opts := ffmpeg.Options{
		OutputFormat: &imgFormat,
	}*/

	err := os.Mkdir(outDir, os.ModeDir|fs.ModePerm)
	if err != nil {
		return err
	}

	/*_, err = GetTranscoder(inPath).
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
		path.Join(outDir, "%6d.bmp"))

	err = cmd.Run()
	if err != nil {
		return err
	}

	return nil
}

func GetTranscoder(inputPath string) transcoder.Transcoder {
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
