package url

import (
	"math/rand"
	netUrl "net/url"
	"time"
)

const (
	size    = 5
	symbols = "abcdefghijklmnopqr...STUVWXYZ1234567890_-+"
)

var repo Repository

func init() {
	rand.NewSource(time.Now().UnixNano())
}

type Url struct {
	Id          string    `json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	Destination string    `json:"destination"`
}

type Stats struct {
	Url    *Url `json:"url"`
	Clicks int  `json:"clicks"`
}

type Repository interface {
	IdExists(id string) bool
	FindById(id string) *Url
	FindByUrl(url string) *Url
	Save(url Url) error
	RegisterClick(id string)
	FindClicks(id string) int
}

func ConfigRepository(r Repository) {
	repo = r
}

func FindOrCreateUrl(destination string) (u *Url, new bool, err error) {
	if u = repo.FindByUrl(destination); u != nil {
		return u, false, nil
	}

	if _, err = netUrl.ParseRequestURI(destination); err != nil {
		return nil, false, err
	}

	url := Url{
		Id:          idGenerator(),
		CreatedAt:   time.Now(),
		Destination: destination,
	}

	repo.Save(url)

	return &url, true, nil
}

func idGenerator() string {
	newId := func() string {
		id := make([]byte, size)
		for i := range id {
			id[i] = symbols[rand.Intn(len(symbols))]
		}
		return string(id)
	}

	for {
		if id := newId(); !repo.IdExists(id) {
			return id
		}
	}
}

func Find(id string) *Url {
	return repo.FindById(id)
}

func RegisterClick(id string) {
	repo.RegisterClick(id)
}

func (u *Url) Stats() *Stats {
	clicks := repo.FindClicks(u.Id)

	return &Stats{u, clicks}
}
