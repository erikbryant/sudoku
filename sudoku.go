package sudoku

import "fmt"

// Board implements a width x height grid of runes
type Board struct {
	name  string
	grid  [9][9]int
	ticks [9][9]map[int]bool
}

// init initializes a new board
func (b *Board) init() {
	for col := 0; col < 9; col++ {
		for row := 0; row < 9; row++ {
			b.ticks[col][row] = make(map[int]bool)
			for digit := 1; digit <= 9; digit++ {
				b.ticks[col][row][digit] = true
			}
		}
	}
}

// New returns a new sudoku board populated with the given starting lines
func New(name string, lines []string) Board {
	var b Board
	b.init()

	b.name = name

	for row := 0; row < 9; row++ {
		for col := 0; col < 9; col++ {
			digit := int(lines[row][col] - 48)
			if digit == 0 {
				continue
			}
			b.setDigit(col, row, digit)
		}
	}

	return b
}

// Name returns the name of the board
func (b *Board) Name() string {
	return b.name
}

// Print prints a board
func (b *Board) Print() {
	fmt.Println(b.Name())
	for row := 0; row < 9; row++ {
		for col := 0; col < 9; col++ {
			if b.grid[col][row] == 0 {
				fmt.Printf(" _ ")
			} else {
				fmt.Printf(" %d ", b.grid[col][row])
			}
		}
		fmt.Printf(" | ")
		for col := 0; col < 9; col++ {
			fmt.Printf(" %d ", len(b.ticks[col][row]))
		}
		fmt.Printf("\n")
	}

	if !b.valid() {
		fmt.Println("Board is not valid!!!")
	}
}

// equal tests whether two boards are equal
func (b *Board) equal(b2 Board) bool {
	for col := 0; col < 9; col++ {
		for row := 0; row < 9; row++ {
			if b.grid[col][row] != b2.grid[col][row] {
				return false
			}
			if !equal(b.ticks[col][row], b2.ticks[col][row]) {
				return false
			}
		}
	}

	return true
}

// copy copies one board into another
func (b *Board) copy() Board {
	var b2 Board

	for col := 0; col < 9; col++ {
		for row := 0; row < 9; row++ {
			b2.setDigit(col, row, b.grid[col][row])
		}
	}

	for col := 0; col < 9; col++ {
		for row := 0; row < 9; row++ {
			b2.ticks[col][row] = make(map[int]bool)
			for digit := range b.ticks[col][row] {
				b2.ticks[col][row][digit] = true
			}
		}
	}

	if !b.equal(b2) {
		fmt.Println("ERROR!!!!!!!! boards are not equal")
	}

	return b2
}

// square returns the bounding indices of the 3x3 square that i is in
func square(i int) (int, int) {
	min := 0
	max := 2
	if i >= 3 && i <= 5 {
		min = 3
		max = 5
	}
	if i >= 6 {
		min = 6
		max = 8
	}

	return min, max
}

// untick removes a tick
func (b *Board) untick(col, row, digit int) {
	delete(b.ticks[col][row], digit)
}

// untickCol removes a tick from an entire column
func (b *Board) untickCol(col, digit int) {
	for row := 0; row < 9; row++ {
		b.untick(col, row, digit)
	}
}

// untickRow removes a tick from an entire row
func (b *Board) untickRow(row, digit int) {
	for col := 0; col < 9; col++ {
		b.untick(col, row, digit)
	}
}

// untickSquare removes a tick from an entire 3x3 square
func (b *Board) untickSquare(col, row, digit int) {
	colMin, colMax := square(col)
	rowMin, rowMax := square(row)

	for c := colMin; c <= colMax; c++ {
		for r := rowMin; r <= rowMax; r++ {
			b.untick(c, r, digit)
		}
	}
}

// setDigit sets the given digit in the solution grid and unticks all relevant ticks
func (b *Board) setDigit(col, row, digit int) {
	if digit == 0 {
		return
	}

	b.grid[col][row] = digit

	// Empty the map now that it is solved
	b.ticks[col][row] = make(map[int]bool)

	b.untickRow(row, digit)
	b.untickCol(col, digit)
	b.untickSquare(col, row, digit)
}

// singleTick finds any cell that can have only one value
func (b *Board) singleTick() {
	for col := 0; col < 9; col++ {
		for row := 0; row < 9; row++ {
			if len(b.ticks[col][row]) == 1 {
				// There's got to be a simpler way to get the one key from the map...
				for digit := range b.ticks[col][row] {
					b.setDigit(col, row, digit)
				}
			}
		}
	}
}

// equal returns true if two maps are identical
func equal(a, b map[int]bool) bool {
	if len(a) != len(b) {
		return false
	}

	for val := range a {
		if !b[val] {
			return false
		}
	}

	return true
}

// doubleTickTwin finds the twin of a given doubletick
func (b *Board) doubleTickTwin(col, row int) {
	// Check the same row
	for c := 0; c < 9; c++ {
		if c == col {
			continue
		}
		if equal(b.ticks[col][row], b.ticks[c][row]) {
			// untick all others in this row
			for c2 := 0; c2 < 9; c2++ {
				if c2 == c || c2 == col {
					continue
				}
				for digit := range b.ticks[col][row] {
					b.untick(c2, row, digit)
				}
			}
			break
		}
	}

	// Check the same col
	for r := 0; r < 9; r++ {
		if r == row {
			continue
		}
		if equal(b.ticks[col][row], b.ticks[col][r]) {
			// untick all others in this col
			for r2 := 0; r2 < 9; r2++ {
				if r2 == r || r2 == row {
					continue
				}
				for digit := range b.ticks[col][row] {
					b.untick(col, r2, digit)
				}
			}
			break
		}
	}

	// Check the same square
	colMin, colMax := square(col)
	rowMin, rowMax := square(row)

	for c := colMin; c <= colMax; c++ {
		for r := rowMin; r <= rowMax; r++ {
			if c == col && r == row {
				continue
			}
			if equal(b.ticks[col][row], b.ticks[c][r]) {
				// untick all others in this square
				for c2 := colMin; c2 <= colMax; c2++ {
					for r2 := rowMin; r2 <= rowMax; r2++ {
						if c2 == col && r2 == row {
							continue
						}
						if c2 == c && r2 == r {
							continue
						}
						for digit := range b.ticks[col][row] {
							b.untick(c2, r2, digit)
						}
					}
				}
				break
			}
		}
	}
}

// doubleTick finds pairs of duoble ticks (like 2/3 and 3/2) in the same
// row, col, or square. If found, it then eliminates those two possibilities
// from every other cell in the row, col, and square.
func (b *Board) doubleTick() {
	for col := 0; col < 9; col++ {
		for row := 0; row < 9; row++ {
			if len(b.ticks[col][row]) == 2 {
				b.doubleTickTwin(col, row)
			}
		}
	}
}

// valid returns whether the given board is in a valid state
func (b *Board) valid() bool {
	valid := true

	for col := 0; col < 9; col++ {
		for row := 0; row < 9; row++ {
			if b.grid[col][row] > 9 || b.grid[col][row] < 0 {
				fmt.Printf("Invalid grid value %d at [%d][%d]\n", b.grid[col][row], col, row)
				valid = false
			}
			if len(b.ticks[col][row]) > 9 {
				fmt.Printf("Too many ticks %d at [%d][%d]\n", len(b.ticks[col][row]), col, row)
				valid = false
			}
			if b.grid[col][row] == 0 && len(b.ticks[col][row]) <= 0 {
				// fmt.Printf("Too few ticks for unsolved cell [%d][%d]\n", col, row)
				valid = false
			}
		}
	}

	for col := 0; col < 9; col++ {
		digits := make(map[int]bool)
		for row := 0; row < 9; row++ {
			d := b.grid[col][row]
			if d == 0 {
				continue
			}
			if digits[d] {
				fmt.Printf("Too many %d's in col %d\n", d, col)
				valid = false
			}
			digits[d] = true
		}
	}

	for row := 0; row < 9; row++ {
		digits := make(map[int]bool)
		for col := 0; col < 9; col++ {
			d := b.grid[col][row]
			if d == 0 {
				continue
			}
			if digits[d] {
				fmt.Printf("Too many %d's in row %d\n", d, row)
				valid = false
			}
			digits[d] = true
		}
	}

	return valid
}

// countSolved returns the number of grid positions that have been solved
func (b *Board) countSolved() int {
	count := 0

	for col := 0; col < 9; col++ {
		for row := 0; row < 9; row++ {
			if b.grid[col][row] != 0 {
				count++
			}
		}
	}

	return count
}

// Solved returns whether a given board is Solved
func (b *Board) Solved() bool {
	if !b.valid() {
		return false
	}

	for col := 0; col < 9; col++ {
		for row := 0; row < 9; row++ {
			if b.grid[col][row] < 1 || b.grid[col][row] > 9 {
				return false
			}
		}
	}

	return true
}

// shortestTicks returns the shortest sequence of ticks on the board
func (b *Board) shortestTicks() (int, int) {
	minLen := 10
	minCol := -1
	minRow := -1

	for col := 0; col < 9; col++ {
		for row := 0; row < 9; row++ {
			if b.grid[col][row] == 0 && len(b.ticks[col][row]) < minLen {
				minLen = len(b.ticks[col][row])
				minCol = col
				minRow = row
			}
		}
	}

	return minCol, minRow
}

// guess makes a guess and tries to solve that new board
func (b *Board) guess() Board {
	col, row := b.shortestTicks()

	for digit := range b.ticks[col][row] {
		b2 := b.copy()
		b2.setDigit(col, row, digit)
		b2.Solve()
		if b2.Solved() {
			return b2
		}
	}

	return *b
}

// Solve tries to solve a given board
func (b *Board) Solve() {
	for {
		before := b.countSolved()

		b.singleTick()
		b.doubleTick()

		if b.Solved() {
			return
		}

		if !b.valid() {
			return
		}

		after := b.countSolved()
		if before == after {
			break
		}
	}

	if !b.Solved() {
		*b = b.guess()
	}
}
