package blogic

import (
	"OnlineSchool/internal/structs"
	"context"
	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"io/ioutil"
	"net/http"
)

func (b *BLogic) addUserCourse(userId int64, addCourse []structs.PayCourseType) (bool, error) {
	userCourse, err := b.DBUser.GetCourses(context.TODO(), userId)
	if err != nil {
		return false, err
	}
	var pushCourse []structs.UserCourse
	for _, val := range addCourse {
		searchCourse := true
		for _, v := range userCourse {
			if v.CourseId == val.CourseId && v.Active {
				searchCourse = false
				v.BuyPeriod = append(v.BuyPeriod, val.Periods...)
			}
		}
		if searchCourse {
			var vr structs.UserCourse
			vr.Active = true
			vr.BuyPeriod = val.Periods
			vr.FreezingDay = 31
			vr.Freeze = true
			vr.CourseId = val.CourseId
			vr.HeartCount = 5
			pushCourse = append(pushCourse, vr)
		}
	}
	userCourse = append(userCourse, pushCourse...)

	modifCound, er := b.DBUser.EditUserCourses(context.TODO(), userId, userCourse)
	if er != nil {
		return false, er
	}
	if modifCound == 0 {
		return false, fmt.Errorf("modifCound is zero")
	}
	return true, nil
}

func (b *BLogic) CheckConnectingCourseGroups(userID int64) (int, []byte) {
	user, erro := b.DBUser.GetUser(context.TODO(), userID)
	if erro != nil {
		if erro == mongo.ErrNoDocuments {
			return 400, []byte("user not found")
		}
		return 500, []byte("Server error")
	}

	type CourseConnect struct {
		NameCourse string `json:"name_course"`
		CourseId   int    `json:"course_id"`
	}
	var retData []CourseConnect
	for _, val := range user.BuyCourses {
		fmt.Println("(for check connect): course =", val.CourseId)
		if !val.Active {
			continue
		}
		course, er := b.DBCourse.GetCourse(context.TODO(), val.CourseId)
		if er != nil {
			continue
		}
		necessityAddingUser := b.checkNecessityAddingUserInGroup(user.VkId, course.VkGroupId, course.VkSecretKey)
		fmt.Println("(for check connect): necessityAddingUser =", necessityAddingUser)
		if necessityAddingUser {
			var vr CourseConnect
			vr.NameCourse = course.NameCourse
			vr.CourseId = course.CourseId
			retData = append(retData, vr)

		}
	}

	ret, errorr := json.Marshal(&retData)
	if errorr != nil {
		return 500, []byte("server error marshal")
	}
	fmt.Println("ret:", string(ret), retData)
	return 200, ret
}

func (b *BLogic) checkNecessityAddingUserInGroup(vkIdUser int64, groupId string, accessToken string) bool {
	url := fmt.Sprintf("https://api.vk.com/method/groups.isMember?group_id=%v&user_id=%v&access_token=%v&&v=5.131", groupId, vkIdUser, accessToken)
	res, errReq := http.Get(url)
	if errReq != nil {
		fmt.Println("[checkUserInGroup] (get to VK):", errReq.Error())
		return false
	}
	byteData, e := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if e != nil {
		fmt.Println("[checkUserInGroup] (reading request data):", errReq.Error())
		return false
	}
	fmt.Println("[checkUserInGroup] (byteData):", string(byteData))
	var resReq struct {
		Response int `json:"response"`
	}
	eor := json.Unmarshal(byteData, &resReq)
	if eor != nil {
		fmt.Println("[checkUserInGroup] (json Unmarshal):", eor.Error())
		return false
	}
	if resReq.Response == 0 {
		return true
	}
	return false
}

func (b *BLogic) GetInvitationLinkVkGroup(userId int64, courseId int) (int, []byte) {
	user, err := b.DBUser.GetUser(context.TODO(), userId)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return 400, []byte("Bad Request")
		}
		return 500, []byte("Server error")
	}
	course, er := b.DBCourse.GetCourse(context.TODO(), courseId)
	if er != nil {
		if er == mongo.ErrNoDocuments {
			return 400, []byte("Bad Request")
		}
		return 500, []byte("Server error")
	}
	necessity := b.checkNecessityAddingUserInGroup(user.VkId, course.VkGroupId, course.VkSecretKey)
	if !necessity {
		return 400, []byte("Bad Request")
	}

	var ret struct {
		VkGroupLink string `json:"vk_group_link"`
	}
	ret.VkGroupLink = course.VkGroupLink
	data, e := json.Marshal(&ret)
	if e != nil {
		return 500, []byte("Server error")
	}
	return 200, data
}
