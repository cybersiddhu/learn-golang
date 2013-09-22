package serializer

import (
	"bitbucket.org/ww/goraptor"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	f, err := os.Create(normalizeFileName(os.Args[1]) + ".json")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	parser := goraptor.NewParser("rdfxml")
	defer parser.Free()

	serializer := goraptor.NewSerializer("json")
	defer serializer.Free()

	err = serializer.SetFile(f, "")
	if err != nil {
		panic(err)
	}

	ch := parser.ParseFile(os.Args[1], "")
	serializer.AddN(ch)
}

func normalizeFileName(name string) string {
	base := filepath.Base(name)
	all := strings.Split(base, ".")
	return all[0]
}
