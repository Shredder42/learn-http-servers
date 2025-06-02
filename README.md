# Chirpy Server

This project creates a server for Chirpy. Chirpy is a social network similar to Twitter. It allows users to create accounts and login, post, view, and delete Chirps. Users can also upgrade to Chirpy Red status.

The server runs on a local machine. HTTP requests are made to the server over the [localhost](https://www.hostinger.com/tutorials/what-is-localhost) in port 8080.

## Running the Chirpy Server

Chirpy requires [Go](https://go.dev/doc/install) version 1.22 or higher and a postgres database

## Environment and Config

You will need a .env file that contains private environment variables necessary for the Chirpy server to function.

* DB_URL - url string to connect to the postgres database
    It will look something like this `postgres://postgres:postgres@localhost:5432/chirpy?sslmode=disable`
* PLATFORM - specifies what platform is accessing the server
* SECRET - an internal JWT used for authentication and authorization (Don't share this!!!)
* POLKA_KEY - an api key the polka client uses to update to chirpy red status

You can load these variables into your main file with: 
```go
    godotenv.Load()
```
Then an example of loading the **dbURL**:
```go
    dbURL := os.Getenv("DB_URL")
    if dbURL == "" {
        log.Fatalf("DB_URL must be set")
    }
```

Then create an **apiConfig** in your main file like so:
```go
	apiCfg := apiConfig{
		polkaKey:       polkaKey,
		jwtSecret:      jwtSecret,
		platform:       platform,
		db:             dbQueries,
		fileserverHits: atomic.Int32{},
	}
```

## Database

You will need to create the database in Postgres. Once you start the postgres server `sudo service postgresql start` on linux or `brew services start postgresql@15` on Mac, you can create the database using this command 
```sql
    CREATE DATABASE chirpy;
```

Then you will need to run the database migrations in the sql/schema folder. To do this, you will need [Goose](https://github.com/pressly/goose), which can be installed using `go install github.com/pressly/goose/v3/cmd/goose@latest`
Then you can use a command like below to continue through the database migrations and set up the database.:
```bash
goose postgres "postgres://postgres:postgres@localhost:5432/chirpy" up
```


## Start the server

Run `go run main.go` to start the server.

This is a guided project from boot.dev.

