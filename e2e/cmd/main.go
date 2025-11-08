package main 

import(
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/jfb0301/golang-testing-reference/e2e/db"
	"github.com/jfb0301/golang-testing-reference/e2e/handlers"
	migrate "github.com/golang-migrate/migrate/v4"
	mpostgres "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/file"
	gormpostgres "gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	port, ok := os.LookupEnv("BOOKSWAP_PORT")
	if !ok {
		log.Fatal("$BOOKSWAP_PORT not found")
	}

	postgresURL, ok := os.LookupEnv("BOOKSWAP_DB_URL")
	if !ok {
		log.Fatal("$BOOKSWAP_DB_URL not found")
	}
		m, err := migrate.New("file://e2e/db/migrations", postgresURL)
		if err != nil {
			log.Fatal("migrate:%v", err)
		}
		if err := m.Up(); err != nil && err != migrate.ErrNoChange {
			log.Fatalf("migration up:%v", err)
		}

		// defer func() {
		// m.Down()
		// }
		dbconn, err := gorm.Open(gormpostgres.Open(postgresURL), &gorm.Config{})
		if err != nil {
			log.Fatal("db open:%v", err)
		}
		ps := db.NewPostingService()
		b := db.NewBookService(dbconn, ps)
		u := db.NewUserService(dbconn, b)
		h := handlers.NewHandler(b, u)

		router := handlers.ConfigureServer(h)
		log.Printf("Listening on: %s...\n" , port)
		log.Fatal(http.ListenAndServe(fmt.Sprint(":", port), router))

}