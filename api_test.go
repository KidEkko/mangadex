package mangodex

import (
	"fmt"
	"os"
	"testing"
)

var client = NewDexClient()

func TestLogin(t *testing.T) {
	err := client.Auth.Login(os.Getenv("USERNAME"), os.Getenv("PASSWORD"))
	if err != nil {
		t.Error("Login failed.")
	}
	fmt.Printf("%v\n", client)
}

func TestGetLoggedUser(t *testing.T) {
	user, err := client.User.GetLoggedUser()
	if err != nil {
		t.Error("Getting user failed.")
	}
	t.Log(user)
}

func TestGetMangaList(t *testing.T) {
	params := &ListMangaParams{
		Limit:    100,
		Offset:   0,
		Includes: []string{AuthorRel},
	}
	// If it is a search, then we add the search term.
	_, err := client.Manga.GetMangaList(params)
	if err != nil {
		t.Errorf("Getting manga failed: %s\n", err.Error())
	}
}
