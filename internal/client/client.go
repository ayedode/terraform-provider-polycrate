// Copyright (c) Ayedo Cloud GmbH
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// Client talks to the Polycrate REST API with a Bearer token.
// Same auth as polycrate-sdk: Authorization: Bearer <token>.
type Client struct {
	Endpoint   string
	Token      string
	HTTPClient *http.Client
}

func New(endpoint, token string) *Client {
	return &Client{
		Endpoint: strings.TrimRight(endpoint, "/"),
		Token:    token,
		HTTPClient: &http.Client{
			Timeout: 60 * time.Second,
		},
	}
}

type Host struct {
	ID              string  `json:"id"`
	Name            string  `json:"name"`
	Kind            string  `json:"kind"`
	Role            string  `json:"role"`
	BootstrapStatus string  `json:"bootstrap_status"`
	DefaultIPv4     *string `json:"default_ipv4"`
	ProductID       string  `json:"-"`
	WorkspaceID     string  `json:"-"`
	Raw             map[string]any
}

type Product struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Kind   string `json:"kind"`
	Active bool   `json:"active"`
}

type CatalogObject struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	DisplayName string `json:"display_name"`
	SalesReady  bool   `json:"sales_ready"`
	SCSIaaS     string `json:"scs_iaas_status"`
	ISO27001    string `json:"iso27001_status"`
}

func (c *Client) CreateHost(ctx context.Context, body map[string]any) (*Host, error) {
	var raw map[string]any
	if err := c.do(ctx, http.MethodPost, "/api/v1/hosts/", body, &raw); err != nil {
		return nil, err
	}
	return hostFromMap(raw)
}

func (c *Client) GetHost(ctx context.Context, id string) (*Host, error) {
	var raw map[string]any
	if err := c.do(ctx, http.MethodGet, "/api/v1/hosts/"+url.PathEscape(id)+"/", nil, &raw); err != nil {
		return nil, err
	}
	return hostFromMap(raw)
}

func (c *Client) UpdateHost(ctx context.Context, id string, body map[string]any) (*Host, error) {
	var raw map[string]any
	if err := c.do(ctx, http.MethodPatch, "/api/v1/hosts/"+url.PathEscape(id)+"/", body, &raw); err != nil {
		return nil, err
	}
	return hostFromMap(raw)
}

func (c *Client) DeleteHost(ctx context.Context, id string) error {
	return c.do(ctx, http.MethodDelete, "/api/v1/hosts/"+url.PathEscape(id)+"/", nil, nil)
}

func (c *Client) GetProduct(ctx context.Context, id string) (*Product, error) {
	var p Product
	if err := c.do(ctx, http.MethodGet, "/api/v1/products/"+url.PathEscape(id)+"/", nil, &p); err != nil {
		return nil, err
	}
	return &p, nil
}

func (c *Client) GetProvider(ctx context.Context, id string) (*CatalogObject, error) {
	var o CatalogObject
	if err := c.do(ctx, http.MethodGet, "/api/v1/providers/"+url.PathEscape(id)+"/", nil, &o); err != nil {
		return nil, err
	}
	return &o, nil
}

func (c *Client) GetPop(ctx context.Context, id string) (*CatalogObject, error) {
	var o CatalogObject
	if err := c.do(ctx, http.MethodGet, "/api/v1/pops/"+url.PathEscape(id)+"/", nil, &o); err != nil {
		return nil, err
	}
	return &o, nil
}

func hostFromMap(raw map[string]any) (*Host, error) {
	h := &Host{Raw: raw}
	h.ID, _ = raw["id"].(string)
	h.Name, _ = raw["name"].(string)
	h.Kind, _ = raw["kind"].(string)
	h.Role, _ = raw["role"].(string)
	h.BootstrapStatus, _ = raw["bootstrap_status"].(string)
	if v, ok := raw["default_ipv4"].(string); ok && v != "" {
		h.DefaultIPv4 = &v
	}
	if p, ok := raw["product"].(map[string]any); ok {
		h.ProductID, _ = p["id"].(string)
	}
	if w, ok := raw["workspace"].(map[string]any); ok {
		h.WorkspaceID, _ = w["id"].(string)
	}
	if h.ID == "" {
		return nil, fmt.Errorf("host response missing id")
	}
	return h, nil
}

func (c *Client) do(ctx context.Context, method, path string, body any, dest any) error {
	var rdr io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return err
		}
		rdr = bytes.NewReader(b)
	}
	req, err := http.NewRequestWithContext(ctx, method, c.Endpoint+path, rdr)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+c.Token)
	req.Header.Set("Accept", "application/json")
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	payload, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if resp.StatusCode == http.StatusNotFound {
		return &NotFoundError{Path: path}
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("%s %s: HTTP %d: %s", method, path, resp.StatusCode, truncate(payload, 512))
	}
	if dest == nil || resp.StatusCode == http.StatusNoContent || resp.StatusCode == http.StatusAccepted {
		return nil
	}
	if len(payload) == 0 {
		return nil
	}
	return json.Unmarshal(payload, dest)
}

type NotFoundError struct {
	Path string
}

func (e *NotFoundError) Error() string {
	return "not found: " + e.Path
}

func truncate(b []byte, n int) string {
	s := strings.TrimSpace(string(b))
	if len(s) > n {
		return s[:n] + "…"
	}
	return s
}
