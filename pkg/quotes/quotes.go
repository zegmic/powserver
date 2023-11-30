package quotes

import (
	"bufio"
	"log"
	"math/rand"
	"os"
)

type Quotes struct {
	quotes []string
}

func New() Quotes {
	file, err := os.Open("quotes.txt")
	if err != nil {
		log.Panic(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	return Quotes{quotes: lines}
}

func (q Quotes) Quote() string {
	n := rand.Int()
	return q.quotes[n%len(q.quotes)]
}
