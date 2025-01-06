# Chirpy

Chirpy is a fictional Twitter clone, built to practice and enhance my GoLang and HTTP server skills and knowledge.

## Tools & Tech

- **Language:** Go
- **Database:** Postgres
- **DB Migration:** Goose
- **SQL Generation:** SQLC

## API Spec

Chirpy uses a RESTful style of API to serve its data.

### App

- **/app/** is a simple file server

### Admin

#### /metrics

- **GET /admin/metrics** returns details about the number of times the file server has been visited

#### /reset

- **GET /admin/reset** simply resets the metrics held in server memory

### API

#### /healthz

- **GET /healthz** returns the health status of the service

#### /chirps

- **GET /api/chirps** serves all existing Chirps
- **GET /api/chirps/{chirpID}** serves an existing Chirp
- **POST /api/chirps** accepts the creation of a new Chirp
- **DELETE /api/chirps/{chirpID}** deletes an existing Chirp

#### /users

- **GET /api/users** serves all existing users
- **POST /api/users** accepts the creation of a new user
- **PUT /api/users** updates an existing user's details

#### /login

- **POST /api/login** allows a user to log in

#### /refresh

- **POST /api/refresh** accepts a refresh token and returns a new access token

#### /revoke

- **POST /api/revoke** revokes a user's refresh token
