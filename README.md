# gorm_generate

Gorm model auto generate scaffold

Helpers:
```json
  -c string
        Special config file, format: .yml if empty, the default db connection file is .yml
  -d string
        Generated directory
  -db string
        DB connect dns
  -name string
        Model name
  -t string
        Table name of generated model
```

Default config file format:
```yaml
db: username:password@tcp(host.mysql.rds.com:3306)/mplive?charset=utf8mb4&collation=utf8mb4_general_ci&parseTime=True&loc=Local
```

####Example: 
    
Generate a model with special table and model name and special config file
    
```json
     ./generate -name PopularAnchor -t mplive.t_popular_anchor -c .yml 
```
and result

```json
package models

type Anchor struct {
	Id                  uint32 `gorm:"column:id"`
	UserId              uint32 `gorm:"column:user_id"`
	State               int8   `gorm:"column:state"`
	BlockExpiredAt      int32  `gorm:"column:block_expired_at"`
	Room                string `gorm:"column:room"`
	Coin                int32  `gorm:"column:coin"`
	TotalCoin           uint32 `gorm:"column:total_coin"`
	Level               int32  `gorm:"column:level"`
	LevelScore          int32  `gorm:"column:level_score"`
	Certificated        int8   `gorm:"column:certificated"`
	CertificatedAt      int32  `gorm:"column:certificated_at"`
	CertificatedDetail  string `gorm:"column:certificated_detail"`
	SignStatus          uint8  `gorm:"column:sign_status"`
	SignCompany         string `gorm:"column:sign_company"`
	SignDetail          string `gorm:"column:sign_detail"`
	LastWithdrawAccount string `gorm:"column:last_withdraw_account"`
	TotalWithdrawCoin   uint32 `gorm:"column:total_withdraw_coin"`
	TotalWithdrawAmount uint32 `gorm:"column:total_withdraw_amount"`
	SignedAt            uint32 `gorm:"column:signed_at"`
	UpdatedAt           int32  `gorm:"column:updated_at"`
	CreatedAt           int32  `gorm:"column:created_at"`
}

func (Anchor) TableName() string {
	return "mplive.t_anchor"
}

```
    
