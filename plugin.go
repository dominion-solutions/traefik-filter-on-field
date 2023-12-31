// This is the filter_on_field plugin.  It contains the code that implements the
// Traefik interfaces, and filters requests based on the parameters passed in the configuration, allowing all values
// except those that match the disallowed content.
package traefik_filter_on_field

import (
	"context"
	"net/http"
	"regexp"
)

// A configuration struct for the filter.  This is populated from the Traefik configuration, regardless
// of whether it's used via the dynamic configuration or the static configuration.
type Config struct {

	//  The name of the field to filter on
	FieldName string `yaml:"fieldName" json:"fieldName" toml:"fieldName" redis:"fieldName"`

	// The message to return when disallowed content is found.  This is returned as the Error Response body.
	ResponseMessage string `yaml:"responseMessage" json:"responseMessage" toml:"responseMessage" redis:"responseMessage"`

	// A group of regular expressions representing content that is not allowed
	DisallowedContent []string `yaml:"disallowedContent" json:"disallowedContent" toml:"disallowedContent" redis:"disallowedContent"`
}

// The Filter On Field Struct that implements all of the Traefik interfaces in order to be used as a Traefik
// middleware plugin.
type FilterOnField struct {
	next   http.Handler
	config *Config
	name   string
	ctx    context.Context
}

// Create a default version of the configuration.
func CreateConfig() *Config {
	return &Config{
		ResponseMessage: "Disallowed content",
	}
}

// Create a new instance of the filter.  This is called by Traefik when the plugin is loaded and returns a new instance of the filter.
// It takes a configuration struct, the next handler in the chain, and the name of the filter.
func New(ctx context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	return &FilterOnField{
		next:   next,
		config: config,
		name:   name,
		ctx:    ctx,
	}, nil
}

// The ServeHTTP method is called by Traefik when a request is received.  It takes a response writer and a request.  When the checks pass,
// it calls the next handler in the chain.  When the checks fail, it returns an error response.
func (f *FilterOnField) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	parameter := req.FormValue(f.config.FieldName)
	if parameter == "" {
		f.next.ServeHTTP(rw, req)
	}
	for _, pattern := range f.config.DisallowedContent {
		regex := regexp.MustCompile(pattern)

		if regex.MatchString(parameter) {
			rw.WriteHeader(http.StatusBadRequest)
			rw.Write([]byte(f.config.ResponseMessage))
			return
		}
	}
	f.next.ServeHTTP(rw, req)
}
