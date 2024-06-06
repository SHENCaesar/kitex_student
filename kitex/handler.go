package main

import (
	"context"
	"strings"
	"sync"

	"github.com/SHENCaesar/kitex_student/kitex/dal/mysql"
	demo "github.com/SHENCaesar/kitex_student/kitex/kitex_gen/demo"
	"github.com/SHENCaesar/kitex_student/kitex/model"
	"github.com/cloudwego/kitex/tool/internal_pkg/log"
)

var students sync.Map
// StudentServiceImpl implements the last service interface defined in the IDL.
type StudentServiceImpl struct{}

// Register implements the StudentServiceImpl interface.
func (s *StudentServiceImpl) Register(ctx context.Context, student *demo.Student) (resp *demo.RegisterResp, err error) {
	log.Infof("Register req: %v", student)
	req := &demo.QueryReq{Id: student.Id}
	stu, _ := s.Query(ctx, req)
	if stu == nil {
		students.Store(student.Id, student)
		emails := ""
		if len(student.GetEmail()) > 0 {
			emails = strings.Join(student.Email, ",") 
		}
		user := &model.User{
			Id:				student.GetId(),
			Name: 			student.GetName(),
			CollegeName: 	student.GetCollege().Name,
			CollegeAddress: student.GetCollege().Address,
			Emails: 		emails,
		}
		//数据库
		err := mysql.CreateUser(user)
		if err != nil {
			return nil,err
		}
	}
	resp = &demo.RegisterResp{
		Success: true,
		Message: "",
	}
	return
}

// Query implements the StudentServiceImpl interface.
func (s *StudentServiceImpl) Query(ctx context.Context, req *demo.QueryReq) (resp *demo.Student, err error) {
	log.Infof("query req: %v", req)
	if stu, ok := students.Load(req.Id); ok {
		resp = stu.(*demo.Student)
		return
	}
	//数据库
	user, err := mysql.QueryUser(req.Id)
	if err != nil {
		return nil, err
	}
	resp = &demo.Student{
		Id: user.Id,
		Name: user.Name,
		College: &demo.College{
			Name: user.CollegeName,
			Address: user.CollegeAddress,
		},
		Email: strings.Split(user.Emails, ","),
	}

	return
}
