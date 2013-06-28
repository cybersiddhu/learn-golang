package main

import (
	dbi "database/sql"
	"flag"
	"fmt"
	_ "github.com/mattn/go-oci8"
	"os"
	"sync"
)

func reportError(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

var wg sync.WaitGroup

func main() {
	var datasource string
	flag.StringVar(&datasource, "datasource", "", "Name of oracle datasource")
	flag.Parse()
	if datasource == "" {
		fmt.Fprintln(os.Stderr, "no *datasource* defined")
		os.Exit(1)
	}

	//database connection
	os.Setenv("NLS_LANG", "")
	orasource := "DPUR_CHADO/DPUR_CHADO@" + datasource
	dbh, err := dbi.Open("oci8", orasource)
	reportError(err)
	defer dbh.Close()

	//query to make
	const listQuery = `select dbxref.accession gene_id, feature.uniquename locus_tag
    FROM feature JOIN organism ON organism.organism_id = feature.organism_id
    JOIN cvterm ON cvterm.cvterm_id = feature.type_id
    JOIN dbxref ON dbxref.dbxref_id = feature.dbxref_id
    WHERE cvterm.name = 'gene'
    AND organism.common_name = :foo
    `

	//statement handle
	stmt, err := dbh.Prepare(listQuery)
	reportError(err)
	defer stmt.Close()

	//organisms to query
	organisms := []string{"pallidum", "fasciculatum"}
	colname := []string{"GeneID", "locus_tag"}

	//fetching data
	for _, name := range organisms {
		wg.Add(1)
		go fetchMapping(stmt, name, colname)
	}

	const purListQuery = `
SELECT gene.uniquename,polypeptide.uniquename FROM feature gene
JOIN organism ON organism.organism_id = gene.organism_id
JOIN cvterm ftype ON ftype.cvterm_id = gene.type_id
JOIN feature_relationship frel ON frel.object_id = gene.feature_id
JOIN feature transcript ON transcript.feature_id = frel.subject_id
JOIN cvterm ttype ON ttype.cvterm_id = transcript.type_id
JOIN feature_relationship frel2 ON frel2.object_id = transcript.feature_id
JOIN feature polypeptide ON polypeptide.feature_id = frel2.subject_id
JOIN cvterm ptype ON ptype.cvterm_id = polypeptide.type_id
WHERE
ftype.name = 'gene'
AND ttype.name = 'mRNA'
AND ptype.name = 'polypeptide'
AND organism.common_name = :foo
    `
	stmt2, err := dbh.Prepare(purListQuery)
	reportError(err)
	wg.Add(1)
	go fetchMapping(stmt2, "purpureum", []string{"GeneID", "ProteinID"})

	wg.Wait()
}

func fetchMapping(stmt *dbi.Stmt, name string, colname []string) {

	f, err := os.Create(name + ".txt")
	reportError(err)
	defer f.Close()

	rows, err := stmt.Query(name)
	reportError(err)
	defer wg.Done()

	fmt.Fprintf(f, "%s\t%s\n", colname[0], colname[1])
	for rows.Next() {
		var acc string
		var locus string
		err = rows.Scan(&acc, &locus)
		reportError(err)
		fmt.Fprintf(f, "%s\t%s\n", acc, locus)
	}
	rows.Close()
	fmt.Printf("done writing %s\n", name+".txt")
}
