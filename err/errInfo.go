package err

import ()

type Exception struct {
	Code         uint
	Message      string
	ExtraMessage string
}

func newException(code uint, message string) func(string) Exception {
	e := Exception{Code: code, Message: message}

	f := func(extra string) Exception {
		e.ExtraMessage = extra
		return e
	}
	return f
}

var (
	NoExcetion=newException(0,"well")("")
	//souce not found 1XX

	//not Found In DataBase
	NotFoundInDataBase = newException(101, "Not Found Item In DataBase")
	//Not Found User
	NotFoundUser=newException(102,"Not Found User In User Set")
	//Todo Not Found
	NotFoundToDo=newException(103,"Target Todo Not Found")
	//API Port Not Found
	NotFoundAPI=newException(104,"Target API Port Not Exist")

	
	//Bad Request 2XX
	TargetParmsNotExist=newException(201,"target Parms Not provide")
	
	//authentication 3XX
	AuthenticationFailure=newException(301 ,"Failure authentication User")
	PermissionDenied=newException(302,"Permission Denied")
	AccessDenied=newException(303,"Access Denied")
	TokenInvaild =newException(304,"Token Invail")
	TokenFailure=newException(305,"Handling Token Failure")
	GenerateTokenFailure =newException(306,"Generate Toke Failure")

	//User Error
	CreateNewUserFailure=newException(401,"Creat New User Failure")
	
	//IOC Error
	UnSupportData=newException(501,"given Data Not Match")
	FailureGenerateFunctionParm=newException(502,"Failure generate Func Parm List")
)
