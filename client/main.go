package main

import (
	"flag"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"ticket/client/model"
	"time"
)

var zoneStringInvalidErr = fmt.Errorf("provided zone string isn't valid")

var serverIP string
var serverPort int

var username string = "testuser0"
var password string = "12345678"

var action string
var zoneString string // ex: A,5,5|B,5,5|C,5,5
var attempts int
var eventID int
var sameUser bool

// Check Makefile for usage

func init() {
	flag.StringVar(&serverIP, "ip", "127.0.0.1", "server ip [default:127.0.0.1]")
	flag.IntVar(&serverPort, "port", 8080, "server port [default:8080]")

	flag.StringVar(&action, "action", "simulate", "event | user | simulate")

	flag.StringVar(&zoneString, "zone", "A,10,10|B,10,10|C,10,10", "zone string [zone/rows/seats|zone/rows/seats ...]")
	flag.IntVar(&attempts, "attempt", 100, "attempt [default 100 times]")
	flag.IntVar(&eventID, "event_id", 1, "eventID [default 1]")
	flag.BoolVar(&sameUser, "same_user", false, "use same user [default false]")

}

func parseZoneString() ([]model.EventZone, error) {
	zones := make([]model.EventZone, 0)

	if len(zoneString) < 5 {
		// best validation BTW
		return zones, zoneStringInvalidErr
	}

	for _, zone := range strings.Split(zoneString, "|") {
		section := strings.Split(zone, ",")
		if len(section) != 3 {
			return zones, zoneStringInvalidErr
		}
		rows, err1 := strconv.ParseInt(section[1], 10, 32)
		seats, err2 := strconv.ParseInt(section[2], 10, 32)

		if err1 != nil || err2 != nil {
			return zones, zoneStringInvalidErr
		}

		zones = append(zones, model.EventZone{
			Zone:  section[0],
			Rows:  int32(rows),
			Seats: int32(seats),
			Price: 1000, // hard code
		})
	}

	return zones, nil
}

func createEvent() (*model.Event, error) {

	client := NewClient(serverIP, serverPort)
	err := client.Login(username, password)

	if err != nil {
		return nil, err
	}

	name := fmt.Sprintf("Event-%s", time.Now().String())

	zones, err := parseZoneString()
	if err != nil {
		return nil, err
	}
	event, err := client.Request.CreateEvent(name, zones)
	if err != nil {
		return nil, err
	}

	return event, err
}

func createUser() {

	client := NewClient(serverIP, serverPort)
	const totalUsers = 5000
	const userPerSecond = 100
	var wg sync.WaitGroup

	for i := 0; i < totalUsers; i += userPerSecond {
		start := time.Now()

		for j := 0; j < userPerSecond && i+j < totalUsers; j++ {
			wg.Add(1)
			go func(id int) {
				defer wg.Done()

				testuser := fmt.Sprintf("testuser%d", id)
				err := client.Register(testuser, password)

				if err != nil {
					fmt.Printf("Create User %d Failed with %v \n", id, err)
				}

			}(i + j)
		}

		wg.Wait()
		elapsed := time.Since(start)
		if elapsed < time.Second {
			time.Sleep(time.Second - elapsed)
		}
	}

}

func simulate() {

	client := NewClient(serverIP, serverPort)
	detail, err := client.Request.GetEventZoneByID(int64(eventID))
	if err != nil {
		fmt.Printf("Can't Find Event ID: %v \n", eventID)
		return
	}

	detailLength := len(detail)

	if detailLength == 0 {
		fmt.Println("No zone detect, skip simulation")
		return
	}

	if sameUser {
		err := client.Login(username, password)
		if err != nil {
			fmt.Printf("Auth Error: %v \n", err)
			return
		}
	}

	var wg sync.WaitGroup
	var success int32
	var failed int32
	var requestErr int32

	seed := time.Now().UnixNano()
	source := rand.NewSource(seed)
	r := rand.New(source)

	for i := 0; i < attempts; i++ {
		wg.Add(1)

		go func(dummy *Client, id int) {
			defer wg.Done()

			if !sameUser {
				dummy = NewClient(serverIP, serverPort)
				testuser := fmt.Sprintf("testuser%d", id)
				err := dummy.Login(testuser, password)

				if err != nil {
					atomic.AddInt32(&requestErr, 1)
					// fmt.Printf("Auth Error: %v \n", err)
					return
				} else {
					time.Sleep(1 * time.Millisecond)
				}
			}

			zone := detail[r.Intn(detailLength)]
			row := r.Int31n(zone.Rows) + 1
			quantity := r.Intn(4) + 1
			res, err := dummy.Request.ClaimTicket(zone.EventID, zone.ID, row, quantity, sameUser)

			if err != nil {
				atomic.AddInt32(&requestErr, 1)
				// fmt.Printf("Cliam Ticket Error: %v \n", err)
				return
			}

			noTicketsClaimed := len(res.ClaimedTickets) == 0
			if noTicketsClaimed {
				atomic.AddInt32(&failed, 1)
				// fmt.Printf("No Tickets Claimed \n")
				return
			}

			time.Sleep(1 * time.Millisecond)

			_, err = dummy.Request.CreateOrder(res.ClaimedTickets, int64(eventID), sameUser)
			if err != nil {
				atomic.AddInt32(&requestErr, 1)
				fmt.Printf("Create Order Error: %v \n", err)
				fmt.Printf("Claimed Tickets: %v \n", res.ClaimedTickets)
				return
			}

			atomic.AddInt32(&success, 1)
		}(client, i)

		time.Sleep(1 * time.Millisecond)
	}

	wg.Wait()

	fmt.Println()
	fmt.Println("----- [Simulation Result] -----")
	fmt.Printf("Get Ticket: %d \n", success)
	fmt.Printf("Didn't get Ticket: %d \n", failed)
	fmt.Printf("HTTP request error: %d \n", requestErr)
	fmt.Printf("Total request: %d \n", success+failed+requestErr)
}

func main() {
	flag.Parse()

	if action == "event" {
		fmt.Println("Start Event Creation")
		event, err := createEvent()
		if err != nil {
			fmt.Printf("Create Event Failed: %v \n", err)
			return
		}
		fmt.Printf("Event Create Success: %v \n", event)
		return
	}

	if action == "user" {
		fmt.Println("Start TestUser Creation")
		createUser()
		fmt.Println("TestUser Creation End")
		return
	}

	if action == "simulate" {
		// Simulate Buying Ticket
		fmt.Printf("Start Ticket Buying Simulation, Event: %d \n", eventID)
		fmt.Printf("Using Same User: %v \n", sameUser)
		simulate()
		fmt.Println("Simulation Complete")
	}

}
