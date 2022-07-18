package restart

import (
	"io"
	"log"
	"os"
)

var logger *log.Logger

func init() {
	logger = log.New(os.Stderr, "[Restart]", log.Llongfile|log.Lmicroseconds|log.Ldate)
}

func NewLogger(out io.Writer) {
	logger = log.New(out, "[Restart]", log.Llongfile|log.Lmicroseconds|log.Ldate)
}
