package ui

import (
	"fmt"

	sptfy "github.com/Harshitk-cp/termtify/internal/spotify"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/sirupsen/logrus"
	"github.com/zmb3/spotify/v2"
)

type App struct {
	*tview.Application
	spotifyClient *spotify.Client
	layout        *tview.Flex
	searchBar     *tview.InputField
	sidebar       *tview.List
	content       *tview.List
	statusBar     *tview.TextView
	trackIDs      map[int]spotify.ID
}

func NewApp(client *spotify.Client) *App {
	app := &App{
		Application:   tview.NewApplication(),
		spotifyClient: client,
		layout:        tview.NewFlex().SetDirection(tview.FlexRow),
		searchBar:     tview.NewInputField().SetLabel("Search: "),
		sidebar:       tview.NewList().ShowSecondaryText(false),
		content:       tview.NewList().ShowSecondaryText(true),
		statusBar:     tview.NewTextView().SetTextAlign(tview.AlignCenter),
		trackIDs:      make(map[int]spotify.ID),
	}

	app.setupLayout()
	app.setupSidebar()
	app.setupContent()
	app.setupSearchBar()
	app.setupNavigation()

	return app
}
func (a *App) setupLayout() {
	a.layout.AddItem(a.searchBar, 1, 0, false)
	mainArea := tview.NewFlex().
		AddItem(a.sidebar, 0, 1, true).
		AddItem(a.content, 0, 4, false)
	a.layout.AddItem(mainArea, 0, 1, true)
	a.layout.AddItem(a.statusBar, 1, 0, false)

	a.SetRoot(a.layout, true)
}

func (a *App) setupSidebar() {
	a.sidebar.SetBorder(true).SetTitle("Your Library")
	a.loadPlaylists()
}

func (a *App) setupContent() {
	a.content.SetBorder(true).SetTitle("Welcome")
	a.content.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEnter {
			a.playSelectedTrack()
		}
		return event
	})
	a.updateWelcomeMessage()
}

func (a *App) playSelectedTrack() {
	if a.content.GetItemCount() == 0 {
		return
	}

	index := a.content.GetCurrentItem()
	mainText, _ := a.content.GetItemText(index)

	if trackID, ok := a.trackIDs[index]; ok {
		err := sptfy.PlayTrack(a.spotifyClient, trackID)
		if err != nil {
			logrus.Errorf("Failed to play track: %v", err)
			a.SetStatus(fmt.Sprintf("Failed to play: %s", mainText))
			return
		}
		a.SetStatus(fmt.Sprintf("Now playing: %s", mainText))
	}
}

func (a *App) setupSearchBar() {
	a.searchBar.SetFieldBackgroundColor(tcell.ColorWhite).
		SetFieldTextColor(tcell.ColorWhite).
		SetLabelColor(tcell.ColorGreen)
	a.searchBar.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEnter {
			query := a.searchBar.GetText()
			a.performSearch(query)
		}
	})
}

func (a *App) loadPlaylists() {
	playlists, err := sptfy.GetUserPlaylists(a.spotifyClient)
	if err != nil {
		logrus.Errorf("Failed to load playlists: %v", err)
		return
	}

	a.sidebar.Clear()
	for _, playlist := range playlists.Playlists {
		playlistID := playlist.ID
		a.sidebar.AddItem(playlist.Name, "", 0, func() {
			a.showPlaylistContent(playlistID)
		})
	}
}

func (a *App) updateWelcomeMessage() {
	a.content.Clear()
	welcomeMsg := "Welcome to Termtify!\n\nUse the sidebar to navigate or the search bar to find tracks."
	user, err := sptfy.GetCurrentUser(a.spotifyClient)
	if err == nil {
		welcomeMsg = fmt.Sprintf("Welcome, %s!\n\nUse the sidebar to navigate or the search bar to find tracks.", user.DisplayName)
	} else {
		logrus.Errorf("Failed to get current user: %v", err)
	}
	a.content.AddItem(welcomeMsg, "", 0, nil)
}

func (a *App) performSearch(query string) {
	results, err := sptfy.Search(a.spotifyClient, query)
	if err != nil {
		logrus.Errorf("Failed to perform search: %v", err)
		return
	}

	a.content.Clear()
	a.trackIDs = make(map[int]spotify.ID)
	a.content.SetTitle(fmt.Sprintf("Search results for '%s'", query))

	for i, track := range results.Tracks.Tracks {
		title := fmt.Sprintf("%s - %s", track.Name, track.Artists[0].Name)
		a.content.AddItem(title, track.Album.Name, 0, nil)
		a.trackIDs[i] = track.ID
	}
}

func (a *App) showPlaylistContent(playlistID spotify.ID) {
	tracks, err := sptfy.GetPlaylistTracks(a.spotifyClient, playlistID)
	if err != nil {
		logrus.Errorf("Failed to load playlist tracks: %v", err)
		return
	}

	a.content.Clear()
	a.trackIDs = make(map[int]spotify.ID)
	a.content.SetTitle("Tracks in playlist")

	for i, item := range tracks.Items {
		track := item.Track.Track
		title := fmt.Sprintf("%s - %s", track.Name, track.Artists[0].Name)
		a.content.AddItem(title, track.Album.Name, 0, nil)
		a.trackIDs[i] = track.ID
	}
}

func (a *App) SetStatus(message string) {
	a.statusBar.SetText(message)
}

func (a *App) Run() error {
	return a.Application.Run()
}
