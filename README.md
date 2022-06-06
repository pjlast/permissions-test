# Permissions PoC

This repo aims to serve as a proof of concept for the proposed permission system.

## Setup

In order to run this project, you'll have to create a Postgres DB for the permissions table. You should be able to do this with `psql`:

```
psql

CREATE DATABASE permissions_test;
```

Next you need to point the project to this Database with a postgres connection URL.

Create a `.env` file in the repository root directory, and give it the following entry:

```
TEST_DATABASE_URL=postgresql://username:password@localhost:5432/permissions_test?sslmode=disable
```

obviously substituting the username/password/db host/database name as necessary.

## Running

You should now be ready to just run the project:

```
go run .
```