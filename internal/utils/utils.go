package utils

func CalculatePosibbleNextTile(
	board [31][28]int,
	xPosition int,
	yPosition int,
	direction int,
) [][3]int {
	possibleMoves := [][3]int{}

	if xPosition < 30 {
		if newPos := board[xPosition+1][yPosition]; newPos == 1 || newPos == 3 || newPos == 2 {
			possibleMoves = append(possibleMoves, [3]int{xPosition + 1, yPosition, 1})
		}
	}
	if xPosition > 0 {
		if newPos := board[xPosition-1][yPosition]; newPos == 1 || newPos == 3 || newPos == 2 {
			possibleMoves = append(possibleMoves, [3]int{xPosition - 1, yPosition, 2})
		}
	}

	if yPosition < 27 {
		if newPos := board[xPosition][yPosition+1]; newPos == 1 || newPos == 3 || newPos == 2 {
			possibleMoves = append(possibleMoves, [3]int{xPosition, yPosition + 1, 4})
		}
	}

	if yPosition > 0 {
		if newPos := board[xPosition][yPosition-1]; newPos == 1 || newPos == 3 || newPos == 2 {
			possibleMoves = append(possibleMoves, [3]int{xPosition, yPosition - 1, 3})
		}
	}

	if len(possibleMoves) == 0 {
		possibleMoves = append(possibleMoves, [3]int{xPosition, yPosition})
	}

	preferredMoves := [][3]int{}
	for _, v := range possibleMoves {
		if direction == 1 && v[2] == 2 ||
			direction == 2 && v[2] == 1 ||
			direction == 3 && v[2] == 4 ||
			direction == 4 && v[2] == 3 {
		} else {
			preferredMoves = append(preferredMoves, v)
		}
	}

	if len(preferredMoves) == 0 {
		return possibleMoves
	}

	return preferredMoves
}
