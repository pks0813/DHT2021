package main

import (
	"flag"
	"fmt"
	"os"
	"time"
)

var (
	help     bool=true
	testName string="basic"
)

func init() {
	//flag.BoolVar(&help, "help", false, "help")
	//flag.StringVar(&testName, "test", "", "which test(s) do you want to run: basic/advance/all")
	//
	//flag.Usage = usage
	//flag.Parse()
	//
	//if help || (testName != "basic" && testName != "advance" && testName != "all") {
	//	flag.Usage()
	//	os.Exit(0)
	//}
	//
	//rand.Seed(time.Now().UnixNano())
}

func main() {
	_, _ = fmt.Println("Welcome to DHT-2020 Test Program!\n")

	var basicFailRate float64
	var forceQuitFailRate float64
	var QASFailRate float64

	switch testName {
	case "all":
		fallthrough
	case "basic":
		_, _ = fmt.Println("Basic Test Begins:")
		basicPanicked, basicFailedCnt, basicTotalCnt := basicTest()
		if basicPanicked {
			_, _ = fmt.Printf("Basic Test Panicked.")
			os.Exit(0)
		}

		basicFailRate = float64(basicFailedCnt) / float64(basicTotalCnt)
		if basicFailRate > basicTestMaxFailRate {
			_, _ = fmt.Printf("Basic test failed with fail rate %.4f\n", basicFailRate)
		} else {
			_, _ = fmt.Printf("Basic test passed with fail rate %.4f\n", basicFailRate)
		}

		if testName == "basic" {
			break
		}
		time.Sleep(afterTestSleepTime)
		fallthrough
	case "advance":
		_, _ = fmt.Println("Advance Test Begins:")

		/* ------ Force Quit Test Begins ------ */
		forceQuitPanicked, forceQuitFailedCnt, forceQuitTotalCnt := forceQuitTest()
		if forceQuitPanicked {
			_, _ = fmt.Printf("Force Quit Test Panicked.")
			os.Exit(0)
		}

		forceQuitFailRate = float64(forceQuitFailedCnt) / float64(forceQuitTotalCnt)
		if forceQuitFailRate > forceQuitMaxFailRate {
			_, _ = fmt.Printf("Force quit test failed with fail rate %.4f\n", forceQuitFailRate)
		} else {
			_, _ = fmt.Printf("Force quit test passed with fail rate %.4f\n", forceQuitFailRate)
		}
		time.Sleep(afterTestSleepTime)
		/* ------ Force Quit Test Ends ------ */

		/* ------ Quit & Stabilize Test Begins ------ */
		QASPanicked, QASFailedCnt, QASTotalCnt := quitAndStabilizeTest()
		if QASPanicked {
			_, _ = fmt.Printf("Quit & Stabilize Test Panicked.")
			os.Exit(0)
		}

		QASFailRate = float64(QASFailedCnt) / float64(QASTotalCnt)
		if QASFailRate > QASMaxFailRate {
			_, _ = fmt.Printf("Quit & Stabilize test failed with fail rate %.4f\n", QASFailRate)
		} else {
			_, _ = fmt.Printf("Quit & Stabilize test passed with fail rate %.4f\n", QASFailRate)
		}
		/* ------ Quit & Stabilize Test Ends ------ */
	}

	_, _ = fmt.Println("\nFinal print:")
	if basicFailRate > basicTestMaxFailRate {
		_, _ = fmt.Printf("Basic test failed with fail rate %.4f\n", basicFailRate)
	} else {
		_, _ = fmt.Printf("Basic test passed with fail rate %.4f\n", basicFailRate)
	}
	if forceQuitFailRate > forceQuitMaxFailRate {
		_, _ = fmt.Printf("Force quit test failed with fail rate %.4f\n", forceQuitFailRate)
	} else {
		_, _ = fmt.Printf("Force quit test passed with fail rate %.4f\n", forceQuitFailRate)
	}
	if QASFailRate > QASMaxFailRate {
		_, _ = fmt.Printf("Quit & Stabilize test failed with fail rate %.4f\n", QASFailRate)
	} else {
		_, _ = fmt.Printf("Quit & Stabilize test passed with fail rate %.4f\n", QASFailRate)
	}
}

func usage() {
	flag.PrintDefaults()
}
