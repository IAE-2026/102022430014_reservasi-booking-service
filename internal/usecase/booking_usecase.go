package usecase

import (
	"reservasi/internal/domain"
)

type bookingUsecase struct {
	bookingRepo domain.BookingRepository
}

// NewBookingUsecase membuat instance usecase baru
func NewBookingUsecase(bookingRepo domain.BookingRepository) domain.BookingUsecase {
	return &bookingUsecase{
		bookingRepo: bookingRepo,
	}
}

func (u *bookingUsecase) GetAllBookings() ([]*domain.Booking, error) {
	return u.bookingRepo.GetAllBookings()
}

func (u *bookingUsecase) GetBookingByID(bookingID string) (*domain.Booking, error) {
	return u.bookingRepo.GetBookingByID(bookingID)
}
