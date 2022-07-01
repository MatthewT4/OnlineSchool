package blogic

import (
	"context"
	"encoding/json"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

func (b *BLogic) GetHomework(userId int, courseId int, homeworkId int) (int, []byte) {
	/*res, err := b.DBSaveHomework.GetHomework(context.TODO(), userId, courseId, homeworkId)
	if err != nil {
		return 404, []byte("not found")
	}
	if res.PublicDate.After(time.Now()) {
		fmt.Println("hw before time.now")
		return 404, []byte("not found")
	}
	var hw struct {
		HomeworkName string                 `json:"homework_name"`
		Deadline     time.Time              `json:"deadline"`
		HomeworkId   int                    `json:"homework_id"`
		Tasks        []structs.HomeworkTask `json:"tasks"`
		Result       int                    `json:"result,omitempty"`
		MaxPoints    int                    `json:"max_points"`
		Delivered    time.Time              `json:"delivered,omitempty"`
	}

	hw.HomeworkName = res.HomeworkName
	hw.Deadline = res.Deadline
	hw.HomeworkId = res.HomeworkId
	hw.Tasks = res.Tasks
	hw.MaxPoints = res.MaxPoints
	hw.Result = res.Result
	hw.Delivered = res.Delivered

	js, er := json.Marshal(&hw)
	if er != nil {
		return 500, []byte("Server error")
	}
	return 200, js*/
	return 500, []byte("Server error")
}

func (b *BLogic) GetNextCourseHomeworks(userId, courseId int) (int, []byte) {
	res, err := b.DBUser.GetCourses(context.TODO(), userId)
	if err != nil {
		return 404, []byte("not found")
	}
	if !b.checkUserCourse(res, courseId) {
		return 404, []byte("not found")
	}
	HwTemp, er := b.DBTempHomework.GetNextTempHomeworks(context.TODO(), courseId)
	if er != nil {
		if er == mongo.ErrNoDocuments {
			return 200, []byte("[]")
		}
		return 404, []byte("not found")
	}
	HwSave, erro := b.DBSaveHomework.GetNextSaveHomeworks(context.TODO(), courseId, userId)
	if erro != nil && erro != mongo.ErrNoDocuments {
		return 500, []byte("Server error")
	}

	var IdDelHw []int
	for i := 0; i < len(HwSave); i++ {
		IdDelHw = append(IdDelHw, HwSave[i].HomeworkId)
	}
	type returnHws struct {
		HomeworkName string    `json:"homework_name"`
		Deadline     time.Time `json:"deadline"`
		CourseId     int       `json:"course_id"`
		HomeworkId   int       `json:"homework_id"`
	}
	var arrRet []returnHws
	for i := 0; i < len(HwTemp); i++ {
		flag := true
		for _, val := range IdDelHw {
			if HwTemp[i].HomeworkId == val {
				flag = false
				break
			}
		}
		if flag {
			var v returnHws
			v.CourseId = HwTemp[i].CourseId
			v.HomeworkId = HwTemp[i].HomeworkId
			v.HomeworkName = HwTemp[i].HomeworkName
			v.Deadline = HwTemp[i].Deadline
			arrRet = append(arrRet, v)
		}
	}

	if len(arrRet) == 0 {
		return 404, []byte("not found")
	}

	js, er := json.Marshal(&arrRet)
	if er != nil {
		return 500, []byte("Server error")
	}
	return 200, js
}

func (b *BLogic) GetNextHomeworks(userId int) (int, []byte) {
	type returnHws struct {
		HomeworkName string    `json:"homework_name"`
		Deadline     time.Time `json:"deadline"`
		CourseId     int       `json:"course_id"`
		HomeworkId   int       `json:"homework_id"`
	}
	var arrRet []returnHws
	res, err := b.DBUser.GetCourses(context.TODO(), userId)
	if err != nil {
		return 404, []byte("not found")
	}
	for _, course := range res {
		HwTemp, er := b.DBTempHomework.GetNextTempHomeworks(context.TODO(), course.CourseId)
		if er != nil {
			if er == mongo.ErrNoDocuments {
				return 200, []byte("[]")
			}
			return 404, []byte("not found")
		}
		HwSave, erro := b.DBSaveHomework.GetNextSaveHomeworks(context.TODO(), course.CourseId, userId)
		if erro != nil && erro != mongo.ErrNoDocuments {
			return 500, []byte("Server error")
		}

		var IdDelHw []int
		for i := 0; i < len(HwSave); i++ {
			IdDelHw = append(IdDelHw, HwSave[i].HomeworkId)
		}
		for i := 0; i < len(HwTemp); i++ {
			flag := true
			for _, val := range IdDelHw {
				if HwTemp[i].HomeworkId == val {
					flag = false
					break
				}
			}
			if flag {
				var v returnHws
				v.CourseId = HwTemp[i].CourseId
				v.HomeworkId = HwTemp[i].HomeworkId
				v.HomeworkName = HwTemp[i].HomeworkName
				v.Deadline = HwTemp[i].Deadline
				arrRet = append(arrRet, v)
			}
		}
	}

	if len(arrRet) == 0 {
		return 404, []byte("not found")
	}

	js, er := json.Marshal(&arrRet)
	if er != nil {
		return 500, []byte("Server error")
	}
	return 200, js
}
