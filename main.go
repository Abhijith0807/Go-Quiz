package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
)

type problem struct {
	q string
	a string
}

func main() {
	csvFile := flag.String("file", "problems.csv", "a csv file of format 'question,answer'")
	timecount := flag.Int("time", 15, "Timer duration for the Quiz")
	flag.Parse()
	file, err := os.Open(*csvFile)
	if err != nil {
		fmt.Println("Error while opening the file", err)
		return
	}
	defer file.Close()

	reader := csv.NewReader(file)

	records, err := reader.ReadAll()
	if err != nil {
		fmt.Println("Error while reading the csv", err)
		return
	}
	var correct_cnt int
	var incorrect_cnt int
	problems := parseRecords(records)
	answers := make(chan string)
	go func() {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			answers <- scanner.Text()
		}
		if err := scanner.Err(); err != nil {
			fmt.Println("Error while reading input", err)
		}
		close(answers)
	}()
	fmt.Printf("Get ready for the QUIZZZZZ, Your Timer of %v starts now!!\n", *timecount)

	timer := time.NewTimer(time.Duration(*timecount) * time.Second)
	defer timer.Stop()
	for _, prob := range problems {
		fmt.Printf("What is %v ?\n", prob.q)
		select {
		case answer, ok := <-answers:
			if !ok {
				fmt.Println("Input channel is closed")
				return
			}
			if answer == prob.a {
				correct_cnt += 1
			} else {
				incorrect_cnt += 1
			}
		case <-timer.C:
			fmt.Println("Time's up!")
			fmt.Printf("You have answered %v correctly and %v incorrectly with %v pending questions", correct_cnt, incorrect_cnt, len(problems)-correct_cnt-incorrect_cnt)
			return
		}
	}
	fmt.Printf("You have answered %v correctly and %v incorrectly", correct_cnt, incorrect_cnt)
}

func parseRecords(records [][]string) []problem {
	ret := make([]problem, len(records))
	for i, row := range records {
		ret[i] = problem{
			q: row[0],
			a: strings.TrimSpace(row[1]),
		}
	}
	return ret
}
