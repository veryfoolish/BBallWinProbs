// GrabbingBBallProbs
package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

type Quarter int

const (
	Q1 Quarter = iota
	Q2
	Q3
	Q4
	OT // overtime
)

func (q Quarter) String() string {
	switch q {
	case Q1, Q2, Q3, Q4:
		return fmt.Sprintf("Q%d", q+1)
	default:
		return "OT"
	}
}

type GameState struct {
	Quarter
	TimeRemaining time.Duration
	ScoreDiff     uint
	Poss          bool
}

func (g GameState) MinutesRemaining() int { return int(g.TimeRemaining.Minutes()) }
func (g GameState) SecondsRemaining() int { return int(g.TimeRemaining.Seconds()) }

func NewGameState(q Quarter, diff uint, tRemain time.Duration, pos bool) GameState {
	return GameState{
		Quarter:       q,
		TimeRemaining: tRemain,
		ScoreDiff:     diff,
		Poss:          pos,
	}
}

func GameStateToUrlVals(g GameState) url.Values {
	v := url.Values{}
	v.Set("qtr", g.Quarter.String())
	v.Set("mintm", fmt.Sprint(g.MinutesRemaining()))
	v.Set("sectm", fmt.Sprint(g.SecondsRemaining()))
	v.Set("scr", fmt.Sprint(g.ScoreDiff))
	posstr := "N"
	if g.Poss {
		posstr = "Y"
	}
	v.Set("poss", posstr)
	return v
}

const postURL = "http://stats.inpredictable.com/nba/wpCalc.php"

func dieError(e error) {
	if e != nil {
		panic(e)
	}
}

func getPercentage(g GameState) (percent float64) {
	resp, err := http.PostForm(postURL, GameStateToUrlVals(g))
	dieError(err)

	body, err := ioutil.ReadAll(resp.Body)
	dieError(err)

	resp.Body.Close()
	percentagePlus := strings.Split(string(body), "Probability: ")
	percentage := strings.Split(percentagePlus[1], "%")
	fmt.Sscanf(percentage[0], "%f", &percent)
	return percent
}

func minSecToDuration(min int, sec int) time.Duration {
	t, err := time.ParseDuration(fmt.Sprintf("%dm%ds", min, sec))
	dieError(err)
	return t
}

const outFile = "output.txt"

func DefaultGameState(q Quarter, min int, sec int, scorediff int) GameState {
	return GameState{
		Quarter:       q,
		TimeRemaining: minSecToDuration(min, sec),
		ScoreDiff:     uint(scorediff),
		Poss:          true,
	}
}

func main() {
	fmt.Println("Attempting POST")
	quarter := Q3

	theOutput, err := os.Create(outFile)
	dieError(err)
	defer theOutput.Close()

	for score := 10; score < 41; score++ {
		fmt.Println("Score: ", score)
		for mintm := 0; mintm < 12; mintm++ {
			for sectm := 0; sectm < 60; sectm++ {
				percent := getPercentage(DefaultGameState(quarter, mintm, sectm, score))
				fmt.Printf("%2.2f\n", percent)
				time.Sleep(10 * time.Second)
			}
		}
	}
}
