package main

func main() {
	// init consumers and producers
	go initConsumer()
	initProducer()
}
