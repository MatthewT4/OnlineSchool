package blogic

import (
	"OnlineSchool/internal/structs"
	"context"
	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type retTask struct {
	Number      int      `json:"number"`
	Text        string   `json:"text"`
	File        []string `json:"file,omitempty"`
	Answers     []string `json:"answers,omitempty"`
	UserAnswers string   `json:"user_answers"`
	Solution    string   `json:"solution,omitempty"`
	Written     bool     `json:"written"`
	TypeAnswers []string `json:"type_answers,omitempty"`
	UserPoint   int      `json:"user_point,omitempty"`
	MaxPoint    int      `json:"max_point"`
}

func (b *BLogic) GetHomework(userId int, homeworkId int) (int, []byte) {
	tempRes, er := b.DBTempHomework.GetHomework(context.TODO(), homeworkId)
	if er != nil {
		if er == mongo.ErrNoDocuments {
			return 404, []byte("not found")
		}
	}
	if tempRes.PublicDate.After(time.Now()) {
		fmt.Println("hw after time.now")
		return 404, []byte("not found")
	}
	save, err := b.DBSaveHomework.GetHomework(context.TODO(), userId, homeworkId)
	if err != nil && err != mongo.ErrNoDocuments {
		return 500, []byte("Server error")
	}
	saveMap := make(map[int]structs.HomeworkTask)
	for _, val := range save.Tasks {
		saveMap[val.TaskId] = val
	}
	var handed bool
	if err == mongo.ErrNoDocuments {
		handed = false
	} else {
		handed = save.Handed
	}

	var tasksId []int
	for _, val := range tempRes.Tasks {
		tasksId = append(tasksId, val.TaskId)
	}
	mapTasks, ok := b.getTasks(tasksId, handed)
	if !ok {
		return 500, []byte("Server error")
	}
	var returnTasks []retTask
	for key, value := range tempRes.Tasks {
		task := mapTasks[value.TaskId]
		var vr retTask
		vr.Number = key + 1
		vr.Text = task.Text
		vr.TypeAnswers = task.TypeAnswers
		vr.File = task.File
		vr.Written = task.Written
		vr.MaxPoint = task.MaxPoint
		vr.UserAnswers = saveMap[value.TaskId].UserAnswer
		if handed {
			vr.UserPoint = saveMap[value.TaskId].Point
			vr.Answers = task.Answers
			vr.Solution = task.Solution
		}
		returnTasks = append(returnTasks, vr)
	}

	js, er := json.Marshal(&returnTasks)
	if er != nil {
		return 500, []byte("Server error")
	}
	return 200, js
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
	var NumberHws []int
	for _, val := range HwTemp {
		NumberHws = append(NumberHws, val.HomeworkId)
	}
	HwSave, erro := b.DBSaveHomework.GetSaveHomeworks(context.TODO(), courseId, userId, NumberHws, true)
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
	fmt.Println(res)
	if err != nil {
		return 404, []byte("not found")
	}
	for _, course := range res {
		fmt.Println(1222223)
		HwTemp, er := b.DBTempHomework.GetNextTempHomeworks(context.TODO(), course.CourseId)
		if er != nil {
			continue
		}
		fmt.Println("Hw.Temp", HwTemp)
		var NumberHws []int
		for _, val := range HwTemp {
			NumberHws = append(NumberHws, val.HomeworkId)
		}
		if len(NumberHws) == 0 {
			continue
		}
		HwSave, erro := b.DBSaveHomework.GetSaveHomeworks(context.TODO(), course.CourseId, userId, NumberHws, true)
		fmt.Println("Hw.Temp", HwTemp)
		fmt.Println("Hw.Save", HwSave)
		if erro != nil && erro != mongo.ErrNilDocument {
			fmt.Println(erro.Error())
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

func (b *BLogic) GetPastCourseHomeworks(userId, courseId int) (int, []byte) {
	res, err := b.DBUser.GetCourses(context.TODO(), userId)
	if err != nil {
		return 404, []byte("not found")
	}
	if !b.checkUserCourse(res, courseId) {
		return 404, []byte("not found")
	}
	HwTemp, er := b.DBTempHomework.GetPastTempHomeworks(context.TODO(), courseId)
	if er != nil {
		if er == mongo.ErrNoDocuments {
			return 200, []byte("[]")
		}
		return 404, []byte("not found")
	}

	var NumberHws []int
	for _, val := range HwTemp {
		NumberHws = append(NumberHws, val.HomeworkId)
	}
	HwSave, erro := b.DBSaveHomework.GetSaveHomeworks(context.TODO(), courseId, userId, NumberHws, false)
	if erro != nil && erro != mongo.ErrNoDocuments {
		return 500, []byte("Server error")
	}

	type returnHws struct {
		HomeworkName string    `json:"homework_name"`
		Result       int       `json:"result,omitempty"`
		MaxPoints    int       `json:"max_points,omitempty"`
		Delivered    time.Time `json:"delivered,omitempty"`
		Handed       bool      `json:"handed"`
		HomeworkId   int       `json:"homework_id,omitempty"`
		Deadline     time.Time `json:"deadline,omitempty"`
	}
	mapRes := make(map[int]returnHws)
	for _, val := range HwTemp {
		var vr returnHws
		vr.HomeworkName = val.HomeworkName
		vr.HomeworkId = val.HomeworkId
		vr.Deadline = val.Deadline
		fmt.Println("до")
		mapRes[val.HomeworkId] = vr
		fmt.Println("после")
	}
	for _, val := range HwSave {
		vr := mapRes[val.HomeworkId]
		vr.Handed = val.Handed
		if val.Handed == true {
			vr.Delivered = val.Delivered
			vr.Result = val.Result
			vr.MaxPoints = val.MaxPoints
		}
		mapRes[val.HomeworkId] = vr
	}

	var arrRet []returnHws
	for _, val := range mapRes {
		if val.Handed == true || val.Deadline.Before(time.Now()) {
			arrRet = append(arrRet, val)
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
