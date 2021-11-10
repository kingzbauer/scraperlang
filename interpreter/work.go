package interpreter

import (
	"errors"
	"net/http"
	"strings"
)

// This package contains the set of functions, structs that are related to creating http request jobs

// A set of errors defining the different runtime errors
var (
	ErrMissingURLScheme = errors.New("Missing a valid URL scheme")
)

type getWorkConfig struct {
	tag     string
	url     string
	headers map[string]interface{}
}

// newGetWork returns a unit of work that is created when we encounter a get expression. It is
// then dispatched to it's own goroutine
// The unit of work is responsible of calling wg.Done when it's done executing so as to allow the main
// interpreter goroutine to exit when all work is complete
//
// TODO: There are a couple todos here especially with regards to handling error paths
func (i *Interpreter) newGetWork(cfg getWorkConfig) func() {
	return func() {
		defer i.wg.Done()

		// Make sure the url has a valid scheme
		parts := strings.SplitN(cfg.url, ":", 2)
		if len(parts) != 2 {
			// TODO: Have an error handling mechanism
			return
		} else if !in(parts[0], []string{"http", "https"}) {
			// TODO: Have an error handling mechanism
			return
		}
		req, err := http.NewRequest(http.MethodGet, cfg.url, nil)
		if err != nil {
			return
		}
		headers := map[string][]string{}
		if cfg.headers != nil {
			for key, value := range cfg.headers {
				switch t := value.(type) {
				case string:
					headers[key] = []string{t}
				case []string:
					headers[key] = t
				}
			}
		}
		req.Header = headers

		if res, err := http.DefaultClient.Do(req); err == nil {
			env := NewEnvironment(map[string]interface{}{
				"status": res.StatusCode,
			}, nil)
			// TODO: This will be handled by the Resolver by doing a pre-semantic analysis
			if closure, ok := i.taggedClosures[cfg.tag]; ok {
				closure.Accept(i, env)
			}
		}
	}
}

func in(val string, array []string) bool {
	for _, entry := range array {
		if val == entry {
			return true
		}
	}
	return false
}
