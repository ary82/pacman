package utils

func CalculatePosibbleNextTile(
	board [31][28]int,
	xPosition int,
	yPosition int,
) [][]int {
	possibleMoves := [][]int{}

	if xPosition < 30 {
		if newPos := board[xPosition+1][yPosition]; newPos == 1 || newPos == 3 || newPos == 2 {
			possibleMoves = append(possibleMoves, []int{xPosition + 1, yPosition})
		}
	}
	if xPosition > 0 {
		if newPos := board[xPosition-1][yPosition]; newPos == 1 || newPos == 3 || newPos == 2 {
			possibleMoves = append(possibleMoves, []int{xPosition - 1, yPosition})
		}
	}

	if yPosition < 27 {
		if newPos := board[xPosition][yPosition+1]; newPos == 1 || newPos == 3 || newPos == 2 {
			possibleMoves = append(possibleMoves, []int{xPosition, yPosition + 1})
		}
	}

	if yPosition > 0 {
		if newPos := board[xPosition][yPosition-1]; newPos == 1 || newPos == 3 || newPos == 2 {
			possibleMoves = append(possibleMoves, []int{xPosition, yPosition - 1})
		}
	}

	if len(possibleMoves) == 0 {
		possibleMoves = append(possibleMoves, []int{xPosition, yPosition})
	}

	return possibleMoves
}
