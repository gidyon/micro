package errs

// Panic on error
func Panic(err error) {
	if err != nil {
		panic(err)
	}
}
