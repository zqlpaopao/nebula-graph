package v1

import (
	"context"
	"database/sql"
	"fmt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	v1 "test/grpc/api/proto"
)

const (
	apiVersion = "v1"
)

type toDoServiceServer struct {
	db *sql.DB
}

func NewToDoServiceServer(db *sql.DB) v1.ToDoServiceServer {
	return &toDoServiceServer{
		db: db,
	}
}

func (s *toDoServiceServer) checkAPI(api string) error {
	if len(api) > 0 {
		if apiVersion != api {
			return status.Error(codes.Unimplemented, "unsupported API version:service implements API '%s',but given '%s'")
		}
	}
	return nil
}

func (s *toDoServiceServer) connect(ctx context.Context) (*sql.Conn, error) {
	c, err := s.db.Conn(ctx)
	if nil != err {
		return nil, status.Error(codes.Unknown, "连接数据库失败"+err.Error())
	}

	return c, nil
}

func (s *toDoServiceServer) Create(ctx context.Context, req *v1.CreateRequest) (resp *v1.CreateResponse, err error) {
	if err := s.checkAPI(req.Api); nil != err {
		return nil, err
	}
	c, err := s.connect(ctx)
	if err != nil {
		return nil, err
	}
	defer c.Close()
	//时间检查
	//reminder ,err := pytes.timesstamp(req.ToDo.)
	res, err := c.ExecContext(ctx, "INSERT INTO  ToDo(`Title`,`Description`) VALUES(?,?)", req.ToDo.Title, req.ToDo.Description)
	if err != nil {
		return nil, status.Error(codes.Unknown, "插入数据失败")
	}
	id, err := res.LastInsertId()
	if nil != err {
		return nil, status.Error(codes.Unknown, "获取最新IO失败")
	}
	resp = &v1.CreateResponse{
		Api: apiVersion,
		Id:  id,
	}
	return
}

func (s *toDoServiceServer) Read(ctx context.Context, req *v1.ReadRequest) (resp *v1.ReadResponse, err error) {
	if err = s.checkAPI(req.Api); nil != err {
		return nil, status.Error(codes.Unauthenticated, "apiCheck is fail")
	}
	c, err := s.connect(ctx)
	if nil != err {
		return nil, status.Error(codes.Unknown, "获取mysqlClient失败")
	}
	defer c.Close()
	rows, err := c.QueryContext(ctx, "select * from ToDO where `ID`=?", req.Id)
	defer rows.Close()
	if nil != err {
		return nil, status.Error(codes.Unknown, "查询失败")
	}

	if !rows.Next() {
		if err := rows.Err(); nil != err {
			return nil, status.Error(codes.Unknown, "获取数据失败"+err.Error())
		}
		return nil, status.Error(codes.NotFound, fmt.Sprint("找不到数据%d", req.Id))
	}

	var td v1.ToDo
	if err := rows.Scan(&td.Id, &td.Description); nil != err {
		return nil, status.Error(codes.Unknown, "查找数据失败"+err.Error())
	}
	return &v1.ReadResponse{Api: apiVersion, ToDo: &v1.ToDo{Id: td.Id, Description: td.Description}}, nil
}
func (s *toDoServiceServer) Update(ctx context.Context, req *v1.UpdateRequest) (resp *v1.UpdateResponse, err error) {

	return
}
func (s *toDoServiceServer) Delete(ctx context.Context, req *v1.DeleteRequest) (resp *v1.DeleteResponse, err error) {
	return
}
func (s *toDoServiceServer) ReadAll(ctx context.Context, req *v1.ReadAllRequest) (resp *v1.ReadAllResponse, err error) {
	return
}
