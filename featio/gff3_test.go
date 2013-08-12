package gff3

import (
	"bufio"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
	"github.com/cybersiddhu/golang-set"
)

func TestDirective(t *testing.T) {
	dir, err := os.Getwd()
	if err != nil {
		t.Error("Could not get current direcotry")
	}

	dirgff := filepath.Join(dir, "data", "directives.gff3")
	f, err := os.Open(dirgff)
	if err != nil {
		t.Error("Could not open test file")
	}

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		m := ParseDirective(scanner.Text())
		if _, ok := m["directive"]; ok != true {
			t.Fatal("Unable to parse directive")
		}
		if _, ok := m["content"]; ok != true {
			t.Fatal("Unable to parse content")
		}

		if v, ok := m["directive"]; ok {
			//test for gff3-version
			if v == "gff3-version" {
				if m["content"] != 3 {
					t.Fatal("Unable to parse gff3-version")
				}
			}
			//more directive
			if v == "feature-ontology" {
				if m["content"] != "bar" {
					t.Fatal("Unable to parse feature-ontology")
				}
			}

			//test for sequence-region
			if v == "sequence-region" {
				for _, key := range []string{"seqid", "start", "end"} {
					if val, ok := m[key]; !ok {
						t.Errorf("Unable to parse %s from sequence-region", key)
						switch val {
						case "seqid":
							if val != "foo" {
								t.Fatal("Unable to parse value of seqid")
							}
						case "start":
							if val != 1 {
								t.Fatal("Unable to parse value of start")
							}
						case "end":
							if val != 100 {
								t.Fatal("Unable to parse value of end")
							}
						}
					}
				}
			}
		}
	}
}

func TestParseFeature(t *testing.T) {
	dir, err := os.Getwd()
	if err != nil {
		t.Error("Could not get current direcotry")
	}

	dirgff := filepath.Join(dir, "data", "spec_eden.gff3")
	f, err := os.Open(dirgff)
	if err != nil {
		t.Error("Could not open test file")
	}

	gffkeys := []string{"seq_id", "source", "type", "start", "end", "score", "strand", "phase", "attributes"}

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "##") {
			continue
		}
		fm := ParseFeature(line)
		for _, key := range gffkeys {
			if _, ok := fm[key]; !ok {
				t.Errorf("Unable to parse %s feature key", key)
			}

			val := fm[key]
			switch key {
			case "seq_id":
				if val != "ctg123" {
					t.Error("Unable to parse seq_id")
				}
			case "type":
				gtype, ok := val.(string)
				if ok {
					m, err := regexp.MatchString("gene|mRNA|CDS|exon", gtype)
					if !m && err != nil {
						t.Error("Unable to parse type")
					}
				} else {
					t.Error("Unable to get a proper type assertion")
				}
			case "start", "end":
				gloc, ok := val.(string)
				if ok {
					m, err := regexp.MatchString(`\d+`, gloc)
					if !m && err != nil {
						t.Error("Unable to parse location")
					}
				} else {
					t.Error("Unable to get a proper type assertion")
				}
			case "attributes":
				attr, ok := val.(map[string][]string)
				if ok {
					set := mapset.NewSet()
					for key := range attr {
						m, err := regexp.MatchString("ID|Parent|Name", key)
						if !m && err != nil {
							t.Error("Unable to parse attr")
						}
						set.Add(key)
						if set.ContainsAll([]string{"ID","Name"}){
							 if attr["ID"][0] != "gene00001" && attr["Name"][0] != "EDEN" {
							 		t.Error("Did not get correct value of attributes")
							 }
						}
					}
				} else {
					t.Error("Unable to get a proper type assertion")
				}
			}
		}
	}
}
