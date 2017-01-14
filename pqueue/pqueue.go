// pqueue implements a generic priority queue, for use in the search
// problems.

package pqueue

type Item struct {
	Priority int
	Value interface{}
	Index int
}

func NewItem(priority int, value interface{}) *Item {
	return &Item{priority, value, 0}
}

type Queue []*Item

func (q Queue) Len() int { return len(q) }

func (q Queue) Less(i, j int) bool {
	return q[i].Priority < q[j].Priority
}

func (q Queue) Swap(i, j int) {
	q[i], q[j] = q[j], q[i]
	q[i].Index = i
	q[j].Index = j
}

func (q *Queue) Push(x interface{}) {
	n := len(*q)
	item := x.(*Item)
	item.Index = n
	*q = append(*q, item)
}

func (q *Queue) Pop() interface{} {
	old := *q
	n := len(old)
	item := old[n-1]
	item.Index = -1 // for safety
	*q = old[0 : n-1]
	return item
}
