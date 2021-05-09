package args

import (
	"os"
	"strconv"
)

func ProcessArgs() Args {
	processed := Args{}

	for i := range os.Args {
		switch os.Args[i] {
		case "-i":
		case "--input":
			processed.InputPath = os.Args[i+1]
			i++
			break

		case "-t":
		case "--temp":
			processed.TempFolderPath = os.Args[i+1]
			i++
			break

		case "-h":
		case "--height":
			var err error
			processed.Height, err = strconv.Atoi(os.Args[i+1])
			if err != nil {
				println("Invalid height value")
				os.Exit(1)
			}
			i++
			break

		case "-w":
		case "--width":
			var err error
			processed.Width, err = strconv.Atoi(os.Args[i+1])
			if err != nil {
				println("Invalid width value")
				os.Exit(1)
			}
			i++
			break

		case "-s":
		case "--asciiSave":
			processed.AsciiSavePath = os.Args[i+1]
			i++
			break

		case "-a":
		case "--savedFrames":
			processed.UseSavedFrames = true
			break

		case "-v":
		case "--viu":
			processed.UseViu = true
			break
		}
	}

	if len(processed.InputPath) == 0 {
		println("Please supply an input file with -i")
		os.Exit(1)
	}

	return processed
}
