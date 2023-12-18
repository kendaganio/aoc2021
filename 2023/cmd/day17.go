package cmd

import (
	"fmt"
	"math"
	"slices"
	"strings"
	"time"

	"github.com/kendaganio/aoc/2023/magic"
	"github.com/spf13/cobra"
	"golang.org/x/exp/maps"
)

type Move struct {
	p Point
	d Direction
	s int
	h int
}

func (m Move) HashKey() string {
	return fmt.Sprintf("%d,%d,%s,%d", m.p.X, m.p.Y, m.d, m.s)
}

func (m Move) next(heatLoss int, minStraight int, maxStraight int) []Move {
	neighbors := []Move{}

	switch m.d {
	case Up:
		if m.s >= minStraight {
			neighbors = append(neighbors,
				Move{Point{m.p.X - 1, m.p.Y}, Left, 1, m.h + heatLoss},
				Move{Point{m.p.X + 1, m.p.Y}, Right, 1, m.h + heatLoss},
			)
		}
		if m.s < maxStraight {
			neighbors = append(neighbors,
				Move{Point{m.p.X, m.p.Y - 1}, Up, m.s + 1, m.h + heatLoss},
			)
		}
	case Down:
		if m.s >= minStraight {
			neighbors = append(neighbors,
				Move{Point{m.p.X - 1, m.p.Y}, Left, 1, m.h + heatLoss},
				Move{Point{m.p.X + 1, m.p.Y}, Right, 1, m.h + heatLoss},
			)
		}
		if m.s < maxStraight {
			neighbors = append(neighbors,
				Move{Point{m.p.X, m.p.Y + 1}, Down, m.s + 1, m.h + heatLoss},
			)
		}
	case Left:
		if m.s >= minStraight {
			neighbors = append(neighbors,
				Move{Point{m.p.X, m.p.Y - 1}, Up, 1, m.h + heatLoss},
				Move{Point{m.p.X, m.p.Y + 1}, Down, 1, m.h + heatLoss},
			)
		}
		if m.s < maxStraight {
			neighbors = append(neighbors,
				Move{Point{m.p.X - 1, m.p.Y}, Left, m.s + 1, m.h + heatLoss},
			)
		}
	case Right:
		if m.s >= minStraight {
			neighbors = append(neighbors,
				Move{Point{m.p.X, m.p.Y - 1}, Up, 1, m.h + heatLoss},
				Move{Point{m.p.X, m.p.Y + 1}, Down, 1, m.h + heatLoss},
			)
		}
		if m.s < maxStraight {
			neighbors = append(neighbors,
				Move{Point{m.p.X + 1, m.p.Y}, Right, m.s + 1, m.h + heatLoss},
			)
		}
	}

	return neighbors
}

func djikstra2(weights map[Point]int, src Point, minStraight int, maxStraight int) (map[string]int, map[string]Move) {
	prev := make(map[string]Move)
	visited := make(map[string]bool)
	dist := make(map[string]int)

	q := []Move{
		{Point{src.X + 1, src.Y}, Right, 2, 0},
		{Point{src.X, src.Y + 1}, Down, 2, 0},
	}

	for len(q) != 0 {
		move := q[0]
		q = q[1:]
		slices.SortFunc(q, func(a Move, b Move) int {
			return a.h - b.h
		})

		if _, ok := weights[move.p]; !ok {
			continue
		}

		if visited[move.HashKey()] {
			continue
		}

		neighbors := move.next(weights[move.p], minStraight, maxStraight)
		q = append(q, neighbors...)

		for _, next := range neighbors {

			if _, ok := weights[next.p]; !ok {
				continue
			}

			if _, ok := dist[next.HashKey()]; !ok {
				dist[next.HashKey()] = math.MaxInt
			}

			alt := dist[move.HashKey()] + weights[move.p]
			if dist[next.HashKey()] > alt {
				dist[next.HashKey()] = alt
				prev[next.HashKey()] = move
			}
		}

		visited[move.HashKey()] = true
	}

	return dist, prev
}

func SolveD17P1(grid map[Point]int, m int, n int, minStraight int, maxStraight int) (total int) {
	total = math.MaxInt
	start := Point{0, 0}

	dist, _ := djikstra2(grid, start, minStraight, maxStraight)

	endKeys := []string{}
	for _, move := range maps.Keys(dist) {
		if strings.Contains(move, fmt.Sprintf("%d,%d,", m-1, n-1)) {
			fmt.Println("DIST", move, dist[move])
			total = min(total, dist[move])
			endKeys = append(endKeys, move)
		}
	}

	return total + grid[Point{n - 1, m - 1}]
}

func SolveD17P2(grid map[Point]int) (total int) {
	return total
}

var day17Cmd = &cobra.Command{
	Use:   "day17 [path/to/input]",
	Short: "Solver for day 17",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		lines := magic.Lines(args[0])
		m, n := len(lines), len(lines[0])

		grid := map[Point]int{}
		for y := range lines {
			for x := range lines[0] {
				grid[Point{x, y}] = magic.ParseInt(string(lines[y][x]))
			}
		}

		start := time.Now()
		a := SolveD17P1(grid, m, n, 0, 3)
		fmt.Println("Part 1:", a, time.Since(start))

		start = time.Now()
		b := SolveD17P1(grid, m, n, 4, 10)
		fmt.Println("Part 2:", b, time.Since(start))
	},
}

func init() {
	rootCmd.AddCommand(day17Cmd)
}