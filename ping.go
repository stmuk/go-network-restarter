package main

import (
	"fmt"
	"net"
	"os"
	"time"

	"github.com/tatsushid/go-fastping"
)

func main() {

	const pings = 10

	testHost := "8.8.8.8"

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

	if avg > 500 {
		fmt.Printf("network broked")
	}

}
