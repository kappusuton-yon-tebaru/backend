package job

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kappusuton-yon-tebaru/backend/internal/job"
	"github.com/kappusuton-yon-tebaru/backend/internal/models"
)

type Handler struct {
	service *job.Service
}

func NewHandler(service *job.Service) *Handler {
	return &Handler{
		service,
	}
}

func (h *Handler) GetAllJobParents(ctx *gin.Context) {
	pagination := models.NewPaginationWithDefault(1, 10)
	err := ctx.ShouldBindQuery(&pagination)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "pagination should be integer"})
		return
	}

	jobs, err := h.service.GetAllJobParents(ctx, pagination.WithMinimum(1, 10))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, jobs)
}

func (h *Handler) GetAllJobsByParentId(ctx *gin.Context) {
	id := ctx.Param("id")
	if len(id) == 0 {
		ctx.JSON(http.StatusBadRequest, map[string]any{
			"message": "empty id",
		})
		return
	}

	pagination := models.NewPaginationWithDefault(1, 10)
	err := ctx.ShouldBindQuery(&pagination)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "pagination should be integer"})
		return
	}

	jobs, werr := h.service.GetAllJobsByParentId(ctx, id, pagination.WithMinimum(1, 10))
	if werr != nil {
		ctx.JSON(werr.GetCodeOr(http.StatusInternalServerError), map[string]any{
			"message": werr.GetMessageOr("internal server error"),
		})
		return
	}

	ctx.JSON(http.StatusOK, jobs)
}

func (h *Handler) CreateJob(ctx *gin.Context) {
	var jobDTO job.CreateJobDTO

	if err := ctx.ShouldBindJSON(&jobDTO); err != nil {
		ctx.JSON(http.StatusBadRequest, map[string]any{
			"message": "invalid input",
			"error":   err.Error(),
		})
		return
	}

	id, err := h.service.CreateJob(ctx, jobDTO)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, map[string]any{
			"message": "failed to create job",
			"error":   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, map[string]any{
		"message": "job created successfully",
		"id":      id,
	})
}

func (h *Handler) DeleteJob(ctx *gin.Context) {
	id := ctx.Param("id")
	if len(id) == 0 {
		ctx.JSON(http.StatusBadRequest, map[string]any{
			"message": "empty id",
		})
		return
	}

	err := h.service.DeleteJob(ctx, id)
	if err != nil {
		ctx.JSON(err.GetCodeOr(http.StatusInternalServerError), map[string]any{
			"message": err.GetMessageOr("internal server error"),
		})
		return
	}

	ctx.JSON(http.StatusOK, map[string]any{
		"message": "deleted",
	})
}
