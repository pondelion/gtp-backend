package main

import (
	"context"
	"gtp/graph"
	"gtp/graph/generated"
	"log"
	"net/http"
	"os"

	firebase "firebase.google.com/go"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"google.golang.org/api/option"
)

const defaultPort = "8080"

func main() {

	saFilepath := os.Getenv("GCP_SA_CREDENTIAL_FILEPATH")
	if saFilepath == "" {
		panic("GCP_SA_CREDENTIAL_FILEPATH must be set")
	}
	ctx := context.Background()
	sa := option.WithCredentialsFile(saFilepath)
	app, err := firebase.NewApp(ctx, nil, sa)
	if err != nil {
		log.Fatalln(err)
	}

	client, err := app.Firestore(ctx)
	if err != nil {
		log.Fatalln(err)
	}
	defer client.Close()

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{}}))

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
