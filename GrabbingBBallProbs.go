// GrabbingBBallProbs
package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

var wg sync.WaitGroup

func getPercentage(quarter string, score int, mintm int, sectm int, theOutput *os.File) {
	defer wg.Done()
	mintmString := strconv.Itoa(mintm)
	sectmString := strconv.Itoa(sectm)
	scoreString := strconv.Itoa(score)
	resp, err := http.PostForm("http://stats.inpredictable.com/nba/wpCalc.php", url.Values{"qtr": {quarter}, "mintm": {mintmString}, "sectm": {sectmString}, "scr": {scoreString}, "poss": {"Y"}})
	if err != nil {
		theOutput.WriteString(quarter + "," + scoreString + "," + mintmString + "," + sectmString + "," + "ERROR,ERROR" + "\r\n")
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		theOutput.WriteString(quarter + "," + scoreString + "," + mintmString + "," + sectmString + "," + "ERROR,ERROR" + "\r\n")
	}
	resp.Body.Close()
	percentagePlus := strings.Split(string(body), "Probability: ")
	percentage := strings.Split(percentagePlus[1], "%")
	fmt.Println(scoreString + " " + mintmString + " " + sectmString + " " + percentage[0])
	theOutput.WriteString(quarter + "," + scoreString + "," + mintmString + "," + sectmString + "," + percentage[0] + "\r\n")
}

func main() {
	fmt.Println("Trying to POST")
	quarter := "Q3"
	theOutput, err := os.Create("output.txt")
	if err != nil {
		panic(err)
	}
	defer theOutput.Close()
	//theOutput.WriteString("Quarter,Score,Minutes,Seconds,WinPercentage\r\n")
	for score := 10; score < 41; score++ {
		fmt.Println("Score: ", score)
		for mintm := 0; mintm < 12; mintm++ {
			for sectm := 0; sectm < 60; sectm++ {
				wg.Add(1)
				go getPercentage(quarter, score, mintm, sectm, theOutput)
				time.Sleep(10000 * time.Millisecond)
			}
		}
	}
	wg.Wait()
}
