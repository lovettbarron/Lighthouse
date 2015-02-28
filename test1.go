package main

import (
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"log"
	rtl "github.com/jpoirier/gortlsdr"
	// plot "code.google.com/p/plotinum/plot"
	// util "code.google.com/p/plotinum/plotutil"

)

const (
	shipFreq float32 = 	88500
	shipChannel int = 2
)

var send_ping bool = false

func main() {
	runtime.GOMAXPROCS(3) // Limits number of threads
	var err error // Declares error
	_ = err // Because
	var dev *rtl.Context // Declares context
	_ = dev // Yeah

	//---------- Open Device ----------
	if dev, err = rtl.Open(0); err != nil {
		log.Fatal("\tOpen Failed, exiting\n")
	}
	defer dev.Close()
	go sig_abort(dev)

	dev.ResetBuffer(); 
	dev.SetCenterFreq(int(shipFreq))
	
	IQch := make(chan bool)
	var userctx rtl.UserCtx = IQch
	go async_stop(dev, IQch)

	fmt.Println(dev.ReadAsync(rtlsdr_cb, &userctx, rtl.DefaultAsyncBufNumber, rtl.DefaultBufLength))


}


// from https://github.com/jpoirier/gortlsdr/blob/master/rtlsdr_example.go
func rtlsdr_cb(buf []byte, userctx *rtl.UserCtx) {
	if send_ping {
		send_ping = false
		// send a ping to async_stop
		if c, ok := (*userctx).(chan bool); ok {
			c <- true // async-read done signal
		}
	}

	var s uint32
	var i, rate int
	i, rate = 0, 5
	
	for sample := range buf {
		s = s << uint32(sample)
		i++
		if i>rate {
			i=0
			s /= uint32(rate)
			fmt.Printf("%x",s)
		}
	}



	fmt.Printf("%x",buf)
}

func async_stop(dev *rtl.Context, c chan bool) {
	log.Println("async_stop running...")
	<-c
	log.Println("Received ping from rtlsdr_cb, calling CancelAsync")
	if err := dev.CancelAsync(); err != nil {
		log.Printf("CancelAsync failed - %s\n", err)
	} else {
		log.Printf("CancelAsync successful\n")
	}

	os.Exit(0)
}

func sig_abort(dev *rtl.Context) {
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT)
	<-ch
	_ = dev.CancelAsync()
	dev.Close()
	os.Exit(0)
}