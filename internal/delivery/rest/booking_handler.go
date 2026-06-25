package rest

import (
	"net/http"
	"os"
	"reservasi/internal/domain"

	"github.com/gin-gonic/gin"
)

type BookingHandler struct {
	bookingUsecase domain.BookingUsecase
}

func serviceMeta() *domain.Meta {
	serviceName := os.Getenv("IAE_SERVICE_NAME")
	if serviceName == "" {
		serviceName = "Reservasi-Service"
	}
	apiVersion := os.Getenv("IAE_API_VERSION")
	if apiVersion == "" {
		apiVersion = "v1"
	}
	return &domain.Meta{
		ServiceName: serviceName,
		ApiVersion:  apiVersion,
	}
}

// NewBookingHandler mendaftarkan endpoint untuk Booking Service
func NewBookingHandler(r gin.IRouter, us domain.BookingUsecase) {
	handler := &BookingHandler{
		bookingUsecase: us,
	}

	r.GET("/bookings", handler.GetAllBookings)
	r.GET("/bookings/:id", handler.GetBooking)
	r.POST("/bookings", handler.CreateBooking)
}

// CreateBooking godoc
// @Summary Membuat pesanan awal (mengunci kamar)
// @Description Endpoint untuk membuat reservasi awal dengan status LOCKED. Terintegrasi dengan SOAP Audit dan Message Broker.
// @Tags Bookings
// @Accept json
// @Produce json
// @Param X-IAE-KEY header string true "API Key"
// @Param Authorization header string true "Bearer JWT Token dari SSO"
// @Param Idempotency-Key header string false "Kunci unik untuk mencegah duplikasi request"
// @Param request body domain.CreateBookingRequest true "Payload Reservasi"
// @Success 201 {object} domain.SuccessResponse{data=domain.Booking} "Created"
// @Failure 400 {object} domain.ErrorResponse "Bad Request"
// @Failure 401 {object} domain.ErrorResponse "Unauthorized"
// @Router /api/v1/bookings [post]
func (h *BookingHandler) CreateBooking(c *gin.Context) {
	var req domain.CreateBookingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{
			Status:  "error",
			Message: "Payload tidak valid: " + err.Error(),
		})
		return
	}

	// Ambil Idempotency-Key dari header untuk mencegah duplikasi request
	req.IdempotencyKey = c.GetHeader("Idempotency-Key")

	booking, err := h.bookingUsecase.CreateBooking(&req)
	if err != nil {
		// Menggunakan StatusInternalServerError atau BadRequest tergantung err
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{
			Status:  "error",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, domain.SuccessResponse{
		Status:  "success",
		Message: "Pesanan awal berhasil dibuat (Kamar dikunci)",
		Data:    booking,
		Meta:    serviceMeta(),
	})
}





// GetAllBookings godoc
// @Summary Mengambil semua pesanan (Collection)
// @Description Mendapatkan daftar seluruh reservasi (kebutuhan kontrak API)
// @Tags Bookings
// @Produce json
// @Param X-IAE-KEY header string true "API Key"
// @Param Authorization header string true "Bearer JWT Token dari SSO"
// @Success 200 {object} domain.SuccessResponse{data=[]domain.Booking} "Success"
// @Failure 401 {object} domain.ErrorResponse "Unauthorized"
// @Security ApiKeyAuth
// @Router /api/v1/bookings [get]
func (h *BookingHandler) GetAllBookings(c *gin.Context) {
	bookings, err := h.bookingUsecase.GetAllBookings()
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrorResponse{
			Status:  "error",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, domain.SuccessResponse{
		Status:  "success",
		Message: "Data daftar reservasi berhasil diambil",
		Data:    bookings,
		Meta:    serviceMeta(),
	})
}

// GetBooking godoc
// @Summary Mengambil detail pesanan (Resource)
// @Description Mendapatkan rincian spesifik satu reservasi berdasarkan ID
// @Tags Bookings
// @Produce json
// @Param X-IAE-KEY header string true "API Key"
// @Param Authorization header string true "Bearer JWT Token dari SSO"
// @Param id path string true "Booking ID (UUID)"
// @Success 200 {object} domain.SuccessResponse{data=domain.Booking} "Success"
// @Failure 404 {object} domain.ErrorResponse "Not Found"
// @Failure 401 {object} domain.ErrorResponse "Unauthorized"
// @Security ApiKeyAuth
// @Router /api/v1/bookings/{id} [get]
func (h *BookingHandler) GetBooking(c *gin.Context) {
	bookingID := c.Param("id")

	booking, err := h.bookingUsecase.GetBookingByID(bookingID)
	if err != nil {
		c.JSON(http.StatusNotFound, domain.ErrorResponse{
			Status:  "error",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, domain.SuccessResponse{
		Status:  "success",
		Message: "Data reservasi berhasil diambil",
		Data:    booking,
		Meta:    serviceMeta(),
	})
}
