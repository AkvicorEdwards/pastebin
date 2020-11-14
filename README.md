# PasteBin

**PasteBin** is a webapp for code share.

After specifying the working directory in `config.yaml`, you can browse all the folders and files in the directory through the web. You can also upload files to this folder via the web

## Build and Run

```bash
git clone https://github.com/AkvicorEdwards/pastebin.git
cd pastebin
go mod download
go build pastebin.go
./pastebin

# ListenAndServe: localhost:8080
```

## Usage

```bash
# Full content
localhost:8080/1

# Text content
localhost:8080/raw/1
# if password == 12345
localhost:8080/raw/1?pwd=12345

```

