
package api

import (
	"encoding/json"
	"fmt"
	"os"
)

// Paginator keeps track of pagination configuration for an API operation.
type Paginator struct {
	InputTokens  interface{} `json:"input_token"`
	OutputTokens interface{} `json:"output_token"`
	LimitKey     string      `json:"limit_key"`
	MoreResults  string      `json:"more_results"`
}

// InputTokensString returns output tokens formatted as a list
func (p *Paginator) InputTokensString() string {
	str := p.InputTokens.([]string)
	return fmt.Sprintf("%#v", str)
}

// OutputTokensString returns output tokens formatted as a list
func (p *Paginator) OutputTokensString() string {
	str := p.OutputTokens.([]string)
	return fmt.Sprintf("%#v", str)
}

// used for unmarshaling from the paginators JSON file
type paginationDefinitions struct {
	*API
	Pagination map[string]Paginator
}

// AttachPaginators attaches pagination configuration from filename to the API.
func (a *API) AttachPaginators(filename string) error {
	p := paginationDefinitions{API: a}

	f, err := os.Open(filename)
	defer f.Close()
	if err != nil {
		return err
	}
	err = json.NewDecoder(f).Decode(&p)
	if err != nil {
		return fmt.Errorf("failed to decode %s, err: %v", filename, err)
	}

	return p.setup()
}

// setup runs post-processing on the paginator configuration.
func (p *paginationDefinitions) setup() error {
	for n, e := range p.Pagination {
		if e.InputTokens == nil || e.OutputTokens == nil {
			continue
		}
		if _, ok := p.Operations[n]; !ok {
			continue
		}
		paginator := e

		switch t := paginator.InputTokens.(type) {
		case string:
			paginator.InputTokens = []string{t}
		case []interface{}:
			toks := []string{}
			for _, e := range t {
				s := e.(string)
				toks = append(toks, s)
			}
			paginator.InputTokens = toks
		}
		switch t := paginator.OutputTokens.(type) {
		case string:
			paginator.OutputTokens = []string{t}
		case []interface{}:
			toks := []string{}
			for _, e := range t {
				s := e.(string)
				toks = append(toks, s)
			}
			paginator.OutputTokens = toks
		}

		p.Operations[n].Paginator = &paginator
	}

	return nil
}

func enableStopOnSameToken(service string) bool {
	switch service {
	case "cloudwatchlogs":
		return true
	default:
		return false
	}
}
