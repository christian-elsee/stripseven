
@test "can reach the echo server" {
  nc -zv echo $ECHO_SERVICE_PORT
}

@test "can verify the echo" {
  echo hello | nc echo $ECHO_SERVICE_PORT | grep hello
}

@test "can proxy the echo server" {
  timeout 5 dist/build run -t localhost:2000 &>/dev/null &
  sleep 2
  curl -s http://localhost:8080 -d 'hello' | grep -q hello
}
