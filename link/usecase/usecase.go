package usecase

import (
	"context"
	"fmt"
	"github.com/kneumoin/go-clean-architecture/link"
	"github.com/kneumoin/go-clean-architecture/models"
	"github.com/spf13/viper"
)

type LinkUseCase struct {
	linkRepo link.Repository
}

func NewLinkUseCase(linkRepo link.Repository) *LinkUseCase {
	return &LinkUseCase{
		linkRepo: linkRepo,
	}
}

func (b LinkUseCase) CreateLink(ctx context.Context, user *models.User, secret string) (*models.Link, error) {
	bm, err := b.linkRepo.CreateLink(ctx, user, secret)
	if err != nil {
		return nil, err
	}
	url := fmt.Sprintf("http://%s:%s/api/links/%s", viper.GetString("hostname"), viper.GetString("port"), bm.ID)
	bm.URL = url
	return bm, nil
}

func (b LinkUseCase) GetLink(ctx context.Context, user *models.User, id string) (*models.Link, error) {
	bm, err := b.linkRepo.GetLink(ctx, user, id)
	if err != nil {
		return nil, err
	}
	if err := b.linkRepo.DeleteLink(ctx, user, id); err != nil {
		return nil, err
	}
	return bm, nil
}
