
@test "can reach the echoserver" {
  nc -zv localhost 2000
}

@test "can verify the echo" {
  echo hello | nc localhost 2000 | grep hello
}
