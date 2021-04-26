package state

import (
	"log"

	"github.com/aretext/aretext/menu"
)

type MenuStyle int

const (
	MenuStyleCommand = MenuStyle(iota)
	MenuStyleFile
)

// MenuState represents the menu for searching and selecting items.
type MenuState struct {
	// visible indicates whether the menu is currently displayed.
	visible bool

	// style controls how the menu is displayed.
	style MenuStyle

	// search controls which items are visible based on the user's current search query.
	search *menu.Search

	// selectedResultIdx is the index of the currently selected search result.
	// If there are no results, this is set to zero.
	// If there are results, this must be less than the number of results.
	selectedResultIdx int
}

func (m *MenuState) Visible() bool {
	return m.visible
}

func (m *MenuState) Style() MenuStyle {
	return m.style
}

func (m *MenuState) SearchQuery() string {
	if m.search == nil {
		return ""
	}
	return m.search.Query()
}

func (m *MenuState) SearchResults() (results []menu.Item, selectedResultIdx int) {
	if m.search == nil {
		return nil, 0
	}
	return m.search.Results(), m.selectedResultIdx
}

// ShowMenu displays the menu with the specified style and items.
func ShowMenu(state *EditorState, style MenuStyle, loadItems func() []menu.Item) {
	emptyQueryShowAll := bool(style == MenuStyleFile)
	items := loadItems()
	if style == MenuStyleCommand {
		items = append(items, state.customMenuItems...)
	}
	search := menu.NewSearch(items, emptyQueryShowAll)
	state.menu = &MenuState{
		visible:           true,
		style:             style,
		search:            search,
		selectedResultIdx: 0,
	}
	state.inputMode = InputModeMenu
}

// HideMenu hides the menu.
func HideMenu(state *EditorState) {
	state.menu = &MenuState{}
	state.inputMode = InputModeNormal
}

// ExecuteSelectedMenuItem executes the action of the selected menu item and closes the menu.
func ExecuteSelectedMenuItem(state *EditorState) {
	search := state.menu.search
	results := search.Results()
	if len(results) == 0 {
		// If there are no results, then there is no action to perform.
		SetStatusMsg(state, StatusMsg{
			Style: StatusMsgStyleError,
			Text:  "No menu item selected",
		})
		HideMenu(state)
		return
	}

	idx := state.menu.selectedResultIdx
	selectedItem := results[idx]
	log.Printf("Executing menu item '%s' at result index %d\n", selectedItem.Name, idx)
	HideMenu(state)

	actionFunc, ok := selectedItem.Action.(func(*EditorState))
	if !ok {
		log.Printf("Invalid action for menu item '%s'\n", selectedItem.Name)
		return
	}
	actionFunc(state)
}

// MoveMenuSelection moves the menu selection up or down with wraparound.
func MoveMenuSelection(state *EditorState, delta int) {
	numResults := len(state.menu.search.Results())
	if numResults == 0 {
		return
	}

	newIdx := (state.menu.selectedResultIdx + delta) % numResults
	if newIdx < 0 {
		newIdx = numResults + newIdx
	}

	state.menu.selectedResultIdx = newIdx
}

// AppendRuneToMenuSearch appends a rune to the menu search query.
func AppendRuneToMenuSearch(state *EditorState, r rune) {
	menu := state.menu
	menu.search.AppendRuneToQuery(r)
	menu.selectedResultIdx = 0
}

// DeleteRuneFromMenuSearch deletes a rune from the menu search query.
func DeleteRuneFromMenuSearch(state *EditorState) {
	menu := state.menu
	menu.search.DeleteRuneFromQuery()
	menu.selectedResultIdx = 0
}
