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

var dataPath = flag.String("f", "", "Location of flashcard data in csv format [question, answer]")
var specificLine = flag.Int64("l", 0, "Use a specific line instead of selecting one at random")

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

	if *dataPath == "" {
		fmt.Print("Missing data source!\n")
		flag.PrintDefaults()
		os.Exit(1)
	}

	infile, err := os.Open(*dataPath)
	if err != nil {
		log.Fatalf("data file (%s) failed to open: %s", *dataPath, aurora.Red(err.Error()))
	}
	defer infile.Close()

	var randomLine int64

	if *specificLine == 0 {
		//count lines in file
		var numLines int64
		scanner := bufio.NewScanner(infile)

		for scanner.Scan() {
			numLines++
		}
		infile.Seek(0, 0)

		//pick a random line
		rand.Seed(int64(time.Now().Nanosecond()))

		if randomLine = rand.Int63n(numLines); randomLine == 0 {
			randomLine = 1;
		}
	} else {
		randomLine = *specificLine
	}



	var card *Card

	//scan to random line and read it
	reader := csv.NewReader(infile)

	var pos int64 = 1 //
	for {

		line, err := reader.Read()
		if err != nil {
			log.Fatalf("Failed to read CSV (or no data found in file): %s (looking for line %d)", aurora.Red(err.Error()), randomLine)
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
