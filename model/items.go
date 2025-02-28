package model

import (
	"time"

	sqlite_vec "github.com/asg017/sqlite-vec-go-bindings/cgo"
	"github.com/spf13/viper"
)

// 声明一个符合标准的结构体
type Item struct {
	PostID    string
	Embedding []float32
	Title     string
	Content   string
	CreatedAt string
	UpdatedAt string
	Distance  float64
}

// 插入向量数据
func InsertDocument(item Item) error {
	sqlite_vec.Auto()

	// 进行blob序列化
	v, err := sqlite_vec.SerializeFloat32(item.Embedding)
	if err != nil {
		return err
	}
	now := time.Now().Format(time.RFC3339)
	_, err = VecDB.Exec("INSERT INTO items(post_id, embedding, title,content, created_at, updated_at) VALUES (?, ?, ?,?, ?, ?)",
		item.PostID, v, item.Title, item.Content, now, now)
	if err != nil {
		return err
	}
	return nil
}

// 查询向量数据
// 查询向量数据
func GetDocument(input []float32) ([]Item, error) {
	sqlite_vec.Auto()
	query, err := sqlite_vec.SerializeFloat32(input)
	if err != nil {
		return nil, err
	}

	doc_limit := viper.GetInt("app.doc_limit")

	rows, err := VecDB.Query(`
		SELECT
			post_id,
			title,
			content,
			vec_distance_L2(embedding, ?) AS distance
		FROM items
		WHERE embedding IS NOT NULL
		ORDER BY distance
		LIMIT ?
	`, query, doc_limit)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// 数据列表
	var items []Item
	for rows.Next() {
		var post_id string
		var title string
		var content string
		var distance float64
		err = rows.Scan(&post_id, &title, &content, &distance)
		if err != nil {
			return nil, err
		}
		items = append(items, Item{PostID: post_id, Title: title, Content: content, Distance: distance})
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return items, nil
}

// 清空整个表
func TruncateItems() error {
	_, err := VecDB.Exec("DELETE FROM items")
	if err != nil {
		return err
	}
	return nil
}

// 根据postid删除单行数据
func DeleteItem(post_id string) error {
	_, err := VecDB.Exec("DELETE FROM items WHERE post_id = ?", post_id)
	if err != nil {
		return err
	}
	return nil
}
