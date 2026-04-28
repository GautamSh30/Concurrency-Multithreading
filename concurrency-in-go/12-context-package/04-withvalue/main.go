package main

import (
	"context"
	"fmt"
)

type contextKey string

const userIDKey contextKey = "userID"

func fetchUser(ctx context.Context) {
	userID := ctx.Value(userIDKey)
	if userID == nil {
		fmt.Println("No user ID in context")
		return
	}
	fmt.Printf("Fetching data for user: %v\n", userID)
}

func middleware(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, userIDKey, userID)
}

func main() {
	ctx := context.Background()

	// Simulate middleware adding user ID to context
	ctx = middleware(ctx, "user-42")

	// Downstream function extracts value
	fetchUser(ctx)

	// Without value
	fetchUser(context.Background())
}
