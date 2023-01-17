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
	EnableFileAuth  = env.Bool("ENABLE_FILE_AUTH", false, false, "Enable authentication through docker configuration file 'config.json'")
	FileAuthDir     = env.String("FILE_AUTH_DIR", false, "/config", "Directory to load 'config.json' docker configuration from")
)
