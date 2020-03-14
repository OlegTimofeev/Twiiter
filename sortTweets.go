package main

type byTime []Tweet
type byId []Tweet

func (t byTime) Len() int {
	return len(t)
}

func (t byTime) Swap(i, j int) {
	t[i], t[j] = t[j], t[i]
}

func (t byTime) Less(i, j int) bool {
	return t[i].Time > t[j].Time
}

func (t byId) Len() int {
	return len(t)
}

func (t byId) Swap(i, j int) {
	t[i], t[j] = t[j], t[i]
}

func (t byId) Less(i, j int) bool {
	return t[i].ID > t[j].ID
}
