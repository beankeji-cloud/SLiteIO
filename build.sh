#!/bin/bash

go build -ldflags '-s -w ' -o _build/node-disk-controller cmd/controller/main.go
go build -ldflags '-s -w ' -o _build/scheduler-plugin cmd/scheduler/main.go
ls -l _build
scp _build/node-disk-controller _build/scheduler-plugin 10.244.68.65:/opt/lei/build_image_liteio/
