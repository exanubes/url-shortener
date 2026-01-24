package dynamodb

type Repository struct {
	client *client
}

func NewRepository(client *client) *Repository {
	return &Repository{client: client}
}
