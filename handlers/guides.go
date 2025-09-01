package handlers

import (
	"backend/helpers"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type GuideHandler struct {
	DB *gorm.DB
}

func NewGuideHandler(db *gorm.DB) *GuideHandler {
	return &GuideHandler{DB: db}
}

func (h *GuideHandler) ListGuides(context *gin.Context) {
	page := context.DefaultQuery("page", "1")
	size := context.DefaultQuery("size", "10")
	expertise := context.Query("expertise")

	guides, expertiseList, pagination, err := helpers.ListGuides(h.DB, page, size, expertise)
	if err != nil {
		helpers.SendErrorResponse(context, http.StatusInternalServerError, "Failed to fetch guides or expertise")
		return
	}

	response := gin.H{
		"guides":     guides,
		"expertises": expertiseList,
		"pagination": pagination,
	}

	helpers.SendSuccessResponse(context, http.StatusOK, response)
}
