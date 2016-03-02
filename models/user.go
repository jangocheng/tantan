package models

type UserInfo struct {
	UserId     int64  `xorm:"user_id INT NOT NULL autoincr" json:"user_id"`
	UserName   string `xorm:"user_name VARHCAR(128) NOT NULL" json:"user_name"`
	Type       string `xorm:"type VARHCAR(128) NOT NULL" json:"type"`
	CreatedUTC int    `xorm:"created_utc INT NOT NULL" json:"created_utc"`
	Status     int    `xorm:"status INT NOT NULL" json:"status"`
}

func (*UserInfo) TableName() string {
	return "user_info"
}

func (m *UserInfo) UniqueCond() (string, []interface{}) {
	return "user_id=?", []interface{}{m.UserId}
}

func GetAllUsers(s *ModelSession, page, count int) ([]*UserInfo, error) {
	var (
		err error
		res = make([]*UserInfo, 0)
	)

	if s == nil {
		s = newAutoCloseModelsSession()
	}
	err = s.OrderBy("user_id desc").Limit(page, (page-1)*count).Find(&res)

	return res, err
}

func GetUserByName(s *ModelSession, name string) (int64, error) {
	var (
		err   error
		count int64
		res   = new(UserInfo)
	)
	if s == nil {
		s = newAutoCloseModelsSession()
	}

	count, err = s.Where("user_name =?", name).Count(res)
	return count, err
}
