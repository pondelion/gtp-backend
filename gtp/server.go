package main

import (
	"fmt"
	"gtp/graph"
	"gtp/graph/generated"
	"gtp/graph/model"
	"gtp/utils/gcp"
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const defaultPort = "8080"

func main() {

	gcp_sa_filepath := os.Getenv("GCP_SA_CREDENTIAL_FILEPATH")
	if gcp_sa_filepath == "" {
		panic("GCP_SA_CREDENTIAL_FILEPATH must be set")
	}
	gcp_project_id := os.Getenv("GCP_PROJECT_ID")
	if gcp_project_id == "" {
		panic("GCP_PROJECT_ID must be set")
	}

	// ctx := context.Background()

	// opt := option.WithCredentialsFile(saFilepath)
	// fmt.Println(opt)
	// app, err := firebase.NewApp(ctx, nil, opt)
	// if err != nil {
	// 	fmt.Errorf("error initializing app: %v", err)
	// }

	// client, err := app.Firestore(ctx)
	// if err != nil {
	// 	log.Fatalln(err)
	// }
	// defer client.Close()

	supabase_db_uri := os.Getenv("SUPABASE_DB_URI")
	if supabase_db_uri == "" {
		var err error = nil
		fmt.Println("Fetching secret settings SUPABASE_DB_URI from GCP...")
		supabase_db_uri, err = gcp.GetSecret("SUPABASE_DB_URI", gcp_sa_filepath, gcp_project_id)
		if err != nil {
			log.Fatalln(err)
		}
	}
	fmt.Println(supabase_db_uri)

	supabase_dsn := os.Getenv("SUPABASE_DB_CONNECTION_STRING_GO")
	if supabase_dsn == "" {
		var err error = nil
		fmt.Println("Fetching secret settings SUPABASE_DB_CONNECTION_STRING_GO from GCP...")
		supabase_dsn, err = gcp.GetSecret("SUPABASE_DB_CONNECTION_STRING_GO", gcp_sa_filepath, gcp_project_id)
		if err != nil {
			log.Fatalln(err)
		}
	}
	fmt.Println(supabase_dsn)

	db, err := gorm.Open(postgres.Open(supabase_dsn), &gorm.Config{})
	if err != nil {
		fmt.Errorf("error initializing app: %v", err)
	}
	fmt.Println(db)

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	err = db.AutoMigrate(&model.NewTodo{}, &model.Todo{}, &model.User{})
	if err != nil {
		log.Fatalln(err)
	}

	// srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{}}))
	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{DB: db}}))

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
