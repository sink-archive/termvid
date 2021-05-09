package args

import (
	"os"
	"strconv"
)

func ProcessArgs() Args {
	processed := Args{}

	rawArgs := os.Args[1:]

	for i, arg := range rawArgs {
		switch arg {
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
