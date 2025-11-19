package apa

import (
	"fmt"
)

type Author struct {
	Given  string `json:"given"`
	Family string `json:"family"`
}

type Message struct {
	Authors        []Author `json:"author"`
	Title          []string `json:"title"`
	ContainerTitle []string `json:"container-title"`
	Volume         string   `json:"volume"`
	Issue          string   `json:"issue"`
	Page           string   `json:"page"`
	PublishedPrint struct {
		DateParts [][]int `json:"date-parts"`
	} `json:"published-print"`
}

// CrossRefData ...
type CrossRefData struct {
	Status  string  `json:"status"`
	Message Message `json:"message"`
}

func PrintCitation(data CrossRefData) {
	msg := data.Message

	// Year
	year := "n.d."
	if len(msg.PublishedPrint.DateParts) > 0 && len(msg.PublishedPrint.DateParts[0]) > 0 {
		year = fmt.Sprintf("%d", msg.PublishedPrint.DateParts[0][0])
	}

	// Authors (in-text APA)
	// If 1 author: (Smith, 2020)
	// If 2 authors: (Smith & Doe, 2020)
	// If 3+ authors: (Smith et al., 2020)

	citation := ""

	if len(msg.Authors) == 0 {
		citation = fmt.Sprintf("(Unknown, %s)", year)
	} else if len(msg.Authors) == 1 {
		a := msg.Authors[0]
		citation = fmt.Sprintf("(%s, %s)", a.Family, year)
	} else if len(msg.Authors) == 2 {
		a1 := msg.Authors[0].Family
		a2 := msg.Authors[1].Family
		citation = fmt.Sprintf("(%s & %s, %s)", a1, a2, year)
	} else {
		a1 := msg.Authors[0].Family
		citation = fmt.Sprintf("(%s et al., %s)", a1, year)
	}

	fmt.Println("Cite:")
	fmt.Println(citation)
}

func PrintReference(data CrossRefData) {
	// Build APA reference
	msg := data.Message

	// Authors: "Family, G."
	authorsAPA := ""
	for i, a := range msg.Authors {
		part := fmt.Sprintf("%s, %s.", a.Family, string([]rune(a.Given)[0]))
		if i == 0 {
			authorsAPA += part
		} else if i == len(msg.Authors)-1 {
			authorsAPA += " & " + part
		} else {
			authorsAPA += ", " + part
		}
	}

	// Year
	year := ""
	if len(msg.PublishedPrint.DateParts) > 0 && len(msg.PublishedPrint.DateParts[0]) > 0 {
		year = fmt.Sprintf("%d", msg.PublishedPrint.DateParts[0][0])
	} else {
		year = "n.d."
	}

	// Title
	title := ""
	if len(msg.Title) > 0 {
		title = msg.Title[0]
	}

	// Container title (journal)
	journal := ""
	if len(msg.ContainerTitle) > 0 {
		journal = msg.ContainerTitle[0]
	}

	// Pages
	pages := msg.Page

	// APA format
	apa := fmt.Sprintf("%s (%s). %s. %s, %s(%s), %s.",
		authorsAPA,
		year,
		title,
		journal,
		msg.Volume,
		msg.Issue,
		pages,
	)

	fmt.Println("APA:")
	fmt.Println(apa)
}
