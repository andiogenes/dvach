package api

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"io"
	"log"
	"net/http"
)

type Board struct {
	// Board id, like "pr" or "b"
	Id string `json:"id"`
	// Board's category
	Category string `json:"category"`
	// Name of the board
	Name string `json:"name"`
}

// ListBoards prints all available boards.
//
// Function will panic if an error occurs.
func ListBoards() {
	boards, err := GetBoards()
	if err != nil {
		log.Panic(err)
	}

	for _, v := range boards {
		fmt.Printf("%-10s %20s %s\n", v.Id, v.Category, v.Name)
	}
}

// GetBoards returns lists of all available boards.
func GetBoards() ([]Board, error) {
	const url = "https://2ch.hk/makaba/mobile.fcgi?task=get_boards"

	response, err := http.Get(url)
	if err != nil {
		return nil, errors.Wrap(err, "cannot get boards")
	}
	defer response.Body.Close()

	return parseBoards(response.Body)
}

// parseBoards parses boards list from JSON response.
func parseBoards(reader io.Reader) ([]Board, error) {
	var rawBoards map[string][]Board

	decoder := json.NewDecoder(reader)
	if err := decoder.Decode(&rawBoards); err != nil {
		return nil, errors.Wrap(err, "cannot parse boards")
	}

	boards := make([]Board, 0)
	for _, v := range rawBoards {
		boards = append(boards, v...)
	}

	return boards, nil
}
