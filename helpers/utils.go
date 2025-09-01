package helpers

import (
	"backend/models"
	"errors"
	"strconv"
	"strings"
	"time"

	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
)

// ParseRFC3339Slot parses datetime string to time.Time
func ParseRFC3339Slot(datetime string) (time.Time, error) {
	t, err := time.Parse(time.RFC3339, datetime)
	if err != nil {
		return time.Time{}, errors.New("datetime must be RFC3339")
	}
	return t, nil
}

// Check if requested slot exists in guide availability
func IsSlotAvailable(guide *models.Guide, slot time.Time) bool {
	for _, s := range guide.Availability {
		if t, _ := time.Parse(time.RFC3339, s); t.Equal(slot) {
			return true
		}
	}
	return false
}

// DB Helpers

func GetUserByEmail(db *gorm.DB, email string) (*models.User, error) {
	var user models.User
	if err := db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func GetGuideByID(db *gorm.DB, guideID uint) (*models.Guide, error) {
	var guide models.Guide
	if err := db.First(&guide, guideID).Error; err != nil {
		return nil, err
	}
	return &guide, nil
}

func IsSlotBooked(db *gorm.DB, guideID uint, slot time.Time) bool {
	var count int64
	db.Model(&models.Booking{}).Where("guide_id = ? AND datetime = ? AND status <> ?", guideID, slot, "cancelled").Count(&count)
	return count > 0
}

func CreateBooking(db *gorm.DB, booking *models.Booking) error {
	return db.Create(booking).Error
}

func ListUserBookings(db *gorm.DB, userID uint) ([]models.BookingResponse, error) {
	var bookings []models.Booking
	if err := db.
		Preload("Guide").
		Where("user_id = ?", userID).
		Order("datetime desc").
		Find(&bookings).Error; err != nil {
		return nil, err
	}

	var responses []models.BookingResponse
	for _, b := range bookings {
		responses = append(responses, models.BookingResponse{
			ID:       b.ID,
			UserID:   b.UserID,
			GuideID:  b.GuideID,
			Datetime: b.Datetime,
			Status:   b.Status,
			Guide: models.GuideResponse{
				ID:           b.Guide.ID,
				Name:         b.Guide.Name,
				Expertise:    b.Guide.Expertise,
				Availability: b.Guide.Availability,
			},
			Notes:     b.Notes,
			CreatedAt: b.CreatedAt,
		})
	}

	return responses, nil
}

func ListGuides(db *gorm.DB, pageParam, sizeParam, expertise string) ([]models.Guide, []string, gin.H, error) {
	var (
		guides        []models.Guide
		expertiseList []string
		total         int64
	)

	page, err := strconv.Atoi(pageParam)
	if err != nil || page < 1 {
		page = 1
	}
	size, err := strconv.Atoi(sizeParam)
	if err != nil || size < 1 {
		size = 10
	}

	offset := (page - 1) * size

	if strings.ToLower(expertise) == "all" {
		expertise = ""
	}

	query := db.Model(&models.Guide{}).Select("id", "name", "expertise", "availability")

	if expertise != "" {
		query = query.Where("LOWER(expertise) = ?", strings.ToLower(expertise))
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, nil, nil, err
	}

	if err := query.Offset(offset).Limit(size).Find(&guides).Error; err != nil {
		return nil, nil, nil, err
	}

	for i := range guides {
		avail := guides[i].Availability
		newAvail := []string{}
		for j := range avail {
			slotTime, err := time.Parse(time.RFC3339, avail[j])
			if err != nil {
				continue
			}
			if !IsSlotBooked(db, guides[i].ID, slotTime) {
				newAvail = append(newAvail, avail[j])
			}
		}
		guides[i].Availability = newAvail
	}

	// Fetch distinct expertise list for filter options
	if err := db.Model(&models.Guide{}).Distinct().Pluck("expertise", &expertiseList).Error; err != nil {
		return nil, nil, nil, err
	}

	// Build pagination metadata
	pagination := gin.H{
		"total": total,
		"page":  page,
		"size":  size,
	}

	return guides, expertiseList, pagination, nil
}

func SendErrorResponse(c *gin.Context, status int, errMsg string) {
	response := gin.H{
		"status":  "error",
		"message": errMsg,
	}
	c.IndentedJSON(status, response)
}

func SendSuccessResponse(c *gin.Context, status int, data interface{}) {
	response := gin.H{
		"status": "success",
		"data":   data,
	}
	c.IndentedJSON(status, response)
}
