package main

import (
	"bufio"
	"fmt"
	"os"
	"time"

	"gopkg.in/mgo.v2"
)

var (
	attempts = 0
	found    = 0
	file     *os.File
)

func main() {
	var err error
	file, err = os.Create("result.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	ips := ReadAllLines("ips.txt")
	ipChannel := make(chan string)
	fmt.Print("Enter Threads Counts: ")
	t := 0
	fmt.Scanln(&t)
	for i := 0; i != t; i++ {
		go thread(ipChannel)
	}
	go printer()
	for _, line := range ips {
		ipChannel <- line
	}
	for attempts + found != len(ips) {
		time.Sleep(time.Second)
	}
	fmt.Println("Done")
}

func printer() {
	for {
		fmt.Printf("\rattempts: %d, found: %d   ", attempts, found)
		time.Sleep(time.Millisecond * 100)
	}
}

func thread(ipChannel chan string) {
	for {
		ip := <-ipChannel
		if checkMongo(ip) {
			found++
			fmt.Println("\r" + ip + " UnAuthinticated MongoDB\t\t\t\t")
			file.WriteString(ip + "\r\n")
		} else {
			attempts++
		}
	}
}

func checkMongo(ip string) bool {
	sess, err := mgo.Dial(fmt.Sprintf("mongodb://%s:27017", ip))
	if err != nil {
		return false
	}
	defer sess.Close()
	_, err = sess.DatabaseNames()
	return err == nil
}

func ReadAllLines(path string) (lines []string) {
	file, err := os.Open(path)
	if err != nil {
		return nil
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines
}
