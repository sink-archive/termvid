package player

import (
	"fmt"
	"math"
	"time"
)

func playGeneric(iterator []string, renderFunc func(string), framerate float64) {
	frameTime := 1000 / framerate
	timeDebt := 0.0

	for _, item := range iterator {
		startTime := time.Now()

		print("\033[0;0H") // reset to 0,0 in console
		renderFunc(item)

		// measure time rendering took
		renderTime := float64(time.Since(startTime).Nanoseconds()) * 1000000
		// amount of time we need to compensate
		makeupTarget := renderTime + timeDebt
		// timedebt is made up for, clear it
		timeDebt = 0
		// max possible correction
		correction := math.Min(makeupTarget, frameTime)
		// if cant make up now, try later
		if makeupTarget > frameTime {
			timeDebt += makeupTarget - frameTime
		}

		duration := time.Duration(int64(correction * 1000000))
		time.Sleep(duration)
	}
}

func PlayAscii(frames []string, framerate float64) {
	playGeneric(frames, playAsciiRenderFunc, framerate)
}
func playAsciiRenderFunc(str string) {
	fmt.Printf(str)
}