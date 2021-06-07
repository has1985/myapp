package svc

import (
	"context"
	"github.com/infobloxopen/atlas-app-toolkit/gateway"
	gorm2 "github.com/infobloxopen/atlas-app-toolkit/gorm"
	query1 "github.com/infobloxopen/atlas-app-toolkit/query"
	errors1 "github.com/infobloxopen/protoc-gen-gorm/errors"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/has1985/myapp/myapp/pkg/pb"
	"github.com/jinzhu/gorm"
)

const (
	// version is the current version of the service
	version = "0.0.1"
)

// Default implementation of the Myapp server interface
type server struct {
	*pb.MyappDefaultServer
	db *gorm.DB
}

// GetVersion returns the current version of the service
func (server) GetVersion(context.Context, *empty.Empty) (*pb.VersionResponse, error) {
	return &pb.VersionResponse{Version: version}, nil
}

// NewBasicServer returns an instance of the default server interface
func NewBasicServer(database *gorm.DB) (pb.MyappServer, error) {
	return &server{db: database}, nil
}

// Create ...
func (m *server) Create(ctx context.Context, in *pb.CreateBookRequest) (*pb.CreateBookResponse, error) {
	txn, ok := gorm2.FromContext(ctx)
	if !ok {
		return nil, errors1.NoTransactionError
	}
	db := txn.Begin()
	if db.Error != nil {
		return nil, db.Error
	}

	res, err := pb.DefaultCreateBook(ctx, in.GetPayload(), db)
	if err != nil {
		return nil, err
	}
	out := &pb.CreateBookResponse{Id: res.Id}

	return out, nil
}

// Delete ...
func (m *server) Delete(ctx context.Context, in *pb.DeleteBookRequest) (*pb.DeleteBookResponse, error) {
	txn, ok := gorm2.FromContext(ctx)
	if !ok {
		return nil, errors1.NoTransactionError
	}
	db := txn.Begin()
	if db.Error != nil {
		return nil, db.Error
	}

	err := pb.DefaultDeleteBook(ctx, &pb.Book{Id: in.GetId()}, db)
	if err != nil {
		return nil, err
	}
	out := &pb.DeleteBookResponse{}
	gateway.SetDeleted(ctx, "Deleted")
	return out, nil
}

// List ...
func (m *server) List(ctx context.Context, in *pb.ListBookRequest) (*pb.ListBookResponse, error) {
	txn, ok := gorm2.FromContext(ctx)
	if !ok {
		return nil, errors1.NoTransactionError
	}
	db := txn.Begin()
	if db.Error != nil {
		return nil, db.Error
	}
	db = db.Where("amount > ?", 0)

	pagedRequest := false
	if in.GetPaging().GetLimit() >= 1 {
		in.Paging.Limit++
		pagedRequest = true
	}
	res, err := pb.DefaultListBook(ctx, db, in.Filter, in.OrderBy, in.Paging, in.Fields)
	if err != nil {
		return nil, err
	}
	var resPaging *query1.PageInfo
	if pagedRequest {
		var offset int32
		var size int32 = int32(len(res))
		if size == in.GetPaging().GetLimit() {
			size--
			res = res[:size]
			offset = in.GetPaging().GetOffset() + size
		}
		resPaging = &query1.PageInfo{Offset: offset}
	}
	out := &pb.ListBookResponse{Results: res, Page: resPaging}

	return out, nil
}

// Read ...
func (m *server) Read(ctx context.Context, in *pb.ReadBookRequest) (*pb.ReadBookResponse, error) {
	txn, ok := gorm2.FromContext(ctx)
	if !ok {
		return nil, errors1.NoTransactionError
	}
	db := txn.Begin()
	if db.Error != nil {
		return nil, db.Error
	}

	db = db.Where("amount > ?", 0)

	res, err := pb.DefaultReadBook(ctx, &pb.Book{Id: in.GetId()}, db)
	if err != nil {
		return nil, err
	}
	out := &pb.ReadBookResponse{Result: res}
	return out, nil
}
