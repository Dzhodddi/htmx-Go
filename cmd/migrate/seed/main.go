package main

import (
	"log"
	"project/internal/db"
	"project/internal/env"
	"project/internal/store"
)

func main() {
	addr := env.GetString("DB_ADDR",
		"postgresql://admin:adminpasswrod@localhost/social?sslmode=disable")
	conn, err := db.New(addr, 3, 3, "15m")
	if err != nil {
		log.Fatal(err)
	}

	defer conn.Close()
	
	store := store.NewStorage(conn)

	db.Seed(store)
}
