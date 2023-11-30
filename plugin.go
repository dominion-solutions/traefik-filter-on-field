package filter_on_field

import (
	"net/http"
	"regexp"
)

type Error struct {
	msg  string
	code int
}

type Config struct {
	/**
	 * The name of the field to filter on
	 */
	fieldName string
	/**
	 * The message to return when disallowed content is found.  Defaults to "Disallowed content"
	 */
	responseMessage string
	/**
	 * A group of regular expressions representing content that is not allowed
	 */
	disallowedContent []string
}

type FilterOnField struct {
	next   http.Handler
	config *Config
	name   string
}

/**
 * Create a new instance of the filter
 * @param ctx The configuration for the filter
 * @param next The next handler in the chain
 * @param name The name of the filter
 * @return A new instance of the filter
 */
func New(ctx *Config, next http.Handler, name string) (http.Handler, error) {
	return &FilterOnField{
		next:   next,
		config: ctx,
		name:   name,
	}, nil
}

/**
 * Apply the middleware to the request
 * @param rw The response writer
 * @param req The request
 */
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
