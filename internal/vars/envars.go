package vars

import (
	"github.com/nicholasjackson/env"
)

var (
	BindAddress     = env.String("EXPORTER_PORT", false, "2121", "Bind address for server")
	Username        = env.String("DOCKERHUB_USER", false, "", "dockerhub username")
	Password        = env.String("DOCKERHUB_PASSWORD", false, "", "Dockerhub password")
	DockerRepoImage = env.String("DOCKERHUB_REPO_IMAGE", false, "ratelimitpreview/test", "Repository image")
	EnableUserAuth  = env.Bool("ENABLE_USER_AUTH", false, false, "Enable metrics for users")
)
