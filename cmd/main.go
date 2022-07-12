package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"sort"

	"github.com/Masterminds/semver"
)

func main() {

	raw, err := parseStdinJSON()

	if err != nil {
		log.Println("input is not the expect structure, input should be a JSON array of strings")
		return
	}

	vs, err := parseVersions(raw)

	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	sort.Sort(semver.Collection(vs))

	json := convertToJSON(vs)

	fmt.Println(json)

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
