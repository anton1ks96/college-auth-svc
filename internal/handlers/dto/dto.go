package dto

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required,min=6"`
}

type UserInfo struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Role     string `json:"role"`
}

type AppUserInfo struct {
	ID            string `json:"id"`
	Username      string `json:"username"`
	Role          string `json:"role"`
	AcademicGroup string `json:"academic_group,omitempty"`
	Profile       string `json:"profile,omitempty"`
	Subgroup      string `json:"subgroup,omitempty"`
	EnglishGroup  string `json:"english_group,omitempty"`
}

type AppSignInResponse struct {
	AccessToken      string      `json:"access_token"`
	RefreshToken     string      `json:"refresh_token"`
	AccessExpiresIn  int         `json:"access_expires_in"`
	RefreshExpiresIn int         `json:"refresh_expires_in"`
	User             AppUserInfo `json:"user"`
}

type AppRefreshRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type AppRefreshResponse struct {
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
}

type AppSignOutRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type AppValidateResponse struct {
	Valid bool        `json:"valid"`
	User  AppUserInfo `json:"user,omitempty"`
}

type StudentSearchRequest struct {
	Query string `json:"query"`
}

type AppGetAccessRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type AppGetAccessTokenResponse struct {
	AccessToken string      `json:"access_token"`
	ExpiresIn   int         `json:"expires_in"`
	User        AppUserInfo `json:"user"`
}
