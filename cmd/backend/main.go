package main

import (
	"context"
	"fmt"

	"github.com/kappusuton-yon-tebaru/backend/cmd/backend/backend"
	"github.com/kappusuton-yon-tebaru/backend/cmd/backend/internal/router"
)

func main() {
	app, err := backend.Initialize()
	if err != nil {
		panic(err)
	}

	defer app.MongoClient.Disconnect(context.Background())

	r := router.New()

	r.RegisterRoutes(app)

	if err := r.Run(fmt.Sprintf("0.0.0.0:%v", app.Config.Backend.Port)); err != nil {
		panic(err)
	}
}
