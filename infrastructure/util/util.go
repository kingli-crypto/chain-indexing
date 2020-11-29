package util
import (
	"os"
	"time"
	"fmt"
)
func WriteLog(filename string, text string) {
	title := fmt.Sprintf("%v-%s", time.Now(), text)

	f, err := os.OpenFile(filename,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	if _, err := f.WriteString(title); err != nil {
		panic(err)
	}
}
