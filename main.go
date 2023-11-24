package main

import (
	"fmt"
	"math/rand"
	"net"
	"os"
	"strconv"
	"time"
)

func philosopher(hungry chan<- string, done chan<- string, eat <-chan string, msgSenderCh chan<- string) {
	s := os.Getenv("EAT_AMOUNT")
	eatcount, _ := strconv.Atoi(s)
	for {

		if eatcount > 0 {
			am := rand.Intn(4-1+1) + 1
			fmt.Println("Thinking for: ", am, "seconds")
			time.Sleep(time.Duration(am) * time.Second)
			hungry <- "hungry"
		} else {
			fmt.Println("Telling my neighbours that I'm done")
			msgSenderCh <- "finish"
			return
		}

		if <-eat == "eat" {
			time.Sleep(time.Second)
			done <- "done"
			eatcount--
			fmt.Println("I HAVE ATE: ", 5-eatcount, " times")
		}
	}

}

func waiter(hungry <-chan string, done <-chan string, eat chan<- string, msgSenderCh chan<- string, rcvStickR <-chan string, rcvStickL <-chan string, passR <-chan string, passL <-chan string) {
	id := os.Getenv("id")
	haveL, haveR := true, false
	dirtyL, dirtyR := true, false
	if id == "5" {
		haveL, haveR = true, true
		dirtyL, dirtyR = true, true
	} else if id == "1" {
		haveL, haveR = false, false
		dirtyL, dirtyR = false, false
	}

	requestedL, requestedR := false, false
	for {
		select {
		case <-hungry:
			{
				fmt.Println("the philosopher is hungry")
				fmt.Println("my sticks are: ", haveL, haveR)
				if haveL && haveR {
					fmt.Println("EATING")
					eat <- "eat"
					dirtyR, dirtyL = true, true
					<-done
					fmt.Println("Done Eating")

				} else {

					if !haveL {
						fmt.Println("pls pass me my LEFT")
						msgSenderCh <- "passMyLStickPls"
					}
					if !haveR {
						fmt.Println("pls pass me my RIGHT")
						msgSenderCh <- "passMyRStickPls"
					}
					for !haveL || !haveR {
						fmt.Println("waiting for my sticks")
						select {
						case <-rcvStickR:
							fmt.Println("Received stick from RIGHT")
							haveR = true
							dirtyR = false
						case <-rcvStickL:
							fmt.Println("Received stick from LEFT")
							haveL = true
							dirtyL = false
						case <-passR:
							{
								requestedR = true
							}
						case <-passL:
							{
								requestedL = true
							}
						}

					}
					if haveL && haveR {
						fmt.Println("EATING")
						eat <- "eat"
						dirtyR, dirtyL = true, true
						<-done
						fmt.Println("Done Eating")
					}
				}
			}
		case <-passR:
			{
				requestedR = true
			}
		case <-passL:
			{
				requestedL = true
			}
		}
		fmt.Println("have they asked me for my sticks?: ", requestedL, requestedR)
		if requestedR {

			if haveR && dirtyR {
				fmt.Println("passing my RIGHT fork")
				dirtyR = false
				haveR = false
				msgSenderCh <- "takeMyRstick"
				requestedR = false
			}
		}
		if requestedL {

			if haveL && dirtyL {
				fmt.Println("passing my LEFT fork")
				dirtyL = false
				haveL = false
				msgSenderCh <- "takeMyLstick"
				requestedL = false
			}
		}

	}

}

func sender(destAddrR string, destAddrL string, msgSenderCh <-chan string, finish_iter chan<- string) {

	connR, _ := net.Dial("udp", destAddrR)
	defer connR.Close()
	connL, _ := net.Dial("udp", destAddrL)
	defer connL.Close()

	for {
		message := <-msgSenderCh
		switch message {
		case "passMyLStickPls":
			_, _ = connL.Write([]byte(message))

		case "passMyRStickPls":
			_, _ = connR.Write([]byte(message))

		case "takeMyLstick":
			_, _ = connL.Write([]byte(message))
		case "takeMyRstick":
			_, _ = connR.Write([]byte(message))

		case "finish":
			_, _ = connL.Write([]byte("finishL"))
			_, _ = connR.Write([]byte("finishR"))
			finish_iter <- "finishL"
		}
	}
}
func receiver(port int, rcvStickR chan<- string, rcvStickL chan<- string, passR chan<- string, passL chan<- string, finish_iter chan<- string) {
	addr1 := fmt.Sprintf(":%d", port)
	addr, _ := net.ResolveUDPAddr("udp", addr1)

	listener, _ := net.ListenUDP("udp", addr)
	defer listener.Close()

	buffer := make([]byte, 1024)
	for {
		n, _, _ := listener.ReadFromUDP(buffer)
		message := string(buffer[:n])

		switch message {
		case "takeMyLstick":
			rcvStickR <- message
		case "takeMyRstick":
			rcvStickL <- message
		case "passMyLStickPls":
			passR <- message
		case "passMyRStickPls":
			passL <- message
		case "finishL":
			finish_iter <- "finishR"
		case "finishR":
			finish_iter <- "finishmi"
		}
	}
}

func main() {

	hungry := make(chan string)
	done := make(chan string)
	eat := make(chan string)
	rcvStickR := make(chan string)
	rcvStickL := make(chan string)
	passR := make(chan string)
	passL := make(chan string)
	msgSenderCh := make(chan string, 1024)

	rightIP := os.Getenv("R_IP")
	leftIP := os.Getenv("L_IP")
	port := 8080
	rightAdress := fmt.Sprintf("%s:%d", rightIP, port)
	leftAdress := fmt.Sprintf("%s:%d", leftIP, port)

	finish_iter := make(chan string)

	go sender(rightAdress, leftAdress, msgSenderCh, finish_iter)
	go receiver(port, rcvStickR, rcvStickL, passR, passL, finish_iter)
	go waiter(hungry, done, eat, msgSenderCh, rcvStickR, rcvStickL, passR, passL)
	go philosopher(hungry, done, eat, msgSenderCh)

	left_finish := false
	right_finish := false
	finish_mi := false
	for !(left_finish && right_finish && finish_mi) {
		message := <-finish_iter
		switch message {
		case "finishR":
			right_finish = true
		case "finishL":
			left_finish = true
		case "finishmi":
			finish_mi = true
		}

	}
	fmt.Println("I ate all my food, and both my neighbours too, I'm leaving now and remember kids never stop thinkin!")

}
