package stats

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"mentat/internal/core"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func showAddDeckModal(pages *tview.Pages, dirPath string, m *core.Mentat, nav Navigator, footer *tview.TextView) {
	dismiss := func() {
		pages.RemovePage("modal")
		footer.SetText("r:review - a:add deck - q:quit")
	}

	form := tview.NewForm()
	form.AddInputField("Name", "", 30, nil, nil)
	form.AddButton("Create", func() {
		name := form.GetFormItemByLabel("Name").(*tview.InputField).GetText()
		if name == "" {
			return
		}
		if !strings.HasSuffix(name, ".md") {
			name += ".md"
		}
		absPath := filepath.Join(m.VaultDir, dirPath, name)
		content := fmt.Sprintf("---\ncreated: %s\n---\n", time.Now().Format(time.DateOnly))
		os.WriteFile(absPath, []byte(content), 0644)
		nav.ShowStats()
	})
	form.AddButton("Cancel", func() {
		dismiss()
	})
	form.SetBorder(true).SetTitle("Add deck to " + dirPath)
	form.SetCancelFunc(dismiss)

	modal := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(nil, 0, 1, false).
		AddItem(tview.NewFlex().
			AddItem(nil, 0, 1, false).
			AddItem(form, 50, 1, true).
			AddItem(nil, 0, 1, false), 9, 1, true).
		AddItem(nil, 0, 1, false)

	pages.AddPage("modal", modal, true, true)
	footer.SetText("esc:cancel")
}

func drawDecklist(deckTable *tview.Table, m *core.Mentat, nav Navigator, pages *tview.Pages, footer *tview.TextView) tview.Primitive {
	rootNode := tview.NewTreeNode(".")
	nodePathsMap, err := addDirEntriesToNode(rootNode, m.VaultDir, m.VaultDir)
	if err != nil {
		panic(err)
	}

	tree := tview.NewTreeView().SetRoot(rootNode).SetCurrentNode(rootNode)
	tree.SetBorder(true).SetTitle("vault")

	tree.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyRune {
			switch event.Rune() {
			case 'j':
				// move down
				tree.Move(1)
				return nil
			case 'k':
				// move up
				tree.Move(-1)
				return nil
			case 'a':
				node := tree.GetCurrentNode()
				dirPath := "."
				if node.GetText() != "." {
					relPath, ok := nodePathsMap[node]
					if !ok {
						return nil
					}
					if strings.HasSuffix(relPath, ".md") {
						dirPath = filepath.Dir(relPath)
					} else {
						dirPath = relPath
					}
				}
				showAddDeckModal(pages, dirPath, m, nav, footer)
				return nil
			case 'r':
				// review the selected decks
				node := tree.GetCurrentNode()
				if node.GetText() == "." {
					// show all decks at the root node
					nav.ShowReview(m.Decks)
				} else {
					relPath, ok := nodePathsMap[node]
					if !ok {
						return nil
					}
					nodeDecks := findDecksAtPath(relPath, m.Decks)
					nav.ShowReview(nodeDecks)
				}
				return nil
			case 'q':
				nav.App().Stop()
				return nil
			}
		}
		return event
	})

	tree.SetChangedFunc(func(node *tview.TreeNode) {
		if node.GetText() == "." {
			drawAllDecksInfo(deckTable, m.Decks)
			return
		}
		relPath, ok := nodePathsMap[node]
		if !ok {
			return
		}
		nodeDecks := findDecksAtPath(relPath, m.Decks)

		if len(nodeDecks) > 1 {
			drawAllDecksInfo(deckTable, nodeDecks)
		} else if len(nodeDecks) == 1 {
			drawDeckInfo(deckTable, nodeDecks[0])
		}
	})

	return tree
}
