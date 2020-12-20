package gscan

import (
	"bufio"
	"os"
)

func GscanFromTerminal() string {
	gscan := bufio.NewScanner(os.Stdin)
	gscan.Scan()
	return gscan.Text()
}
