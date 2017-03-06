package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

// NginxClient lets you update HTTP/Stream servers in NGINX Plus via its upstream_conf API
type NginxClient struct {
	httpClient   *Client
	streamClient *Client
}

// NewNginxClient creates an NginxClient.
func NewNginxClient(upstreamConfEndpoint string, statusEndpoint string, timeout time.Duration) (*NginxClient, error) {
	httpClient, err := NewHTTPClient(upstreamConfEndpoint, statusEndpoint, timeout)
	if err != nil {
		return nil, err
	}

	streamClient, err := NewStreamClient(upstreamConfEndpoint, statusEndpoint, timeout)
	if err != nil {
		return nil, err
	}

	return &NginxClient{httpClient, streamClient}, nil

}

// CheckIfHTTPUpstreamExists checks if the HTTP upstream exists in NGINX. If the upstream doesn't exist, it returns an error.
func (client *NginxClient) CheckIfHTTPUpstreamExists(upstream string) error {
	return client.httpClient.CheckIfUpstreamExists(upstream)
}

// CheckIfStreamUpstreamExists checks if the Stream upstream exists in NGINX. If the upstream doesn't exist, it returns an error.
func (client *NginxClient) CheckIfStreamUpstreamExists(upstream string) error {
	return client.streamClient.CheckIfUpstreamExists(upstream)
}

// UpdateHTTPServers updates the servers of the HTTP upstream.
// Servers that are in the slice, but don't exist in NGINX will be added to NGINX.
// Servers that aren't in the slice, but exist in NGINX, will be removed from NGINX.
func (client *NginxClient) UpdateHTTPServers(upstream string, servers []string) ([]string, []string, error) {
	return client.httpClient.UpdateServers(upstream, servers)
}

// UpdateStreamServers updates the servers of the Stream upstream.
// Servers that are in the slice, but don't exist in NGINX will be added to NGINX.
// Servers that aren't in the slice, but exist in NGINX, will be removed from NGINX.
func (client *NginxClient) UpdateStreamServers(upstream string, servers []string) ([]string, []string, error) {
	return client.streamClient.UpdateServers(upstream, servers)
}

// Client lets you add/remove servers to/from NGINX Plus via its upstream_conf API
type Client struct {
	upstreamConfEndpoint string
	statusEndpoint       string
	httpClient           *http.Client
}

type peers struct {
	Peers []peer
}

type peer struct {
	ID     int
	Server string
}

// NewHTTPClient creates a new HTTP client.
func NewHTTPClient(upstreamConfEndpoint string, statusEndpoint string, timeout time.Duration) (*Client, error) {
	httpClient := &http.Client{Timeout: timeout}

	err := checkIfUpstreamConfIsAccessible(httpClient, upstreamConfEndpoint)
	if err != nil {
		return nil, err
	}

	err = checkIfStatusIsAccessible(httpClient, statusEndpoint)
	if err != nil {
		return nil, err
	}

	client := &Client{upstreamConfEndpoint: upstreamConfEndpoint + "?", statusEndpoint: statusEndpoint, httpClient: httpClient}
	return client, nil
}

// NewStreamClient creates a new Stream client.
func NewStreamClient(upstreamConfEndpoint string, statusEndpoint string, timeout time.Duration) (*Client, error) {
	httpClient := &http.Client{Timeout: timeout}

	err := checkIfUpstreamConfIsAccessible(httpClient, upstreamConfEndpoint)
	if err != nil {
		return nil, err
	}

	err = checkIfStatusIsAccessible(httpClient, statusEndpoint)
	if err != nil {
		return nil, err
	}

	client := &Client{upstreamConfEndpoint: upstreamConfEndpoint + "?stream=", statusEndpoint: statusEndpoint + "/stream", httpClient: httpClient}
	return client, nil
}

func checkIfUpstreamConfIsAccessible(httpClient *http.Client, endpoint string) error {
	resp, err := httpClient.Get(endpoint)
	if err != nil {
		return fmt.Errorf("upstream_conf endpoint %v is not accessible: %v", endpoint, err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("upstream_conf endpoint %v is not accessible: %v", endpoint, err)
	}

	if resp.StatusCode != http.StatusNotFound {
		return fmt.Errorf("upstream_conf endpoint %v is not accessible: expected 404 response, got %v", endpoint, resp.StatusCode)
	}

	bodyStr := string(body)
	expected := "missing \"upstream\" argument\n"
	if bodyStr != expected {
		return fmt.Errorf("upstream_conf endpoint %v is not accessible: expected %q body, got %q", endpoint, expected, bodyStr)
	}

	return nil
}

func checkIfStatusIsAccessible(httpClient *http.Client, endpoint string) error {
	resp, err := httpClient.Get(endpoint)
	if err != nil {
		return fmt.Errorf("status endpoint is %v not accessible: %v", endpoint, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("status endpoint is %v not accessible: expected 200 response, got %v", endpoint, resp.StatusCode)
	}

	return nil
}

// CheckIfUpstreamExists checks if the upstream exists in NGINX. If the upstream doesn't exist, it returns an error.
func (client *Client) CheckIfUpstreamExists(upstream string) error {
	_, err := client.getUpstreamPeers(upstream)
	return err
}

func (client *Client) getUpstreamPeers(upstream string) (*peers, error) {
	request := fmt.Sprintf("%v/upstreams/%v", client.statusEndpoint, upstream)

	resp, err := client.httpClient.Get(request)
	if err != nil {
		return nil, fmt.Errorf("Failed to connect to the status api to get upstream %v info: %v", upstream, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("Upstream %v is not found", upstream)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Failed to read the response body with upstream %v info: %v", upstream, err)
	}
	var prs peers
	err = json.Unmarshal(body, &prs)
	if err != nil {
		return nil, fmt.Errorf("Error unmarshaling upstream %v: got %q response: %v", upstream, string(body), err)
	}

	return &prs, nil
}

// AddServer adds the server to the upstream.
func (client *Client) AddServer(upstream string, server string) error {
	id, err := client.getIDOfServer(upstream, server)

	if err != nil {
		return fmt.Errorf("Failed to add %v server to %v upstream: %v", server, upstream, err)
	}
	if id != -1 {
		return fmt.Errorf("Failed to add %v server to %v upstream: server already exists", server, upstream)
	}

	request := fmt.Sprintf("%v&upstream=%v&add=&server=%v", client.upstreamConfEndpoint, upstream, server)

	resp, err := client.httpClient.Get(request)
	if err != nil {
		return fmt.Errorf("Failed to add %v server to %v upstream: %v", server, upstream, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Failed to add %v server to %v upstream: expected 200 response, got %v", server, upstream, resp.StatusCode)
	}

	return nil
}

// DeleteServer the server from the upstream.
func (client *Client) DeleteServer(upstream string, server string) error {
	id, err := client.getIDOfServer(upstream, server)
	if err != nil {
		return fmt.Errorf("Failed to remove %v server from  %v upstream: %v", server, upstream, err)
	}
	if id == -1 {
		return fmt.Errorf("Failed to remove %v server from %v upstream: server doesn't exists", server, upstream)
	}

	request := fmt.Sprintf("%v&upstream=%v&remove=&id=%v", client.upstreamConfEndpoint, upstream, id)

	resp, err := client.httpClient.Get(request)
	if err != nil {
		return fmt.Errorf("Failed to remove %v server from %v upstream: %v", server, upstream, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("Failed to remove %v server to %v upstream: expected 200 or 204 response, got %v", server, upstream, resp.StatusCode)
	}

	return nil
}

// UpdateServers updates the servers of the upstream.
// Servers that are in the slice, but don't exist in NGINX will be added to NGINX.
// Servers that aren't in the slice, but exist in NGINX, will be removed from NGINX.
func (client *Client) UpdateServers(upstream string, servers []string) ([]string, []string, error) {
	serversInNginx, err := client.GetServers(upstream)
	if err != nil {
		return nil, nil, fmt.Errorf("Failed to update servers of %v upstream: %v", upstream, err)
	}

	toAdd, toDelete := determineUpdates(servers, serversInNginx)

	for _, server := range toAdd {
		err := client.AddServer(upstream, server)
		if err != nil {
			return nil, nil, fmt.Errorf("Failed to update servers of %v upstream: %v", upstream, err)
		}
	}

	for _, server := range toDelete {
		err := client.DeleteServer(upstream, server)
		if err != nil {
			return nil, nil, fmt.Errorf("Failed to update servers of %v upstream: %v", upstream, err)
		}
	}

	return toAdd, toDelete, nil
}

func determineUpdates(updatedServers []string, nginxServers []string) (toAdd []string, toRemove []string) {
	for _, server := range updatedServers {
		found := false
		for _, serverNGX := range nginxServers {
			if server == serverNGX {
				found = true
				break
			}
		}
		if !found {
			toAdd = append(toAdd, server)
		}
	}

	for _, serverNGX := range nginxServers {
		found := false
		for _, server := range updatedServers {
			if serverNGX == server {
				found = true
				break
			}
		}
		if !found {
			toRemove = append(toRemove, serverNGX)
		}
	}

	return
}

// GetServers returns the servers of the upsteam from NGINX.
func (client *Client) GetServers(upstream string) ([]string, error) {
	peers, err := client.getUpstreamPeers(upstream)
	if err != nil {
		return nil, fmt.Errorf("Error getting servers of %v upstream: %v", upstream, err)
	}

	var servers []string
	for _, peer := range peers.Peers {
		servers = append(servers, peer.Server)
	}

	return servers, nil
}

func (client *Client) getIDOfServer(upstream string, name string) (int, error) {
	peers, err := client.getUpstreamPeers(upstream)
	if err != nil {
		return -1, fmt.Errorf("Error getting id of server %v of upstream %v:", name, upstream)
	}

	for _, p := range peers.Peers {
		if p.Server == name {
			return p.ID, nil
		}
	}

	return -1, nil
}
