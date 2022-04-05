package student

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	pb "go-student-management/proto/student"
	"math"
	"net/http"
	"strings"
)

// GetStudents godoc
// @Summary      Get all students
// @Description  Get all students' information
// @Tags         Student
// @Accept       json
// @Produce      json
// @Success      200  {object}  map[string][]string
// @Router       /students [get]
func GetStudents(ctx *gin.Context, client pb.StudentServiceClient) {
	response, err := client.GetStudents(ctx, &pb.EmptyRequest{})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"result": response,
	})
}

// GetStudentById godoc
// @Summary      Get student by student id
// @Description  Get student information based on student id
// @Tags         Student
// @Accept       json
// @Produce      json
// @Param        id path string true "Student ID" minLength(8)
// @Success      200  {object}  map[string][]string
// @Failure      400  {object}  map[string][]string
// @Failure      500  {object}  map[string][]string
// @Router       /students/{id} [get]
func GetStudentById(ctx *gin.Context, client pb.StudentServiceClient) {
	id := ctx.Param("id")
	if validateId(id) {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "ID length must be 8",
		})
		return
	}

	response, err := client.GetStudentById(ctx, &pb.StudentId{Id: id})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"result": response,
	})
}

// GetStudentsByName godoc
// @Summary      Get student by name
// @Description  Get student information based on student name
// @Tags         Student
// @Accept       json
// @Produce      json
// @Param        name query string true "Student name" pattern("^[a-zA-Z]+(([',. -][a-zA-Z ])?[a-zA-Z]*)*$)
// @Success      200  {object}  map[string][]string
// @Failure      400  {object}  map[string][]string
// @Failure      500  {object}  map[string][]string
// @Router       /students/name/{name} [get]
func GetStudentsByName(ctx *gin.Context, client pb.StudentServiceClient) {
	name := ctx.Query("name")
	if validateName(name) {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid name pattern",
		})
		return
	}
	preprocessName(&name)

	response, err := client.GetStudentsByName(ctx, &pb.StudentName{Name: name})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"result": response,
	})
}

// UpdateStudent godoc
// @Summary      Update student by ID
// @Description  Update student information based on student ID
// @Tags         Student
// @Accept       json
// @Produce      json
// @Param        id path string true "Student ID to update" minLength(8)
// @Param        student body student.Student true "Student information to update"
// @Success      200  {object}  map[string][]string
// @Failure      400  {object}  map[string][]string
// @Failure      500  {object}  map[string][]string
// @Router       /students/{id} [put]
func UpdateStudent(ctx *gin.Context, client pb.StudentServiceClient) {
	var student *pb.Student
	if err := ctx.ShouldBindBodyWith(&student, binding.JSON); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	student.Id = ctx.Param("id")

	if validateId(student.Id) {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "ID length must be 8",
		})
		return
	} else if validateName(student.Name) {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid name pattern",
		})
		return
	} else if validateAge(student.Age) {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Age must be an integer larger than 0",
		})
		return
	} else if validateClass(student.Class) {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid class",
		})
		return
	} else if validateCpa(student.Cpa) {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "CPA must be between 0.0 and 4.0",
		})
		return
	}
	preprocessName(&student.Name)
	preprocessCpa(&student.Cpa)
	preprocessClass(&student.Class)

	response, err := client.UpdateStudent(ctx, student)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"result": response,
	})
}

// DeleteStudent godoc
// @Summary      Delete student by ID
// @Description  Delete student record based on student ID
// @Tags         Student
// @Accept       json
// @Produce      json
// @Param        id path string true "Student ID to delete"
// @Success      200  {object}  map[string][]string
// @Failure      400  {object}  map[string][]string
// @Failure      500  {object}  map[string][]string
// @Router       /students/{id} [delete]
func DeleteStudent(ctx *gin.Context, client pb.StudentServiceClient) {
	id := ctx.Param("id")

	if validateId(id) {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "ID length must be 8",
		})
		return
	}

	response, err := client.DeleteStudent(ctx, &pb.StudentId{Id: id})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"result": response,
	})
}

// AddStudent godoc
// @Summary      Add new student
// @Description  Get student information based on student id
// @Tags         Student
// @Accept       json
// @Produce      json
// @Param        student body student.Student true "Student information"
// @Success      200  {object}  map[string][]string
// @Failure      400  {object}  map[string][]string
// @Failure      500  {object}  map[string][]string
// @Router       /students [post]
func AddStudent(ctx *gin.Context, client pb.StudentServiceClient) {
	var student *pb.Student
	if err := ctx.ShouldBindBodyWith(&student, binding.JSON); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
	}

	if validateId(student.Id) {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "ID length must be 8",
		})
		return
	} else if validateName(student.Name) {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid name pattern",
		})
		return
	} else if validateAge(student.Age) {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Age must be an integer larger than 0",
		})
		return
	} else if validateClass(student.Class) {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid class",
		})
		return
	} else if validateCpa(student.Cpa) {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "CPA must be between 0.0 and 4.0",
		})
		return
	}
	preprocessName(&student.Name)
	preprocessCpa(&student.Cpa)
	preprocessClass(&student.Class)

	response, err := client.AddStudent(ctx, student)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"result": response,
	})
}

func validateId(id string) bool {
	return len(id) != 8
}

func validateName(name string) bool {
	return false
}

func validateAge(age int32) bool {
	return age < 0
}

func validateClass(class string) bool {
	return false
}

func validateCpa(Cpa float32) bool {
	return Cpa < 0.0 || Cpa > 4.0
}

func preprocessName(name *string) {
	*name = strings.Title(strings.ToLower(*name))
}

func preprocessClass(class *string) {
	*class = strings.ToUpper(*class)
}

func preprocessCpa(Cpa *float32) {
	*Cpa = float32(math.Ceil(float64(*Cpa*100)) / 100)
}
