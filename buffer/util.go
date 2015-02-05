package buffer
import(
	"strconv"
)

type String string

func (this String) String() string {
	return string(this)
}

type Int int

func (this Int) String() string {
	return strconv.Itoa(int(this))
}
