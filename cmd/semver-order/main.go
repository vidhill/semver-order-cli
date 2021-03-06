package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sort"

	"github.com/Masterminds/semver"
)

var logger = DefaultLogger{}

func main() {

	if checkIfEmptyStdin() {
		logger.Info("Empty input")
		os.Exit(2)
	}

	raw, err := parseStdinJSON()

	if err != nil {
		logger.Info("The input is not in the expect structure, input should be a JSON array of strings")
		return
	}

	vs, invalid, err := parseVersions(raw)

	if err != nil {
		logger.Infof("Error: %v\n", err)
		os.Exit(2)
		return
	}

	sort.Sort(semver.Collection(vs))

	ordered := getOriginalNames(vs)

	// append the invalid tag names to the end
	ordered = append(ordered, invalid...)

	json := convertToJSON(ordered)

	logger.Info(json)

}

func parseStdinJSON() ([]string, error) {
	record := []string{}

	dec := json.NewDecoder(os.Stdin)
	for {
		err := dec.Decode(&record)
		if err == io.EOF {
			return record, nil
		}
		if err != nil {
			return []string{}, err
		}
	}
}

func parseVersions(s []string) ([]*semver.Version, []string, error) {

	vs := []*semver.Version{}
	invalid := []string{}
	for _, r := range s {
		v, err := semver.NewVersion(r)
		if err != nil {
			invalid = append(invalid, r)
		} else {
			vs = append(vs, v)
		}
	}

	return vs, invalid, nil
}

// return ordered as original values as slice of string
func getOriginalNames(vs []*semver.Version) []string {

	s := make([]string, len(vs))

	for i, v := range vs {
		s[i] = v.Original()
	}

	return s
}

func convertToJSON(s []string) string {
	b, _ := json.MarshalIndent(s, "", "\t")

	return string(b)
}

func checkIfEmptyStdin() bool {
	stat, _ := os.Stdin.Stat()
	if (stat.Mode() & os.ModeCharDevice) == 0 {
		return false
	} else {
		return true
	}
}

type Logger interface {
	Info(...interface{})
	Infof(string, ...interface{})
}

type DefaultLogger struct{}

func (l DefaultLogger) Info(v ...interface{}) {
	fmt.Println(v...)
}

func (l DefaultLogger) Infof(s string, v ...interface{}) {
	fmt.Printf(s, v...)
}
