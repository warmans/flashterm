package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"math/rand"
	"os"

	"time"

	"log"

	"github.com/logrusorgru/aurora"
)

var dataPath = flag.String("data.path", "data.csv", "Location of flashcard data in csv format [question, answer]")

type Card struct {
	A    string
	B    string
	Line int64
}

func (c *Card) SideA() string {
	return fmt.Sprintf("%s", aurora.Blue(c.A))
}

func (c *Card) SideB() string {
	return fmt.Sprintf("%s", aurora.Green(c.B))
}

func main() {

	flag.Parse()

	infile, err := os.Open(*dataPath)
	if err != nil {
		log.Fatalf("data file (%s) failed to open: %s", *dataPath, aurora.Red(err.Error()))
	}
	defer infile.Close()

	//count lines in file
	var numLines int64
	scanner := bufio.NewScanner(infile)

	for scanner.Scan() {
		numLines++
	}
	infile.Seek(0, 0)

	//pick a random line
	rand.Seed(int64(time.Now().Nanosecond()))
	randomLine := rand.Int63n(numLines)

	var card *Card

	//scan to random line and read it
	reader := csv.NewReader(infile)
	var pos int64
	for {

		line, err := reader.Read()
		if err != nil {
			log.Fatalf("Failed to read CSV (or no data found in file): %s", aurora.Red(err.Error()))
		}
		if len(line) != 2 {
			log.Fatalf("Malformed CSV. All rows should have two columns, %d has %d", randomLine, len(line))
		}

		if randomLine != pos {
			pos++
			continue
		}

		//create card
		card = &Card{A: line[0], B: line[1], Line: randomLine}
		break
	}

	//do output
	fmt.Printf("%s", card.SideA())
	bufio.NewReader(os.Stdin).ReadString(byte(10))
	fmt.Printf("\033[1A%s > %s\n", card.SideA(), card.SideB())
}
