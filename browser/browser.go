package browser

// Family ...
type Family int

// Browser Families
const (
	Chrome Family = iota
	Firefox
	Edge
	Safari
	InternetExplorer
)

var names = []string{"chrome", "firefox", "Microsoft Edge", "Safari", "internet explorer"}

// Parse ...
func Parse(b string) Family {
	idx := indexOf(b, names)
	return Family(idx)
}

// String returns the name of the day
func (b Family) String() string {
	return names[b]
}

func indexOf(element string, data []string) int {
	for k, v := range data {
		if element == v {
			return k
		}
	}
	return -1 //not found.
}
