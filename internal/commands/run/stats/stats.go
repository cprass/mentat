package stats

import (
	"fmt"
	"maps"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"mentat/internal/core"

	"github.com/rivo/tview"
)

type Navigator interface {
	ShowStats()
	ShowReview(decks []*core.Deck)
	RenderView(view tview.Primitive)
	App() *tview.Application
}

func RenderStats(m *core.Mentat, nav Navigator) {
	decks := m.Decks

	pages := tview.NewPages()

	deckTable := tview.NewTable()

	deckTable.
		SetBorder(true).
		SetBorderPadding(1, 1, 2, 2)

	drawAllDecksInfo(deckTable, decks)

	footer := tview.NewTextView()
	footer.SetText("r:review - a:add deck - q:quit")

	deckList := drawDecklist(deckTable, m, nav, pages, footer)

	content := tview.NewFlex()
	content.AddItem(deckList, 0, 1, true)
	content.AddItem(deckTable, 0, 1, false)

	flex := tview.NewFlex()
	flex.SetDirection(tview.FlexRow)
	flex.AddItem(content, 0, 1, true)
	flex.AddItem(footer, 1, 0, false)

	pages.AddPage("main", flex, true, true)

	nav.RenderView(pages)
}

func setTableRow(t *tview.Table, row int, field1, field2 string) {
	t.SetCell(row, 0, tview.NewTableCell(field1))
	t.SetCell(row, 1, tview.NewTableCell(field2))
}

func drawTable(t *tview.Table, title string, cards, new, learning, due int) {
	t.Clear()
	t.SetTitle(title)
	setTableRow(t, 0, "Cards", strconv.Itoa(cards))
	setTableRow(t, 1, "New", strconv.Itoa(new))
	setTableRow(t, 2, "Learning", strconv.Itoa(learning))
	setTableRow(t, 3, "Due", strconv.Itoa(due))
}

func drawDeckInfo(t *tview.Table, deck *core.Deck) {
	drawTable(t, deck.Name, len(deck.Cards), 0, 0, deck.CardsDue())
}

func drawAllDecksInfo(t *tview.Table, decks []*core.Deck) {
	t.SetTitle(fmt.Sprintf("%d decks", len(decks)))
	t.Clear()

	var stats core.DeckStats

	for _, deck := range decks {
		stats = stats.Add(deck.Stats())
	}

	drawTable(t, "all cards", stats.Cards.Count, stats.Cards.New, stats.Cards.Learning, stats.Cards.Due)
}

func addDirEntriesToNode(parent *tview.TreeNode, path string, vaultPath string) (map[*tview.TreeNode]string, error) {
	relPathsMap := make(map[*tview.TreeNode]string)
	entries, err := os.ReadDir(path)
	if err != nil {
		return relPathsMap, err
	}
	for _, e := range entries {
		if e.IsDir() {
			// ignore specific directories
			if e.Name() == ".git" || e.Name() == ".jj" {
				continue
			}
			dirPath := filepath.Join(path, e.Name())
			relDirPath, err := filepath.Rel(vaultPath, dirPath)
			if err != nil {
				return relPathsMap, err
			}
			node := tview.NewTreeNode(e.Name())
			childPathsMap, err := addDirEntriesToNode(node, dirPath, vaultPath)
			if err != nil {
				return relPathsMap, err
			}
			if len(childPathsMap) == 0 {
				// Ignore node because there are no children
				continue
			}
			// Add child paths to the paths map
			maps.Copy(relPathsMap, childPathsMap)
			// Add dir path to the paths map
			relPathsMap[node] = relDirPath
			parent.AddChild(node)
			continue
		}

		if !strings.HasSuffix(e.Name(), ".md") {
			// ignore node
			continue
		}
		node := tview.NewTreeNode(e.Name())
		relPath, err := filepath.Rel(vaultPath, filepath.Join(path, e.Name()))
		if err != nil {
			return relPathsMap, err
		}
		relPathsMap[node] = relPath
		parent.AddChild(node)
	}

	return relPathsMap, nil
}

func findDecksAtPath(relPath string, decks []*core.Deck) []*core.Deck {
	result := make([]*core.Deck, 0)
	isSingle := strings.HasSuffix(relPath, ".md")

	for _, d := range decks {
		if isSingle && d.Location == relPath {
			result = append(result, d)
			break
		}
		if !isSingle && strings.HasPrefix(d.Location, relPath) {
			result = append(result, d)
			continue
		}
	}

	return result
}
