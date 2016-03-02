package models

type RelationShip struct {
	Id         int64 `xorm:"id INT NOT NULL autoincr" json:"id"`
	Master     int64 `xorm:"master INT NOT NULL" json:"master"`
	Liker      int64 `xorm:"liker INT NOT NULL" json:"liker"`
	Type       int   `xorm:"type INT NOT NULL" json:"type"`
	State      int   `xorm:"state INT NOT NULL" json:"state"`
	CreatedUTC int   `xorm:"created_utc INT" json:"created_utc"`
	Status     int   `xorm:"status INT NOT NULL" json:"status"`
}

func (*RelationShip) TableName() string {
	return "relationship"
}

func (m *RelationShip) UniqueCond() (string, []interface{}) {
	return "master=? and liker=?", []interface{}{m.Master, m.Liker}
}

func GetAllRelationByUserId(s *ModelSession, uid, page, count int) ([]*RelationShip, error) {
	var (
		err error
		res = make([]*RelationShip, 0)
	)

	if s == nil {
		s = newAutoCloseModelsSession()
	}
	err = s.Where("master=?", uid).OrderBy("id desc").Limit(page, (page-1)*count).Find(&res)

	return res, err
}

func GetRelationById(s *ModelSession, master, linker int64) ([]*RelationShip, error) {
	var (
		err error
		res = make([]*RelationShip, 0)
	)

	if s == nil {
		s = newAutoCloseModelsSession()
	}
	err = s.Where("master=? and liker=?", master, linker).Find(&res)

	return res, err
}

func DelRelation(s *ModelSession, master, linker int64) (int64, error) {
	var (
		err      error
		affected int64
		res      RelationShip
	)
	if s == nil {
		s = newAutoCloseModelsSession()
	}

	res.Master = master
	res.Liker = linker
	affected, err = s.Delete(&res)

	return affected, err
}
