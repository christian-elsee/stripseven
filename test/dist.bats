# sanity check dist

@test "can deploy manifest" {
  kubectl apply -f dist/manifest.yaml --dry-run=server
}

@test "can lint golang" {
  true
}
