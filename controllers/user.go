package controllers

import (
	"acme/models"
	pk "acme/utilities/pbkdf2"
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/alexcesaro/mail/gomail"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/validation"
	"github.com/twinj/uuid"
	"time"
)

type hStruct struct {
	Hash [32]byte
	Salt [16]byte
}

func (this *MainController) Login() {
	this.activeContent("user/login")

	back := this.Ctx.Input.Param(":back")
	if this.Ctx.Input.Method() == "POST" {
		flash := beego.NewFlash()
		email := this.GetString("email")
		password := this.GetString("password")
		valid := validation.Validation{}
		valid.Email(email, "email")
		valid.Required(password, "password")
		if valid.HasErrors() {
			errormap := []string{}
			for _, err := range valid.Errors {
				errormap = append(errormap, "Validation failed on "+err.Key+": "+err.Message+"\n")
			}
			this.Data["Errors"] = errormap
			return
		}
		fmt.Println("Authorization is", email, ":", password)

		//******** Read password hash from database into temp struct y
		var y hStruct

		o := orm.NewOrm()
		o.Using("default")
		user := models.AuthUser{Email: email}
		err := o.Read(&user, "Email")
		if err == nil {
			if user.Reg_key != "" {
				flash.Error("Account not verified")
				flash.Store(&this.Controller)
				return
			}

			ibuf := bytes.NewReader([]byte(user.Password))
			err = binary.Read(ibuf, binary.LittleEndian, &y.Hash)
			if err != nil {
				flash.Error("Internal error")
				flash.Store(&this.Controller)
				return
			}
			err = binary.Read(ibuf, binary.LittleEndian, &y.Salt)
			if err != nil {
				flash.Error("Internal error")
				flash.Store(&this.Controller)
				return
			}
			fmt.Println("password hash y is", y)
		} else {
			flash.Error("No such user/email")
			flash.Store(&this.Controller)
			return
		}

		//******** Compare submitted password with database
		var x pk.PasswordHash

		x.Hash = make([]byte, 32)
		copy(x.Hash, y.Hash[:32])
		x.Salt = make([]byte, 16)
		copy(x.Salt, y.Salt[:16])
		if !pk.MatchPassword(password, &x) {
			flash.Error("Bad password")
			flash.Store(&this.Controller)
			return
		}

		//******** Create session and go back to previous page
		m := make(map[string]interface{})
		m["first"] = user.First
		m["username"] = email
		m["timestamp"] = time.Now()
		this.SetSession("acme", m)
		this.Redirect("/"+back, 302)
	}
}

func (this *MainController) Logout() {
	this.activeContent("logout")
	this.DelSession("acme")
	this.Redirect("/home", 302)
}

func (this *MainController) Register() {
	this.activeContent("user/register")

	if this.Ctx.Input.Method() == "POST" {
		flash := beego.NewFlash()
		first := this.GetString("first")
		last := this.GetString("last")
		email := this.GetString("email")
		password := this.GetString("password")
		password2 := this.GetString("password2")
		valid := validation.Validation{}
		valid.Required(first, "first")
		valid.Email(email, "email")
		valid.MinSize(password, 6, "password")
		valid.Required(password2, "password2")
		if valid.HasErrors() {
			errormap := []string{}
			for _, err := range valid.Errors {
				errormap = append(errormap, "Validation failed on "+err.Key+": "+err.Message+"\n")
			}
			this.Data["Errors"] = errormap
			return
		}
		if password != password2 {
			flash.Error("Passwords don't match")
			flash.Store(&this.Controller)
			return
		}
		h := pk.HashPassword(password)

		//******** Save user info to database
		o := orm.NewOrm()
		o.Using("default")

		user := new(models.AuthUser)
		user.First = first
		user.Last = last
		user.Email = email

		//******** Convert password hash to string
		buf := new(bytes.Buffer)
		err := binary.Write(buf, binary.LittleEndian, h.Hash)
		if err != nil {
			flash.Error("Internal error")
			flash.Store(&this.Controller)
			return
		}
		err = binary.Write(buf, binary.LittleEndian, h.Salt)
		if err != nil {
			flash.Error("Internal error")
			flash.Store(&this.Controller)
			return
		}
		b := buf.Bytes()
		fmt.Printf("password hash/salt is: %x\n", b)
		user.Password = string(b)

		//******** Add user to database with new uuid and send verification email
		u := uuid.NewV4()
		user.Reg_key = u.String()
		id, err := o.Insert(user) // BUG: the input parameter MUST be passed by value, contradicting the documentation
		if err != nil {
			flash.Error(email + " already registered")
			flash.Store(&this.Controller)
			return
		}
		fmt.Println("Id =", id)

		if !sendVerification(email, u.String()) {
			flash.Error("Unable to send verification email")
			flash.Store(&this.Controller)
			return
		}
		flash.Notice("Your account has been created. You must verify the account in your email.")
		flash.Store(&this.Controller)
		this.Redirect("/notice", 302)
	}
}

func sendVerification(email, u string) bool {
	link := "http://localhost:8080/user/verify/" + u
	host := "smtp.gmail.com"
	port := 587
	msg := gomail.NewMessage()
	msg.SetAddressHeader("From", "acmecorp@gmail.com", "ACME Corporation")
	msg.SetHeader("To", email)
	msg.SetHeader("Subject", "Account Verification for ACME Corporation")
	msg.SetBody("text/html", "To verify your account, please click on the link: <a href=\""+link+
		"\">"+link+"</a><br><br>Best Regards,<br>ACME Corporation")
	m := gomail.NewMailer(host, "youraccount@gmail.com", "YourPassword", port)
	if err := m.Send(msg); err != nil {
		return false
	}
	return true
}

func (this *MainController) Verify() {
	this.activeContent("user/verify")

	u := this.Ctx.Input.Param(":uuid")
	o := orm.NewOrm()
	o.Using("default")
	user := models.AuthUser{Reg_key: u}
	err := o.Read(&user, "Reg_key")
	if err == nil {
		this.Data["Verified"] = 1
		user.Reg_key = ""
		if _, err := o.Update(&user); err != nil {
			delete(this.Data, "Verified")
		}
	}
}

func (this *MainController) Profile() {
	this.activeContent("user/profile")

	//******** This page requires login
	sess := this.GetSession("acme")
	if sess == nil {
		this.Redirect("/login/home", 302)
		return
	}
	m := sess.(map[string]interface{})

	flash := beego.NewFlash()

	//******** Read password hash from database into temp struct y
	var y hStruct

	o := orm.NewOrm()
	o.Using("default")
	user := models.AuthUser{Email: m["username"].(string)}
	err := o.Read(&user, "Email")
	if err == nil {
		ibuf := bytes.NewReader([]byte(user.Password))
		err = binary.Read(ibuf, binary.LittleEndian, &y.Hash)
		if err != nil {
			flash.Error("Internal error")
			flash.Store(&this.Controller)
			return
		}
		err = binary.Read(ibuf, binary.LittleEndian, &y.Salt)
		if err != nil {
			flash.Error("Internal error")
			flash.Store(&this.Controller)
			return
		}
		fmt.Println("password hash y is", y)
	} else {
		flash.Error("Internal error")
		flash.Store(&this.Controller)
		return
	}

	defer func(this *MainController, user *models.AuthUser) {
		this.Data["First"] = user.First
		this.Data["Last"] = user.Last
		this.Data["Email"] = user.Email
	}(this, &user)

	if this.Ctx.Input.Method() == "POST" {
		first := this.GetString("first")
		last := this.GetString("last")
		email := this.GetString("email")
		current := this.GetString("current")
		password := this.GetString("password")
		password2 := this.GetString("password2")
		valid := validation.Validation{}
		valid.Required(first, "first")
		valid.Email(email, "email")
		valid.Required(current, "current")
		if valid.HasErrors() {
			errormap := []string{}
			for _, err := range valid.Errors {
				errormap = append(errormap, "Validation failed on "+err.Key+": "+err.Message+"\n")
			}
			this.Data["Errors"] = errormap
			return
		}

		if password != "" {
			valid.MinSize(password, 6, "password")
			valid.Required(password2, "password2")
			if valid.HasErrors() {
				errormap := []string{}
				for _, err := range valid.Errors {
					errormap = append(errormap, "Validation failed on "+err.Key+": "+err.Message+"\n")
				}
				this.Data["Errors"] = errormap
				return
			}

			if password != password2 {
				flash.Error("Passwords don't match")
				flash.Store(&this.Controller)
				return
			}
			h := pk.HashPassword(password)

			//******** Convert password hash to string
			buf := new(bytes.Buffer)
			err = binary.Write(buf, binary.LittleEndian, h.Hash)
			if err != nil {
				flash.Error("Internal error")
				flash.Store(&this.Controller)
				return
			}
			err = binary.Write(buf, binary.LittleEndian, h.Salt)
			if err != nil {
				flash.Error("Internal error")
				flash.Store(&this.Controller)
				return
			}
			b := buf.Bytes()
			fmt.Printf("password hash/salt is: %x\n", b)
			user.Password = string(b)
		}

		//******** Compare submitted password with database
		var x pk.PasswordHash

		x.Hash = make([]byte, 32)
		copy(x.Hash, y.Hash[:32])
		x.Salt = make([]byte, 16)
		copy(x.Salt, y.Salt[:16])
		if !pk.MatchPassword(current, &x) {
			flash.Error("Bad current password")
			flash.Store(&this.Controller)
			return
		}

		//******** Save user info to database
		user.First = first
		user.Last = last
		user.Email = email

		_, err := o.Update(&user)
		if err == nil {
			flash.Notice("Profile updated")
			flash.Store(&this.Controller)
			m["username"] = email
		} else {
			flash.Error("Internal error")
			flash.Store(&this.Controller)
			return
		}
	}
}

func (this *MainController) Remove() {
	this.activeContent("user/remove")

	//******** This page requires login
	sess := this.GetSession("acme")
	if sess == nil {
		this.Redirect("/login/home", 302)
		return
	}
	m := sess.(map[string]interface{})

	if this.Ctx.Input.Method() == "POST" {
		current := this.GetString("current")
		valid := validation.Validation{}
		valid.Required(current, "current")
		if valid.HasErrors() {
			errormap := []string{}
			for _, err := range valid.Errors {
				errormap = append(errormap, "Validation failed on "+err.Key+": "+err.Message+"\n")
			}
			this.Data["Errors"] = errormap
			return
		}

		flash := beego.NewFlash()

		//******** Read password hash from database into temp struct y
		var y hStruct

		o := orm.NewOrm()
		o.Using("default")
		user := models.AuthUser{Email: m["username"].(string)}
		err := o.Read(&user, "Email")
		if err == nil {
			ibuf := bytes.NewReader([]byte(user.Password))
			err = binary.Read(ibuf, binary.LittleEndian, &y.Hash)
			if err != nil {
				flash.Error("Internal error")
				flash.Store(&this.Controller)
				return
			}
			err = binary.Read(ibuf, binary.LittleEndian, &y.Salt)
			if err != nil {
				flash.Error("Internal error")
				flash.Store(&this.Controller)
				return
			}
			fmt.Println("password hash y is", y)
		} else {
			flash.Error("Internal error")
			flash.Store(&this.Controller)
			return
		}

		//******** Compare submitted password with database
		var x pk.PasswordHash

		x.Hash = make([]byte, 32)
		copy(x.Hash, y.Hash[:32])
		x.Salt = make([]byte, 16)
		copy(x.Salt, y.Salt[:16])
		if !pk.MatchPassword(current, &x) {
			flash.Error("Bad current password")
			flash.Store(&this.Controller)
			return
		}

		//******** Delete user record
		_, err = o.Delete(&user)
		if err == nil {
			flash.Notice("Your account is deleted.")
			flash.Store(&this.Controller)
			this.DelSession("acme")
			this.Redirect("/notice", 302)
		} else {
			flash.Error("Internal error")
			flash.Store(&this.Controller)
			return
		}
	}
}
