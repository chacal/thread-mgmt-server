#!/usr/bin/env bash
function stop_server {
  echo -n "Stopping server with PID $MAIN_PID.. "
  kill "$MAIN_PID"
  echo "done."
}
trap stop_server EXIT

go build -o bin cmd/mgmt-server/*
./bin/main --mqtt-broker "" --mqtt-username "" --mqtt-password pass &
MAIN_PID=$!

npx webpack serve

