package api

import (
	mainv1 "github.com/yashbek/y2j/api/main/v1"
)

type MainServer struct{
	mainv1.UnimplementedMainServiceServer
}