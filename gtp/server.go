package main

import (
	"fmt"
	"gtp/db/rdb"
	"gtp/db/rdb/model"
	"gtp/graph"
	"gtp/graph/generated"
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
)

const defaultPort = "8080"

func main() {

	// gcp_sa_filepath := os.Getenv("GCP_SA_CREDENTIAL_FILEPATH")
	// if gcp_sa_filepath == "" {
	// 	panic("GCP_SA_CREDENTIAL_FILEPATH must be set")
	// }
	// gcp_project_id := os.Getenv("GCP_PROJECT_ID")
	// if gcp_project_id == "" {
	// 	panic("GCP_PROJECT_ID must be set")
	// }

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

	// supabase_db_uri := os.Getenv("SUPABASE_DB_URI")
	// if supabase_db_uri == "" {
	// 	var err error = nil
	// 	fmt.Println("Fetching secret settings SUPABASE_DB_URI from GCP...")
	// 	supabase_db_uri, err = gcp.GetSecret("SUPABASE_DB_URI", gcp_sa_filepath, gcp_project_id)
	// 	if err != nil {
	// 		log.Fatalln(err)
	// 	}
	// }
	// fmt.Println(supabase_db_uri)

	db, err := rdb.SupabaseDB()
	if err != nil {
		log.Fatalln(err)
	}

	err = db.AutoMigrate(&model.NewTodo{}, &model.Todo{}, &model.User{})
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

	// srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{}}))
	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{DB: db}}))

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
