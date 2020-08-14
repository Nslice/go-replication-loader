package replication

import (
	"io/ioutil"
	"log"
	"os"
	"strings"
)

// LineBreak is the windows line break
const LineBreak = "\r\n"

// DescriptionLoader provides an ability
// to load info about issues included in the replication
type DescriptionLoader struct {
	FileLoader
}

// GetDescriptionContent looks for files with *.desc pattern and combines it in a string
func (loader *DescriptionLoader) GetDescriptionContent(removeAfterRead bool) string {
	var result []string

	files := loader.GetFiles("*.desc")
	for _, file := range files {
		dat, err := ioutil.ReadFile(file)

		if err != nil {
			log.Fatalf("%v couln't read due to error %v", file, err)
		} else {
			result = append(result, string(dat))
		}

		if removeAfterRead {
			err = os.Remove(file)
			if err != nil {
				log.Fatalf("%v couln't be removed due to error %v", file, err)
			}
		}
	}

	return strings.Join(result, LineBreak)
}
