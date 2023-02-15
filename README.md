# SnippetBox using Let's Go
- [SnippetBox using Let's Go](#snippetbox-using-lets-go)
- [Introduction](#introduction)
  - [Pre-requisites](#pre-requisites)
  - [Running the program](#running-the-program)
  - [Running Code Coverage](#running-code-coverage)
  - [Appendix](#appendix)
    - [Setting up a MySQL Server using GitPod](#setting-up-a-mysql-server-using-gitpod)

# Introduction

This is my own implementation SnippetBox using [Let's Go](https://lets-go.alexedwards.net)

## Pre-requisites

Download the CSS and Javascript files by running this in your terminal `$ curl https://www.alexedwards.net/static/sb.v120.tar.gz | tar -xvz -C ./ui/static`

## Running the program
1. Run the Web Server using this command `go run cmd/web/* -port=":4000"`
2. Curl to the server using this command `curl -iL -X POST http://localhost:4000/snippet/create`
3. See the contents of mysql using these commands
    - Start MySQL: `mysql -D snippetbox -u root -p `
    - Check its contents: `SELECT id, title, expires FROM snippets;`

## Running Code Coverage
Execute the following statements to generate a Code Coverage Report
```
go test -coverpkg=./... -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

## Appendix

### Setting up a MySQL Server using GitPod
1. Open the MySQL configuration file in a text editor of your choice. In this example, we use nano:
    ```
    sudo nano /etc/mysql/my.cnf
    ```
2. Then, add the following lines at the end of the MySQL configuration file:
    ```
    [mysqld]
    socket=[path to mysqld.sock]
    [client]
    socket=[path to mysqld.sock]
    ```

    E.g.
    ```
    [mysqld]
    socket=/run/mysqld/mysqld.sock
    [client]
    socket=/run/mysqld/mysqld.sock
    ```

    Note:
    * I already placed the my.cnf here. You can just copy the file to `/etc/mysql`

3. Go to `./snippetbox` directory
4. Run the SQL scripts by typing the following commands:
  ```
  sudo mysql
  source ./db/snippetbox.sql
  source ./db/newUser.sql
  ```