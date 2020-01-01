package threads

import (
	"encoding/json"
	"fmt"
	"github.com/andiogenes/dvach/pkg/utils"
	"github.com/microcosm-cc/bluemonday"
	"github.com/pkg/errors"
	"io"
	"log"
	"net/http"
)

type Thread struct {
	// Thread id
	Id string `json:"num"`
	// Thread subject
	Subject string `json:"subject"`
	// Beginning of the first threads post
	Comment string `json:"comment"`
}

// ListThreads prints all available threads.
//
// Function will panic if an error occurs.
func ListThreads(board string) {
	threads, err := GetThreads(board)
	if err != nil {
		log.Panic(err)
	}

	p := bluemonday.StrictPolicy()

	for _, v := range threads {
		id := utils.BlueColor(v.Id)
		subject := v.Subject
		comment := p.Sanitize(v.Comment)

		fmt.Printf("%s %s\n\t%s\n", id, subject, comment)
	}
}

// GetThreads returns list of all available threads in given board.
func GetThreads(board string) ([]Thread, error) {
	url := fmt.Sprintf("https://2ch.hk/%s/catalog.json", board)

	response, err := http.Get(url)
	if err != nil {
		return nil, errors.Wrapf(err, "cannot get threads for %s", board)
	}
	defer response.Body.Close()

	return parseThreads(response.Body)
}

// parseThreads parses threads from JSON response.
func parseThreads(reader io.Reader) ([]Thread, error) {
	var rawThreads struct {
		Threads []Thread `json:"threads"`
	}

	decoder := json.NewDecoder(reader)
	if err := decoder.Decode(&rawThreads); err != nil {
		return nil, errors.Wrap(err, "cannot parse threads")
	}

	return rawThreads.Threads, nil
}
