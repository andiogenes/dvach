package ui

import (
	"fmt"
	"github.com/andiogenes/dvach/pkg/api/boards"
	"github.com/andiogenes/dvach/pkg/api/posts"
	"github.com/andiogenes/dvach/pkg/api/threads"
	"github.com/andiogenes/dvach/pkg/utils"
	"github.com/marcusolsson/tui-go"
	"image"
	"strconv"
)

// Run creates and runs instance of TUI client.
func Run() error {
	ui, err := tui.New(nil)
	if err != nil {
		return err
	}

	root := newMainView(ui)
	ui.SetWidget(root)

	ui.SetKeybinding("Esc", func() { ui.Quit() })

	if err := ui.Run(); err != nil {
		return err
	}

	return nil
}

// newMainView creates main view of application.
func newMainView(ui tui.UI) tui.Widget {
	// Create quit button
	quit := tui.NewButton("[Выйти]")
	quit.OnActivated(func(b *tui.Button) {
		ui.Quit()
	})

	// Create layout for view
	boardsList := newBoardsList(ui)
	boardsList.SetFocused(true)

	boardsWrapper := tui.NewHBox(boardsList)
	boardsWrapper.SetBorder(true)

	window := tui.NewVBox(
		tui.NewPadder(10, 1, tui.NewLabel("Добро пожаловать. Снова.")),
		boardsWrapper,
		tui.NewHBox(
			tui.NewSpacer(),
			tui.NewPadder(5, 2, quit),
			tui.NewSpacer(),
		),
	)
	window.SetBorder(true)

	container := tui.NewHBox(
		tui.NewSpacer(),
		tui.NewVBox(
			tui.NewSpacer(),
			window,
			tui.NewSpacer(),
		),
		tui.NewSpacer(),
	)

	// Set up focus chain
	focusChain := &tui.SimpleFocusChain{}
	focusChain.Set(boardsList, quit)

	ui.SetFocusChain(focusChain)

	return container
}

// newBoardsList creates list filled with board names and binds routing action on it.
func newBoardsList(ui tui.UI) tui.Widget {
	// Create and set up list widget
	list := tui.NewList()

	// Move to corresponding board when list item is activated
	list.OnItemActivated(func(l *tui.List) {
		ui.SetWidget(newThreadsView(ui, list.SelectedItem()))
	})

	retrievedBoards, _ := boards.GetBoards()

	// Fill the list widget
	items := make([]string, len(retrievedBoards))
	for i, v := range retrievedBoards {
		items[i] = v.Id
	}

	list.AddItems(items...)
	list.SetSelected(0)

	return list
}

// newThreadsView creates view with threads of given board.
func newThreadsView(ui tui.UI, board string) tui.Widget {
	retrievedThreads, _ := threads.GetThreads(board)

	// Create, fill and set up list of threads
	list := newThreadsList(retrievedThreads)
	list.SetFocused(true)
	list.Select(0)

	// Move to corresponding thread when list item is activated
	list.OnItemActivated(func(l *tui.List) {
		thread, err := strconv.ParseUint(list.threads[l.Selected()].Id, 10, 64)
		if err != nil {
			return
		}

		ui.SetWidget(newPostsView(ui, board, thread))
	})

	// Wrap list with border
	listWrapper := tui.NewHBox(list)
	listWrapper.SetBorder(true)

	// Create and set up quit button
	quit := tui.NewButton("[Выйти]")
	quit.OnActivated(func(b *tui.Button) {
		ui.SetWidget(newMainView(ui))
	})

	// Create layout for view
	window := tui.NewVBox(
		tui.NewHBox(
			tui.NewSpacer(),
			tui.NewPadder(10, 1, tui.NewLabel(fmt.Sprintf("/%s/", board))),
			tui.NewSpacer(),
		),
		listWrapper,
		tui.NewHBox(
			tui.NewSpacer(),
			tui.NewPadder(5, 2, quit),
		),
	)
	window.SetBorder(true)

	// Set up focus chain
	focusChain := &tui.SimpleFocusChain{}
	focusChain.Set(list, quit)

	ui.SetFocusChain(focusChain)

	return window
}

// newPostsView creates view with posts of given thread on a given board.
func newPostsView(ui tui.UI, board string, thread uint64) tui.Widget {
	retrievedPosts, _ := posts.GetPosts(board, thread)

	// Create slice of widgets that represent posts in thread.
	comments := make([]tui.Widget, len(retrievedPosts))
	for i, v := range retrievedPosts {
		label := tui.NewLabel(utils.FormatHTML(v.Comment, false))
		label.SetWordWrap(true)

		container := tui.NewHBox(
			label,
		)
		container.SetBorder(true)
		container.SetTitle(fmt.Sprintf("%v %s", v.Id, v.Date))

		comments[i] = container
	}

	// Put widgets in container
	postsContainer := tui.NewVBox(comments...)
	// Hack to fix layout in scroll area, use precisely
	postsContainer.Resize(image.Pt(20, 0))

	// Create and set up quit button
	quit := tui.NewButton("[Выйти]")
	quit.OnActivated(func(b *tui.Button) {
		ui.SetWidget(newThreadsView(ui, board))
	})
	quit.SetFocused(true)

	// Wrap container with posts in scroll area
	scrollArea := tui.NewScrollArea(postsContainer)
	scrollArea.SetSizePolicy(tui.Expanding, tui.Expanding)

	// Hope this will not cause NPE or memory leak in the future
	ui.SetKeybinding("Up", func() { scrollArea.Scroll(0, -4) })
	ui.SetKeybinding("Down", func() { scrollArea.Scroll(0, 4) })

	// Create layout for view
	window := tui.NewVBox(
		scrollArea,
		tui.NewHBox(
			tui.NewSpacer(),
			quit,
		),
	)
	window.SetBorder(true)

	// Set up focus chain
	focusChain := &tui.SimpleFocusChain{}
	focusChain.Set(quit)

	ui.SetFocusChain(focusChain)

	return window
}
