package rdb

import (
	"fmt"
	"gtp/db/rdb/model"
	gcp_util "gtp/utils/gcp"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var supabase_db *gorm.DB

func SupabaseDB() (*gorm.DB, error) {
	if supabase_db != nil {
		return supabase_db, nil
	}

	supabase_db, err := InitDB()
	if err != nil {
		log.Fatalln(err)
	}

	Migrate(supabase_db)

	return supabase_db, nil
}

func InitDB() (*gorm.DB, error) {
	supabase_dsn := os.Getenv("SUPABASE_DB_CONNECTION_STRING_GO")
	if supabase_dsn == "" {
		fmt.Println("Environment variable SUPABASE_DB_CONNECTION_STRING_GO is not set, fetching secret settings SUPABASE_DB_CONNECTION_STRING_GO from GCP...")

		gcp_project_id := os.Getenv("GCP_PROJECT_ID")
		if gcp_project_id == "" {
			panic("GCP_PROJECT_ID must be set")
		}
		var err error = nil
		gcp_sa_filepath := os.Getenv("GCP_SA_CREDENTIAL_FILEPATH")
		if gcp_sa_filepath == "" {
			supabase_dsn, err = gcp_util.GetSecret("SUPABASE_DB_CONNECTION_STRING_GO", gcp_project_id, nil)
		} else {
			supabase_dsn, err = gcp_util.GetSecret("SUPABASE_DB_CONNECTION_STRING_GO", gcp_project_id, &gcp_sa_filepath)
		}
		if err != nil {
			log.Fatalln(err)
		}
	}
	fmt.Println(supabase_dsn)

	var err error
	supabase_db, err = gorm.Open(postgres.Open(supabase_dsn), &gorm.Config{})
	if err != nil {
		fmt.Errorf("error initializing app: %v", err)
	}
	fmt.Println(supabase_db)

	return supabase_db, nil
}

func Migrate(db *gorm.DB) {
	err := db.AutoMigrate(&model.NewTodo{}, &model.Todo{}, &model.User{})
	if err != nil {
		log.Fatalln(err)
	}
}
