# vama-api

# API Architecture

This API is built with the Clean architecture pattern in mind and is based on entities. To read more, go to this website: https://github.com/bxcodec/go-clean-arch

**Each entity lives in internal/entities/ and contains its own HTTP endpoints as well as the following:**

- usecase.go - **PUBLIC** interface for communicating with other entites and getting access via any transports like HTTP/GRPC/etc...
- repository.go - **INTERNAL** interface for accessing data at DB/JSON/STORAGE/etc
- structure.go - **INTERNAL** models for communicating in usecase
- controller - configurations for transport
  - controller/decoders.go - decoders input model
  - controller/endpoints.go - map transport to usecase
  - controller/index.go - configuration for transport
- models - **PUBLIC** models for incoming requests and outgoing responses
- repositories - repository **Each func at separate file**
- repositories - usecase **Each func at separate file**

# Migrations

At one point, we migrated from MySQL to PostgreSQL. The SQL statements that weused for this migrations live at `./migrations`

The SQL files in this folder correspond to PostgreSQL schemas in our DB.

# Cloud Tasks & Cloud Scheduler Jobs
Any endpoints that have the following `Version` with the following values in their `index.go` file are not user facing endpoints and are executed by some external system (eg. Cloud Task, Cloud Scheduler, etc.).

```
constants.CLOUD_SCHEDULER_V1
constants.CLOUD_TASK_V1
```

As such, if the URL is changed in the `index.go`, one must also change the URL in terraform, since this is where we instantiate and adjust our GCP cloud infrastructure.

**_NOTE: All scripts should be idempotent_**

# Installation

# Installing local database
1. Install Postgres 13 `https://www.postgresql.org/download/`.  Select `macOS` and then `Download the installer`
2. Run the installer.  Default port `5432` is fine.
3. Skip the stack builder option.
4. Add `export PATH="/Library/PostgreSQL/13/bin:$PATH"` to your `~/.bash_profile`
5. Install brew `https://brew.sh/`
6. Install libpq `brew doctor; brew update; brew install libpq; brew link --force libpq`
7. Ask for an engineer for a fresh dump of the database table schemas
8. Go to the directory of the dump and execute `psql -f <nameOfDump.sql> -d postgres -U postgres`

# Installing Redis
There are many instances in which we cache results in order to improve performance, hence, we have a SaaS Redis caching our entries.

1. Download Docker `https://docs.docker.com/desktop/mac/install/` and input the default params in the installer
2. Run the following commands to run a Docker container that already has all of the local Redis dependencies (including `redis-cli`) installed
  ```
    docker pull redis
    docker run -d --name vama-redis redis
  ```

3. Run `docker exec -it vama-redis bash` to start the container
4. To connect to a redis instance, run the command `redis-cli -h <endpoint> -p <port> -a <password>`
5. Run the command `ping` to prove to yourself that you've connected successfully

# Installing API dependencies
1. Install Go `https://golang.org/doc/install`
2. Clone `https://github.com/VamaSingapore/vama-api`
3. Place the `creds.json` and `.env.local.rc` files in the root of your directory (ask an engineer to securely send these files to you)
4. Execute `source .env.local.rc` in the root directory

# Running API locally
1. Run `make build` to compile and then `./api` to run the server locally
2. Run `make test` to execute local tests
3. (Optional) If you are getting a password error after running `make test` and then `make build`
  1. `sudo -u postgres psql`
  2. type in your computer password
  3. once the postgres command line presents itself, type `\password postgres`
  4. enter the password `pass`
  5. verify the password `pass`
  6. quit with `\q`
  7. try running `make build` and `make test` again

# Flags
Command line flags:

| Flag    | Type     | Default  | Description |
|   ---   | ---      | ---      | ---         |
| -silent  | bool     | false   | When passed, turn off request/response logging  |

For example: 
`make build; ./api -silent`

# Tasks

Tasks live in https://finchat.atlassian.net/secure/RapidBoard.jspa?rapidView=3&projectKey=API

# Documentation

API documentation lives in https://documenter.getpostman.com/view/18713237/UVRBnmg5

**_NOTE: After ANY change to API, you should run it in Postmann and click "Save", because the same documentation is used by mobile developers_**

# Process

- Create and assign tasks in Jira
- Move tickets to `In Progress`
- Create branch off of `staging` with the following naming conventions
  - feature/API-XXX - for new features
  - hotfix/20210621-1 - for hotfixes
  - release release/X.Y.Z
- After a task is finished, you should create a `Pull Request`, and request Code Review
- Update DB (migrations, in the future this will be automated)
- Update GitBook documentation
- Move ticket to QA

# Testing

Manual

In order to manually test the API, you will need to download Postman and be added to the `vama-apis` collection. This will give you access to our HTTP requests that we use to test the API.

**_NOTE: Each API folder at Postman contains a Pre-request Script for authorization_**

Once you are added to the collection, you will then need to update your Postman environment so that you have the correct tokens and identification when you make HTTP requests to the API.

We have four environments in Postman: `prod`, `staging`, `dev`, `local`. `dev` is meant for sending HTTP requests to our `api-dev` deployment in GCP, whereas, `local` is meant for hitting your locally running server at `localhost:8080`.

When you navigate to the “Environments” tab in Postman, you will need to fill out the following fields for all environments:

```
baseUrl
apiKey
idToken
idTokenExpiresAt
testUserEmail
testUserPassword
```

Ask one of the developers for the values for each of the environments listed above and they will make sure to send them over to you in a secure fashion.

Once the values are set, you can start making HTTP requests via Postman by navigating to the endpoint to which you want to send an HTTP request to and press the “Send” button. You should then start receiving responses from the API.

# Test Suite

You may need to `chmod +x ./cmd/test/runtests.sh`
Run `make test`

You must have psql installed.
