package collector

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/m47ik/drl-exporter/internal/vars"
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
	l             = log.New(os.Stdout, "drl-exporter ", log.LstdFlags)
)

// GetMetrics will get the token from dockerhub and send requests to get metrics headers
func GetMetrics() {

	tokenUrl := "https://auth.docker.io/token?service=" +
		"registry.docker.io&scope=repository:" + *vars.DockerRepoImage + ":pull"
	repoUrl := "https://registry-1.docker.io/v2/" +
		"registry.docker.io&scope=repository:" + *vars.DockerRepoImage + "/manifests/latest"
	var auth = Authentication{}

	c := &http.Client{Timeout: 10 * time.Second}

	tokenReq, err := http.NewRequest("GET", tokenUrl, nil)
	if err != nil {
		l.Println("unable to send request for token")
	}

	if *vars.EnableUserAuth == true {
		tokenReq.SetBasicAuth(*vars.Username, *vars.Password)
	}

	tokenReq.Close = true

	tokenResp, err := c.Do(tokenReq)
	if err != nil {
		l.Println("unable to get token", err)
	}
	defer tokenResp.Body.Close()

	tokenBody, err := ioutil.ReadAll(tokenResp.Body)
	if err != nil {
		l.Println("unable to get token body")
	}

	err = json.Unmarshal([]byte(tokenBody), &auth)
	if err != nil {
		l.Println("unable to save token")
	}

	limitReq, err := http.NewRequest("HEAD", repoUrl, nil)
	if err != nil {
		l.Println(err)
	}

	limitReq.Header.Add("Authorization", "Bearer "+auth.Token)
	limitResp, err := c.Do(limitReq)
	if err != nil {
		l.Println("unable to get response from dockerhub")
	}

	limitHeader := limitResp.Header.Get("RateLimit-Limit")
	remainHeader := limitResp.Header.Get("RateLimit-Remaining")

	dockerLimit, err := convertHeaders(limitHeader)
	if err != nil {
		l.Println("no header data for limit")
	}
	dockerLimitRemain, err := convertHeaders(remainHeader)
	if err != nil {
		l.Println("no header data for remaining limit")
	}

	switch {
	case len(dockerLimit) <= 0:
		return
	case len(dockerLimitRemain) <= 0:
		return
	default:
		DockerMetrics["maxRequestTotal"] = dockerLimit[0]
		DockerMetrics["maxRequestTotalTime"] = dockerLimit[1]
		DockerMetrics["remainingRequestTotal"] = dockerLimitRemain[0]
		DockerMetrics["remainingRequestTotalTime"] = dockerLimitRemain[1]
	}
}

// convert and normalize the recieved headers
func convertHeaders(data string) ([]float64, error) {

	rs := strings.Replace(data, "w=", "", 2)
	ss := strings.Split(rs, ";")
	xFloat, err := convertToFloat(ss)
	if err != nil {
		return nil, err
	}
	return xFloat, nil
}

// takes slice of strings and convert them to float64
func convertToFloat(xs []string) ([]float64, error) {
	var xFloat []float64
	for i := range xs {
		str := xs[i]
		in, err := strconv.ParseFloat(str, 64)
		if err != nil {
			return nil, err
		}
		xFloat = append(xFloat, in)
	}

	return xFloat, nil
}
