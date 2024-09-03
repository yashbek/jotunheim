package matchmaking

import (
	"log"
	"math"

	"github.com/yashbek/jotunheim/models"
	"github.com/yashbek/jotunheim/utils"
)

type PlayerTreeRecord struct {
	profile     models.Profile
	rightKey  string
	leftKey   string
	parentKey string
	height      int
}

var queue map[string]PlayerTreeRecord
var rootKey string

func Init() {
	queue = make(map[string]PlayerTreeRecord, 0)
}

func AddRecord(profile models.Profile) {
	record := wrapProfile(profile)

	if len(queue) == 0 {
		queue[profile.ID] = record
		rootKey = profile.ID
		return
	}

	recursiveTraverse(record, queue[rootKey])
}

func recursiveTraverse(new, curr PlayerTreeRecord) {
	diff := new.profile.Elo - curr.profile.Elo
	if math.Abs(float64(diff)) <= utils.DefaultEloInterval {
		err := startGame(new.profile.ID, curr.profile.ID)
		if err != nil {
			log.Default().Print("couldn't start game ")
		}
		// Pair players, remove record
		// @TODO implement deletion
		return 
	}

	if diff > 0 {
		if curr.rightKey != "" {
			new.parentKey = curr.profile.ID
			curr.rightKey = new.profile.ID
		} else {
			recursiveTraverse(new, queue[curr.rightKey])
		}
	} else {
		if curr.leftKey != "" {
			new.parentKey = curr.profile.ID
			curr.leftKey = new.profile.ID
		} else {
			recursiveTraverse(new, queue[curr.rightKey])
		}
	}

	updateHeight(curr.profile.ID)

	balance := queue[curr.rightKey].height - queue[curr.leftKey].height

	if math.Abs(float64(balance)) < 2 {
		return
	}

	if balance > 0 {
		if new.profile.Elo > queue[curr.rightKey].profile.Elo {
			// right right 
			counterClockWiseRotate(curr.profile.ID)
		} else {
			// right left
			clockWiseRotate(queue[curr.rightKey].profile.ID)
			counterClockWiseRotate(curr.profile.ID)
		}
	} else {
		if new.profile.Elo > queue[curr.rightKey].profile.Elo {
			// left right
			counterClockWiseRotate(queue[curr.leftKey].profile.ID)
			clockWiseRotate(curr.profile.ID)
		} else {
			// left left 
			clockWiseRotate(curr.profile.ID)
		}
	}

}

func wrapProfile(profile models.Profile) PlayerTreeRecord {
	return PlayerTreeRecord{
		profile:     profile,
		rightKey:  "",
		leftKey:   "",
		parentKey: "",
		height:      1,
	}
}

func updateHeight(recordKey string) {
	copy := queue[recordKey]
	copy.height =  1 + utils.Max(queue[copy.rightKey].height, queue[copy.leftKey].height)
	queue[recordKey] = copy
}

func setRecord(recordKey string, record PlayerTreeRecord) {
	queue[recordKey] = record
}

func updateRecord(record PlayerTreeRecord) {
	setRecord(record.profile.ID, record)
}

func counterClockWiseRotate(subTreeRootkey string) {
	oldRoot := queue[subTreeRootkey]
	firstRight := queue[oldRoot.rightKey]
	firstLeft := queue[firstRight.leftKey]

	firstRight.leftKey = oldRoot.profile.ID
	firstRight.parentKey = oldRoot.parentKey
	updateRecord(firstRight)

    oldRoot.rightKey = firstLeft.profile.ID
	oldRoot.parentKey = firstRight.profile.ID
	updateRecord(oldRoot)

	firstLeft.parentKey = oldRoot.profile.ID
	updateRecord(firstLeft)


    updateHeight(oldRoot.profile.ID)
	updateHeight(firstRight.profile.ID)
}

func clockWiseRotate(subTreeRootkey string) {
	oldRoot := queue[subTreeRootkey]
	firstLeft := queue[oldRoot.leftKey]
	firstRight := queue[firstLeft.rightKey]

	firstLeft.rightKey = oldRoot.profile.ID
	firstLeft.parentKey = oldRoot.parentKey
	updateRecord(firstLeft)

    oldRoot.leftKey = firstRight.profile.ID
	oldRoot.parentKey = firstLeft.profile.ID
	updateRecord(oldRoot)

	firstRight.parentKey = oldRoot.profile.ID
	updateRecord(firstRight)


    updateHeight(oldRoot.profile.ID)
	updateHeight(firstRight.profile.ID)
}

func startGame(p1, p2 string) error {
	for _, p := range []string{p1, p2} {
		_, exists := queue[p]
		if exists {
			deleteRecord(p)
		}
	}

	return nil
}

func deleteRecord(recordKey string) {
	toBeRemoved := queue[recordKey]
	var candidateKey string

	switch getNumberOfChildren(recordKey) {
	case 2:
		rightMostRecordInLeftSubtree := getRightMostChild(toBeRemoved.leftKey)
		// condition for an edge case
		RightMostIsSubTreeRoot := rightMostRecordInLeftSubtree.profile.ID == toBeRemoved.leftKey
		// to avoid an edge case where the right most element in the left subtree is the left subtree root itself
		if RightMostIsSubTreeRoot {
			rightMostRecordInLeftSubtree.parentKey = toBeRemoved.parentKey
			rightMostRecordInLeftSubtree.rightKey = toBeRemoved.rightKey
			candidateKey = rightMostRecordInLeftSubtree.profile.ID
			updateRecord(rightMostRecordInLeftSubtree)
			break
		} 
		parent := queue[rightMostRecordInLeftSubtree.parentKey]
		parent.rightKey = ""

		if rightMostRecordInLeftSubtree.leftKey != "" {
			left := queue[rightMostRecordInLeftSubtree.leftKey]			
			left.parentKey = parent.profile.ID
			parent.rightKey = left.profile.ID
			updateRecord(left)
		}
		updateRecord(parent)

		rightMostRecordInLeftSubtree.rightKey = toBeRemoved.rightKey
		rightMostRecordInLeftSubtree.parentKey = toBeRemoved.parentKey
		rightMostRecordInLeftSubtree.leftKey = toBeRemoved.leftKey
		candidateKey = rightMostRecordInLeftSubtree.profile.ID
		updateRecord(rightMostRecordInLeftSubtree)
	case 1:
		if toBeRemoved.leftKey != "" {
			candidateKey = toBeRemoved.leftKey
		} else {
			candidateKey = toBeRemoved.rightKey
		}
	case 0:
		candidateKey = ""
	} 

	parent, exists := queue[toBeRemoved.parentKey]
	if exists {
		if toBeRemoved.profile.Elo > parent.profile.Elo {
			parent.rightKey = candidateKey
		} else {
			parent.leftKey = candidateKey
		}
		updateRecord(parent)
	}
	
	delete(queue, toBeRemoved.profile.ID)
}

func getRightMostChild(key string) PlayerTreeRecord {
	curr := queue[key]
	if curr.rightKey != "" {
		return getRightMostChild(curr.rightKey)
	}
	return curr
}

func getNumberOfChildren(key string) int {
	record := queue[key]
	count := 0
	if record.leftKey != "" {count++}
	if record.rightKey != "" {count++}
	return count
}
