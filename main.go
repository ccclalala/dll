package main

/*
extern int helloFromC();
*/
import "C"
import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"runtime"
	"strconv"

	"github.com/phayes/freeport"
)

var port int

//export HelloFromGo
func HelloFromGo() {
	fmt.Println(port)
	port, _ = freeport.GetFreePort()
	m := http.NewServeMux()
	s := http.Server{Addr: ":" + strconv.Itoa(port), Handler: m}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	m.HandleFunc("/shutdown", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
		// Cancel the context on request
		cancel()
	})
	go func() {
		if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()
	open("http://localhost:" + strconv.Itoa(port))
	select {
	case <-ctx.Done():
		// Shutdown the server when the context is canceled
		s.Shutdown(ctx)
		C.helloFromC()
	}
	log.Printf("Finished")
}

func main() {
}

func open(url string) error {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start"}
	case "darwin":
		cmd = "open"
	default: // "linux", "freebsd", "openbsd", "netbsd"
		cmd = "xdg-open"
	}
	args = append(args, url)
	return exec.Command(cmd, args...).Start()
}
