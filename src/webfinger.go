//###################################################################//
//#  (c) 2025 Frederic Roggon                                       #//
//#                                                                 #//
//#  Licensed under the terms of GNU AFFERO GENERAL PUBLIC LICENSE. #//
//#  The full terms are provided via LICENSE file which is based    #//
//#  in the root of the code repository.                            #//
//#                                                                 #//
//#  Author: Frederic Roggon <frederic.roggon@codeadmin.de>         #//
//###################################################################//
package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
)

// Regex to extract the domain from an acct: URI
var acctDomainRegex = regexp.MustCompile(`^acct:[a-zA-Z0-9._%+-]+@([a-zA-Z0-9.-]+\.[a-zA-Z]{2,})$`)

// extractDomain extracts the domain from an acct: URI.
// Returns the domain and nil on success, or an empty string and an error if the format is invalid.
func extractDomain(acctResource string) (string, error) {
	matches := acctDomainRegex.FindStringSubmatch(acctResource)
	if len(matches) < 2 {
		return "", fmt.Errorf("invalid acct resource format: %s", acctResource)
	}
	return matches[1], nil
}

// JRD is a JSON Resource Descriptor, as defined in RFC 7033.
// The structure is kept minimal for this application's purpose.
type JRD struct {
	Subject string   `json:"subject"`
	Links   []Link   `json:"links"`
}

// Link is a link within an JRD.
type Link struct {
	Rel  string `json:"rel"`
	Type string `json:"type,omitempty"`
	Href string `json:"href,omitempty"`
}

// webfingerHandler handles acutal webfinger requests. The user-provided `resource` query
// parameter is only used for string comparison and therefore not vulnerable to injection
// attacks. The response Content-Type is set to `application/jrd+json` to mitigate risk of XSS.
func webfingerHandler(w http.ResponseWriter, r *http.Request, cfg *Config) {
	resource := r.URL.Query().Get("resource")
	if resource == "" {
		http.Error(w, "Missing resource parameter", http.StatusBadRequest)
		return
	}

	if cfg.AllowDomainWildcard {
		requestedDomain, err := extractDomain(resource)
		if err != nil {
			http.Error(w, "Invalid resource format", http.StatusBadRequest)
			return
		}

		configuredDomain, err := extractDomain(cfg.Resource)
		if err != nil {
			// This should ideally not happen if cfg.Resource is validated at startup
			http.Error(w, "Internal server error: invalid configured resource", http.StatusInternalServerError)
			return
		}

		if requestedDomain != configuredDomain {
			http.Error(w, "Resource not found", http.StatusNotFound)
			return
		}
	} else {
		if resource != cfg.Resource {
			http.Error(w, "Resource not found", http.StatusNotFound)
			return
		}
	}

	jrd := JRD{
		Subject: resource,
		Links: []Link{
			{
				Rel:  "http://openid.net/specs/connect/1.0/issuer",
				Href: cfg.IssuerURL,
			},
		},
	}

	w.Header().Set("Content-Type", "application/jrd+json")
	json.NewEncoder(w).Encode(jrd)
}
