package pkg

// func /////////////////////////////////////////

// Convenience function to return a single value
func Must[T any](value T, err error) T {
	if err != nil {
		panic(err)
	}

	return value
}
