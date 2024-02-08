package ws

const boardSize = 15
const marksToWin = 5

func fitsVertically(i int) bool {
	return i+boardSize*(marksToWin-1) < boardSize*boardSize
}

func fitsHorizontally(i int) bool {
	return i%boardSize+marksToWin <= boardSize
}

func fitsBackHorizontally(i int) bool {
	return i%boardSize-(marksToWin-1) >= 0
}

func getHorizontalRow(i int, cell string, board []string) []int {
	if !fitsHorizontally(i) {
		return []int{}
	}

	currentWinnerRow := make([]int, 0)
	currentRowCount := 1

	currentWinnerRow = append(currentWinnerRow, i)

	for j := i + 1; currentRowCount < boardSize; j++ {
		if board[j] == cell {
			currentWinnerRow = append(currentWinnerRow, j)
			currentRowCount++
		} else {
			break
		}
	}

	if len(currentWinnerRow) >= marksToWin {
		return currentWinnerRow
	}

	return []int{}
}

func getVerticalRow(i int, cell string, board []string) []int {
	if !fitsVertically(i) {
		return []int{}
	}

	currentWinnerRow := make([]int, 0)
	currentRowCount := 1

	currentWinnerRow = append(currentWinnerRow, i)

	for j := i + boardSize; currentRowCount < boardSize; j += boardSize {
		if board[j] == cell {
			currentWinnerRow = append(currentWinnerRow, j)
			currentRowCount++
		} else {
			break
		}
	}

	if len(currentWinnerRow) >= marksToWin {
		return currentWinnerRow
	}

	return []int{}
}

func getDiagonalLTRRow(i int, cell string, board []string) []int {
	if !fitsHorizontally(i) || !fitsVertically(i) {
		return []int{}
	}

	currentWinnerRow := make([]int, 0)
	currentRowCount := 1

	currentWinnerRow = append(currentWinnerRow, i)

	for j := i + boardSize + 1; currentRowCount < boardSize; j += boardSize + 1 {
		if board[j] == cell {
			currentWinnerRow = append(currentWinnerRow, j)
			currentRowCount++
		} else {
			break
		}
	}

	if len(currentWinnerRow) >= marksToWin {
		return currentWinnerRow
	}

	return []int{}
}

func getDiagonalRTLRow(i int, cell string, board []string) []int {
	if !fitsHorizontally(i) || !fitsVertically(i) {
		return []int{}
	}

	currentWinnerRow := make([]int, 0)
	currentRowCount := 1

	currentWinnerRow = append(currentWinnerRow, i)

	for j := i + boardSize - 1; currentRowCount < boardSize; j += boardSize - 1 {
		if board[j] == cell {
			currentWinnerRow = append(currentWinnerRow, j)
			currentRowCount++
		} else {
			break
		}
	}

	if len(currentWinnerRow) >= marksToWin {
		return currentWinnerRow
	}

	return []int{}
}

func getWinnerRow(board []string) []int {
	emptyCellCount := len(board)

	for i := 0; i < len(board); i++ {
		cell := board[i]

		if cell != "" {
			emptyCellCount--

			horizontalRow := getHorizontalRow(i, cell, board)
			if len(horizontalRow) > 0 {
				return horizontalRow
			}

			verticalRow := getVerticalRow(i, cell, board)
			if len(verticalRow) > 0 {
				return verticalRow
			}

			diagonalLTRRow := getDiagonalLTRRow(i, cell, board)
			if len(diagonalLTRRow) > 0 {
				return diagonalLTRRow
			}

			diagonalRTLRow := getDiagonalRTLRow(i, cell, board)
			if len(diagonalRTLRow) > 0 {
				return diagonalRTLRow
			}
		}
	}

	return []int{}
}
