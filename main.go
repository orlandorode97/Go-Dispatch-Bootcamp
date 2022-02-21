package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/orlandorode97/go-disptach/infrastructure/api"
	"github.com/orlandorode97/go-disptach/infrastructure/router"
	"github.com/orlandorode97/go-disptach/registry"
)

func main() {
	var (
		imaggaApiKey       = flag.String("imagga-api-key", envString("IMAGGA_API_KEY", ""), "Imagga api key")
		imaggaApiKeySecret = flag.String("imagga-api-key-secret", envString("IMAGGA_API_KEY_SECRET", ""), "Imagga api key")
		imaggaPort         = flag.String("imagga-port", envString("IMAGGA_PORT", "8080"), "Imagga port to listen to")
	)
	flag.Parse()

	imaggaClient := api.New(*imaggaApiKey, *imaggaApiKeySecret)

	r := registry.NewRegistry(imaggaClient)

	router := router.NewRouter(r.NewAppController())
	log.Printf("listening on http://localhost:%s", *imaggaPort)
	http.ListenAndServe(fmt.Sprintf(":%s", *imaggaPort), router)
}

func envString(env, fallback string) string {
	e := os.Getenv(env)
	if e == "" {
		return fallback
	}
	return e
}
