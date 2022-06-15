package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"cloud.google.com/go/spanner"
	database "cloud.google.com/go/spanner/admin/database/apiv1"
	"google.golang.org/api/iterator"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", helloSpanner)

	log.Printf("Server listening on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, mux))
}

func createClients(ctx context.Context, db string) (*database.DatabaseAdminClient, *spanner.Client) {
	adminClient, err := database.NewDatabaseAdminClient(ctx)
	if err != nil {
		log.Fatal(err)
	}

	dataClient, err := spanner.NewClient(ctx, db)
	if err != nil {
		log.Fatal(err)
	}

	return adminClient, dataClient
}

func helloSpanner(writer http.ResponseWriter, request *http.Request) {
	host, _ := os.Hostname()
	fmt.Fprintf(writer, "Hostname: %s\n", host)

	db := "projects/" + os.Getenv("SPINNER_PROJECT_ID") + "/instances/" + os.Getenv("DB_INSTANCE") + "/databases/" + os.Getenv("DB_NAME")
	ctx := context.Background()
	adminClient, dataClient := createClients(ctx, db)
	defer adminClient.Close()
	defer dataClient.Close()

	stmt := spanner.Statement{SQL: `SELECT SingerId, AlbumId, AlbumTitle FROM Albums`}
	iter := dataClient.Single().Query(ctx, stmt)
	defer iter.Stop()
	for {
		row, err := iter.Next()
		if err == iterator.Done {
			fmt.Fprintln(writer, "Done")
			return
		}
		if err != nil {
			fmt.Fprintf(writer, "Query failed with %v", err)
			return
		}
		var singerID, albumID int64
		var albumTitle string
		if err := row.Columns(&singerID, &albumID, &albumTitle); err != nil {
			fmt.Fprintf(writer, "Failed to parse row %v", err)
			return
		}
		fmt.Fprintf(writer, "%d %d %s\n", singerID, albumID, albumTitle)
	}
}
