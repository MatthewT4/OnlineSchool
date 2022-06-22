package blogic

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"
)

func (b *BLogic) GetUserCourses(user_id int) (int, string) {
	res, err := b.DBUser.GetCourses(context.TODO(), user_id)
	if err != nil {
		fmt.Println(err.Error())
		return 404, "not found"
	}
	type resCourses struct {
		NameCourse string    `json:"name_course"`
		PaymentEnd time.Time `json:"payment_end"`
	}

	var mas []resCourses
	for i := 0; i < len(res); i++ {
		course, er := b.DBCourse.GetCourse(context.TODO(), res[i].CourseId)
		if er == nil {
			//find max payment period
			max := 0
			for j := 0; j < len(res[i].BuyPeriod); j++ {
				if max < res[i].BuyPeriod[j] {
					max = res[i].BuyPeriod[j]
				}
			}
			var c resCourses
			c.NameCourse = course.NameCourse
			c.PaymentEnd = course.PaymentPeriod[max]
			fmt.Println(c)
			mas = append(mas, c)
		} else {
			return 404, "not found"
			fmt.Println(err.Error())
		}
	}
	if len(mas) == 0 {
		return 404, "not found"
	}
	jr, err := json.Marshal(mas)
	if err != nil {
		return 404, "not found"
		log.Fatal(err)
	}
	return 200, string(jr)
}
