// Create and maintain by Chaiyapong Lapliengtrakul (chaiyapong@3dsinteractive.com), All right reserved (2021 - Present)
package gcloud

const (
	GCloudOAuthState string = "pam-oauth"
)

const (
	ErrInvalidOAuthToken   string = "invalid oauth token"
	ErrInvalidOAuthConfig  string = "invalid oauth config"
	ErrInvalidGCloudConfig string = "invalid gcloud config"
)

var (
	SheetsColumnLetters []string = []string{
		"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z",
	}
	SheetsMaxColumnLen int = len(SheetsColumnLetters) + (len(SheetsColumnLetters) * len(SheetsColumnLetters))
)
