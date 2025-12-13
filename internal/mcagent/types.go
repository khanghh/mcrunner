package mcagent

type CASClient struct {
	LoginURL     string
	ValidateURL  string
	ClientID     string
	ClientSecret string
}

type UserInfo struct {
	UserID   string `json:"userId"`
	Username string `json:"username"`
	FullName string `json:"fullName"`
	Email    string `json:"email"`
	Picture  string `json:"picture,omitempty"`
}

type authenticationResponse struct {
	AuthenticationSuccess *struct {
		User *UserInfo `json:"user"`
	} `json:"authenticationSuccess,omitempty"`
	AuthenticationFailure *struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	} `json:"authenticationFailure,omitempty"`
}

type casErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type ServerInfo struct {
	Name          string    `json:"name"`
	Version       string    `json:"version"`
	TPS           []float64 `json:"tps"`
	PlayersOnline int       `json:"playersOnline"`
	PlayersMax    int       `json:"playersMax"`
}
