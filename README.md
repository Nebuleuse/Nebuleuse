# Nebuleuse
Nebuleuse backend written in Go using a REST like API  

# API
GET /  
GET /status  
Input  
Response : JSON {
	Maintenance      bool  
	NebuleuseVersion int  
	GameVersion      int  
	UpdaterVersion   int  
	Motd             string  
}  

POST /connect  
Input : username, password  
Response : JSON {
	SessionId string  
}  

POST /getUserInfos  
Input : sessionid  
Response : JSON {  
	Username     string  
	SessionId    string  
	Rank         int  
	Avatar       string  
	Achievements [ {  
			Id       int  
			Name     string  
			Progress uint  
			Value    uint  
		}  
	]  
	Stats [	{  
			Name  string  
			Value int64  
		}  
	]  
}  

POST /updateAchievements  
Input : sessionid, data  
Data format : [ {  
		Id    int  
		Value int  
	}  
]  
Response : Generic response

POST /updateStats  
Input : sessionid, data  
Data Format : [ {  
		Name  string  
		Value int64  
	}  
]  
Response : Generic response  
  
POST /addComplexStats  
Input : sessionid, data  
Data Format : [ {  
		Name 	String  
		Values 	[ {  
				Name	string  
				Value 	string  
			}  
		]  
	}  
]  
Response : Generic response  

## Generic response format  
{  
	Code    int  
	Message string  
}  

# Enums
NebErrorNone         = 0  
	No error to report  
NebError             = 1
	Unknown error  
NebErrorDisconnected = 2  
	User is not connected  
NebErrorLogin        = 3  
	Could not login  
NebErrorPartialFail  = 4  
	One or more operation failed  