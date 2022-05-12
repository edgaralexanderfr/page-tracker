package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/edgaralexanderfr/page-tracker/pkg/color"
)

const (
	URL           string        = "https://sedeelectronica.antioquia.gov.co/pasaporte/user/pago/"
	CHECK_TIME    time.Duration = 10
	TIME_INTERVAL time.Duration = time.Minute
	ALARM_URL     string        = "https://youtu.be/1bpE1IpXcfs"
)

var (
	unexpectedTexts [3]string = [3]string{
		"Comunicado pago de pasaporte",
		"límite diario de transacciones",
		"Consultar el pago",
	}
)

func main() {
	_, hErr := getParam("-h")
	_, helpErr := getParam("--help")
	_, taErr := getParam("-ta")
	_, tryAlarmErr := getParam("--try-alarm")

	if hErr == nil || helpErr == nil {
		fmt.Println("Page Tracker 1.0.0")
		fmt.Println("")
		fmt.Println("Usage:")
		fmt.Println("  command [arguments]:")
		fmt.Println("")
		fmt.Println("Options:")
		fmt.Println("  -h, --help       Display this help message")
		fmt.Println("  -ta, --try-alarm Try the alarm")
		fmt.Println("")

		return
	}

	if taErr == nil || tryAlarmErr == nil {
		alarm()

		return
	}

	for true {
		check()

		time.Sleep(CHECK_TIME * TIME_INTERVAL)
	}
}

func check() {
	fmt.Print(color.Yellow, "[CHECKING] ")
	fmt.Print(color.White, "Verifying ")
	fmt.Print(color.Cyan, URL)
	fmt.Println(color.White, "...")
	fmt.Println()

	resp, err := http.Get(URL)

	if err != nil {
		log.Fatalln(err)
	}

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Fatalln(err)
	}

	responseText := string(body)
	valid := true

	for _, unexpectedText := range unexpectedTexts {
		if strings.Index(responseText, unexpectedText) != -1 {
			valid = false

			break
		}
	}

	if resp.StatusCode != 200 || !valid {
		fmt.Print(color.Red, "[ERROR] ")
		fmt.Print(color.White, "Page has returned an unexpected response...")
		fmt.Println()
	} else {
		fmt.Print(color.Green, "[SUCCESS] ")
		fmt.Print(color.White, "Page has returned a new response! Please visit ")
		fmt.Print(color.Cyan, URL)
		fmt.Print(color.White, " to check it out.")
		fmt.Println()

		alarm()
	}
}

func getParam(name string) (string, error) {
	args := os.Args[1:]
	last := len(args) - 1

	for i, param := range args {
		if param == name {
			if i == last {
				return "true", nil
			} else {
				return args[i+1], nil
			}
		}
	}

	return "", fmt.Errorf("Param not found")
}

func alarm() {
	var err error

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", ALARM_URL).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", ALARM_URL).Start()
	default:
		err = fmt.Errorf("Unable to run the alarm")
	}

	if err != nil {
		log.Fatal(err)
	}
}
