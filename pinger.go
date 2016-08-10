package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/ThomasRooney/gexpect"
	"github.com/tatsushid/go-fastping"
)

const pings = 1
const routerUrl = "http://192.168.0.1/setup.cgi?todo=debug"

const testHost = "8.8.8.8" // Google DNS as ping test

func main() {

	user := os.Getenv("netgear_user")
	pw := os.Getenv("netgear_pw")

	p := fastping.NewPinger()

	var pingTot float32

	for i := 0; i < pings; i++ {
		ra, err := net.ResolveIPAddr("ip4:icmp", testHost)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		p.AddIPAddr(ra)
		p.OnRecv = func(addr *net.IPAddr, rtt time.Duration) {
			fmt.Printf("IP Addr: %s receive, RTT: %v\n", addr.String(), rtt)
			pingTot += float32(rtt)
		}

		err = p.Run()
		if err != nil {
			fmt.Println(err)
		}
	}

	avg := (pingTot / (pings * 1000000))

	fmt.Printf("avg: %f", avg)

	if avg > 1 { // XXX
		fmt.Printf(httpReq(user, pw, routerUrl))
		sendCmd()
	}

}

func httpReq(user string, pw string, routerUrl string) string {

	client := &http.Client{}

	req, err := http.NewRequest("GET", routerUrl, nil)
	req.SetBasicAuth(user, pw)
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	bodyText, err := ioutil.ReadAll(resp.Body)
	s := string(bodyText)
	return s
}

func sendCmd() {
	child, err := gexpect.Spawn("telnet 192.168.0.1")
	if err != nil {
		panic(err)
	}
	child.Expect("D7000 login:")
	child.SendLine("root")
	child.Expect("#")
	child.SendLine("echo $USER")
	child.Expect("root")
	child.SendLine("")
	child.Interact()
	child.Close()
}
