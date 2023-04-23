// Create and maintain by Chaiyapong Lapliengtrakul (chaiyapong@3dsinteractive.com), All right reserved (2021 - Present)
package microservices

import "time"

// IContext is the context for service
type IContext interface {
	Log(message string)
	Param(name string) string
	QueryParam(name string) string
	Response(responseCode int, responseData interface{})
	ReadInput() string
	ReadInputs() []string

	// Time
	Now() time.Time

	Requester(baseURL string, timeout time.Duration) IRequester
	WrapError(errIn error, errOut error) error
}
