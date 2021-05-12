package frameio

import (
	"github.com/vmihailenco/msgpack/v5"
	"io/ioutil"
)

func PackFrames(frames SavedFrames) ([]byte, error) {
	b, err := msgpack.Marshal(frames)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func UnpackFrames(b []byte) (SavedFrames, error) {
	var frames SavedFrames
	err := msgpack.Unmarshal(b, &frames)
	if err != nil {
		return SavedFrames{}, err
	}
	return frames, nil
}

func SaveFrames(frames SavedFrames, path string) error {
	packed, err := PackFrames(frames)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(path, packed, 644)
	if err != nil {
		return err
	}
	return nil
}

func LoadFrames(path string) (SavedFrames, error) {
	raw, err := ioutil.ReadFile(path)
	if err != nil {
		return SavedFrames{}, err
	}
	parsed, err := UnpackFrames(raw)
	if err != nil {
		return SavedFrames{}, err
	}
	return parsed, nil
}
