// Copyright (c) Ayedo Cloud GmbH
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreateAndGetHost(t *testing.T) {
	t.Parallel()
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v1/hosts/", func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer test-token" {
			t.Errorf("missing bearer token")
		}
		if r.Method == http.MethodPost {
			var body map[string]any
			if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
				t.Fatal(err)
			}
			if body["name"] != "web-1" {
				t.Errorf("name=%v", body["name"])
			}
			_ = json.NewEncoder(w).Encode(map[string]any{
				"id":               "11111111-1111-1111-1111-111111111111",
				"name":             "web-1",
				"kind":             "vm",
				"role":             "generic",
				"bootstrap_status": "pending",
				"default_ipv4":     nil,
				"product":          map[string]any{"id": "p1"},
				"workspace":        map[string]any{"id": "w1"},
			})
			return
		}
		http.NotFound(w, r)
	})
	srv := httptest.NewServer(mux)
	t.Cleanup(srv.Close)

	c := New(srv.URL, "test-token")
	h, err := c.CreateHost(context.Background(), map[string]any{
		"name":         "web-1",
		"workspace_id": "w1",
		"product_id":   "p1",
		"kind":         "vm",
		"role":         "generic",
	})
	if err != nil {
		t.Fatal(err)
	}
	if h.ID == "" || h.Name != "web-1" || h.Role != "generic" {
		t.Fatalf("unexpected host %+v", h)
	}
}

func TestNotFound(t *testing.T) {
	t.Parallel()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	t.Cleanup(srv.Close)
	c := New(srv.URL, "t")
	_, err := c.GetHost(context.Background(), "missing")
	if _, ok := err.(*NotFoundError); !ok {
		t.Fatalf("want NotFoundError, got %v", err)
	}
}
