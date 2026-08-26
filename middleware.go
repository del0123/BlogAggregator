package main

import (
	"context"
	"gator/internal/database"
)

func middlewareLoggedIn(
	handler func(s *state, cmd command, user database.User) error,
) func(*state, command) error {

	return func(s *state, cmd command) error {
		ctx := context.Background()
		user, err := s.Db.GetUser(ctx, s.Config.CurrentUserName)
		if err != nil {
			return err
		}
		return handler(s, cmd, user)
	}
}
