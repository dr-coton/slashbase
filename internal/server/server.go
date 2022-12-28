package server

import (
	"fmt"
	"log"
	"mime"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/icza/gox/osx"
	"github.com/knadh/stuffbin"
	"github.com/slashbaseide/slashbase/internal/config"
	"github.com/slashbaseide/slashbase/internal/utils"
)

// Init server
func Init() {
	fmt.Println("Running Slashbase IDE at http://localhost:" + config.GetServerPort())
	if config.IsLive() {
		gin.SetMode(gin.ReleaseMode)
	}
	router := NewRouter()
	serveStaticFiles(router)
	if config.GetConfig().BuildName == config.BUILD_DOCKER_PROD {
		router.Run(":" + config.GetServerPort())
	} else {
		osx.OpenDefault("https://app.slashbase.com")
		go router.Run(":" + config.GetServerPort())
	}
}

func serveStaticFiles(router *gin.Engine) {
	// Serving the Frontend files in Production
	if config.IsLive() {
		fs := initFS()
		router.NoRoute(func(c *gin.Context) {
			if file, err := fs.Read("web/" + c.Request.URL.Path); err == nil {
				contentType := mime.TypeByExtension("." + utils.FileExtensionFromPath(c.Request.URL.Path))
				c.Data(http.StatusOK, contentType, file)
				return
			}
			indexFileData, _ := fs.Read("web/index.html")
			c.Data(http.StatusOK, "text/html", indexFileData)
		})
	}
}

// initFS initializes the stuffbin FileSystem to provide
// access to bunded static assets to the app.
func initFS() stuffbin.FileSystem {
	path, err := os.Executable()
	if err != nil {
		log.Fatalf("error getting executable path: %v", err)
	}

	fs, err := stuffbin.UnStuff(path)
	if err == nil {
		return fs
	}

	// Running in local mode. Load the required static assets into
	// the in-memory stuffbin.FileSystem.
	// unable to initialize embedded filesystem
	// using local filesystem for static assets

	files := []string{
		"web",
	}

	fs, err = stuffbin.NewLocalFS("/", files...)
	if err != nil {
		log.Fatalf("failed to load local static files: %v", err)
	}

	return fs
}
