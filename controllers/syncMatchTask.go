package controllers

import "github.com/tantan/models"

var (
	syncMatchTaskChan = make(chan *models.RelationShip, 1024)
	queueLen          = 1
	taskQueue         = make([]*models.RelationShip, queueLen)
)

func init() {
	go syncTask()
}

func doMatch2db(q []*models.RelationShip) (err error) {

	ms := models.NewModelSession()
	defer ms.Close()
	if err = ms.Begin(); err != nil {
		return
	}

	for _, relation := range q {
		relation.State = 2
		if err = models.UpdateDBModel(ms, relation); err != nil {
			ms.Rollback()
			return
		}

		relation.Master, relation.Liker = relation.Liker, relation.Master
		if err = models.UpdateDBModel(ms, relation); err != nil {
			ms.Rollback()
			return
		}
	}

	if err = ms.Commit(); err != nil {
		return
	}

	return
}

func syncTask() {
	for {
		for i := 0; i < queueLen; i++ {
			taskQueue[i] = <-syncMatchTaskChan
		}
		doMatch2db(taskQueue)
	}
}
