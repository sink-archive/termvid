package args

import (
	"fmt"
	"os"
	"strconv"
)

func ProcessArgs() Args {
	processed := Args{
		Width:          128,
		Height:         72,
		UseViu:         false,
		UseSavedFrames: false,
	}

	rawArgs := os.Args[1:]

	if len(rawArgs) == 0 {
		help()
	}

	for i, arg := range rawArgs {
		switch arg {
		case "--help":
			help()
		case "-i", "--input":
			processed.InputPath = rawArgs[i+1]
			i++

		case "-t", "--temp":
			processed.TempFolderPath = rawArgs[i+1]
			i++

		case "-h", "--height":
			var err error
			processed.Height, err = strconv.Atoi(rawArgs[i+1])
			if err != nil {
				println("Invalid height value")
				os.Exit(1)
			}
			i++

		case "-w", "--width":
			var err error
			processed.Width, err = strconv.Atoi(rawArgs[i+1])
			if err != nil {
				println("Invalid width value")
				os.Exit(1)
			}
			i++

		case "-s", "--asciiSave":
			processed.AsciiSavePath = rawArgs[i+1]
			i++

		case "-a", "--savedFrames":
			processed.UseSavedFrames = true

		case "-v", "--viu":
			processed.UseViu = true
		}
	}

	if len(processed.InputPath) == 0 {
		println("Please supply an input file with -i")
		os.Exit(1)
	}

	if processed.UseViu && processed.UseSavedFrames {
		println("Cannot use viu and read saved frames together")
		os.Exit(1)
	}

	if processed.UseViu && len(processed.AsciiSavePath) != 0 {
		println("Cannot use viu and save frames together")
		os.Exit(1)
	}

	return processed
}

func help() {
	fmt.Println("TermVid by Cain Atkinson (Yellowsink)\nLicensed under GPL-3.0-or-later\n\n" +
		"   --help        This screen\n" +
		"-i --input       Provide an input file\n" +
		"-t --temp        Choose a custom temp folder\n" +
		"-h --height      Provide a custom height for the image\n" +
		"-w --width       Provide a custom width for the image\n" +
		"-s --asciiSave   Save to a file for later playback\n" +
		"-a --savedFrames Use saved frames from -s\n" +
		"-v --viu         Display frames with viu")
	os.Exit(0)
}
