package v1

import (
	"context"
	"database/sql"
	vs "../../api/proto/v1"

)

const (
	apiVersion = "v1"
)

type toDoServicesServer struct {
	db *sql.DB
}

func NewToDoServiceServer(db * sql.DB) v1.ToDoServiceServer{
	return &toDoServicesServer{
		db: db,
	}
}

func (s *toDoServicesServer) checkAPI(api string)error{
	if len(api) > 0{
		if apiVersion != api{
			return status.Error(codes.Uniplemented,"unsupported API version:service implements API '%s',but given '%s'",api,apiVersion)
		}
	}
	return nil
}

func(s *toDoServicesServer) connect(ctx context.Context)(*sql.Conn,error){
	c,err := s.db.Conn(ctx)
	if nil != err{
		return nil ,status.Error(codes.Unkonwn,"连接数据库失败"+err.Error())
	}

	return c,nil
}

//func (s *toDoServicesServer) Creat(ctx context.Context,req *v1.Creat)