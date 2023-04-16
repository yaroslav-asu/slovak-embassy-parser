package user

import (
	"go.uber.org/zap"
	"main/internal/session"
	"main/internal/utils/funcs"
	gorm_models "main/models/gorm"
	"main/models/gorm/datetime"
	"net/url"
	"strings"
)

type User struct {
	Session session.Session
	DB      gorm_models.User
}

func (u *User) LogIn() {
	u.Session.LogIn(u.DB.UserName, u.DB.Password)
}

func (u *User) LogOut() {
	zap.L().Info("Starting logout user: " + u.DB.UserName)
	u.Session.LogOut()
	zap.L().Info("Finished logout user: " + u.DB.UserName)
}

func NewUser(username, password string) User {
	newUser := User{
		Session: session.NewLoggedInSession(username, password),
		DB: gorm_models.User{
			UserName: username,
			Password: password,
		},
	}
	return newUser
}

func (u *User) ReserveDatetime(city gorm_models.City, date datetime.Date) {
	zap.L().Info("Starting to reserve date in: " + city.Name + " at: " + date.Format(datetime.DateTime))
	res := u.Session.Get(funcs.Linkify("calendarDay.do?day=", date.Format(datetime.DateOnly), "&timeSlotId=&calendarId=&consularPostId=", city.Id))
	defer res.Body.Close()
	funcs.Sleep()
	u.Session.DownloadCaptcha()
	captchaSolve := session.SolveCaptcha()
	res = u.Session.PostForm(
		funcs.Linkify("calendarDay.do?day=", date.Format(datetime.DateOnly), "&consularPostId=", city.Id),
		url.Values{
			"calendar.timeOfVisit":               {date.Format(datetime.FormDateTime)},
			"calendar.sequenceNo":                {"1"},
			"calendar.consularPost.consularPost": {city.Id},
			"captcha":                            {captchaSolve},
		},
	)
	defer res.Body.Close()
	res = u.Session.Get(funcs.Linkify("logout.do"))
	defer res.Body.Close()
}

func (u *User) IsReserved() bool {
	doc := u.Session.GetParsedSoup(funcs.Linkify("dateOfVisitDecision.do"))
	return strings.Contains(doc.Find("td", "class", "infoTableInformationText").Text(), "have reservation")
}