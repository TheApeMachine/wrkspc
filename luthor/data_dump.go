package luthor

import (
	"bufio"
	"bytes"
	"log"
	"regexp"

	"github.com/acarl005/stripansi"
)

type DataLine string

type DataDump struct {
	content *bytes.Buffer
}

func NewDataDump(content []byte) *DataDump {
	return &DataDump{
		content: bytes.NewBuffer(content),
	}
}

func (dump DataDump) GenerateLines() chan DataLine {
	out := make(chan DataLine)
	scanner := bufio.NewScanner(dump.content)
	// detector := chardet.NewTextDetector()

	go func() {
		defer close(out)

		for scanner.Scan() {
			latin := regexp.MustCompile(`[^(\x20-\x7F)\x0A\x0D]*`)
			line := latin.ReplaceAllString(scanner.Text(), "")

			space := regexp.MustCompile(`\s+`)
			line = space.ReplaceAllString(stripansi.Strip(line), " ")

			out <- DataLine(line)
		}

		if err := scanner.Err(); err != nil {
			log.Fatal(err)
		}
	}()

	return out
}
