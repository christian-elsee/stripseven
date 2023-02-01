
@test "can run the build" {
  docker run -it --rm local/stripseven
}

@test "can run the build's run command" {
  docker run -it --rm local/stripseven run -h
}

@test "can run the build as a proxy" {
  name=test-local-stripseven
  docker network create "$name" ||:
  docker rm -f "$name" ||:
  docker run \
    -d \
    --rm \
    --name "$name" \
    --network "$name" \
    --  local/stripseven run \
          -l 0.0.0.0:1221 \
          -t localhost:1234
  docker run -it --rm --network "$name" \
    -- alpine nc -zv "$name" 1221

  docker rm -f "$name"
}

teardown_file() {
  docker rm -f test-local-stripseven ||:
  docker network rm test-local-stripseven ||:
}