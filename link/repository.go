package link

import (
	"context"
	"github.com/kneumoin/go-clean-architecture/models"
)

type Repository interface {
	CreateLink(ctx context.Context, user *models.User, secret string) (*models.Link, error)
	GetLink(ctx context.Context, user *models.User, id string) (*models.Link, error)
	DeleteLink(ctx context.Context, user *models.User, id string) error
}
