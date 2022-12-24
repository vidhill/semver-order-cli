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

	if checkIfEmptyStdin(os.Stdin) {
		logger.Info("Empty input")
		os.Exit(2)
	}

	raw, err := parseStdinJSON(os.Stdin)

	if err != nil {
		logger.Info("The input is not in the expected structure, input should be a JSON array of strings")
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

	json := convertToJSON(ordered, true)

	logger.Info(json)

}

func parseStdinJSON(r io.Reader) ([]string, error) {
	record := []string{}

	err := json.NewDecoder(r).Decode(&record)

	return record, err
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

func convertToJSON(s []string, pretty bool) string {
	b := marshallIgnoreError(s, pretty)
	return string(b)
}

func marshallIgnoreError(i interface{}, pretty bool) []byte {
	if pretty {
		b, _ := json.MarshalIndent(i, "", "\t")
		return b
	}
	b, _ := json.Marshal(i)
	return b
}

func checkIfEmptyStdin(f *os.File) bool {
	stat, _ := f.Stat()
	emptyInput := (stat.Mode() & os.ModeCharDevice) == 0

	return !emptyInput
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
