package repositories

//package repositories

import (
	"encoding/json"
	"fmt"
	"time"
	. "zyg/datamodels"

	"github.com/go-redis/redis"
	//"zyg/datamodels"
)

type redisConfings struct {
	Addr   string
	Pass   string
	Tag    string
	BookDb int
	SessDb int
}

var RedisConfig redisConfings

func RedisNewDB(db int) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     RedisConfig.Addr,
		Password: RedisConfig.Pass, // no password set
		DB:       db,               // use default DB
	})
	_, err := client.Ping().Result()
	if err != nil {
		return client, err
	}
	return client, err
}

func GetRedisBookData(value string) (BookData, error) {
	var book BookData
	client, err := RedisNewDB(RedisConfig.BookDb)
	if err != nil {
		return book, redis.Nil
	}
	defer client.Close()
	value = RedisConfig.Tag + "Book_" + value
	bookJson, err := client.Get(value).Result()
	if err == redis.Nil {
		return book, err
	} else if err != nil {
		return book, err
	} else {
		err = json.Unmarshal([]byte(bookJson), &book)
	}
	return book, err
}

func WriteRedisBookData(value string, Book BookData, times time.Duration) {
	client, err := RedisNewDB(RedisConfig.BookDb)
	defer client.Close()
	if err != nil {

	} else {
		data, err := json.Marshal(Book)
		if err != nil {
			fmt.Printf("\n转换为json错误%v", err)
		}
		value = RedisConfig.Tag + "Book_" + value
		err = client.Set(value, data, times).Err()
		if err != nil {
			fmt.Printf("\n写入出错%v", err)
		}
	}

}

func GetRedisRecBookData(value string) (ArticleStructSet, error) {
	var bookSet ArticleStructSet
	client, err := RedisNewDB(RedisConfig.BookDb)
	if err != nil {
		return bookSet, redis.Nil
	}
	defer client.Close()
	value = RedisConfig.Tag + "Rec_" + value
	bookJson, err := client.Get(value).Result()
	if err == redis.Nil {
		return bookSet, err
	} else if err != nil {
		return bookSet, err
	} else {
		err = json.Unmarshal([]byte(bookJson), &bookSet)
	}
	return bookSet, err
}

func WriteRedisRecBookData(value string, bookSet ArticleStructSet, times time.Duration) {
	client, err := RedisNewDB(RedisConfig.BookDb)
	defer client.Close()
	if err != nil {

	} else {
		data, err := json.Marshal(bookSet)
		if err != nil {
			fmt.Printf("\n转换为json错误%v", err)
		}
		value = RedisConfig.Tag + "Rec_" + value
		err = client.Set(value, data, times).Err()
		if err != nil {
			fmt.Printf("\n写入出错%v", err)
		}
	}
}
func SessSet(user UserInfo, times time.Duration) error {
	client, err := RedisNewDB(RedisConfig.SessDb)
	if err != nil {
		return redis.Nil
	}
	defer client.Close()
	data, err := json.Marshal(user)
	if err != nil {
		return err
	}
	value := RedisConfig.Tag + "sess_" + user.Userid
	err = client.Set(value, data, times).Err()
	if err != nil {
		return err
	}
	return nil
}
func SessGet(userid string) (UserInfo, error) {
	var user UserInfo
	client, err := RedisNewDB(RedisConfig.SessDb)
	if err != nil {
		return user, redis.Nil
	}
	defer client.Close()
	value := RedisConfig.Tag + "sess_" + userid
	userJson, err := client.Get(value).Result()
	if err == redis.Nil {
		return user, err
	} else if err != nil {
		return user, err
	} else {
		err = json.Unmarshal([]byte(userJson), &user)
	}
	return user, err
}
func SessDel(userid string) error {
	client, err := RedisNewDB(RedisConfig.SessDb)
	if err != nil {
		return redis.Nil
	}
	defer client.Close()
	value := RedisConfig.Tag + "sess_" + userid
	err = client.Del(value).Err()
	return err
}
