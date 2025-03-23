# go-not-safecli

> [!WARNING]
>THIS PROJECT IS FOR LEARNING PURPOSE ONLY DO NOT USE IT TO STORE SENSITIVE INFORMATION

go-not-safe-cli is a simple command-line application designed for storing email and password credentials.

## Installation

To install and run go-not-safecli, follow these steps:

### Prerequisites

- Go version 1.20 and above
- PostgreSQL

### Steps

1. Clone the repository:

    ```bash
    git clone https://github.com/disikoX/go-not-safe-cli.git
    cd go-not-safe-cli
    ```

2. Create a `.env` file that looks like this:

    ```bash
    host     = ""
    port     = 
    user     = ""
    password = ""
    dbname   = ""
    database_url = "postgres://username:dbname@host:port/password"
    ```

3. Install all the Go dependencies with :

    ```bash
    go mod download
    ```

3. Build the app with :

    ```bash
    go build -o go-not-safecli
    
    ```

3. Run the app with :

    ```bash
    go build -o go-not-safecli
    chmod +x go-not-safecli 
    
    ```

## Usage

Once the application is running, you can use the following commands:

### Add Credentials

To add a new email and password, use the following command:

```bash
./go-not-safe-cli add <your_email> <your_password>
```

### Show Credentials

To show credentials in a table, use the following command:
```bash
./go-not-safe-cli all 
```

### Remove Credentials

To remove an existing email and password, use the following command:

```bash
./go-not-safe-cli rm <your_id> 
```

### Modify Credentials

To modify an existing email and password, use the following command:

```bash
./go-not-safe-cli md <your_id> <your_email> <your_password>
```

## Contributing

Contributions are welcome! Please open an issue or submit a pull request for any improvements or bug fixes.

## License

This project is licensed under the MIT License. See the [MIT](https://www.google.com/url?sa=t&source=web&rct=j&opi=89978449&url=https://opensource.org/license/mit&ved=2ahUKEwiso6WWhoiMAxVzT0EAHYcyHhAQFnoECBYQAQ&usg=AOvVaw0JouoMsOReC1lXVEak9dPg) file for details.
