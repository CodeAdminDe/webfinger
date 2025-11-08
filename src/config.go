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
	"fmt"
	"os"
	"regexp"
)

type Config struct {
	Resource  string
	IssuerURL string
	AllowDomainWildcard bool
}

func NewConfig() (*Config, error) {
	resource := os.Getenv("WEBFINGER_RESOURCE")
	if resource == "" {
		return nil, fmt.Errorf("WEBFINGER_RESOURCE environment variable not set")
	}

	issuerURL := os.Getenv("WEBFINGER_ISSUER_URL")
	if issuerURL == "" {
		return nil, fmt.Errorf("WEBFINGER_ISSUER_URL environment variable not set")
	}

	allowDomainWildcardStr := os.Getenv("WEBFINGER_ALLOW_DOMAIN_WILDCARD")
	allowDomainWildcard := (allowDomainWildcardStr == "true") || (allowDomainWildcardStr == "TRUE")

	// Validate resource format: acct:user@domain.com
	resourceRegex := regexp.MustCompile(`^acct:[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	
	if !resourceRegex.MatchString(resource) {
		return nil, fmt.Errorf("WEBFINGER_RESOURCE is not in the format acct:user@domain.com")
	}

	// Validate issuer URL format
	urlRegex := regexp.MustCompile(`^https?://[^\s/$.?#][^\s]*$`)
	if !urlRegex.MatchString(issuerURL) {
		return nil, fmt.Errorf("WEBFINGER_ISSUER_URL is not a valid URL")
	}

	return &Config{
		Resource:            resource,
		IssuerURL:           issuerURL,
		AllowDomainWildcard: allowDomainWildcard,
	}, nil 
}