// Create and maintain by Chaiyapong Lapliengtrakul (chaiyapong@3dsinteractive.com), All right reserved (2021 - Present)
package main

import (
	"fmt"
	"net/http"

	"github.com/martechassociation/martechassociation.github.io/generator"
	"github.com/martechassociation/martechassociation.github.io/microservices"
)

func main() {
	cfg := microservices.NewConfig()
	ms := microservices.NewMicroservice()
	defer ms.Stop()

	startServices(ms, cfg)

	ms.Start()
}

func startServices(ms *microservices.Microservice, cfg microservices.IConfig) {
	serviceID := cfg.ServiceID()
	ms.Log("main", fmt.Sprintf("serviceID=%s", serviceID))
	switch serviceID {
	case "http":
		startHTTPServices(ms, cfg)
	}
}

func startHTTPServices(ms *microservices.Microservice, cfg microservices.IConfig) {
	ms.GET("/gen", func(ctx microservices.IContext) error {
		gen := generator.NewGenerator(ctx, cfg)
		gen.GenerateLandingPages()

		resp := map[string]string{"status": "success"}
		ctx.Response(http.StatusOK, resp)
		return nil
	})
}
