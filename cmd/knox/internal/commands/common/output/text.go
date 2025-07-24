package output

type Text string

func (s Text) String() string {
	return string(s)
}
