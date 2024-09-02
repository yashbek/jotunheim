package matchmaking

import (
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

func updateRecord(recordKey string, record PlayerTreeRecord) {
	queue[recordKey] = record
}

func counterClockWiseRotate(subTreeRootkey string) {
	oldRoot := queue[subTreeRootkey]
	firstRight := queue[oldRoot.rightKey]
	firstLeft := queue[firstRight.leftKey]

	firstRight.leftKey = oldRoot.profile.ID
	firstRight.parentKey = oldRoot.parentKey
	updateRecord(firstRight.profile.ID, firstRight)

    oldRoot.rightKey = firstLeft.profile.ID
	oldRoot.parentKey = firstRight.profile.ID
	updateRecord(oldRoot.profile.ID, oldRoot)

	firstLeft.parentKey = oldRoot.profile.ID
	updateRecord(firstLeft.profile.ID, firstLeft)


    updateHeight(oldRoot.profile.ID)
	updateHeight(firstRight.profile.ID)

}

func clockWiseRotate(subTreeRootkey string) {
	oldRoot := queue[subTreeRootkey]
	firstLeft := queue[oldRoot.leftKey]
	firstRight := queue[firstLeft.rightKey]

	firstLeft.rightKey = oldRoot.profile.ID
	firstLeft.parentKey = oldRoot.parentKey
	updateRecord(firstLeft.profile.ID, firstLeft)

    oldRoot.leftKey = firstRight.profile.ID
	oldRoot.parentKey = firstLeft.profile.ID
	updateRecord(oldRoot.profile.ID, oldRoot)

	firstRight.parentKey = oldRoot.profile.ID
	updateRecord(firstRight.profile.ID, firstRight)


    updateHeight(oldRoot.profile.ID)
	updateHeight(firstRight.profile.ID)

}