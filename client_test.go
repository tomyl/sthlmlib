package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewClientLogsInWithPassword(t *testing.T) {
	var requestBody LoginRequest

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("method = %s, want %s", r.Method, http.MethodPost)
		}
		if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
			t.Fatalf("could not decode request: %v", err)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"data":{"loginPatron":{"status":"OK","statusMessage":null}}}`))
	}))
	defer server.Close()

	previousEndpoint := apiEndpoint
	apiEndpoint = server.URL
	defer func() {
		apiEndpoint = previousEndpoint
	}()

	if _, err := NewClient("198102030405", "secret"); err != nil {
		t.Fatalf("NewClient returned error: %v", err)
	}

	if requestBody.Variables.Operation != "loginPatron" {
		t.Fatalf("operation = %q, want %q", requestBody.Variables.Operation, "loginPatron")
	}
	if requestBody.Variables.LoginInput.CardNumber != "198102030405" {
		t.Fatalf("cardNumber = %q", requestBody.Variables.LoginInput.CardNumber)
	}
	if requestBody.Variables.LoginInput.Password != "secret" {
		t.Fatalf("password = %q", requestBody.Variables.LoginInput.Password)
	}
}

func TestLoginInputKeepsServerPinCodeFieldName(t *testing.T) {
	data, err := json.Marshal(LoginInput{
		CardNumber: "198102030405",
		Password:   "secret",
	})
	if err != nil {
		t.Fatalf("could not marshal login input: %v", err)
	}

	want := `{"cardNumber":"198102030405","pinCode":"secret"}`
	if string(data) != want {
		t.Fatalf("login input JSON = %s, want %s", data, want)
	}
}
