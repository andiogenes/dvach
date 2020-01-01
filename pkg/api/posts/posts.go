package posts

import (
	"encoding/json"
	"fmt"
	"github.com/andiogenes/dvach/pkg/utils"
	"github.com/pkg/errors"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type Post struct {
	Id      uint64  `json:"num"`
	Comment string  `json:"comment"`
	Date    string  `json:"date"`
	Images  []Image `json:"files"`
}

type Image struct {
	Name     string `json:"name"`
	FullName string `json:"fullname"`
	Path     string `json:"path"`
}

// ListPosts prints all posts in given thread on a given board.
//
// Function will panic if an error occurs.
func ListPosts(board string, thread uint64) {
	posts, err := GetPosts(board, thread)
	if err != nil {
		log.Panic(err)
	}

	for _, v := range posts {
		id := utils.BlueColor(strconv.FormatUint(v.Id, 10))
		comment := utils.FormatHTML(v.Comment, true)
		date := utils.GreenColor(v.Date)
		images := formatImages(v.Images)

		fmt.Printf("%s %s %s\n\t%s\n", id, date, images, comment)
	}
}

// GetPosts returns list of all posts in given thread on a given board.
func GetPosts(board string, thread uint64) ([]Post, error) {
	url := fmt.Sprintf("https://2ch.hk/%s/res/%v.json", board, thread)

	response, err := http.Get(url)
	if err != nil {
		return nil, errors.Wrapf(err, "cannot get thread %s/%v", board, thread)
	}
	defer response.Body.Close()

	posts, err := parsePosts(response.Body)
	if err != nil {
		return nil, errors.Wrap(err, "cannot parse posts")
	}

	return posts, nil
}

// parsePosts parses posts from JSON response.
func parsePosts(reader io.Reader) ([]Post, error) {
	type thread struct {
		Posts []Post `json:"posts"`
	}

	var rawPosts struct {
		Threads []thread `json:"threads"`
	}

	decoder := json.NewDecoder(reader)
	if err := decoder.Decode(&rawPosts); err != nil {
		return nil, err
	}

	if len(rawPosts.Threads) == 0 {
		return nil, errors.New("threads must be present")
	}

	return rawPosts.Threads[0].Posts, nil
}

// formatImages represents slice of Image as string.
func formatImages(images []Image) string {
	descriptions := make([]string, len(images))

	for i, v := range images {
		name := v.Name
		path := v.Path
		fullName := utils.YellowColor(v.FullName)

		descriptions[i] = fmt.Sprintf("%s %s %s", fullName, path, name)
	}

	return strings.Join(descriptions, "\n")
}
