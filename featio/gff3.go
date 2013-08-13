package gff3

import (
	"regexp"
	"strings"
)

var DirecRegexp = regexp.MustCompile(`^##(\S+)\s+(.+)$`)
var SpaceRegexp = regexp.MustCompile(`\s+`)

func ParseDirective(line string) map[string]interface{} {
	directive := make(map[string]interface{})
	match := DirecRegexp.FindStringSubmatch(line)
	
	directive["directive"] = match[1]
	directive["content"] = match[2]

	if match[1] == "sequence-region" {
		value := SpaceRegexp.Split(match[2], -1)
		directive["seqid"] = value[0]
		directive["start"] = value[1]
		directive["end"] = value[2]
	} else if match[1] == "genome-build" {
		value := SpaceRegexp.Split(match[2], -1)
		directive["source"] = value[0]
		directive["buildname"] = value[1]
	}
	return directive
}

func ParseFeature(line string) map[string]interface{} {
	feat := make(map[string]interface{})
	arr := strings.Split(line, "\t")

	feat["seq_id"] = arr[0]
	feat["source"] = arr[1]
	feat["type"] = arr[2]
	feat["start"] = arr[3]
	feat["end"] = arr[4]
	feat["score"] = arr[5]
	feat["strand"] = arr[6]
	feat["phase"] = arr[7]

	feat["attributes"] = ParseAttribute(arr[8])
	return feat
}

func ParseAttribute(line string) map[string][]string {
	//handling attributes
	am := make(map[string][]string)
	//Separate tags delimited by ";"
	attr := strings.Split(line, ";")

	//Split each tag delimited by "="
	for _, val := range attr {
		kv := strings.Split(val, "=")
		//Split tag with multiple values if any ","
		if strings.Contains(kv[1], ",") {
			avals := strings.Split(kv[1], ",")
			am[kv[0]] = avals
		} else {
			am[kv[0]] = []string{kv[1]}
		}
	}
	return am
}
