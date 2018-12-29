// cursor movement in byte slice
package cursor

// Index returns the index within the buffer for a given column and line
// If column overflows in either positive or negative, buffer length determins endpoints.
// line starts with 1
// column starts with 1
// returned index starts with 0
func Index(buf []rune, sep rune, line, column int) (i int) {
	if line < 1 { // don't handle negative lines
		return 0
	}
	// Move index forward for each line
FINDLINE:
	for i = 0; i < len(buf); i++ {
		if line == 1 {
			break FINDLINE
		}
		if buf[i] == sep { // line ending found
			line--
		}
	}
	if i == len(buf) { // We reached the end
		i--
	}
	i = i + column
	if column > 0 {
		i--
	}
	if i < 0 {
		i = 0
	}
	return
}

// IndexLeft returns the right most index of the sep character
// to the left of index
func IndexLeft(buf []rune, sep rune, index int) (i int) {
	if index >= len(buf) {
		index = len(buf) - 1
	}
	if index > 0 && index < len(buf) {
		return IndexLast(buf[:index], sep)
	}
	return 0
}

// IndexLast returns the last index of sep inside buf or 0 if none is found
// Similar to IndexLeft but works on buf alone.
func IndexLast(buf []rune, sep rune) (i int) {
	for i := len(buf) - 1; i > 0; i-- {
		if sep == buf[i] {
			return i
		}
	}
	return 0
}

// IndexRune returns the index of the first position of sep in buf
func IndexRune(buf []rune, sep rune) (i int) {
	for i := 0; i < len(buf); i++ {
		if sep == buf[i] {
			return i
		}
	}
	return len(buf) - 1
}

// IndexUp moves the cursor one line up from the given index.
// It tries to keep the column position if possible.
func IndexUp(buf []rune, sep rune, index int) (i int) {
	// make sure we're within the buffer limits
	if index >= len(buf) {
		index = len(buf) - 1
	}
	if index < 0 {
		return 0
	}
	// Move to end of previous line
	end := IndexLeft(buf, sep, index)
	// Which column are we at now?
	currentcol := index - end - 1
	// We are at the first character
	if end == 0 {
		return 0
	}
	// Find beginning of the previous line
	begin := IndexLeft(buf, sep, end)
	// move to the right of the sep so we end up on column 0
	if begin > 0 {
		begin = begin + 1
	}
	endcol := end - begin
	// the above line is empty
	if endcol == 0 {
		return end
	}
	// if the line above is shorter than the current column position
	if endcol < currentcol {
		return begin + endcol
	}
	i = begin + currentcol
	return
}

// IndexDown moves the cursor one line down from the given index.
// It tries to keep the column position if possible.
func IndexDown(buf []rune, sep rune, index int) (i int) {
	// make sure we're within the buffer limits
	if index >= len(buf) {
		return len(buf) - 1
	}
	if index < 0 {
		index = 0
	}
	currentcol := index // first line
	// calculate current column when index is on another line
	newlineleft := IndexLast(buf[:index], sep)
	if newlineleft != 0 {
		currentcol = index - newlineleft - 1
	}
	// index of the next new line
	begin := IndexRune(buf[index:], sep) + index
	if begin >= 0 {
		begin = begin + 1
	}
	endcol := IndexRune(buf[begin:], sep)
	if endcol == -1 { // there is no ending, eg. no newline at EOF
		endcol = len(buf[begin:]) - 1
	}
	if endcol == 0 { // means we're on empty line
		return begin
	}
	end := begin + endcol
	// line below is shorter than current column position
	if currentcol > endcol {
		return end
	}
	if begin == len(buf)-1 {
		return begin
	}
	i = begin + currentcol
	return
}

// Position returns the line and column the given index(i) is on using sep as newline separator
// Lines and columns start with 1
func Position(buf []rune, sep rune, i int) (line, column, index int) {
	line = 1
	column = 1
	if i >= len(buf) {
		i = len(buf) - 1
	}
	if i < 0 {
		return line, column, 0
	}
	index = i
	// On the first line
	line = Count(buf[:i], sep) + 1
	if line == 1 {
		column = i + 1
		return
	}
	column = i - IndexLast(buf[:i], sep)
	return
}

// Count returns the number of occurences of sep in buf
func Count(buf []rune, sep rune) (c int) {
	for i := 0; i < len(buf); i++ {
		if sep == buf[i] {
			c++
		}
	}
	return
}

// LineBefore returns slice from the given index up to the first separator
func LineBefore(buf []rune, sep rune, i int) (r []rune) {
	f := IndexLeft(buf, sep, i)
	if buf[f] == sep {
		k := f + 1
		if k == i {
			return []rune("")
		}
		return buf[k:i]
	}
	r = buf[f:i]

	return
}
