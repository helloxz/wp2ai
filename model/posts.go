package model

import (
	"fmt"
	"time"
)

type Post struct {
	ID          uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	PostID      uint      `gorm:"not null;uniqueIndex" json:"post_id"` // 添加 uniqueIndex
	PostDate    string    `gorm:"type:TEXT" json:"post_date"`          // 使用 *time.Time 处理 TEXT 类型
	PostContent string    `gorm:"type:TEXT" json:"post_content"`
	PostTitle   string    `gorm:"type:TEXT" json:"post_title"`
	Status      int       `gorm:"not null;default:0" json:"status"` // 默认0：未处理，1：处理中，2：待更新，3：已处理,4:存在错误
	CreatedAt   time.Time `gorm:"not null" json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// TableName sets the insert table name for struct
func (Post) TableName() string {
	return "posts"
}

// 写一个函数，得到一批postid，然后一次性插入到数据库中
func InsertPosts(posts []Post) error {
	return DB.Create(&posts).Error
}

// 插入单个数据
func InsertPost(post Post) error {
	return DB.Create(&post).Error
}

// 写一个函数，函数指定查询的数量，然后查询出status=0的数据
func GetPosts(limit int) []Post {
	var posts []Post
	DB.Where("status = ?", 0).Limit(limit).Find(&posts)
	return posts
}

// 写一个函数，函数接收一批postid，然后批量更新status=1
func UpdatePostsStatus(ids []uint, status int) error {
	return DB.Model(&Post{}).Where("post_id IN (?)", ids).Update("status", status).Error
}

// 根据postid，更新单行数据
func UpdatePost(post Post) error {
	return DB.Model(&Post{}).Where("post_id = ?", post.PostID).Updates(&post).Error
}

// 写一个函数，根据页码和每页数量，查询出所有的数据
func GetPostsByPage(page, limitInt int) ([]Post, error) {
	var posts []Post
	err := DB.Offset((page - 1) * limitInt).Limit(limitInt).Order("post_id desc").Find(&posts).Error
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	// fmt.Println("db", posts)
	return posts, nil
}

// 根据状态统计所有文章数量，如果状态为-1，则统计全部文章
func CountPosts(status int) int {
	var count int64
	if status == -1 {
		DB.Model(&Post{}).Count(&count)
	} else {
		DB.Model(&Post{}).Where("status = ?", status).Count(&count)
	}
	return int(count)
}

// 清空posts整个表
func TruncatePosts() error {
	return DB.Exec("DELETE FROM posts").Error
}

// 根据postid删除单行数据
func DeletePost(postId uint) error {
	return DB.Where("post_id = ?", postId).Delete(&Post{}).Error
}
