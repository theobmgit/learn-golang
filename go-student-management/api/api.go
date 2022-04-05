package main

import (
	"flag"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	studentAPI "go-student-management/api/student"
	"go-student-management/docs"
	pb "go-student-management/proto/student"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
)

var (
	clientAddr = flag.String("c-addr", "localhost:50052", "The api address in the format of host:port")
	serverAddr = flag.String("addr", "localhost:50051", "The server address in the format of host:port")
)

// @title Student Management API
// @version 1.0
// @description Sample student management RPC-based API using gRPC server and gin api.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:50052
// @BasePath /api/v1
// @schemes http
func main() {
	const BasePath = "/api/v1"
	connection, err := grpc.Dial(*serverAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Client error: %v\n", err)
	}
	defer connection.Close()

	client := pb.NewStudentServiceClient(connection)

	router := gin.Default()
	router.StaticFile("/favicon.ico", "resources/favicon.ico")

	v1 := router.Group(BasePath)
	{
		student := v1.Group("/students")
		{
			student.GET("", func(ctx *gin.Context) {
				studentAPI.GetStudents(ctx, client)
			})
			student.GET("/:id", func(ctx *gin.Context) {
				studentAPI.GetStudentById(ctx, client)
			})
			student.GET("/name", func(ctx *gin.Context) {
				studentAPI.GetStudentsByName(ctx, client)
			})

			student.POST("", func(ctx *gin.Context) {
				studentAPI.AddStudent(ctx, client)
			})

			student.PUT("/:id", func(ctx *gin.Context) {
				studentAPI.UpdateStudent(ctx, client)
			})

			student.DELETE("/:id", func(ctx *gin.Context) {
				studentAPI.DeleteStudent(ctx, client)
			})
		}
	}

	docs.SwaggerInfo.BasePath = BasePath
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	if err := router.Run(*clientAddr); err != nil {
		log.Fatalf("could not run server: %v", err)
	}
}
