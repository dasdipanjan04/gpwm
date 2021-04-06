package gscan

import (
	"bufio"
	"os"
)

// GscanFromTerminal takes the termial input and returns a string.
func GscanFromTerminal() string {
	gscan := bufio.NewScanner(os.Stdin)
	gscan.Scan()
	return gscan.Text()
}
