# SnippetBox using Let's Go
- [SnippetBox using Let's Go](#snippetbox-using-lets-go)
- [1. Introduction](#1-introduction)
  - [1.1. Pre-requisites](#11-pre-requisites)
  - [1.2 Running the program](#12-running-the-program)

# 1. Introduction

This is my own implementation SnippetBox using [Let's Go](https://lets-go.alexedwards.net)

## 1.1. Pre-requisites

Download the CSS and Javascript files by running this in your terminal `$ curl https://www.alexedwards.net/static/sb.v120.tar.gz | tar -xvz -C ./ui/static`

## 1.2 Running the program
1. Run the Web Server using this command `go run cmd/web/* -port=":4000"`
2. Curl to the server using this command `curl -iL -X POST http://localhost:4000/snippet/create`
3. See the contents of mysql using these commands
    - Start MySQL: `mysql -D snippetbox -u web -p `
    - Check its contents: `SELECT id, title, expires FROM snippets;`