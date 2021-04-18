package invite

import (
	"fmt"
	"math"
	"math/rand"
	"strconv"
	"time"
)

func GenerageCode(value string) string {

	chars := []string{
		"0", "1", "2", "3", "4", "5", "6", "7", "8", "9",
		"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z",
		"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z",
	}
	cl := float64(62)
	b := 5

	codeLen := 8

	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(chars), func(i, j int) { chars[i], chars[j] = chars[j], chars[i] })

	var code string

	for i := 0; i < codeLen; i++ {
		start := 2 + (i * b)
		v := value[start:(b + start)]
		n, err := strconv.ParseUint(v, 16, 32)
		if err != nil {
			fmt.Println(err)
		}
		m := int(math.Mod(float64(n), cl))
		code += chars[m]
	}

	return code
}
