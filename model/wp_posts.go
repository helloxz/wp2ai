package model

// 声明wp_posts结构体
type WpPost struct {
	ID          uint   `gorm:"primaryKey" json:"id"`
	PostDate    string `gorm:"column:post_date" json:"post_date"`
	PostContent string `gorm:"column:post_content;type:longtext" json:"post_content"`
	PostTitle   string `gorm:"column:post_title" json:"post_title"`
	PostStatus  string `gorm:"column:post_status" json:"post_status"`
}

// TableName 设置表名，如果你的表名不是默认的Post的复数形式
func (WpPost) TableName() string {
	return "wp_posts"
}

// 写一个函数，查询出所有post_status='publish'的文章，只查询ID字段即可
func GetPostIds() []WpPost {
	var posts []WpPost
	err := WP.Select("ID").
		Where("post_status = ?", "publish").
		Where("post_type = ?", "post").
		Where("post_password = ?", "").
		Order("id desc").Find(&posts).Error
	if err != nil {
		return []WpPost{}
	}
	return posts
}

// 根据一批ID，查询出匹配ID的数据
func GetPostsByIds(ids []uint) []WpPost {
	var posts []WpPost
	err := WP.Where("id IN (?)", ids).Find(&posts).Error
	if err != nil {
		return []WpPost{}
	}
	return posts
}

// 根据ID，和where条件，检查是否存在
func GetWpPostById(id uint) (*WpPost, error) {
	var post *WpPost
	err := WP.Where("id = ?", id).
		Where("post_status = ?", "publish").
		Where("post_type = ?", "post").
		Where("post_password = ?", "").First(&post).Error
	if err != nil {
		return nil, err
	}
	return post, nil
}
