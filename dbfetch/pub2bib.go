package main

import (
	 "fmt"
	 "encoding/xml"
	 "io/ioutil"
	 "os"
)

type  Article struct {
	 XMLName xml.Name `xml:"PubmedArticleSet"`
	 Pmid int `xml:"PubmedArticle>MedlineCitation>PMID"`
	 Entry Entry `xml:"PubmedArticle>MedlineCitation>Article"`
}

type Author struct {
	 Lastname string `xml:"LastName"`
	 Firstname string `xml:"ForeName"`
	 Initials string
}

type Entry struct {
	 Title string `xml:"ArticleTitle"`
	 Page string `xml:"Pagination>MedlinePgn"`
	 Abstract string `xml:"Abstract>AbstractText"`
	 Affiliation string
	 Journal Journal `xml:"Journal"`
	 Authors []Author `xml:"AuthorList>Author"`
}

type  Journal struct {
	 Issn string `xml:"ISSN"`
	 Volume int `xml:"JournalIssue>Volume"`
	 Issue int `xml:"JournalIssue>Issue"`
	 Year int `xml:"JournalIssue>PubDate>Year"`
	 Day int `xml:"JournalIssue>PubDate>Day"`
	 Month string `xml:"JournalIssue>PubDate>Month"`
	 Name string `xml:"Title"`
	 Abbreviation string `xml:"ISOAbbreviation"`
	
}

func main() {
	 content,err := ioutil.ReadFile(os.Args[1])
	 if err != nil {
	 		panic(err)
	 }

	 var article Article
	 err = xml.Unmarshal(content,&article)
	 if err != nil {
	 		panic(err)
	 }

	 fmt.Printf("pubmed id %d\n",article.Pmid)
	 j := article.Entry.Journal
	 fmt.Printf("journal issn:%s issue:%d year:%d month:%s name:%s\n",j.Issn,j.Issue,j.Year,j.Month,j.Name)
	 for _,author := range article.Entry.Authors  {
			fmt.Printf("Firstname:%s Lastname:%s\n",author.Firstname,author.Lastname)
	 }

}
