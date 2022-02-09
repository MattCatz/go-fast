package fast

import (
	"fmt"
	"testing"
	"time"
	"net"

	"github.com/ddo/go-spin"
)

var f *Fast
var testing_urls []string

func TestNew(t *testing.T) {
	var err error
	f, err = New(nil)

	if err != nil {
		t.Error(err)
		return
	}

	if f.client == nil {
		t.Error(err)
	}
}

func TestInit(t *testing.T) {
	err := f.Init()

	if err != nil {
		t.Error(err)
		return
	}

	if f.url == "" {
		t.Error(err)
	}

	if f.token == "" {
		t.Error(err)
	}

	if f.urlCount == 0 {
		t.Error(err)
	}
}

func TestGetUrls(t *testing.T) {
	urls, err := f.GetUrls()

	if err != nil {
		t.Error(err)
		return
	}

	if len(urls) != f.urlCount {
		t.Error(err)
	}

	testing_urls = urls
}

func TestDownload(t *testing.T) {
	byteLenChan := make(chan int64)
	done := make(chan struct{})

	spinner := spin.New("")

	go func() {
		for range byteLenChan {
			fmt.Printf(" \r %s", spinner.Spin())
		}
	}()

	err := f.download(testing_urls[0], byteLenChan, done)

	if err != nil {
		t.Error(err)
		return
	}

	fmt.Println("done")
}

// stop after 5s
func TestDownloadStop(t *testing.T) {
	byteLenChan := make(chan int64)
	done := make(chan struct{})

	spinner := spin.New("")

	go func() {
		for range byteLenChan {
			fmt.Printf(" \r %s", spinner.Spin())
		}
	}()

	go func() {
		<-time.After(2 * time.Second)
		close(done)
	}()

	err := f.download(testing_urls[0], byteLenChan, done)

	if err != nil {
		t.Error(err)
		return
	}

	fmt.Println("done")
}

func TestMeasure(t *testing.T) {
	KbpsChan := make(chan float64)

	spinner := spin.New("")

	go func() {
		for Kbps := range KbpsChan {
			fmt.Printf(" \r %s %.2f Kbps %.2f Mbps", spinner.Spin(), Kbps, Kbps/1000)
		}
	}()

	err := f.Measure(testing_urls, KbpsChan)

	if err != nil {
		t.Error(err)
		return
	}

	fmt.Println("done")
}

func TestBadBind(t *testing.T) {
	var err error

	// The Download should fail is the server is unreachable from the bound address
	f, err = New(&Option{
		bindAddress: "127.0.0.1",
	})

	err = f.Init()

	if err == nil {
		t.Fatal("Expecting Error")
		return
	}
}

func TestBind(t *testing.T) {
	// https://stackoverflow.com/questions/23558425/how-do-i-get-the-local-ip-address-in-go
	conn, err := net.Dial("udp", "8.8.8.8:80")
    if err != nil {
        t.Fatal("Failed to find output ip")
    }
    defer conn.Close()

    localAddr := conn.LocalAddr().(*net.UDPAddr)
	
	f, err = New(&Option{
		bindAddress: localAddr.IP.String(),
	})

	if err != nil {
		t.Error(err)
		return
	}

	TestInit(t)
	TestDownload(t)
	TestMeasure(t)
}
