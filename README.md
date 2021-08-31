# mtgobananium

## Installation

First of all you need to clone the Tdlib repo and build it:
```bash
git clone git@github.com:tdlib/td.git --depth 1
cd td
mkdir build
cd build
cmake -DCMAKE_BUILD_TYPE=Release ..
cmake --build . -- -j5
make install

# -j5 refers to number of your cpu cores + 1 for multi-threaded build.
cd ../../
# Build project
git clone https://github.com/zombiedevel/mtgobananium.git
cd mtgobananium
go mod download
go build -o mtgobananium cmd/main.go

```

## Docker
You can use prebuilt tdlib with following Docker image:

***Linux:***
``` shell
git clone https://github.com/zombiedevel/mtgobananium.git 
cd mtgobananium
docker build -f.DockerFile -t bananium .
```
##On build success
```docker run -d bananium -app-id <app id> -app-hash <app hash> -token <bot token>```