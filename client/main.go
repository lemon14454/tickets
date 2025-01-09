package main

import (
	"flag"
	"fmt"
	"strconv"
	"strings"
	"ticket/client/model"
	"time"
)

var zoneStringInvalidErr = fmt.Errorf("provided zone string isn't valid")

var serverIP string
var serverPort int
var username string
var password string
var zoneString string // ex: A,5,5|B,5,5|C,5,5
var action string     // create or simulate

// e.g.
// ./client -ip 127.0.0.1 -port	8080 -user user@mail.com -password 12345678

// 1. Shell Script generate multiple user
// 2. Create Event with CLI
// 3. Simulate Massive ticket buying request with CLI

func init() {
	flag.StringVar(&serverIP, "ip", "127.0.0.1", "server ip [default:127.0.0.1]")
	flag.IntVar(&serverPort, "port", 8080, "server port [default:8080]")

	flag.StringVar(&username, "user", "", "username [auto register if not exists]")
	flag.StringVar(&password, "pwd", "", "password [atleast 8 characters]")

	flag.StringVar(&zoneString, "zone", "A,10,10|B,10,10|C,10,10", "zone string [zone/rows/seats|zone/rows/seats ...]")
	flag.StringVar(&action, "action", "create", "action [create or simulate]")
}

func authUser(client *Client) error {
	err := client.Login(username, password)

	if err != nil {
		// User Probably not created
		err := client.Register(username, password)
		if err != nil {
			return err
		}

		// Retry
		err = client.Login(username, password)
		if err != nil {
			return err
		}
	}

	fmt.Println("User Successfully Logged In")
	return nil
}

func create(client *Client, zoneString string) {
	fmt.Println("Start Event Creation")
	name := fmt.Sprintf("Event-%s", time.Now().String())

	zones := make([]model.EventZone, 0)

	if len(zoneString) < 5 {
		// best validation BTW
		fmt.Println(zoneStringInvalidErr)
		return
	}

	for _, zone := range strings.Split(zoneString, "|") {
		section := strings.Split(zone, ",")
		if len(section) != 3 {
			fmt.Println(zoneStringInvalidErr)
			return
		}
		rows, err1 := strconv.ParseInt(section[1], 10, 32)
		seats, err2 := strconv.ParseInt(section[2], 10, 32)

		if err1 != nil || err2 != nil {
			fmt.Println(zoneStringInvalidErr)
			return
		}

		zones = append(zones, model.EventZone{
			Zone:  section[0],
			Rows:  int32(rows),
			Seats: int32(seats),
			Price: 1000, // hard code
		})
	}

	event, err := client.Request.CreateEvent(name, zones)
	if err != nil {
		fmt.Printf("Create Event Request Err: %v", err)
		return
	}

	fmt.Printf("Event Create Success: %v", event)
}

func main() {
	flag.Parse()
	client := NewClient(serverIP, serverPort)

	err := authUser(client)
	if err != nil {
		fmt.Printf("Auth Error: %v \n", err)
		return
	}

	switch action {
	case "create":
		create(client, zoneString)
	case "simulate":
		fmt.Println("Start Ticket Buying Simulation")
	default:
		fmt.Println("Action doesn't exists, please try 'create' or 'simulate' ")
	}
}
