// Custom widgets that extend functionality of basic tui-go components.

package ui

import (
	"fmt"
	threadsapi "github.com/andiogenes/dvach/pkg/api/threads"
	"github.com/andiogenes/dvach/pkg/utils"
	"github.com/marcusolsson/tui-go"
)

const maxCommentWidth = 180

// threadsList is extension of tui.List with hidden data model of board.
type threadsList struct {
	threads []threadsapi.Thread
	*tui.List
}

// newThreadsList creates new threadList filled with given threads.
func newThreadsList(threads []threadsapi.Thread) *threadsList {
	list := tui.NewList()

	for _, v := range threads {
		comment := []rune(utils.FormatHTML(v.Comment, false))

		commentWidth := len(comment)
		if commentWidth > maxCommentWidth {
			commentWidth = maxCommentWidth
		}

		list.AddItems(fmt.Sprintf("%s %s", v.Id, string(comment[:commentWidth])))
	}

	return &threadsList{
		threads: threads,
		List:    list,
	}
}
