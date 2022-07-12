package main

import (
	"encoding/json"
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
		log.Fatalf("one item was an invalid version number:\n\t%s", err)
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
			return vs, NewMyError(`the string "%s" is not a valid Semantic version`, r, err)
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

type MyErr struct {
	msg string
	err error
}

func (e MyErr) Error() string {
	return fmt.Sprintf("%v, %s", e.err, e.msg)
}

func (e MyErr) Unwrap() error {
	return e.err
}

func NewMyError(s, input string, e error) MyErr {
	return MyErr{
		msg: fmt.Sprintf(s, input),
		err: e,
	}
}
