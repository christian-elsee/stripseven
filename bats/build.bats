
@test "can run the build" {
  dist/build
}

@test "can run the build's run command" {
  dist/build run -h
}

@test "can proxy the echo server" {
  timeout 5 dist/build run -t localhost:2000 &>/dev/null &
  sleep 2
  curl -s http://localhost:8080 -d 'hello' | grep -q hello
}