// Create and maintain by Chaiyapong Lapliengtrakul (chaiyapong@3dsinteractive.com), All right reserved (2021 - Present)
package generator

import (
	"github.com/martechassociation/martechassociation.github.io/microservices"
)

// Generator is service to genereate martech page files
type Generator struct {
	ctx microservices.IContext
	cfg microservices.IConfig
}

// NewGenerator return new generator service
func NewGenerator(ctx microservices.IContext, cfg microservices.IConfig) *Generator {
	return &Generator{
		ctx: ctx,
		cfg: cfg,
	}
}

func (svc *Generator) GenerateLandingPages() error {
	ctx := svc.ctx
	ctx.Log("Generate landing pages")

	return nil
}
