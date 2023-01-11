package collector

import (
	"encoding/json"
	"io"
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

var DockerMetrics = make(map[string]float64, 4)
var DockerLabels = make(map[string]string, 1)

// GetMetrics will save metrics in a float64 map
func GetMetrics() {
	l := log.New(os.Stdout, "drl-exporter ", log.LstdFlags)
	tokenUrl := "https://auth.docker.io/token?service=" +
		"registry.docker.io&scope=repository:" + *vars.DockerRepoImage + ":pull"
	repoUrl := "https://registry-1.docker.io/v2/" +
		"registry.docker.io&scope=repository:" + *vars.DockerRepoImage + "/manifests/latest"

	tr, err := tokenRequest(tokenUrl)
	if err != nil {
		l.Printf("unable to send request %v\n",err)
	}

	tb, err := tokenBody(tr)
	if err != nil {
		l.Printf("unable to get token data %v\n", err)
	}


	lh, err := getLimitHeaders(repoUrl, tb)
	if err != nil {
		l.Printf("unexpected response from docker %v\n", err)
	}

	limitHeader := lh.Header.Get("RateLimit-Limit")
	remainHeader := lh.Header.Get("RateLimit-Remaining")
	
	sourceHeader := lh.Header.Get("Docker-RateLimit-Source")
	if sourceHeader == "" {
		l.Println("no header data for docker-ratelimit-source")
	}

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
		DockerLabels["reqsource"] = sourceHeader
	}

}


func convertHeaders(data string) ([]float64, error) {

	rs := strings.Replace(data, "w=", "", 2)
	ss := strings.Split(rs, ";")
	xFloat, err := convertToFloat(ss)
	if err != nil {
		return nil, err
	}
	return xFloat, nil
}

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

func tokenRequest(url string ) (*http.Request, error){
	tr, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	if *vars.EnableUserAuth == true {
		tr.SetBasicAuth(*vars.Username, *vars.Password)
	}
	return tr, nil
}

func tokenBody(req *http.Request) ([]byte, error) {
	c := http.Client{}
	rsp, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	tr, err := io.ReadAll(rsp.Body)
	if err != nil {
		return nil, err
	}
	rsp.Body.Close()
	return tr, nil

}

func getLimitHeaders(url string, td []byte) (*http.Response,error) {
	c := &http.Client{Timeout: 10 * time.Second}
	var auth = Authentication{}

	tkErr := json.Unmarshal(td, &auth)
	if tkErr != nil {
		return nil, tkErr
	}

	lmr, err := http.NewRequest("HEAD", url, nil)
	if err != nil {
		return nil, err
	}

	lmr.Header.Add("Authorization", "Bearer "+auth.Token)
	lr, err := c.Do(lmr)
	if err != nil {
		return nil, err
	}

	return lr, nil
}