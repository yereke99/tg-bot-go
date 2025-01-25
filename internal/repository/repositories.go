package repository

type IRedisRepository interface{}

type Repository struct{}

func NewRepository() *Repository {
	return &Repository{}
}
