
@test "can run the build" {
  dist/build
}

@test "can run the build's run command" {
  dist/build run -h
}
