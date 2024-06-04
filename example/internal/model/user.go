package model

type Person struct {
	ID      uint   `gorm:"primaryKey;autoIncrement;column:id"` // 定义主键 id，并且开启自动递增
	Name    string `gorm:"column:name"`
	Age     int    `gorm:"column:age"`
	Email   string `gorm:"column:email;uniqueIndex"`
	Address string `gorm:"column:address"`
}

func (Person) TableName() string {
	return "person"
}
