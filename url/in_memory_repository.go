package url

type inMemoryRepository struct {
	urls map[string]*Url
}

func NewRopositoryInMemory() *inMemoryRepository {
	return &inMemoryRepository{make(map[string]*Url)}
}

func FindByUrl(destination string) (u *string) {
	return u
}

func (r *inMemoryRepository) IdExists(id string) bool {
	_, exists := r.urls[id]

	return exists
}

func (r *inMemoryRepository) FindById(id string) *Url {
	return r.urls[id]
}

func (r *inMemoryRepository) FindByUrl(url string) *Url {
	for _, u := range r.urls {
		if u.Destination == url {
			return u
		}
	}

	return nil
}

func (r *inMemoryRepository) Save(url Url) error {
	r.urls[url.Id] = &url

	return nil
}
