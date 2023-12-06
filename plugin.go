// This is the filter_on_field plugin.  It contains the code that implements the
// Traefik interfaces, and filters requests based on the parameters passed in the configuration, allowing all values
// except those that match the disallowed content.
package filter_on_field

import (
	"net/http"
	"regexp"
)

// An error type that can be returned by the filter
type Error struct {
	// The error message
	msg string
	// The HTTP status code to return
	code int
}

// A configuration struct for the filter.  This is populated from the Traefik configuration, regardless
// of whether it's used via the dynamic configuration or the static configuration.
type Config struct {

	//  The name of the field to filter on
	fieldName string

	// The message to return when disallowed content is found.  This is returned as the Error Response body.
	responseMessage string

	// A group of regular expressions representing content that is not allowed
	disallowedContent []string
}

// The Filter On Field Struct that implements all of the Traefik interfaces in order to be used as a Traefik
// middleware plugin.
type FilterOnField struct {
	next   http.Handler
	config *Config
	name   string
}

func CreateConfig(fieldName string, responseMessage string, disallowedContent []string) *Config {
	return &Config{
		fieldName:         fieldName,
		responseMessage:   responseMessage,
		disallowedContent: disallowedContent,
	}
}

// Create a new instance of the filter.  This is called by Traefik when the plugin is loaded and returns a new instance of the filter.
// It takes a configuration struct, the next handler in the chain, and the name of the filter.
func New(ctx *Config, next http.Handler, name string) (http.Handler, error) {
	return &FilterOnField{
		next:   next,
		config: ctx,
		name:   name,
	}, nil
}

// The ServeHTTP method is called by Traefik when a request is received.  It takes a response writer and a request.  When the checks pass,
// it calls the next handler in the chain.  When the checks fail, it returns an error response.
func (f *FilterOnField) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	parameter := req.FormValue(f.config.fieldName)
	if parameter == "" {
		f.next.ServeHTTP(rw, req)
	}
	for _, pattern := range f.config.disallowedContent {
		regex := regexp.MustCompile(pattern)

		if regex.MatchString(parameter) {
			rw.WriteHeader(http.StatusBadRequest)
			rw.Write([]byte(f.config.responseMessage))
			return
		}
	}
	f.next.ServeHTTP(rw, req)
}
