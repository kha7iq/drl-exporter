package collector

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/docker/cli/cli/config"
	"github.com/docker/cli/cli/config/types"
	"github.com/kha7iq/drl-exporter/internal/vars"
)

// Authentication struct is used for saving token data returned by dockerhub
type Authentication struct {
	Token       string    `json:"token"`
	AccessToken string    `json:"access_token"`
	ExpiresIn   int       `json:"expires_in"`
	IssuedAt    time.Time `json:"issued_at"`
}

var (
	DockerMetrics = make(map[string]float64, 4)
	DockerLabels  = make(map[string]string, 1)
)

// HTTPClient is an interface for making HTTP requests.
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// GetMetrics will save metrics in a float64 map
func GetMetrics() {
	logger := log.New(os.Stdout, "drl-exporter ", log.LstdFlags)

	tokenURL, repoURL := buildURLs()
	tokenReq, err := createTokenRequest(tokenURL, loadDockerHubAuth)
	if err != nil {
		logger.Printf("unable to create token request: %v\n", err)
		return
	}

	client := &http.Client{}
	tokenBody, err := executeRequest(client, tokenReq)
	if err != nil {
		logger.Printf("unable to get valid token: %v\n", err)
		return
	}

	response, err := fetchLimitHeaders(client, repoURL, tokenBody)
	if err != nil {
		logger.Printf("unexpected response from docker: %v\n", err)
		return
	}

	processHeaders(response, logger)
}

func buildURLs() (string, string) {
	tokenBaseURL := "https://auth.docker.io"
	repoBaseURL := "https://registry-1.docker.io"
	if *vars.EnableIPv6 {
		tokenBaseURL = "https://auth.ipv6.docker.com"
		repoBaseURL = "https://registry.ipv6.docker.com"
	}

	tokenURL := fmt.Sprintf("%s/token?service=registry.docker.io&scope=repository:%s:pull", tokenBaseURL, *vars.DockerRepoImage)
	repoURL := fmt.Sprintf("%s/v2/registry.docker.io&scope=repository:%s/manifests/latest", repoBaseURL, *vars.DockerRepoImage)

	return tokenURL, repoURL
}

// Inject loadDockerHubAuth as a dependency
func createTokenRequest(url string, loadAuthFunc func() (types.AuthConfig, error)) (*http.Request, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	if *vars.EnableFileAuth && !*vars.EnableUserAuth {
		authConfig, err := loadAuthFunc()
		if err != nil {
			return nil, err
		}
		req.SetBasicAuth(authConfig.Username, authConfig.Password)
	} else if *vars.EnableUserAuth {
		req.SetBasicAuth(*vars.Username, *vars.Password)
	}

	return req, nil
}

func loadDockerHubAuth() (types.AuthConfig, error) {
	logger := log.New(os.Stdout, "drl-exporter-dockerauth ", log.LstdFlags)
	dockerConfDir := *vars.FileAuthDir
	dockerConfig, err := config.Load(dockerConfDir)
	if err != nil {
		logger.Printf("Unable to load docker configuration from config file '%v/config.json'\n", dockerConfDir)
		return types.AuthConfig{}, err
	}

	if !dockerConfig.ContainsAuth() {
		logger.Printf("No 'auths' found in configuration file '%v/config.json'\n", dockerConfDir)
	}

	return dockerConfig.GetAuthConfig("https://index.docker.io/v1/")
}

func executeRequest(client HTTPClient, req *http.Request) ([]byte, error) {
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("request failed with status code %d: %s", resp.StatusCode, string(body))
	}

	return body, nil
}

func fetchLimitHeaders(client HTTPClient, url string, tokenBody []byte) (*http.Response, error) {
	var auth Authentication

	if err := json.Unmarshal(tokenBody, &auth); err != nil {
		return nil, err
	}

	req, err := http.NewRequest("HEAD", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", "Bearer "+auth.Token)
	return client.Do(req)
}

func processHeaders(resp *http.Response, logger *log.Logger) {
	limitHeader := resp.Header.Get("RateLimit-Limit")
	remainHeader := resp.Header.Get("RateLimit-Remaining")
	sourceHeader := resp.Header.Get("Docker-RateLimit-Source")

	if sourceHeader == "" {
		logger.Println("no header data for docker-ratelimit-source")
	}

	dockerLimit, err := parseHeaderValues(limitHeader)
	if err != nil {
		logger.Println("no header data for limit")
		return
	}

	dockerLimitRemain, err := parseHeaderValues(remainHeader)
	if err != nil {
		logger.Println("no header data for remaining limit")
		return
	}

	if len(dockerLimit) > 0 && len(dockerLimitRemain) > 0 {
		DockerMetrics["maxRequestTotal"] = dockerLimit[0]
		DockerMetrics["maxRequestTotalTime"] = dockerLimit[1]
		DockerMetrics["remainingRequestTotal"] = dockerLimitRemain[0]
		DockerMetrics["remainingRequestTotalTime"] = dockerLimitRemain[1]
		DockerLabels["reqsource"] = sourceHeader
	}
}

func parseHeaderValues(data string) ([]float64, error) {
	cleanedData := strings.ReplaceAll(data, "w=", "")
	stringValues := strings.Split(cleanedData, ";")
	return convertStringsToFloats(stringValues)
}

func convertStringsToFloats(strings []string) ([]float64, error) {
	var floatValues []float64
	for _, str := range strings {
		value, err := strconv.ParseFloat(str, 64)
		if err != nil {
			return nil, err
		}
		floatValues = append(floatValues, value)
	}
	return floatValues, nil
}
