package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/zmb3/spotify"
	"golang.org/x/oauth2"
)

const (
	SpotifyClientId     = "insert"
	SpotifyClientSecret = "insert"
	RedirectUri         = "http://localhost:3000/callback"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "<html><body>")
		fmt.Fprintf(w, "<h1>Welcome to My App</h1>")
		fmt.Fprintf(w, "<p>Please <a href='/login'>login with Spotify</a> to continue.</p>")
		fmt.Fprintf(w, "</body></html>")
	})

	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, getSpotifyAuthURL(), http.StatusFound)
	})

	http.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		token, err := getAccessToken(r.FormValue("code"))
		if err != nil {
			fmt.Fprintf(w, "Error getting access token: %v", err)
			return
		}

		client := spotify.Authenticator{}.NewClient(token)

		user, err := client.CurrentUser()
		if err != nil {
			fmt.Fprintf(w, "Error getting user info: %v", err)
			return
		}

		followedArtists, err := client.CurrentUsersFollowedArtists()
		if err != nil {
			fmt.Fprintf(w, "Error getting followed artists: %v", err)
			return
		}

		fmt.Fprintf(w, "<html><body>")
		fmt.Fprintf(w, "<h1>User Information</h1>")
		fmt.Fprintf(w, "<p>Display Name: %s</p>", user.DisplayName)
		fmt.Fprintf(w, "<p>Email: %s</p>", user.Email)
		fmt.Fprintf(w, "<h2>Followed Artists</h2>")
		fmt.Fprintf(w, "<ul>")
		for _, artist := range followedArtists.Artists {
			if releasedInLastWeek(client, artist) {
				fmt.Fprintf(w, "<li>%s</li>", getLastReleasedSong(client, artist))
			}
		}
		fmt.Fprintf(w, "</ul>")
		fmt.Fprintf(w, "</body></html>")
	})

	fmt.Println("Server is running at http://localhost:3000")
	http.ListenAndServe(":3000", nil)
}

func getSpotifyAuthURL() string {
	auth := spotify.NewAuthenticator(RedirectUri, spotify.ScopeUserReadPrivate, spotify.ScopeUserReadEmail, spotify.ScopeUserFollowRead)
	auth.SetAuthInfo(SpotifyClientId, SpotifyClientSecret)
	return auth.AuthURL("")
}

func getAccessToken(code string) (*oauth2.Token, error) {
	auth := spotify.NewAuthenticator(RedirectUri, spotify.ScopeUserReadPrivate, spotify.ScopeUserReadEmail, spotify.ScopeUserFollowRead)
	auth.SetAuthInfo(SpotifyClientId, SpotifyClientSecret)
	token, err := auth.Exchange(code)
	if err != nil {
		return nil, fmt.Errorf("error exchanging code for token: %v", err)
	}
	return token, nil
}

func releasedInLastWeek(client spotify.Client, artist spotify.FullArtist) bool {
	// Get all albums by the artist
	albums, err := client.GetArtistAlbums(artist.ID)
	if err != nil {
		fmt.Println("Error getting artist's albums:", err)
		return false
	}

	// Check release date of each album
	for _, album := range albums.Albums {
		releaseDateParse(album)
	}
	return false
}

func getLastReleasedSong(client spotify.Client, artist spotify.FullArtist) *string {
	// Get the artist's top tracks
	tracks, err := client.GetArtistsTopTracks(artist.ID, spotify.CountryUSA)
	if err != nil {
		fmt.Println("Error getting artist's top tracks:", err)
		return nil
	}

	// Iterate through the tracks to find the latest one released within the last week
	for _, track := range tracks {
		// Parse the release date of the track
		releaseDate, err := time.Parse("YY-DD-MM", track.Album.ReleaseDate) //HERE
		if err != nil {
			fmt.Printf("Error parsing release date: %v\n", err)
			continue
		}

		// Check if the track was released within the last week
		if releaseDate.After(time.Now().AddDate(0, 0, -7)) {
			return &track.Album.Name
		}
	}

	return nil
}

func releaseDateParse(album spotify.SimpleAlbum) bool {
	switch album.ReleaseDatePrecision {
	case "year":
		result, err := time.Parse("YYYY", album.ReleaseDate)
		if err != nil {
			fmt.Printf("Error parsing release date: %v\n", err)
		}
		return result.After(time.Now().AddDate(0, 0, -7))
	case "month":
		result, err := time.Parse("YYYY-MM", album.ReleaseDate)
		if err != nil {
			fmt.Printf("Error parsing release date: %v\n", err)
		}
		return result.After(time.Now().AddDate(0, 0, -7))
	case "day":
		result, err := time.Parse("YYYY-DD-MM", album.ReleaseDate)
		if err != nil {
			fmt.Printf("Error parsing release date: %v\n", err)
		}
		return result.After(time.Now().AddDate(0, 0, -7))
	default:
		return false
	}
}

func radarFinal(client spotify.Client, artist spotify.FullArtist) {}
