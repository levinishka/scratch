# scratch
Template Golang service and libraries.

## Getting started
### How to install
Clone scratch
```shell
git clone https://github.com/causelovem/scratch.git
```
You can either compile scratch by yourself:

using make
```shell
cd scratch
make build
```
manual
```shell
cd scratch
go build -o cmd/bin/scratch cmd/scratch/main.go # for Linux
GOOS=windows GOARCH=amd64 go build -o cmd/bin/scratch.exe cmd/scratch/main.go # for Windows
```
Or use precompiled binaries in `cmd/bin` folder:
* `scratch` for **Linux**
* `scratch.exe` for **Windows**

### How to use
Check help
```shell
./cmd/bin/scratch -help
```

Scratch creates template for Golang service:

1. create service from scratch
   1. use `-project` parameter to specify new project's path: last element in a path will be new project's name
   2. use `-repo` parameter to specify git repository path of new project
2. initialize go modules
3. write your own logic using template service
4. test it
5. initialize git and push
6. ???
7. PROFIT

#### Examples
Creates new project with name `testProject`
```shell
# create project from scratch
chmod a+x ./cmd/bin/scratch
./cmd/bin/scratch -project /absolute/path/to/testProject -repo github.com/causelovem

# initialize go modules
cd /absolute/path/to/testProject
go mod init github.com/causelovem/testProject
go mod tidy

# write logic and then test service
make lint
make test
make build
make test-run

# initialize git and push
git init
git add --all
git commit -m "Initial Commit"
git remote add origin github.com/causelovem/testProject.git
git push -u origin master
```

## Libraries
Scratch contains some useful libraries which you can import and use:
* `config` simply reads config file in JSON format and unmarshal it to a structure
* `logger` provides preconfigured zap-logger
* `router` provides mux router with pprof handlers added
* `server` provides http server with graceful shutdown
* `metrics` provides prometheus http server with basic service metrics

To import any of these packages use `"github.com/causelovem/scratch/pkg/PACKAGE"`
