package frameio

type SavedFrames struct {
	Frames    []string `msgpack:"0"`
	Framerate float64  `msgpack:"1"`
	AudioWav  []byte   `msgpack:"2"`
}
