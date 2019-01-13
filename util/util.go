package util

func Runechan(s string) <-chan rune {
	rch := make(chan rune)
	go func() {
		for _, c := range s {
			rch <- c
		}
		close(rch)
	}()
	return rch
}
