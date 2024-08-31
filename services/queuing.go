package services

var queue map[int][]string

const defaultEloInterval = 50

func getQueue() map[int][]string {
	return queue
}

func AddToMatchMakingQueue(userID string, elo int) {
	q := getQueue()
	key := getKey(elo)
	if len(q[key]) == 0 {
		q[key] = make([]string, 0)
	}
	q[key] = append(q[key], userID)
}
 
func getKey(num int) int {
	return num / defaultEloInterval
}