package drafts

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/Gealber/limengo/repositories/models"

	"github.com/gin-gonic/gin"
)

const (
	peopleHeader = "X-People"
)

type draftsRepository interface {
	List(id int) ([]models.Draft, error)
	Get(id int, context string) (*models.Draft, error)
	Create(draft models.Draft) error
	Delete(id int, context string) error
	Update(draft models.Draft) (*models.Draft, error)
}

type draftsController struct {
	repo draftsRepository
}

func New(repo draftsRepository) *draftsController {
	return &draftsController{repo: repo}
}

func (ctr *draftsController) List(c *gin.Context) {
	peopleID := c.Request.Header.Get(peopleHeader)
	id, err := strconv.Atoi(peopleID)
	if err != nil {
		c.JSON(http.StatusPreconditionFailed, gin.H{"message": "people id absent or invalid"})

		return
	}

	drafts, err := ctr.repo.List(id)
	if err != nil {
		if errors.Is(err, models.NotFoundErr) {
			c.JSON(http.StatusNotFound, gin.H{"message": "resource not found"})

			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"message": "internal error"})

		return
	}

	c.JSON(http.StatusOK, drafts)
}

func (ctr *draftsController) Post(c *gin.Context) {
	peopleID := c.Request.Header.Get(peopleHeader)
	id, err := strconv.Atoi(peopleID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "bad request"})

		return
	}

	var request CreateDraftRequest
	err = c.BindJSON(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "bad request"})

		return
	}

	draft := models.Draft{
		ID:        id,
		Type:      request.Type,
		Context:   request.Context,
		Data:      request.Data,
		CreatedAt: time.Now().Format(time.RFC3339),
		Iri:       fmt.Sprintf("/front/drafts/%s", request.Context),
	}

	err = ctr.repo.Create(draft)
	if err != nil {
		if errors.Is(err, models.DuplicateValueErr) {
			c.JSON(422, gin.H{"message": "resource already exists"})

			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"message": "internal error"})

		return
	}

	c.JSON(201, draft)
}

func (ctr *draftsController) Get(c *gin.Context) {
	peopleID := c.Request.Header.Get(peopleHeader)
	id, err := strconv.Atoi(peopleID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "bad request"})
		return
	}

	context := c.Param("context")

	draft, err := ctr.repo.Get(id, context)
	if err != nil {
		if errors.Is(err, models.NotFoundErr) {
			c.JSON(http.StatusNotFound, gin.H{"message": "resource not found"})

			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"message": "internal error"})

		return
	}

	c.JSON(http.StatusOK, draft)
}

func (ctr *draftsController) Put(c *gin.Context) {
	peopleID := c.Request.Header.Get(peopleHeader)
	id, err := strconv.Atoi(peopleID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "bad request"})

		return
	}

	context := c.Param("context")

	var request UpdateDraftRequest
	err = c.BindJSON(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "bad request"})
		return
	}

	draft := models.Draft{
		ID:      id,
		Type:    request.Type,
		Context: context,
		Data:    request.Data,
	}

	updatedDraft, err := ctr.repo.Update(draft)
	if err != nil {
		if errors.Is(err, models.NotFoundErr) {
			c.JSON(http.StatusNotFound, gin.H{"message": "resource not found"})

			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"message": "internal error"})

		return
	}

	c.JSON(http.StatusOK, updatedDraft)
}

func (ctr *draftsController) Delete(c *gin.Context) {
	peopleID := c.Request.Header.Get(peopleHeader)
	id, err := strconv.Atoi(peopleID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "bad request"})

		return
	}

	context := c.Param("context")

	err = ctr.repo.Delete(id, context)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "internal error"})

		return
	}

	c.JSON(http.StatusNoContent, nil)
}
