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
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestNewConfig(t *testing.T) {
	// Test case 1: Valid environment variables
	os.Setenv("WEBFINGER_RESOURCE", "acct:user@example.com")
	os.Setenv("WEBFINGER_ISSUER_URL", "https://example.com")
	cfg, err := NewConfig()
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}
	if cfg.Resource != "acct:user@example.com" {
		t.Errorf("Expected resource to be 'acct:user@example.com', but got: %s", cfg.Resource)
	}
	if cfg.IssuerURL != "https://example.com" {
		t.Errorf("Expected issuer URL to be 'https://example.com', but got: %s", cfg.IssuerURL)
	}
	os.Unsetenv("WEBFINGER_RESOURCE")
	os.Unsetenv("WEBFINGER_ISSUER_URL")

	// Test case 2: Missing WEBFINGER_RESOURCE
	os.Setenv("WEBFINGER_ISSUER_URL", "https://example.com")
	_, err = NewConfig()
	if err == nil {
		t.Error("Expected an error for missing WEBFINGER_RESOURCE, but got nil")
	}
	os.Unsetenv("WEBFINGER_ISSUER_URL")

	// Test case 3: Invalid WEBFINGER_RESOURCE
	os.Setenv("WEBFINGER_RESOURCE", "invalid-resource")
	os.Setenv("WEBFINGER_ISSUER_URL", "https://example.com")
	_, err = NewConfig()
	if err == nil {
		t.Error("Expected an error for invalid WEBFINGER_RESOURCE, but got nil")
	}
	os.Unsetenv("WEBFINGER_RESOURCE")
	os.Unsetenv("WEBFINGER_ISSUER_URL")

	// Test case 4: Invalid WEBFINGER_ISSUER_URL
	os.Setenv("WEBFINGER_RESOURCE", "acct:user@example.com")
	os.Setenv("WEBFINGER_ISSUER_URL", "invalid-url")
	_, err = NewConfig()
	if err == nil {
		t.Error("Expected an error for invalid WEBFINGER_ISSUER_URL, but got nil")
	}
	os.Unsetenv("WEBFINGER_RESOURCE")
	os.Unsetenv("WEBFINGER_ISSUER_URL")

	// Test case 5: More invalid WEBFINGER_ISSUER_URL
	os.Setenv("WEBFINGER_RESOURCE", "acct:user@example.com")
	os.Setenv("WEBFINGER_ISSUER_URL", "http://.example.com")
	_, err = NewConfig()
	if err == nil {
		t.Error("Expected an error for invalid WEBFINGER_ISSUER_URL, but got nil")
	}
	os.Unsetenv("WEBFINGER_RESOURCE")
	os.Unsetenv("WEBFINGER_ISSUER_URL")

	// Test case 6: Invalid characters in WEBFINGER_RESOURCE
	os.Setenv("WEBFINGER_RESOURCE", "acct:us<er@example.com")
	os.Setenv("WEBFINGER_ISSUER_URL", "https://example.com")
	_, err = NewConfig()
	if err == nil {
		t.Error("Expected an error for invalid characters in WEBFINGER_RESOURCE, but got nil")
	}
	// Test case for WEBFINGER_ALLOW_DOMAIN_WILDCARD
	os.Setenv("WEBFINGER_RESOURCE", "acct:user@example.com")
	os.Setenv("WEBFINGER_ISSUER_URL", "https://example.com") // Ensure IssuerURL is set for these tests

	// Test case 7: WEBFINGER_ALLOW_DOMAIN_WILDCARD set to "true"
	os.Setenv("WEBFINGER_ALLOW_DOMAIN_WILDCARD", "true")
	cfg, err = NewConfig()
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}
	if !cfg.AllowDomainWildcard {
		t.Errorf("Expected AllowDomainWildcard to be true, but got: %v", cfg.AllowDomainWildcard)
	}
	os.Unsetenv("WEBFINGER_ALLOW_DOMAIN_WILDCARD")

	// Test case 8: WEBFINGER_ALLOW_DOMAIN_WILDCARD set to "TRUE"
	os.Setenv("WEBFINGER_ALLOW_DOMAIN_WILDCARD", "TRUE")
	cfg, err = NewConfig()
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}
	if !cfg.AllowDomainWildcard {
		t.Errorf("Expected AllowDomainWildcard to be true, but got: %v", cfg.AllowDomainWildcard)
	}
	os.Unsetenv("WEBFINGER_ALLOW_DOMAIN_WILDCARD")

	// Test case 9: WEBFINGER_ALLOW_DOMAIN_WILDCARD set to "false"
	os.Setenv("WEBFINGER_ALLOW_DOMAIN_WILDCARD", "false")
	cfg, err = NewConfig()
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}
	if cfg.AllowDomainWildcard {
		t.Errorf("Expected AllowDomainWildcard to be false, but got: %v", cfg.AllowDomainWildcard)
	}
	os.Unsetenv("WEBFINGER_ALLOW_DOMAIN_WILDCARD")

	// Test case 10: WEBFINGER_ALLOW_DOMAIN_WILDCARD not set
	os.Setenv("WEBFINGER_RESOURCE", "acct:user@example.com") // Ensure these are set
	os.Setenv("WEBFINGER_ISSUER_URL", "https://example.com")   // Ensure these are set
	cfg, err = NewConfig()
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}
	if cfg.AllowDomainWildcard {
		t.Errorf("Expected AllowDomainWildcard to be false, but got: %v", cfg.AllowDomainWildcard)
	}
	os.Unsetenv("WEBFINGER_RESOURCE") // Unset after this test case
	os.Unsetenv("WEBFINGER_ISSUER_URL") // Unset after this test case
}

func TestWebfingerHandlerAllowDomainWildcard(t *testing.T) {
	// Setup: AllowDomainWildcard = true
	cfg := &Config{
		Resource:            "acct:user@example.com",
		IssuerURL:           "https://example.com/issuer",
		AllowDomainWildcard: true,
	}

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		webfingerHandler(w, r, cfg)
	})

	tests := []struct {
		name           string
		resourceParam  string
		expectedStatus int
		expectedSubject string
	}{
		{
			name:           "Matching domain, different user",
			resourceParam:  "acct:anotheruser@example.com",
			expectedStatus: http.StatusOK,
			expectedSubject: "acct:anotheruser@example.com",
		},
		{
			name:           "Matching domain, same user",
			resourceParam:  "acct:user@example.com",
			expectedStatus: http.StatusOK,
			expectedSubject: "acct:user@example.com",
		},
		{
			name:           "Non-matching domain",
			resourceParam:  "acct:user@another.com",
			expectedStatus: http.StatusNotFound,
			expectedSubject: "", // Not applicable for error cases
		},
		{
			name:           "Invalid acct format",
			resourceParam:  "acct:invalid",
			expectedStatus: http.StatusBadRequest,
			expectedSubject: "", // Not applicable for error cases
		},
		{
			name:           "Missing resource parameter",
			resourceParam:  "",
			expectedStatus: http.StatusBadRequest,
			expectedSubject: "", // Not applicable for error cases
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest("GET", "/.well-known/webfinger?resource="+tt.resourceParam, nil)
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v", status, tt.expectedStatus)
			}

			if tt.expectedStatus == http.StatusOK {
				expected := `{"subject":"` + tt.expectedSubject + `","links":[{"rel":"http://openid.net/specs/connect/1.0/issuer","href":"https://example.com/issuer"}]}`
				actual := rr.Body.String()
				if len(actual) > 0 && actual[len(actual)-1] == '\n' {
					actual = actual[:len(actual)-1]
				}
				if actual != expected {
					t.Errorf("handler returned unexpected body: got %v want %v", actual, expected)
				}
			}
		})
	}

	// Test case for AllowDomainWildcard = false (default behavior)
	cfgFalse := &Config{
		Resource:            "acct:user@example.com",
		IssuerURL:           "https://example.com/issuer",
		AllowDomainWildcard: false,
	}
	handlerFalse := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		webfingerHandler(w, r, cfgFalse)
	})

	t.Run("AllowDomainWildcard=false: Non-matching user", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/.well-known/webfinger?resource=acct:anotheruser@example.com", nil)
		if err != nil {
			t.Fatal(err)
		}
		rr := httptest.NewRecorder()
		handlerFalse.ServeHTTP(rr, req)
		if status := rr.Code; status != http.StatusNotFound {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusNotFound)
		}
	})

	t.Run("AllowDomainWildcard=false: Matching user", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/.well-known/webfinger?resource=acct:user@example.com", nil)
		if err != nil {
				t.Fatal(err)
		}
		rr := httptest.NewRecorder()
		handlerFalse.ServeHTTP(rr, req)
		if status := rr.Code; status != http.StatusOK {
				t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}
	})
}

func TestWebfingerHandler(t *testing.T) {
	cfg := &Config{
		Resource:  "acct:user@example.com",
		IssuerURL: "https://example.com/issuer",
	}

	req, err := http.NewRequest("GET", "/.well-known/webfinger?resource=acct:user@example.com", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		webfingerHandler(w, r, cfg)
	})

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	expected := `{"subject":"acct:user@example.com","links":[{"rel":"http://openid.net/specs/connect/1.0/issuer","href":"https://example.com/issuer"}]}`
	// Trim newline from actual response
	actual := rr.Body.String()
	if len(actual) > 0 && actual[len(actual)-1] == '\n' {
		actual = actual[:len(actual)-1]
	}

	if actual != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			actual, expected)
	}

	t.Run("Resource param does not exist and returns 400", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/doesntmatter-but-does-not-have-resource-param", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()

		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				webfingerHandler(w, r, cfg)
			})

		handler.ServeHTTP(rr, req)
			
		if status := rr.Code; status != http.StatusBadRequest {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusBadRequest)
		}

		expected := "Missing resource parameter\n"
		if rr.Body.String() != expected {
			t.Errorf("handler returned unexpected body: got %v want %v",
				rr.Body.String(), expected)
		}
	})
}

func TestWebfingerHandler404(t *testing.T) {
	t.Run("Root returns not found 404", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "404 Not Found", http.StatusNotFound)
		})

		handler.ServeHTTP(rr, req)
		
			

		if status := rr.Code; status != http.StatusNotFound {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusNotFound)
		}

		expected := `404 Not Found`
		// Trim newline from actual response
		actual := rr.Body.String()
		if len(actual) > 0 && actual[len(actual)-1] == '\n' {
			actual = actual[:len(actual)-1]
		}

		if actual != expected {
			t.Errorf("handler returned unexpected body: got %v want %v",
				actual, expected)
		}
	})

	t.Run(".well-known returns not found 404", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/.well-known", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "404 Not Found", http.StatusNotFound)
		})

		handler.ServeHTTP(rr, req)
		
			

		if status := rr.Code; status != http.StatusNotFound {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusNotFound)
		}

		expected := `404 Not Found`
		// Trim newline from actual response
		actual := rr.Body.String()
		if len(actual) > 0 && actual[len(actual)-1] == '\n' {
			actual = actual[:len(actual)-1]
		}

		if actual != expected {
			t.Errorf("handler returned unexpected body: got %v want %v",
				actual, expected)
		}
	})

	t.Run("non-existing returns not found 404", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/anything-random", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "404 Not Found", http.StatusNotFound)
		})

		handler.ServeHTTP(rr, req)
		
		if status := rr.Code; status != http.StatusNotFound {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusNotFound)
		}

		expected := `404 Not Found`
		// Trim newline from actual response
		actual := rr.Body.String()
		if len(actual) > 0 && actual[len(actual)-1] == '\n' {
			actual = actual[:len(actual)-1]
		}

		if actual != expected {
			t.Errorf("handler returned unexpected body: got %v want %v",
				actual, expected)
		}
	})
}

func TestWebfingerHandlerHealth200(t *testing.T) {
	t.Run("_healthz returns 200 ok", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/_healthz", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "200 OK", http.StatusOK)
		})

		handler.ServeHTTP(rr, req)
		
		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusOK)
		}

		expected := `200 OK`
		// Trim newline from actual response
		actual := rr.Body.String()
		if len(actual) > 0 && actual[len(actual)-1] == '\n' {
			actual = actual[:len(actual)-1]
		}

		if actual != expected {
			t.Errorf("handler returned unexpected body: got %v want %v",
				actual, expected)
		}
	})
}
