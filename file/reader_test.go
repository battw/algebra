package file

import (
	"testing"
)

func Test_Read(t *testing.T) {
	//read file
	linechan := make(chan string)
	err := Read("testscript.rw", linechan)
	if err != nil {
		t.Errorf("Error reading file:\n%s", err)
	}
	//check stuff is coming out of the channel
	line1 := <-linechan
	line1str := "one two three four five"
	if line1 != line1str {
		t.Errorf("Error reading first line of file. Is\n%s\nshould be\n%s", line1, line1str)
	}
	line2 := <-linechan
	line2str := "six seven eight nine ten"
	if line2 != line2str {
		t.Errorf("Error reading second line of file. Is\n%s\nshould be\n%s", line2, line2str)
	}
	line3 := <-linechan
	line3str := "eleven twelve thirteen fourteen fifteen"
	if line3 != line3str {
		t.Errorf("Error reading third line of file. Is\n%s\nshould be\n%s", line3, line3str)
	}
}
