package blogic

import (
	"OnlineSchool/internal/structs"
	mongodb "OnlineSchool/pkg/mongoDB"
	"context"
	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"strconv"
	"time"
)

const (
	clientId = 8219136
	secretId = "fobH7n71sa1Hhhl771Ek"
)

func (b *BLogic) Login(VKCode string, redirectUrl string) (int, []byte, string /*cookie*/) {
	fmt.Println(VKCode)
	res, err := http.Get(fmt.Sprintf("https://oauth.vk.com/access_token?code=%v&redirect_uri=%v&client_id=%v&client_secret=%v", VKCode, redirectUrl, clientId, secretId))
	if err != nil {
		log.Fatal(err)
	}
	byteData, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}
	var data struct {
		AccessToken      string `json:"access_token,omitempty"`
		ExpiresIn        int64  `json:"expires_in,omitempty"`
		UserId           int64  `json:"user_id,omitempty"`
		Error            string `json:"error,omitempty"`
		ErrorDescription string `json:"error_description,omitempty"`
	}
	er := json.Unmarshal(byteData, &data)
	if er != nil {

	}
	if data.Error != "" {
		fmt.Println("error:", data.Error, data.ErrorDescription)
		return 400, []byte("VK code is not valid"), ""
	}
	user, erro := b.DBUser.GetUser(context.TODO(), data.UserId)
	if erro != nil {
		if erro != mongo.ErrNoDocuments {
			return 500, []byte("Server error"), ""
		}
		fmt.Println(data.UserId)
		fmt.Println(data.AccessToken)
		id, errr := b.createUser(data.UserId, data.AccessToken)
		if errr != nil {
			return 500, []byte("Server error"), ""
		}
		token, eo := b.JWTManager.NewJWT(string(id), time.Hour*24*30)
		if eo != nil {
			return 500, []byte("Server error"), token
		}
		return 200, []byte(string(id)), token
	}
	token, eo := b.JWTManager.NewJWT(strconv.FormatInt(user.UserId, 10), time.Hour*24*30)
	if eo != nil {
		return 500, []byte("Server error"), token
	}
	return 200, []byte("OK"), token
}

func (b *BLogic) createUser(VkUserId int64, accessToken string) (userId int64, err error) {
	var user structs.User
	user.VkId = VkUserId
	res, err := http.Get(fmt.Sprintf("https://api.vk.com/method/users.get?user_ids=%v&fields=photo_100&access_token=%v&v=5.131", VkUserId, accessToken))
	if err != nil {
		fmt.Println("err func (create user) [vk request]: ", err.Error())
		return 0, err
	}
	byteData, erro := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if erro != nil {
		fmt.Println("err func (create user) [read vk body]: ", erro.Error())
		return 0, erro
	}
	var us struct {
		Response []struct {
			FirstName string `json:"first_name"`
			LastName  string `json:"last_name"`
			Avatar    string `json:"photo_100"`
		}
		Error interface{} `json:"error"`
	}
	er := json.Unmarshal(byteData, &us)
	if er != nil {
		fmt.Println("err func (create user) [JSON Unmarshal]: ", er.Error())
		return 0, er
	}
	if us.Error != nil {
		return 0, fmt.Errorf("func err")
	}
	user.Avatar = us.Response[0].Avatar
	user.FirstName = us.Response[0].FirstName
	user.LastName = us.Response[0].LastName
	user.UserId = VkUserId
	user.BuyCourses = append(user.BuyCourses, structs.UserCourse{CourseId: 0, Active: false})
	fmt.Println(us)
	errorr := b.DBUser.CreateUser(context.TODO(), user)
	if errorr != nil {
		if mongodb.IsDuplicate(errorr) {
			uId := user.UserId
			var i int64 = 0
			for mongodb.IsDuplicate(errorr) {
				var numerator int64 = int64(len(string(i)))
				user.UserId = uId*int64(math.Pow(10, float64(numerator))) + i
				errorr = b.DBUser.CreateUser(context.TODO(), user)
				i++
			}
		} else {
			fmt.Println("err func (create user) [DB CreateUser]", errorr.Error())
		}
	}
	return user.UserId, nil
}

func (b *BLogic) Authentication(token string) (int64, int, error) {
	strId, er := b.JWTManager.Parse(token)
	if er != nil {
		fmt.Println("Authentication [JWT Parse] error: ", er.Error())
		return 0, 401, er
	}
	fmt.Println("strId", strId)
	userId, erro := strconv.ParseInt(strId, 10, 64)
	if erro != nil {
		return 0, 401, erro
	}
	_, err := b.DBUser.GetCourses(context.TODO(), userId)
	if err != nil {
		return 0, 401, err
	}
	return userId, 200, nil
}
