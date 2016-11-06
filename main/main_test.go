package main

func init() {
	go startServer()
	deleteAllBanks()
}

func panicOnErrInTest(err error) {
	if err != nil {
		panic(err)
	}
}
func panicOnHttpErrInTest(err *httpError) {
	if err != nil {
		panic(err)
	}
}
