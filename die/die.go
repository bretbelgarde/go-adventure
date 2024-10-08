package die

import (
	"math/rand"
	"regexp"
	"strconv"
	"strings"
)

// Roll rolls a die in the format XdY and returns the total as an int.
// The format of the string must be a single digit followed by a 'd'
// followed by one or more digits. If the string does not match this
// format, an error is returned.
func Roll(s string) (int, error) {
	matched, err := regexp.MatchString("^[1-9]d[1-9]\\d*", s)
	if matched {
		return -1, err
	}
	split := strings.Split(s, "d")
	total := 0
	times, _ := strconv.Atoi(split[0])
	max, _ := strconv.Atoi(split[1])

	for i := 0; i < times; i++ {
		total += (rand.Intn(max) + 1)
	}

	return total, nil
}
