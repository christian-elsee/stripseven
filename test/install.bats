
@test "can reach the echo server" {
  kubectl exec -it nc -zv localhost 12321
}

@test "can verify the echo" {
  echo hello | nc localhost sq | grep hello
}

@test "can proxy the echo server" {
  timeout 5 dist/build run -t localhost:2000 &>/dev/null &
  sleep 2
  curl -s http://localhost:8080 -d 'hello' | grep -q hello
}