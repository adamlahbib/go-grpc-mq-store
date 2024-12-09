package main

func main() {
	// init consumers and producers
	go initConsumer()
	go initProducer()
	// init api
	initApi()
}
