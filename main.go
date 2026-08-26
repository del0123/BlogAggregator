package main

import (
	"database/sql"
	"fmt"
	"gator/internal/config"
	"gator/internal/database"
	"log"
	"os"

	_ "github.com/lib/pq"
)

func main() {
	fmt.Println("Welcome to Gator!")

	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("error reading config: %v", err)
	}

	db, err := sql.Open("postgres", cfg.DBUrl)
	if err != nil {
		log.Fatalf("error opening database: %v", err)
	}
	defer db.Close()

	dbQueries := database.New(db)

	s := &state{
		Db:     dbQueries,
		Config: &cfg,
	}

	cmds := commands{
		RegisteredCmds: make(map[string]func(*state, command) error),
	}

	cmds.register("login", handlerLogin)
	cmds.register("register", handlerRegister)
	cmds.register("reset", handlerReset)
	cmds.register("users", handlerGetUsers)
	cmds.register("agg", handlerAggregate)
	cmds.register("addfeed", middlewareLoggedIn(handlerAddFeed))
	cmds.register("feeds", handlerGetFeeds)
	cmds.register("follow", middlewareLoggedIn(handlerFollow))
	cmds.register("following", middlewareLoggedIn(handlerFollowing))
	cmds.register("unfollow", middlewareLoggedIn(handlerUnfollow))
	cmds.register("browse", middlewareLoggedIn(handlerBrowse))

	if len(os.Args) < 2 {
		fmt.Println("insufficient arguments, needs command")
		os.Exit(1)
	}

	cmd := command{
		Name: os.Args[1],
		Args: os.Args[2:],
	}

	err = cmds.run(s, cmd)
	if err != nil {
		fmt.Printf("error running command: %v \n", err)
		os.Exit(1)
	}
}
