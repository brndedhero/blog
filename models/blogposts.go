package models

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/brndedhero/blog/config"
	"github.com/brndedhero/blog/helpers"
	"github.com/sirupsen/logrus"
)

type BlogPost struct {
	Id        uint64     `json:"id" gorm:"primary_key"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
	DeletedAt *time.Time `json:"-" sql:"index"`
	Title     string     `json:"title"`
	Body      string     `gorm:"type:text" json:"body"`
	Tags      []Tag      `gorm:"many2many:post_tags;" json:"tags"`
}

func GetAllBlogPosts() (string, error) {
	var blogPosts []BlogPost
	data, err := config.Redis.Get(config.Redis.Context(), "blogPost:all").Result()
	if err != nil {
		config.Log.WithFields(logrus.Fields{
			"app":      "redis",
			"func":     "getAllBlogPosts",
			"redisKey": "blogPost:all",
		}).Warn("key not found")
		res := config.DB.Find(&blogPosts)
		if res.Error != nil {
			config.Log.WithFields(logrus.Fields{
				"app":  "mysql",
				"func": "getAllBlogPosts",
			}).Error(res.Error)
			message, _ := helpers.PrepareString(404, nil)
			return message, res.Error
		}
		json, err := json.Marshal(blogPosts)
		if err != nil {
			config.Log.WithFields(logrus.Fields{
				"app":  "blog",
				"func": "getAllBlogPosts",
			}).Error(err)
			message, _ := helpers.PrepareString(500, nil)
			return message, err
		}
		redisErr := config.Redis.Set(config.Redis.Context(), "blogPost:all", json, 0).Err()
		if redisErr != nil {
			config.Log.WithFields(logrus.Fields{
				"app":      "redis",
				"func":     "getAllBlogPosts",
				"redisKey": "blogPost:all",
			}).Warn(redisErr)
		}
		message, _ := helpers.PrepareString(200, json)
		return message, nil
	}
	message, _ := helpers.PrepareString(200, []byte(data))
	return message, nil
}

func CreateBlogPost(title string, body string) (string, error) {
	blogPost := BlogPost{Title: title, Body: body}
	res := config.DB.Create(&blogPost)
	if res.Error != nil {
		config.Log.WithFields(logrus.Fields{
			"app":  "mysql",
			"func": "createBlogPost",
		}).Error(res.Error)
		message := helpers.PrepareErrorString(500, res.Error)
		return message, res.Error
	}
	redisKey := "blogPost:" + fmt.Sprint(blogPost.Id)
	json, err := json.Marshal(blogPost)
	if err != nil {
		config.Log.WithFields(logrus.Fields{
			"app":  "blog",
			"func": "createBlogPost",
		}).Error(err)
		message := helpers.PrepareErrorString(500, err)
		return message, err
	}
	redisErr := config.Redis.Set(config.Redis.Context(), redisKey, json, 0).Err()
	if redisErr != nil {
		config.Log.WithFields(logrus.Fields{
			"app":      "redis",
			"func":     "createBlogPost",
			"redisKey": redisKey,
		}).Warn(redisErr)
	}
	redisDelErr := config.Redis.Del(config.Redis.Context(), "blogPost:all").Err()
	if redisDelErr != nil {
		config.Log.WithFields(logrus.Fields{
			"app":      "redis",
			"func":     "createBlogPost",
			"redisKey": "blogPost:all",
		}).Warn(redisDelErr)
	}
	data, _ := helpers.PrepareCreateJson(int(blogPost.Id), int(res.RowsAffected))
	message, _ := helpers.PrepareString(201, data)
	return message, nil
}

func GetBlogPost(id uint64) (string, error) {
	var blogPost BlogPost
	redisKey := fmt.Sprintf("blogPost:%d", id)
	data, err := config.Redis.Get(config.Redis.Context(), redisKey).Result()
	if err != nil {
		config.Log.WithFields(logrus.Fields{
			"app":      "redis",
			"func":     "getBlogPost",
			"redisKey": redisKey,
		}).Warn("key not found")
		res := config.DB.First(&blogPost, id)
		if res.Error != nil {
			config.Log.WithFields(logrus.Fields{
				"app":  "mysql",
				"func": "getBlogPost",
			}).Error(res.Error)
			message, _ := helpers.PrepareString(404, nil)
			return message, res.Error
		}
		json, err := json.Marshal(blogPost)
		if err != nil {
			config.Log.WithFields(logrus.Fields{
				"app":  "blog",
				"func": "getBlogPost",
			}).Error(err)
			message := helpers.PrepareErrorString(500, err)
			return message, err
		}
		redisErr := config.Redis.Set(config.Redis.Context(), redisKey, json, 0).Err()
		if redisErr != nil {
			config.Log.WithFields(logrus.Fields{
				"app":      "redis",
				"func":     "getBlogPost",
				"redisKey": redisKey,
			}).Warn(redisErr)
			redisDelErr := config.Redis.Del(config.Redis.Context(), redisKey).Err()
			if redisDelErr != nil {
				config.Log.WithFields(logrus.Fields{
					"app":      "redis",
					"func":     "getBlogPost",
					"redisKey": redisKey,
				}).Warn(redisDelErr)
			}
		}
		message, _ := helpers.PrepareString(200, json)
		return message, nil
	}
	message, _ := helpers.PrepareString(200, []byte(data))
	return message, nil
}

func UpdateBlogPost(id uint64, title string, body string) (string, error) {
	var blogPost BlogPost
	res := config.DB.First(&blogPost, id)
	if res.Error != nil {
		config.Log.WithFields(logrus.Fields{
			"app":  "mysql",
			"func": "updateBlogPost",
		}).Error(res.Error)
		message, _ := helpers.PrepareString(404, nil)
		return message, res.Error
	}
	blogPost.Title = title
	blogPost.Body = body
	res = config.DB.Save(&blogPost)
	if res.Error != nil {
		config.Log.WithFields(logrus.Fields{
			"app":  "mysql",
			"func": "updateBlogPost",
		}).Error(res.Error)
		message, _ := helpers.PrepareString(500, nil)
		return message, res.Error
	}

	jsonData, err := json.Marshal(blogPost)
	if err != nil {
		config.Log.WithFields(logrus.Fields{
			"app":  "blog",
			"func": "updateBlogPost",
		}).Error(err)
		message := helpers.PrepareErrorString(500, err)
		return message, err
	}
	redisKey := fmt.Sprintf("blogPost:%d", blogPost.Id)
	redisErr := config.Redis.Set(config.Redis.Context(), redisKey, jsonData, 0).Err()
	if redisErr != nil {
		config.Log.WithFields(logrus.Fields{
			"app":      "redis",
			"func":     "updateBlogPost",
			"redisKey": redisKey,
		}).Warn(redisErr)
		redisDelErr := config.Redis.Del(config.Redis.Context(), redisKey)
		if redisDelErr != nil {
			config.Log.WithFields(logrus.Fields{
				"app":      "redis",
				"func":     "updateBlogPost",
				"redisKey": redisKey,
			}).Warn(redisDelErr)
		}
	}
	redisDelErr := config.Redis.Del(config.Redis.Context(), "blogPost:all")
	if redisDelErr != nil {
		config.Log.WithFields(logrus.Fields{
			"app":      "redis",
			"func":     "updateBlogPost",
			"redisKey": "blogPost:all",
		}).Warn(redisDelErr)
	}
	data, _ := helpers.PrepareCreateJson(int(blogPost.Id), int(res.RowsAffected))
	message, _ := helpers.PrepareString(200, data)
	return message, nil
}

func DeleteBlogPost(id uint64) (string, error) {
	var blogPost BlogPost
	res := config.DB.First(&blogPost, id)
	if res.Error != nil {
		config.Log.WithFields(logrus.Fields{
			"app":  "mysql",
			"func": "deleteBlogPost",
		}).Error(res.Error)
		message, _ := helpers.PrepareString(404, nil)
		return message, res.Error
	}
	res = config.DB.Delete(&blogPost)
	if res.Error != nil {
		config.Log.WithFields(logrus.Fields{
			"app":  "mysql",
			"func": "deleteBlogPost",
		}).Error(res.Error)
		message, _ := helpers.PrepareString(500, nil)
		return message, res.Error
	}
	redisKey := fmt.Sprintf("blogPost:%d", blogPost.Id)
	redisErr := config.Redis.Del(config.Redis.Context(), redisKey).Err()
	if redisErr != nil {
		config.Log.WithFields(logrus.Fields{
			"app":      "redis",
			"func":     "updateBlogPost",
			"redisKey": redisKey,
		}).Warn(redisErr)
	}
	redisDelErr := config.Redis.Del(config.Redis.Context(), "blogPost:all").Err()
	if redisDelErr != nil {
		config.Log.WithFields(logrus.Fields{
			"app":      "redis",
			"func":     "updateBlogPost",
			"redisKey": "blogPost:all",
		}).Warn(redisDelErr)
	}
	data, _ := helpers.PrepareCreateJson(int(blogPost.Id), int(res.RowsAffected))
	message, _ := helpers.PrepareString(200, data)
	return message, nil
}
