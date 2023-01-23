package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"os"
	"strings"
)

func main() {

	pathToSource := flag.String("source_path", "", "path to source file")
	pathToIni := flag.String("ini_path", "", "path to .ini file")

	flag.Parse()

	if len(*pathToSource) == 0 || len(*pathToIni) == 0 {
		if !replaceDefaultText() {
			os.Exit(-1)
		}
		fmt.Println("Empty argument! Version tag defaultly pathched to 1.0")
		os.Exit(1)
	}

	f, err := os.Open(*pathToIni)
	if err != nil {
		fmt.Println(err)
		if !replaceDefaultText() {
			os.Exit(-1)
		}
		fmt.Println("Open error! Version tag defaultly pathched to 1.0")
		os.Exit(1)
	}
	defer f.Close()

	varsToFind := []string{"VersionMajor", "VersionMinor", "VersionPatch"}
	versions := parseVarsInFile(f, varsToFind)
	var finalVersion []byte
	if len(versions) != 3 {
		finalVersion = append(finalVersion, "1.0"...)
	} else {
		for _, v := range varsToFind {
			finalVersion = append(finalVersion, versions[v]+"."...)
		}
		finalVersion = append(finalVersion[:0], finalVersion[:len(finalVersion)-1]...)
	}
	if !replaceTextInFile(*pathToSource, "{{version}}", string(finalVersion)) {
		if !replaceDefaultText() {
			os.Exit(-1)
		}
		os.Exit(1)
	}
	os.Exit(1)
}

func parseVarsInFile(file *os.File, vars []string) map[string]string {

	scanner := bufio.NewScanner(file)
	versions := make(map[string]string)
	for scanner.Scan() {
		for _, v := range vars {
			fmt.Println(scanner.Text())
			if strings.Contains(scanner.Text(), v) {
				ind := strings.Index(scanner.Text(), "=")
				versions[v] = scanner.Text()[ind+1:]
				break
			}
		}
		if len(vars) == 0 {
			break
		}
	}
	return versions
}

func replaceTextInFile(pathToSource string, textToReplace string, text string) bool {
	input, err := os.ReadFile(pathToSource)
	defer fmt.Println(err)
	if err != nil {
		return false
	}

	output := bytes.Replace(input, []byte(textToReplace), []byte(text), 1)
	if err = os.WriteFile(pathToSource, output, 0666); err != nil {
		return false
	}
	return true
}

func replaceDefaultText() bool {
	return replaceTextInFile("defaultPath", "{{version}}", "1.0")
}
