package main

import (
	"flag"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"sync"
	"ticket/client/model"
	"time"
)

var zoneStringInvalidErr = fmt.Errorf("provided zone string isn't valid")

var serverIP string
var serverPort int
var username string
var password string
var zoneString string // ex: A,5,5|B,5,5|C,5,5
var attemptUser int
var sameUser bool

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
	flag.IntVar(&attemptUser, "attempt", 100, "attempt [default 100 times]")

	flag.BoolVar(&sameUser, "same_user", false, "use same user [default false]")
}

func authUser(client *Client, user, pwd string) error {
	err := client.Login(user, pwd)

	if err != nil {
		// User Probably not created
		err := client.Register(user, pwd)
		if err != nil {
			return err
		}

		// Retry
		err = client.Login(user, pwd)
		if err != nil {
			return err
		}
	}

	// fmt.Println("User Successfully Logged In")
	return nil
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

func create() (*model.Event, *[]model.EventZoneDetail, error) {
	client := NewClient(serverIP, serverPort)

	err := authUser(client, username, password)
	if err != nil {
		return nil, nil, err
	}

	name := fmt.Sprintf("Event-%s", time.Now().String())

	zones, err := parseZoneString()
	if err != nil {
		return nil, nil, err
	}
	event, err := client.Request.CreateEvent(name, zones)
	if err != nil {
		return nil, nil, err
	}

	detail, err := client.Request.GetEventZoneByID(event.ID)

	return event, detail, err
}

func simulate(detail []model.EventZoneDetail) {

	if len(detail) == 0 {
		fmt.Println("No zone detect, skip simulation")
		return
	}

	var wg sync.WaitGroup
	success, failed := 0, 0

	seed := time.Now().UnixNano()
	source := rand.NewSource(seed)
	r := rand.New(source)

	for i := 0; i < attemptUser; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			client := NewClient(serverIP, serverPort)

			simulateUserID := i
			if sameUser {
				simulateUserID = 0
			}

			testuser := fmt.Sprintf("testuser%d", simulateUserID)
			err := authUser(client, testuser, password)
			if err != nil {
				failed++
				fmt.Printf("Auth Error: %v \n", err)
				return
			}
			zone := detail[r.Intn(len(detail))]
			row := r.Int31n(zone.Rows) + 1
			quantity := r.Intn(4) + 1
			res, err := client.Request.ClaimTicket(zone.EventID, zone.ID, row, quantity)

			if err != nil {
				failed++
				fmt.Printf("Cliam Ticket Error: %v \n", err)
				return
			}

			noTicketsClaimed := len(res.ClaimedTickets) == 0
			if noTicketsClaimed {
				failed++
				// fmt.Printf("No Tickets Claimed \n")
				return
			}

			_, err = client.Request.CreateOrder(res.ClaimedTickets)
			if err != nil {
				failed++
				fmt.Printf("Create Order Error: %v \n", err)
				fmt.Printf("Input: %v %v %v %v \n", zone.EventID, zone.ID, row, quantity)
				return
			}

			success++
		}()
	}

	wg.Wait()
	fmt.Printf("[Result] Success: %d, Failed: %d \n", success, failed)
}

func main() {
	flag.Parse()

	// Create Event
	fmt.Println("Start Event Creation")
	event, detail, err := create()
	if err != nil {
		fmt.Printf("Create Event Failed: %v \n", err)
		return
	}
	fmt.Printf("Event Create Success: %v \n", event)

	// Simulate Buying Ticket
	fmt.Println("Start Ticket Buying Simulation")
	simulate(*detail)
	fmt.Println("Simulation Complete")
}
