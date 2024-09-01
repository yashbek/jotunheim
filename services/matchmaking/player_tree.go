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
