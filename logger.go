package restart

import (
	"log"
	"os"
)

var logger *log.Logger

func init() {
	logger = log.New(os.Stderr, "[Restart]", log.Llongfile|log.Lmicroseconds|log.Ldate)
}
