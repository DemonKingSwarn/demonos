package main

/*
void kHalt(void);
*/
import "C"

//go:nosplit
func halt() {
	C.kHalt()
}
