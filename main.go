package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/urfave/cli/v3"
)

var pool *pgxpool.Pool
var ctx = context.Background()

// User represent the users information in database
type User struct {
	ID       int
	Email    string
	Password string
}

func initDB(ctx context.Context) (*pgxpool.Pool, error) {

	// Load the .env file
	_ = godotenv.Load(".env")

	// Initialize the connection Pool
	config, err := pgxpool.ParseConfig(os.Getenv("database_url"))
	if err != nil {
		return nil, fmt.Errorf("unable to connect to the database: %w", err)
	}

	// Pool configuration
	config.MaxConns = 10
	config.MinConns = 2
	config.MaxConnLifetime = time.Hour
	config.MaxConnIdleTime = 30 * time.Minute

	// Create the pool
	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("error creating connection pool: %w", err)
	}

	// Verify the connection
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("unable to ping database: %w", err)
	}

	return pool, nil
}

// Function to add e-mail and password
func createAction(dbPool *pgxpool.Pool, text string) error {
	sql := `
		INSERT INTO users_information (email, password)
		VALUES ($1,$2)
		RETURNING user_id
	`
	var id int
	err := pool.QueryRow(ctx, sql, text, false).Scan(&id)
	if err != nil {
		return fmt.Errorf("error creating e-mail and password: %w", err)
	}

	fmt.Printf("Created e-mail and password: %d\n", id)
	return nil
}

// function to show all the e-mail and password
func getAllAction(dbPool *pgxpool.Pool) ([]User, error) {
	sql := `
	SELECT user_id, email, password
	FROM users_information
	ORDER BY id
	`
	rows, err := pool.Query(ctx, sql) // Query to execute `SELECT` statement that returns multiple rows
	if err != nil {
		return nil, fmt.Errorf("error querying users_information: %w", err)
	}
	defer rows.Close() // Ensure that ressources are cleaned up when func exit

	var users []User
	for rows.Next() {
		var user User
		err := rows.Scan(
			&user.ID,
			&user.Email,
			&user.Password,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning user row: %w", err)
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating User rows: %w", err)
	}

	return users, nil

}

// function to print the email and the password of an user
func printAction(users []User) {
	if len(users) == 0 {
		fmt.Println("No information found")
		return
	}

	for _, user := range users {
		fmt.Printf("%s, %s \n",
			user.Email,
			user.Password,
		)
	}

}

func main() {

	// Initialize context
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Initialize the database
	dbPool, err := initDB(ctx)
	if err != nil {
		log.Fatalf("Database initialization failed: %v", err)
	}
	defer dbPool.Close() // Close connection pool when exiting

	// Handle SIGINT (Ctrl+C) for graceful shutdown
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, os.Interrupt)
		<-sigChan
		cancel()
	}()

	app := &cli.Command{
		Name:  "go-safecli",
		Usage: " simple CLI tool for simply storing and managing passwords and email credentials",
		Action: func(ctx context.Context, cmd *cli.Command) error {
			dbPool, err := initDB(ctx)
			if err != nil {
				log.Fatalf("Database initialization failed: %v", err)
			} else {
				fmt.Println("Success")
			}
			defer dbPool.Close()
			return nil
		},
		Commands: []*cli.Command{
			{
				Name:    "add",
				Aliases: []string{"a"},
				Usage:   "Add new email and password",
				Action: func(c context.Context, cmd *cli.Command) error {
					text := cmd.Args().First()
					if text == "" {
						return fmt.Errorf("email and password cannot be empty")
					}
					return createAction(dbPool, text)
				},
			},

			{
				Name:    "all",
				Aliases: []string{"l"},
				Usage:   "Show email and password",
				Action: func(ctx context.Context, cmd *cli.Command) error {
					users, err := getAllAction(dbPool)
					if err != nil {
						return err
					}
					fmt.Println("All user information:")
					printAction(users)
					return nil
				},
			},
		},
	}
	if err := app.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
