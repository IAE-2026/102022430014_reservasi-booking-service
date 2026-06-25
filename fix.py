import os, re

path = 'internal/delivery/rest/booking_handler.go'
with open(path, 'r', encoding='utf-8') as f:
    content = f.read()

content = content.replace(
    'type BookingHandler struct {\n\tbookingUsecase domain.BookingUsecase\n}',
    'type BookingHandler struct {\n\tbookingUsecase domain.BookingUsecase\n}\n\nfunc serviceMeta() *domain.Meta {\n\tserviceName := os.Getenv(\"IAE_SERVICE_NAME\")\n\tif serviceName == \"\" {\n\t\tserviceName = \"Reservasi-Service\"\n\t}\n\tapiVersion := os.Getenv(\"IAE_API_VERSION\")\n\tif apiVersion == \"\" {\n\t\tapiVersion = \"v1\"\n\t}\n\treturn &domain.Meta{\n\t\tServiceName: serviceName,\n\t\tApiVersion:  apiVersion,\n\t}\n}'
)
content = content.replace(
    '\"reservasi/internal/domain\"',
    '\"os\"\n\t\"reservasi/internal/domain\"'
)

content = content.replace(
    '\tr.POST(\"/bookings\", handler.CreateBooking)',
    '\tr.GET(\"/bookings\", handler.GetAllBookings)\n\tr.GET(\"/bookings/:id\", handler.GetBooking)\n\tr.POST(\"/bookings\", handler.CreateBooking)'
)

content = re.sub(
    r'(c\.JSON\(http\.Status(?:OK|Created),\s*domain\.SuccessResponse\{\n(?:\s+.*\n)+?)\s*\})',
    lambda m: m.group(1) + '\t\tMeta:    serviceMeta(),\n\t}' if 'Meta:' not in m.group(1) else m.group(0),
    content
)

new_handlers = '''
// GetAllBookings godoc
// @Summary Mengambil semua pesanan (Collection)
// @Description Mendapatkan daftar seluruh reservasi (kebutuhan kontrak API)
// @Tags Bookings
// @Produce json
// @Param X-IAE-KEY header string true "API Key"
// @Param Authorization header string false "Bearer JWT Token dari SSO"
// @Success 200 {object} domain.SuccessResponse{data=[]domain.Booking} "Success"
// @Failure 401 {object} domain.ErrorResponse "Unauthorized"
// @Security ApiKeyAuth
// @Router /api/v1/bookings [get]
func (h *BookingHandler) GetAllBookings(c *gin.Context) {
\tbookings, err := h.bookingUsecase.GetAllBookings()
\tif err != nil {
\t\tc.JSON(http.StatusInternalServerError, domain.ErrorResponse{
\t\t\tStatus:  "error",
\t\t\tMessage: err.Error(),
\t\t})
\t\treturn
\t}

\tc.JSON(http.StatusOK, domain.SuccessResponse{
\t\tStatus:  "success",
\t\tMessage: "Data daftar reservasi berhasil diambil",
\t\tData:    bookings,
\t\tMeta:    serviceMeta(),
\t})
}

// GetBooking godoc
// @Summary Mengambil detail pesanan (Resource)
// @Description Mendapatkan rincian spesifik satu reservasi berdasarkan ID
// @Tags Bookings
// @Produce json
// @Param X-IAE-KEY header string true "API Key"
// @Param Authorization header string false "Bearer JWT Token dari SSO"
// @Param id path string true "Booking ID (UUID)"
// @Success 200 {object} domain.SuccessResponse{data=domain.Booking} "Success"
// @Failure 404 {object} domain.ErrorResponse "Not Found"
// @Failure 401 {object} domain.ErrorResponse "Unauthorized"
// @Security ApiKeyAuth
// @Router /api/v1/bookings/{id} [get]
func (h *BookingHandler) GetBooking(c *gin.Context) {
\tbookingID := c.Param("id")

\tbooking, err := h.bookingUsecase.GetBookingByID(bookingID)
\tif err != nil {
\t\tc.JSON(http.StatusNotFound, domain.ErrorResponse{
\t\t\tStatus:  "error",
\t\t\tMessage: err.Error(),
\t\t})
\t\treturn
\t}

\tc.JSON(http.StatusOK, domain.SuccessResponse{
\t\tStatus:  "success",
\t\tMessage: "Data reservasi berhasil diambil",
\t\tData:    booking,
\t\tMeta:    serviceMeta(),
\t})
}
'''
if 'GetAllBookings' not in content[content.rfind('}'):]:
    content += new_handlers

with open(path, 'w', encoding='utf-8') as f:
    f.write(content)
