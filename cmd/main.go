package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"sort"

	"github.com/Masterminds/semver"
)

var logger = DefaultLogger{}

func main() {

	raw, err := parseStdinJSON()

	if err != nil {
		logger.Info("The input is not in the expect structure, input should be a JSON array of strings")
		return
	}

	vs, err := parseVersions(raw)

	if err != nil {
		logger.Infof("Error: %v\n", err)
		os.Exit(2)
		return
	}

	sort.Sort(semver.Collection(vs))

	json := convertToJSON(vs)

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

func parseVersions(s []string) ([]*semver.Version, error) {

	vs := make([]*semver.Version, len(s))
	for i, r := range s {
		v, err := semver.NewVersion(r)
		if err != nil {
			errMsg := fmt.Sprintf(`the string "%s" is not a valid Semantic version`, r)
			return vs, errors.New(errMsg)
		}

		vs[i] = v
	}

	return vs, nil
}

// return ordered as original values as slice of string
func convertToJSON(vs []*semver.Version) string {

	s := make([]string, len(vs))

	for i, v := range vs {
		s[i] = v.Original()
	}

	b, _ := json.MarshalIndent(s, "", "\t")

	return string(b)
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
