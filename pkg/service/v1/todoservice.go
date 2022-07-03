package v1

import (
    "context"
    "database/sql"

    "github.com/golang/protobuf/ptypes"
    "google.golang.org/grpc/codes"
    "google.golang.org/grpc/status"

    v1 "github.com/dansusman/todoservice/pkg/api/v1"
)

const (
    apiVersion = "v1"
)

type todoServiceServer struct {
    db *sql.DB
}

// check if the API version requested by client is supported by server!!
func (s *todoServiceServer) checkAPI(api string) error {
    if len(api) > 0 && apiVersion != api {
        return status.Error(codes.Unimplemented, "unsupported API version")
    }
    return nil // empty string is fine, just use current version of API
}

// connect to DB and return connection
func (s *todoServiceServer) connect(ctx context.Context) (*sql.Conn, error) {
    c, err := s.db.Conn(ctx)
    if err != nil {
        return nil, status.Error(codes.Unknown, "failed to connect to DB" + err.Error())
    }
    return c, nil
}

//////////////////////////// CRUD ////////////////////////////

// Create new task
func (s *todoServiceServer) Create(ctx context.Context, req *v1.CreateRequest) (*v1.CreateResponse, error) {

    // verify API version
    if err := s.checkAPI(req.Api); err != nil {
        return nil, err
    }

    // get connection to DB
    conn, err := s.connect(ctx)
    if err != nil {
        return nil, err
    }

    defer conn.Close()

    reminder, err := ptypes.Timestamp(req.Todo.Reminder)
    if err != nil {
        return nil, status.Error(codes.InvalidArgument, "reminder field invalid: " + err.Error())
    }

    // insert todo data into DB
	res, err := conn.ExecContext(ctx, "INSERT INTO Todo(`Title`, `Description`, `Reminder`) VALUES(?, ?, ?)",
		req.Todo.Title, req.Todo.Description, reminder)
	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to insert into Todo-> " + err.Error())
	}

    // get next id in sequence of creations
    id, err := res.LastInsertId()
    if err != nil {
        return nil, status.Error(codes.Unknown, "error grabbing id for creation" + err.Error())
    }

    return &v1.CreateResponse{
        Api: apiVersion,
        Id: id,
    }, nil
    

}

