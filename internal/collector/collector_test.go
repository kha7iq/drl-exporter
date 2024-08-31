package collector

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/docker/cli/cli/config/types"
	"github.com/kha7iq/drl-exporter/internal/vars"
	"github.com/stretchr/testify/assert"
)

// Mock HTTPClient to simulate HTTP responses
type MockClient struct{}

func (m *MockClient) Do(req *http.Request) (*http.Response, error) {
	return &http.Response{
		StatusCode: http.StatusOK,
		Body:       ioutil.NopCloser(bytes.NewBufferString(`{"token": "test-token"}`)),
	}, nil
}

// TestBuildURLs tests the URL construction logic
func TestBuildURLs(t *testing.T) {
	// Setup the variables
	repoImage := "test/image"
	vars.DockerRepoImage = &repoImage

	// Test with IPv6 disabled
	vars.EnableIPv6 = new(bool)
	tokenURL, repoURL := buildURLs()
	assert.Equal(t, "https://auth.docker.io/token?service=registry.docker.io&scope=repository:test/image:pull", tokenURL)
	assert.Equal(t, "https://registry-1.docker.io/v2/registry.docker.io&scope=repository:test/image/manifests/latest", repoURL)

	// Test with IPv6 enabled
	*vars.EnableIPv6 = true
	tokenURL, repoURL = buildURLs()
	assert.Equal(t, "https://auth.ipv6.docker.com/token?service=registry.docker.io&scope=repository:test/image:pull", tokenURL)
	assert.Equal(t, "https://registry.ipv6.docker.com/v2/registry.docker.io&scope=repository:test/image/manifests/latest", repoURL)
}

// TestCreateTokenRequest tests the token request creation
func TestCreateTokenRequest(t *testing.T) {
	// Mock loadDockerHubAuth function
	mockLoadDockerHubAuth := func() (types.AuthConfig, error) {
		return types.AuthConfig{Username: "fileuser", Password: "filepass"}, nil
	}

	// Set up variables
	username := "user"
	password := "pass"
	vars.Username = &username
	vars.Password = &password
	vars.EnableUserAuth = new(bool)
	*vars.EnableUserAuth = true

	// Test with user authentication enabled
	req, err := createTokenRequest("http://example.com", mockLoadDockerHubAuth)
	assert.NoError(t, err)
	assert.Equal(t, "Basic dXNlcjpwYXNz", req.Header.Get("Authorization"))

	// Test with file-based authentication (mocked)
	*vars.EnableUserAuth = false
	vars.EnableFileAuth = new(bool)
	*vars.EnableFileAuth = true

	req, err = createTokenRequest("http://example.com", mockLoadDockerHubAuth)
	assert.NoError(t, err)
	assert.Equal(t, "Basic ZmlsZXVzZXI6ZmlsZXBhc3M=", req.Header.Get("Authorization"))
}

// TestExecuteRequest tests the execution of HTTP requests
func TestExecuteRequest(t *testing.T) {
	client := &MockClient{}
	req, _ := http.NewRequest("GET", "http://example.com", nil)

	body, err := executeRequest(client, req)
	assert.NoError(t, err)
	assert.Equal(t, `{"token": "test-token"}`, string(body))
}

// TestParseHeaderValues tests the header parsing logic
func TestParseHeaderValues(t *testing.T) {
	header := "100;w=10"
	expected := []float64{100, 10}

	values, err := parseHeaderValues(header)
	assert.NoError(t, err)
	assert.Equal(t, expected, values)
}
