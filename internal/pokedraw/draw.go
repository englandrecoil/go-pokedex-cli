package pokedraw

import (
	"bytes"
	"io"
	"os"
)

func DisplayImage(pokemonName string, data []byte) error {
	fileName := "./images/" + pokemonName + "_image.png"

	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	if _, err = io.Copy(file, bytes.NewReader(data)); err != nil {
		return err
	}

	return nil

}
