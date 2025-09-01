package models

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"
)

type User struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	Name      string    `json:"name"`
	Email     string    `gorm:"uniqueIndex" json:"email"`
	Password  string    `json:"-"`
	CreatedAt time.Time `gorm:"not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"not null;default:CURRENT_TIMESTAMP" json:"updated_at"`
}

type Guide struct {
	ID           uint         `gorm:"primarykey" json:"id"`
	Name         string       `json:"name"`
	Expertise    string       `json:"expertise"`
	Availability Availability `gorm:"type:jsonb"`
	CreatedAt    time.Time    `gorm:"not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt    time.Time    `gorm:"not null;default:CURRENT_TIMESTAMP" json:"updated_at"`
}

type BookingStatus string

const (
	Pending   BookingStatus = "pending"
	Confirmed BookingStatus = "confirmed"
	Cancelled BookingStatus = "cancelled"
	Completed BookingStatus = "completed"
)

func (s *BookingStatus) Scan(value interface{}) error {
	if value == nil {
		*s = Pending
		return nil
	}

	switch v := value.(type) {
	case string:
		*s = BookingStatus(v)
	case []byte:
		*s = BookingStatus(string(v))
	default:
		return fmt.Errorf("cannot scan %T into BookingStatus", value)
	}
	return nil
}

func (s BookingStatus) Value() (driver.Value, error) {
	return string(s), nil
}

type Booking struct {
	ID        uint          `gorm:"primarykey" json:"id"`
	UserID    uint          `json:"user_id"`
	GuideID   uint          `json:"guide_id"`
	Datetime  time.Time     `json:"datetime"`
	Status    BookingStatus `json:"status" gorm:"type:booking_status;default:'pending'"`
	Notes     string        `json:"notes"`
	User      User          `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"-"`
	Guide     *Guide        `gorm:"foreignKey:GuideID;constraint:OnDelete:CASCADE" json:"-"`
	CreatedAt time.Time     `gorm:"not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time     `gorm:"not null;default:CURRENT_TIMESTAMP" json:"updated_at"`
}
type Availability []string

func (a Availability) MarshalJSON() ([]byte, error) {
	formatted := make([]string, len(a))
	for i, s := range a {
		t, err := time.Parse(time.RFC3339, s)
		if err == nil {
			formatted[i] = t.UTC().Format(time.RFC3339)
		} else {
			formatted[i] = s
		}
	}
	return json.Marshal(formatted)
}

func (a Availability) Value() (driver.Value, error) {
	return json.Marshal(a)
}

func (a *Availability) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to scan availability: %v", value)
	}
	return json.Unmarshal(bytes, a)
}

type UserResponse struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type GuideResponse struct {
	ID           uint         `json:"id"`
	Name         string       `json:"name"`
	Expertise    string       `json:"expertise"`
	Availability Availability `json:"availability"`
}

type BookingResponse struct {
	ID        uint          `json:"id"`
	UserID    uint          `json:"user_id"`
	GuideID   uint          `json:"guide_id"`
	Datetime  time.Time     `json:"datetime"`
	Status    BookingStatus `json:"booking_status"`
	Notes     string        `json:"notes"`
	Guide     GuideResponse `json:"guide"`
	CreatedAt time.Time     `json:"created_at"`
}
