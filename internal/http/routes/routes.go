package routes

// ====================
// User Routes
// ====================
const (
	HealthRoute = "/health"

	BooksRoute  = "/books"
	BookRoute   = "/book"
	BookIDRoute = "/book/:id"

	UsersRoute = "/users"
	UserRoute  = "/user/:id"

	ForgotPassword       = "/forgot-password"
	LoginRoute           = "/login"
	RegisterRoute        = "/register"
	VerifyEmailRoute     = "/verify-email"
	VerifyMobileRoute    = "/verify-mobile"
	ResendEmailRoute     = "/resend-email-verification"
	ResendMobileOTPRoute = "/resend-mobile-otp"
	ResetPasswordRoute   = "/reset-password"
)

// ====================
// Publisher Routes
// ====================
const (
	PublisherRoute       = "/publishers"
	PublisherByIDRoute   = "/publishers/:id"
	PublisherStatusRoute = "/publishers/:id/status"
)
