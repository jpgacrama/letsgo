# SnippetBox: Generate Certificate file for GoLang

generate_cert.go is a Symbolic Link to the actual file located within your file system.
To change where the link goes, type:
```
ln -s <Absolute Path to file to be pointed to> <name of the file in the current directory>
```

## How to generate:

Type the following command in terminal
```
 go run generate_cert.go --rsa-bits=2048 --host=localhost
 ```