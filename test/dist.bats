# sanity check dist

@test "can validate manifest yaml" {
  <dist/manifest.yaml yq -re .
}

@test "can lint golang" {
  true
}
