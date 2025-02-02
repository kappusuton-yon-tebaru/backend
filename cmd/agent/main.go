package main

import (
	"fmt"

	"github.com/kappusuton-yon-tebaru/backend/cmd/agent/agent"
	"github.com/kappusuton-yon-tebaru/backend/cmd/agent/internal/router"
)

func main() {
	app, err := agent.Initialize()
	if err != nil {
		panic(err)
	}

	r := router.New()

	r.RegisterRoutes(app)

	if err := r.Run(fmt.Sprintf("0.0.0.0:%v", app.Config.Agent.Port)); err != nil {
		panic(err)
	}
}
