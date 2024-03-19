package http

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/kneumoin/go-clean-architecture/auth"
	"github.com/kneumoin/go-clean-architecture/link"
	"github.com/kneumoin/go-clean-architecture/models"
	"net/http"
)

type Link struct {
	URL string `json:"url"`
}

type Secret struct {
	Secret string `json:"secret"`
}

type Handler struct {
	useCase link.UseCase
}

func NewHandler(useCase link.UseCase) *Handler {
	return &Handler{
		useCase: useCase,
	}
}

type createInput struct {
	Secret string `json:"secret"`
}

func (h *Handler) Create(c *gin.Context) {
	inp := new(createInput)
	if err := c.BindJSON(inp); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	user := c.MustGet(auth.CtxUserKey).(*models.User)

	bm, err := h.useCase.CreateLink(c.Request.Context(), user, inp.Secret)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, toLink(bm))
}

func toLink(b *models.Link) *Link {
	return &Link{
		URL: b.URL,
	}
}

func toSecret(b *models.Link) *Secret {
	return &Secret{
		Secret: b.Secret,
	}
}

func (h *Handler) Get(c *gin.Context) {
	user := c.MustGet(auth.CtxUserKey).(*models.User)

	bm, err := h.useCase.GetLink(c.Request.Context(), user, c.Param("id"))
	if errors.Is(err, link.ErrLinkNotFound) {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, toSecret(bm))
}
