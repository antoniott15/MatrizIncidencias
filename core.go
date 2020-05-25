package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func GetFiles() []string {
	var files []string
	err := filepath.Walk(DIRECTORY, func(path string, info os.FileInfo, err error) error {
		files = append(files, path)
		return nil
	})
	if err != nil {
		panic(err)
	}

	return files
}

func scanStopLisT(path string) (map[string]int, error) {

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)

	scanner.Split(bufio.ScanWords)
	words := make(map[string]int)
	for scanner.Scan() {
		words[scanner.Text()] = 1
	}

	return words, nil
}

func scanWords(path string, stopWords map[string]int, words map[string][]string) (map[string][]string, error) {
	if path == PATHSTOPLIST {
		return nil, nil
	}
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)

	scanner.Split(bufio.ScanWords)

	for scanner.Scan() {
		if _, ok := stopWords[scanner.Text()]; !ok {
			words[scanner.Text()] = append(words[scanner.Text()], path)
		}
	}

	return words, nil
}

func Save(result string) error {
	file, err := os.OpenFile(RESULT, os.O_CREATE|os.O_WRONLY|os.O_APPEND, os.ModeAppend)
	if err != nil {
		return err
	}

	defer file.Close()

	dataWrite := bufio.NewWriter(file)

	_, err = io.WriteString(dataWrite, "\n"+result)
	if err != nil {
		return err
	}

	dataWrite.Flush()
	return file.Sync()
}

func WritingResult(result map[string][]string) {
	for key, value := range result {
		val := fmt.Sprintf("%s: %s \n", key, value)
		if err := Save(val); err != nil {
			fmt.Println(err)
		}
	}
}

func initF() map[string][]string {
	files := GetFiles()
	stopList, err := scanStopLisT(PATHSTOPLIST)
	if err != nil {
		panic(err)
	}
	words := make(map[string][]string)
	for _, elements := range files {
		resultWords, err := scanWords(elements, stopList, words)
		if err != nil {
			panic(err)
		}
		for key, value := range resultWords {
			words[key] = value
		}
	}

	WritingResult(words)
	return words
}

func QueryMachine(query string, toQuery map[string][]string) bool {
	query = strings.ToLower(query)
	str := strings.Split(query, " ")

	var ok, ok2 bool
	if str[0] == "not" {
		_, ok = toQuery[str[1]]
		_, ok2 = toQuery[str[len(str)-1]]
	} else {
		_, ok = toQuery[str[0]]
		_, ok2 = toQuery[str[len(str)-1]]
	}

	if !strings.Contains(query, "not") {
		if strings.Contains(query, "or") {
			return ok || ok2
		} else if strings.Contains(query, "and") {

			return ok && ok2
		}
	} else {
		if str[0] == "not" {
			if strings.Contains(query, "or") {
				return !ok || ok2
			} else if strings.Contains(query, "and") {
				return !ok && ok2
			}
		} else {
			if strings.Contains(query, "or") {
				return ok || !ok2
			} else if strings.Contains(query, "and") {
				return ok && !ok2
			}
		}
	}

	return false
}

func main() {
	words := initF()

	query1 := "huye or asd"
	fmt.Print(QueryMachine(query1, words))

	query2 := "jardinero and llegan"
	fmt.Print(QueryMachine(query2, words))

	query3 := "mithril and not funeral"
	fmt.Print(QueryMachine(query3, words))

}
