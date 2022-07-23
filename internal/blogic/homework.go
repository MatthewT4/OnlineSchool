package blogic

import (
	"OnlineSchool/internal/structs"
	mongodb "OnlineSchool/pkg/mongoDB"
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
	UserAnswer  string   `json:"user_answer"`
	Solution    string   `json:"solution,omitempty"`
	Written     bool     `json:"written"`
	TypeAnswers []string `json:"type_answers,omitempty"`
	Point       int      `json:"point"`
	MaxPoint    int      `json:"max_point"`
}

func (b *BLogic) GetHomework(userId int64, homeworkId int) (int, []byte) {
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
	res, e := b.DBUser.GetCourses(context.TODO(), userId)
	if e != nil {
		return 404, []byte("not found")
	}
	if !b.checkUserCourse(res, tempRes.CourseId) {
		return 404, []byte("not found")
	}

	save, err := b.DBSaveHomework.GetHomework(context.TODO(), userId, homeworkId)
	if err != nil && err != mongo.ErrNoDocuments {
		fmt.Println(err.Error())
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
		vr.UserAnswer = saveMap[value.TaskId].UserAnswer
		if handed {
			vr.Point = saveMap[value.TaskId].Point
			vr.Answers = task.Answers
			vr.Solution = task.Solution
		}
		returnTasks = append(returnTasks, vr)
	}
	var ret struct {
		Tasks     []retTask `json:"tasks"`
		Name      string    `json:"name"`
		Handed    bool      `json:"handed"`
		UserPoint int       `json:"user_point,omitempty"`
		MaxPoint  int       `json:"max_point,omitempty"`
	}
	ret.Name = tempRes.HomeworkName
	ret.Tasks = returnTasks
	ret.Handed = handed
	if handed {
		ret.UserPoint = save.Result
		ret.MaxPoint = save.MaxPoints
	}
	js, er := json.Marshal(&ret)
	if er != nil {
		fmt.Println("marshal error")
		return 500, []byte("Server error")
	}
	return 200, js
}

func (b *BLogic) GetNextCourseHomeworks(userId int64, courseId int) (int, []byte) {
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
	HwSave, erro := b.DBSaveHomework.GetSaveHomeworks(context.TODO(), userId, NumberHws, true)
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

func (b *BLogic) GetNextHomeworks(userId int64) (int, []byte) {
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
		HwSave, erro := b.DBSaveHomework.GetSaveHomeworks(context.TODO(), userId, NumberHws, true)
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

func (b *BLogic) GetPastCourseHomeworks(userId int64, courseId int) (int, []byte) {
	fmt.Println(userId, courseId)
	res, err := b.DBUser.GetCourses(context.TODO(), userId)
	if err != nil {
		return 404, []byte("not found")
	}
	if !b.checkUserCourse(res, courseId) {
		return 404, []byte("not found")
	}
	HwTemp, er := b.DBTempHomework.GetPastTempHomeworks(context.TODO(), courseId)
	fmt.Println("past temp", HwTemp)
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
	HwSave, erro := b.DBSaveHomework.GetSaveHomeworks(context.TODO(), userId, NumberHws, false)
	fmt.Println("past save", HwSave)
	if erro != nil && erro != mongo.ErrNoDocuments {
		return 500, []byte("Server error")
	}

	type returnHws struct {
		HomeworkName string    `json:"homework_name"`
		Result       int       `json:"result,omitempty"`
		MaxPoints    int       `json:"max_points,omitempty"`
		Delivered    time.Time `json:"delivered,omitempty"`
		Handed       bool      `json:"handed"`
		HomeworkId   int       `json:"homework_id"`
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

func (b *BLogic) CreateEmptySaveHw(ownerId int64, temp structs.HomeworkTemplate) error {
	var tasks []structs.HomeworkTask
	maxPoint := 0
	for _, val := range temp.Tasks {
		var vr structs.HomeworkTask
		vr.Number = val.Number
		vr.MaxPoint = val.MaxPoint
		vr.UserAnswer = ""
		tasks = append(tasks, vr)
		maxPoint += val.MaxPoint
	}
	var saveHw structs.HomeworkSave
	saveHw.Handed = false
	saveHw.MaxPoints = maxPoint
	saveHw.HomeworkId = temp.HomeworkId
	saveHw.OwnerId = ownerId
	saveHw.Tasks = tasks
	err := b.DBSaveHomework.CreateSaveHw(context.TODO(), saveHw)
	return err
}
func (b *BLogic) SubmitHomework(userId int64, homeworkId int, answers []structs.HomeworkTask) (int, string) {
	//обработка крайних случаев
	temp, er := b.DBTempHomework.GetHomework(context.TODO(), homeworkId)
	if er != nil {
		if er == mongo.ErrNoDocuments {
			return 404, "not found"
		}
		return 500, "server error"
	}

	res, err := b.DBUser.GetCourses(context.TODO(), userId)
	if err != nil {
		return 400, "This homework not private user"
	}
	if !b.checkUserCourse(res, temp.CourseId) {
		return 400, "This homework not private user"
	}

	save, err := b.DBSaveHomework.GetHomework(context.TODO(), userId, homeworkId)
	if err == nil {
		if save.Handed {
			return 400, "hw submit already"
		}
	} else {
		if err != mongo.ErrNoDocuments {
			return 500, "server error"
		}
		err = b.CreateEmptySaveHw(userId, temp)
		if err != nil && !mongodb.IsDuplicate(err) {
			return 500, "server error"
		}
	}

	//Сам обработчик
	var tasksId []int
	mapObrabot := make(map[int](structs.HomeworkTask)) //key = Number

	//Наполяем задачами из tempHw
	for _, val := range temp.Tasks {
		tasksId = append(tasksId, val.TaskId)
		mapObrabot[val.Number] = val
	}
	//вставляем ответы из запроса
	for _, val := range answers {
		_, ok := mapObrabot[val.Number]
		if !ok {
			return 400, "number response not found template homework"
		}
		vr := mapObrabot[val.Number]
		vr.UserAnswer = val.UserAnswer
		mapObrabot[val.Number] = vr
	}

	//получаем словарь самих задач, key = taskId
	tasks, ok := b.getTasks(tasksId, true)
	if !ok {
		return 500, "Server error"
	}

	// Проверяем ответы и ставим баллы за них
	overallResult := 0
	for key, val := range mapObrabot {
		vr := val
		vr.MaxPoint = tasks[val.TaskId].MaxPoint
		if tasks[val.TaskId].Handler == "" {
			vr.Point = orderly(vr.UserAnswer, tasks[val.TaskId].Answers, tasks[val.TaskId].MaxPoint)
			overallResult += vr.Point
		}
		mapObrabot[key] = vr
		fmt.Println("vr:", mapObrabot[key])
	}

	// Превращаем словарь в массив
	var masTasks []structs.HomeworkTask
	for _, val := range mapObrabot {
		masTasks = append(masTasks, val)
	}
	modifCount, e := b.DBSaveHomework.UpdateTasks(context.TODO(), homeworkId, userId, masTasks, overallResult)
	if e != nil {
		return 500, "update service error"
	}
	if modifCount == 0 {
		return 400, "request incorrect"
	}
	return 200, "ok"

}
