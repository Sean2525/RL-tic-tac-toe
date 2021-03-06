package main

import (
	"fmt"
	"os"
	"time"

	agent "github.com/sean2525/RL-tic-tac-toe/agent"
	board "github.com/sean2525/RL-tic-tac-toe/board"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
)

func nbWins(wins []entry, c string) int {
	count := 0
	for _, e := range wins {
		if e.Value == c {
			count++
		}
	}
	return count
}

type entry struct {
	Timestamp int64  `json:"timestamp"`
	Value     string `json:"value"`
}

func (e entry) String() string {
	return fmt.Sprintf("%v, %v\n", e.Timestamp, e.Value)
}
func learnGame(a1, a2 *agent.Agent, numEpisodes int) {
	var b *board.Game
	for i := 0; i < numEpisodes; i++ {
		if i%2 == 0 {
			b = board.New("X")
		} else {
			b = board.New("O")
		}

		for b.PlayAble() {
			state := b.State
			if b.Player == a1.Sign {
				if err := a1.TrainPlay(b); err != nil {
					fmt.Println(err)
					break
				}
				a1.LearnFromMove(state, b)
				a2.LearnFromMove(state, b)
			} else if b.Player == a2.Sign {
				if err := a2.TrainPlay(b); err != nil {
					fmt.Println(err)
					break
				}
				a1.LearnFromMove(state, b)
				a2.LearnFromMove(state, b)
			} else {
				break
			}

		}
		a1.LearnFromMove(b.State, b)
		a2.LearnFromMove(b.State, b)
	}
}

func demoGameStats(a1, a2 *agent.Agent, numEpisodes int) {
	wins := make([]entry, 0, numEpisodes)
	draw := make([]entry, 0, numEpisodes)
	start := time.Now().UnixNano()
	var b *board.Game
	for i := 0; i < numEpisodes; i++ {
		if i%2 == 0 {
			b = board.New("X")
		} else {
			b = board.New("O")
		}

		for true {
			if b.PlayAble() {
				if b.Player == a1.Sign {
					a1.Play(b)
				} else if b.Player == a2.Sign {
					a2.Play(b)
				} else {
					break
				}
			} else {
				break
			}
		}

		if b.Winner == a1.Sign {
			wins = append(wins, entry{time.Now().UnixNano() - start, b.Winner})
		} else if b.Winner == a2.Sign {
			wins = append(wins, entry{time.Now().UnixNano() - start, b.Winner})
		} else {
			draw = append(draw, entry{time.Now().UnixNano() - start, "Draw"})
		}
	}
	// Display new stats
	fmt.Printf("%v wins %v%% times\n", a1.Sign, float64(nbWins(wins, a1.Sign))/float64(numEpisodes)*100)
	fmt.Printf("%v wins %v%% times\n", a2.Sign, float64(nbWins(wins, a2.Sign))/float64(numEpisodes)*100)
	fmt.Printf("Draws %v%% times\n", float64(len(draw))/float64(numEpisodes)*100)
	if len(os.Args) > 1 && os.Args[1] == "--plot" {
		generateFigure(wins, draw, numEpisodes, a1, a2)
	}

}

func main() {
	agent.Init()
	// Create the board and two agent
	a1 := agent.New(0.5, 1, 0.1, "X")
	a2 := agent.New(0.5, 1, 0.1, "O")

	fmt.Println("before leaning")
	demoGameStats(a1, a2, 3000)

	for i := 1; i <= 10; i++ {
		fmt.Printf("after learning %d times\n", i*3000)
		learnGame(a1, a2, 3000)
		demoGameStats(a1, a2, 3000)
	}

	a1.InteractiveGame()
}

// Generate figure from the wins array
func generateFigure(wins []entry, draw []entry, loopNb int, a1 *agent.Agent, a2 *agent.Agent) {
	// Create plot
	p, err := plot.New()
	if err != nil {
		panic(err)
	}
	// Set plot meta data
	p.Title.Text = "Both learning"
	p.X.Label.Text = "Time"
	p.Y.Label.Text = "Number of wins"
	// Build plot data
	ptsX := make(plotter.XYs, nbWins(wins, a1.Sign)+1)
	ptsO := make(plotter.XYs, nbWins(wins, a2.Sign)+1)
	ptsDraw := make(plotter.XYs, len(draw))
	countX := 0
	countO := 0
	for _, w := range wins[:loopNb] {
		if w.Value == "X" {
			countX++
			ptsX[countX].Y = float64(countX)
			ptsX[countX].X = float64(w.Timestamp)
		} else if w.Value == "O" {
			countO++
			ptsO[countO].Y = float64(countO)
			ptsO[countO].X = float64(w.Timestamp)
		}
	}
	for i, w := range draw {
		ptsDraw[i].Y = float64(i)
		ptsDraw[i].X = float64(w.Timestamp)
	}
	// Add data to plot
	err = plotutil.AddLines(p, "X", ptsX, "O", ptsO, "_", ptsDraw)
	if err != nil {
		panic(err)
	}
	// Save the plot to a PNG file.
	err = p.Save(4*vg.Inch, 4*vg.Inch, "points.png")
	if err != nil {
		panic(err)
	}
}
