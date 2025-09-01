package handlers

import (
	"net/http"

	"backend/helpers"
	"backend/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type BookingHandler struct {
	DB *gorm.DB
}

func NewBookingHandler(db *gorm.DB) *BookingHandler {
	return &BookingHandler{DB: db}
}

func (h *BookingHandler) CreateBooking(context *gin.Context) {
	userID := context.GetUint("userId")

	var req struct {
		GuideID  uint   `json:"guide_id" binding:"required"`
		Datetime string `json:"datetime" binding:"required,datetime=2006-01-02T15:04:05Z07:00"`
		Notes    string `json:"notes"`
	}

	if err := context.ShouldBindJSON(&req); err != nil {
		helpers.SendErrorResponse(context, http.StatusBadRequest, "Invalid request format. Please check your input.")
		return
	}

	if err := helpers.ValidateCreateBookingInput(req.GuideID, req.Datetime, req.Notes); err != nil {
		helpers.SendErrorResponse(context, http.StatusBadRequest, err.Error())
		return
	}

	if req.GuideID == 0 || req.Datetime == "" {
		helpers.SendErrorResponse(context, http.StatusBadRequest, "Guide ID and Datetime are required.")
		return
	}

	guide, err := helpers.GetGuideByID(h.DB, req.GuideID)
	if err != nil {
		helpers.SendErrorResponse(context, http.StatusNotFound, "Guide with the specified ID not found.")
		return
	}

	slot, err := helpers.ParseRFC3339Slot(req.Datetime)
	if err != nil {
		helpers.SendErrorResponse(context, http.StatusBadRequest, "Invalid datetime format. Please use RFC3339 format.")
		return
	}

	if !helpers.IsSlotAvailable(guide, slot) {
		context.JSON(http.StatusBadRequest, gin.H{"error": "The guide is not available at the requested time slot."})
		return
	}

	if helpers.IsSlotBooked(h.DB, guide.ID, slot) {
		helpers.SendErrorResponse(context, http.StatusConflict, "This time slot has already been booked.")
		return
	}

	booking := models.Booking{
		UserID:   userID,
		GuideID:  guide.ID,
		Datetime: slot,
		Status:   "pending",
		Notes:    req.Notes,
		Guide:    guide,
	}

	if err := helpers.CreateBooking(h.DB, &booking); err != nil {
		helpers.SendErrorResponse(context, http.StatusInternalServerError, "Failed to create booking. Please try again later.")
		return
	}

	guideData := models.GuideResponse{
		ID:           guide.ID,
		Name:         guide.Name,
		Expertise:    guide.Expertise,
		Availability: guide.Availability,
	}

	response := models.BookingResponse{
		ID:        booking.ID,
		UserID:    booking.UserID,
		GuideID:   booking.GuideID,
		Datetime:  booking.Datetime,
		Status:    booking.Status,
		Notes:     booking.Notes,
		Guide:     guideData,
		CreatedAt: booking.CreatedAt,
	}

	helpers.SendSuccessResponse(context, http.StatusCreated, response)
}

func (h *BookingHandler) ListUserBookings(context *gin.Context) {
	userID := context.GetUint("userId")
	bookings, err := helpers.ListUserBookings(h.DB, userID)
	if err != nil {
		helpers.SendErrorResponse(context, http.StatusInternalServerError, "Failed to retrieve bookings. Please try again later.")
		return
	}

	if bookings == nil {
		bookings = []models.BookingResponse{}
	}

	helpers.SendSuccessResponse(context, http.StatusOK, bookings)
}
