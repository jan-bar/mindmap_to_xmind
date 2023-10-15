#!/bin/bash
# 第一步: docker build -t xgo .
# 第二步: docker run --rm -v $PWD:/root/xmind xgo
cd /root/xmind
GOOS=linux GOARCH=amd64 CGO_ENABLED=1 go build -ldflags "-s -w" -trimpath -o mindmap.linux
CC=x86_64-w64-mingw32-gcc-posix CXX=x86_64-w64-mingw32-g++-posix GOOS=windows GOARCH=amd64 CGO_ENABLED=1 go build -ldflags "-s -w" -trimpath -o mindmap_x64.exe
CC=i686-w64-mingw32-gcc-posix CXX=i686-w64-mingw32-g++-posix GOOS=windows GOARCH=386 CGO_ENABLED=1 go build -ldflags "-s -w" -trimpath -o mindmap_x86.exe
