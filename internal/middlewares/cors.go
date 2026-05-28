package middlewares

import "github.com/go-chi/cors"

func CORS() *cors.Cors {

	return cors.New(cors.Options{
		AllowedOrigins: []string{
			"http://localhost:3000",
		},

		AllowedMethods: []string{
			"GET",
			"POST",
			"PUT",
			"DELETE",
			"OPTIONS",
		},

		AllowedHeaders: []string{
			"Accept",
			"Authorization",
			"Content-Type",
		},

		AllowCredentials: true,
		MaxAge:           300,
	})
}
