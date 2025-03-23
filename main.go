package main

import (
	"context"
	"fmt"
	"log"
	"net/mail"
	"os"
	"os/signal"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
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
	_ = godotenv.Load()

	// Initialize the connection Pool
	config, err := pgxpool.ParseConfig(os.Getenv("DATABASE_URL"))
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
func createAction(pool *pgxpool.Pool, email, password string) error {
	sql := `
		INSERT INTO users_information (email, password)
		VALUES ($1, $2)
		RETURNING user_id
	`
	e, err := mail.ParseAddress(email)
	if err != nil {
		return fmt.Errorf("error incorrect mail format")
	}

	var id int
	err = pool.QueryRow(ctx, sql, e, password).Scan(&id)
	if err != nil {
		return fmt.Errorf("error creating e-mail and password: %w", err)
	}

	fmt.Println("E-mail and password created")
	return nil
}

// function to show all the e-mail and password
func getAllAction(pool *pgxpool.Pool) ([]User, error) {
	sql := `
	SELECT user_id, email, password
	FROM users_information
	ORDER BY user_id
	`
	rows, err := pool.Query(ctx, sql) // to execute `SELECT` statement that returns multiple rows
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

// function to delete email and password
func deleteAction(pool *pgxpool.Pool, id int) error {
	sql := `DELETE FROM users_information WHERE user_id = $1`

	del, err := pool.Exec(ctx, sql, id)
	if err != nil {
		return fmt.Errorf("error deleting user information: %w", err)
	}

	if del.RowsAffected() == 0 {
		return fmt.Errorf("no user information found with id %d", id)
	}

	fmt.Printf("user %d has been deleted successfully\n", id)
	return nil
}

// function to modify email and password
func changeAction(pool *pgxpool.Pool, id int, email, password string) error {
	sql := `UPDATE users_information
			SET email = $1, password = $2
			WHERE user_id = $3
			RETURNING *
	`
	change, err := pool.Exec(ctx, sql, email, password, id)
	if err != nil {
		return fmt.Errorf("error updating the email and the password: %w", err)
	}

	if change.RowsAffected() == 0 {
		return fmt.Errorf("no user's found with id: %d", id)
	}

	return nil
}

// function to print the email and the password of an user
func printAction(users []User) {
	if len(users) == 0 {
		fmt.Println("No information found")
		return
	}

	t := table.NewWriter()
	t.SetTitle("Users information")
	t.Style().Format.Header = text.FormatTitle

	t.AppendHeader(table.Row{"ID", "Email", "Password"})

	for _, user := range users {
		t.AppendRow(table.Row{user.ID, user.Email, user.Password})
	}

	fmt.Println(t.Render())
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
		Name:  "go-not-safecli",
		Usage: " simple CLI tool for simply storing and managing passwords and email credentials",
		Commands: []*cli.Command{
			{
				Name:    "add",
				Aliases: []string{"a"},
				Usage:   "Add new email and password",
				Action: func(c context.Context, cmd *cli.Command) error {
					args := cmd.Args().Slice()
					if len(args) < 2 {
						return fmt.Errorf("email and password cannot be empty")
					}
					email := args[0]
					password := args[1]
					return createAction(dbPool, email, password)
				},
			},

			{
				Name:    "rm",
				Aliases: []string{"r"},
				Usage:   "Remove the user's information",
				Action: func(ctx context.Context, c *cli.Command) error {
					idStr := c.Args().First()
					if idStr == "" {
						return fmt.Errorf("user's ID required")
					}

					var id int
					if _, err := fmt.Sscanf(idStr, "%d", &id); err != nil {
						return fmt.Errorf("invalid user ID: %s", idStr)
					}

					return deleteAction(dbPool, id)
				},
			},

			{
				Name:    "md",
				Aliases: []string{"m"},
				Usage:   "Modify password and email",
				Action: func(ctx context.Context, cmd *cli.Command) error {
					args := cmd.Args().Slice()
					if len(args) < 2 {
						return fmt.Errorf("email and password cannot be empty")
					}

					var id int
					args_id := cmd.Args().First()
					if args_id == "" {
						return fmt.Errorf("user's ID required")
					}
					if _, err := fmt.Sscanf(args_id, "%d", &id); err != nil {
						return fmt.Errorf("invalid user's ID: %s", args_id)
					}

					email := args[1]
					password := args[2]
					return changeAction(dbPool, id, email, password)
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
