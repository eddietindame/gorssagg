package errors

type HandlerError string

const (
	Nil                HandlerError = ""
	ServerError        HandlerError = "Server error"
	SessionError       HandlerError = "Session error"
	LoginCredentials   HandlerError = "Invalid credentials"
	RegisterPassword   HandlerError = "Passwords do not match"
	RegisterUsername   HandlerError = "Invalid username"
	RegisterEmail      HandlerError = "Invalid email address"
	RegisterUserExists HandlerError = "User already exists"
	ForgotNotFound     HandlerError = "User not found"
	ForgotToken        HandlerError = "Failed to generate reset token"
	ForgotSend         HandlerError = "Failed to send reset email"
	ResetPassword      HandlerError = "Passwords do not match"
	ResetToken         HandlerError = "Invalid or expired token"
	ResetFailed        HandlerError = "Failed to reset password"
	FeedInvalid        HandlerError = "Invalid feed url"
	FeedFetch          HandlerError = "Couldn't fetch feed"
	FeedCreate         HandlerError = "Error creating feed"
	FeedRead           HandlerError = "Error reading feed"
	FollowCreate       HandlerError = "Error creating follow"
	FollowedFeedsRead  HandlerError = "Error fetching followed feeds"
)

func (err HandlerError) ToString() string {
	return string(err)
}

func (err HandlerError) ToFriendlyString() string {
	switch err {
	case ServerError,
		SessionError,
		ForgotToken,
		ForgotSend,
		ResetFailed,
		FeedRead,
		FeedCreate,
		FollowCreate,
		FollowedFeedsRead:
		return "Something went wrong"
	default:
		return err.ToString()
	}
}
