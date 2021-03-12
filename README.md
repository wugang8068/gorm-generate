# gorm_generate

Gorm model auto generate scaffold

Helpers:
```json
    -config string
          Special config file, format: .yml
    -connection string
          DB connect dns
    -dao string
          The directory of dao generate.
    -model-directory string
          Generated model directory
    -model-file string
          Generate model file name
    -model-name string
          Generate model struct name
    -repo string
          The directory of repository generate.
    -table string
          Table name of generated model
```

Default config file format:
```yaml
db: username:password@tcp(host.mysql.rds.com:3306)/mplive?charset=utf8mb4&collation=utf8mb4_general_ci&parseTime=True&loc=Local
```

####Example: 
    
Generate a model with special table and model name and special config file
    
```json
     ./generate -table mplive.t_popular_anchor -model-name PopularAnchor -model-directory models -repo repo -dao dao -model-file popular_anchor -config .yml
```
and result is below: 
popular_anchor.go
```go
package models

type PopularAnchor struct { 
	Id uint32 `json:"id" gorm:"column:id"`
	UserId int32 `json:"user_id" gorm:"column:user_id"`
	VisitorCount int32 `json:"visitor_count" gorm:"column:visitor_count"`
	CreatedAt int32 `json:"created_at" gorm:"column:created_at"`
	UpdatedAt int32 `json:"updated_at" gorm:"column:updated_at"`
}

func(PopularAnchor) TableName() string {
	return "mplive.t_popular_anchor"
}
```

popular_anchor_dao.go
```go
package dao

import (
	"github.com/jinzhu/gorm"
	models "gorm_generate/models"
	"gorm_generate/mysql"
)

type PopularAnchorDao struct { }

func(PopularAnchorDao) List() (l []*models.PopularAnchor) {
	mysql.DefaultConnection().Order("id desc").Find(&l) 
	return
}

func(PopularAnchorDao) GetById(id uint32) (*models.PopularAnchor, error) {
	var m models.PopularAnchor
	e := mysql.DefaultConnection().Where("id = ?", id).First(&m).Error
	if e != nil && e != gorm.ErrRecordNotFound {
		return nil, e 
	}
	return &m, nil
}

func (PopularAnchorDao) Create(m models.PopularAnchor) (*models.PopularAnchor, error)  {
	e := mysql.DefaultConnection().Create(&m).Error
	if e != nil {
		return nil, e
	}
	return &m, nil
}

func (PopularAnchorDao) Update(m models.PopularAnchor, updates map[string]interface{}) (*models.PopularAnchor, error) {
	if len(updates) == 0 {
		return &m, nil
	}
	e := mysql.DefaultConnection().Model(&m).UpdateColumns(updates).Error
	if e != nil {
		return nil, e
	}
	return &m, nil
}

func (PopularAnchorDao) Delete(m models.PopularAnchor) error {
	return mysql.DefaultConnection().Delete(m).Error
}
```

popular_anchor_repo.go
```go
package repo

import (
	models "gorm_generate/models"
	dao "gorm_generate/dao"
)

type PopularAnchorRepository interface {
	List() (l []*models.PopularAnchor)
	GetById(id uint32) (*models.PopularAnchor, error)
	Create(m models.PopularAnchor) (*models.PopularAnchor, error)
	Update(m models.PopularAnchor, updates map[string]interface{}) (*models.PopularAnchor, error)
	Delete(m models.PopularAnchor) error
}

func NewPopularAnchorRepository() PopularAnchorRepository {
	return dao.PopularAnchorDao{}
}
```