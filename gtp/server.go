package main

import (
	"fmt"
	"gtp/db/rdb"
	"gtp/db/rdb/model"
	"gtp/graph"
	"gtp/graph/generated"
	"gtp/middleware/auth"
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/go-chi/chi"
)

const defaultPort = "8080"

func main() {

	db, err := rdb.SupabaseDB()
	if err != nil {
		log.Fatalln(err)
	}

	user := model.User{Name: "user_name1"}
	todo := model.Todo{Text: "todo1_text", Done: false, User: &user}
	result := db.Create(&user)
	fmt.Println(result)
	result = db.Create(&todo)
	fmt.Println(result)

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	router := chi.NewRouter()
	// router.Use(auth.FirebaseAuth())
	router.Use(auth.GCIPAuth())

	// srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{}}))
	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{DB: db}}))

	// http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	// http.Handle("/query", srv)
	router.Handle("/", playground.Handler("GraphQL playground", "/query"))
	router.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	// log.Fatal(http.ListenAndServe(":"+port, nil))
	log.Fatal(http.ListenAndServe(":"+port, router))
}
