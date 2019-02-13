package file

import (
	//	"errors"
	"io/ioutil"
	"strings"
)

//Read - puts the input file into the string channel, a line at a time.
func Read(path string, linechan chan<- string) error {
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	lines := strings.Split(string(bytes), "\n")
	go func() {
		for _, line := range lines {
			linechan <- line
		}
	}()
	return nil
}
