
@test "can run the build" {
  docker run --rm local/stripseven
}

@test "can run the build with arguments" {
  docker run --rm local/stripseven run -h
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
  docker run --rm --network "$name" \
    -- alpine nc -zv "$name" 1221
}

teardown_file() {
  name=test-local-stripseven
  docker rm -f "$name" ||:
  docker network rm "$name" ||:
}