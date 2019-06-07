package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"time"
)

const (
	diceValues    = 6
	winnerPoints  = 100
	looseValue    = 1
	maxIterations = 1 << 10
	gameCount     = 1 << 20
)

var (
	configFile = flag.String("config-file", "game.json", "path to config file")
)

// Player represents a player go figure
type Player struct {
	Name     string `json:"name"`
	Strategy string `json:"strategy"`
	Rolls    int    `json:"rolls"`
	Cap      int    `json:"cap"`
}

// Turn represents a Players turn
type Turn struct {
	p       *Player
	points  int
	winners int
	looses  int
}

// Game Represents a Game
type Game struct {
	turns      []*Turn
	scores     map[*Player]int
	iterations int
	winner     *Player
}

// Tournament is a collection of games
type Tournament struct {
	games  []*Game
	scores map[*Player]int
}

func init() {
	// initialize random generator
	rand.Seed(time.Now().UnixNano())
}

func main() {
	var v map[string]Player
	var players []*Player

	var t Tournament

	jsonData, err := ioutil.ReadFile(*configFile)

	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(jsonData, &v)

	if err != nil {
		log.Panicf("Could not get config data from %d: %d\n", configFile, err)
	}

	for player := range v {
		p := v[player]
		players = append(players, &p)
	}

	t.scores = make(map[*Player]int)

	for i := 0; i < gameCount; i++ {
		g, _ := runGame(players)
		t.scores[g.winner]++
	}

	for player := range t.scores {
		fmt.Printf("%v won %v times\n", player.Name, t.scores[player])
	}
}

func runGame(players []*Player) (Game, error) {

	var g Game

	g.scores = make(map[*Player]int)

	for i := 0; i < maxIterations; i++ {

		g.iterations++

		for _, player := range players {

			t, _ := playTurn(player)

			g.turns = append(g.turns, &t)

			if t.points > 0 {
				g.scores[player] += t.points
			}

			if g.scores[player] >= winnerPoints {
				g.winner = player
				return g, nil
			}

		}
	}

	return g, fmt.Errorf("%v iterations played with no winners", maxIterations)
}

func playTurn(p *Player) (Turn, error) {

	var t Turn
	t.p = p

	switch p.Strategy {

	case "fixedRolls":
		t, _ = playFixedRolls(p.Rolls)
	case "randomRolls":
		r := rand.Intn(p.Rolls) + 1
		t, _ = playFixedRolls(r)
	case "capAt":
		t, _ = playCapAt(p.Cap)
	}

	return t, nil
}

func playFixedRolls(n int) (Turn, error) {

	var rollValue int
	var t Turn

	for i := 0; i < n; i++ {

		rollValue = roll()

		if rollValue == looseValue {
			t.points = 0
			t.looses++
			break
		}

		t.points = t.points + rollValue
		t.winners++
	}

	return t, nil
}

func playCapAt(n int) (Turn, error) {

	var t Turn
	var rollValue int

	for {

		rollValue = roll()

		if rollValue == looseValue {
			t.points = 0
			t.looses++
			break
		}

		t.points = t.points + rollValue
		t.winners++

		if t.points >= n {
			break
		}
	}

	return t, nil
}

func roll() int {
	return rand.Intn(diceValues) + 1
}
