package main

import (
	"fmt"
	"github.com/andiogenes/dvach/pkg/api/boards"
	"github.com/andiogenes/dvach/pkg/api/posts"
	"github.com/andiogenes/dvach/pkg/api/threads"
	"github.com/andiogenes/dvach/pkg/ui"
	"log"
	"os"
	"strconv"
)

func main() {
	args := os.Args

	// Match command line arguments count
	switch len(args) - 1 {
	// no args => run UI
	case 0:
		if err := ui.Run(); err != nil {
			log.Fatal(err)
		}
	// one argument (board) => print info about board(s)
	case 1:
		if board := args[1]; board == "." {
			boards.ListBoards()
		} else {
			threads.ListThreads(board)
		}
	// two arguments (board, thread) => print info about thread
	case 2:
		board := args[1]
		thread, err := strconv.ParseUint(args[2], 10, 64)
		if err != nil {
			fmt.Println("Cannot parse thread")
			break
		}
		posts.ListPosts(board, thread)
	}
}
