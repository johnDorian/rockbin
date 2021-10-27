package vacuum

import (
	"bufio"
	"fmt"
	"time"

	"os"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
)

var defaultCapacity = 60 * 40.

// Bin a type to represent the bin
type Bin struct {
	FilePath string
	Capacity float64
	Seconds  float64
	Unit     string
	Value    string
}

// Update update the bin values
func (b *Bin) Update() {
	file, err := os.Open(b.FilePath)
	if err != nil {
		log.Fatalln(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	line := ""
	for scanner.Scan() {
		line = scanner.Text()
		if strings.Contains(line, "bin_in_time") {
			line = strings.Split(line, "=")[1]
			line = strings.Trim(line, " ;")
			break
		}
	}
	file.Close()
	b.Seconds, err = strconv.ParseFloat(line, 32)
	log.WithFields(log.Fields{"bin_time": b.Seconds}).Info("Parsed bin time")
	if err != nil {
		log.Fatalln(err)
	}
	b.convert()
}

// Convert convert the observed value to desired value
func (b *Bin) convert() {

	switch b.Unit {
	case "%":
		b.Value = fmt.Sprintf("%.2f", b.Seconds/b.Capacity*100.)
	case "sec":
		b.Value = fmt.Sprintf("%.0f", (time.Duration(b.Seconds) * time.Second).Seconds())
	case "min":
		b.Value = fmt.Sprintf("%.0f", (time.Duration(b.Seconds) * time.Second).Minutes())
	default:
		b.Value = fmt.Sprintf("%.2f", b.Seconds/b.Capacity*100.)
	}
	log.WithFields(log.Fields{"bin_time": b.Seconds, "value": b.Value}).Info("Converted Value")

}
