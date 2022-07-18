package blogic

import (
	"context"
	"encoding/json"
	"fmt"
)

func (b *BLogic) GetInfoCourse(userId int64, courseId int) (int, []byte) {
	coursesMas, err := b.DBUser.GetCourses(context.TODO(), userId)
	if err != nil {
		return 404, []byte("not found")
	}
	index := -1
	for key, val := range coursesMas {
		if val.CourseId == courseId {
			index = key
		}
	}
	if index == -1 {
		return 404, []byte("not found")
	}
	if !coursesMas[index].Active {
		return 404, []byte("not found")
	}
	course, er := b.DBCourse.GetCourse(context.TODO(), courseId)
	if er != nil {
		return 404, []byte("not found")
	}

	var RetData struct {
		HeartCount  int    `json:"heart_count"`
		FreezingDay int    `json:"freezing_day"`
		Freeze      bool   `json:"freeze"`
		VkChat      string `json:"vk_chat"`
	}
	RetData.HeartCount = coursesMas[index].HeartCount
	RetData.Freeze = coursesMas[index].Freeze
	RetData.FreezingDay = coursesMas[index].FreezingDay
	RetData.VkChat = course.VkChat

	re, erro := json.Marshal(&RetData)
	if erro != nil {
		fmt.Println(erro)
		return 500, []byte("json marshal fail")
	}
	return 200, re
}
