
@test "can reach the echo server" {
  nc -zv echo $ECHO_SERVICE_PORT
}

@test "can verify the echo" {
  echo hello | nc echo $ECHO_SERVICE_PORT | grep hello
}

@test "can proxy the echo server" {
  curl -s http://stripseven -d 'hello' | grep -q hello
}
