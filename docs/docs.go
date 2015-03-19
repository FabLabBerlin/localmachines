package docs

import (
	"encoding/json"
	"strings"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/swagger"
)

var rootinfo string = `{"apiVersion":"0.1.0","swaggerVersion":"1.2","apis":[{"path":"/user","description":""},{"path":"/machine","description":""}],"info":{"title":"Fabsmith API","description":"Makerspace machine management","contact":"krisjanis.rijnieks@gmail.com"}}`
var subapi string = `{"/machine":{"apiVersion":"0.1.0","swaggerVersion":"1.2","basePath":"","resourcePath":"/machine","produces":["application/json","application/xml","text/plain","text/html"],"apis":[{"path":"/:mid","description":"","operations":[{"httpMethod":"GET","nickname":"Get","type":"","summary":"Get machine by machine ID","parameters":[{"paramType":"path","name":"mid","description":"\"Machine ID\"","dataType":"int","type":"","format":"","allowMultiple":false,"required":true,"minimum":0,"maximum":0}],"responseMessages":[{"code":200,"message":"models.Machine","responseModel":"Machine"},{"code":403,"message":"Failed to get machine","responseModel":""},{"code":401,"message":"Not authorized","responseModel":""}]}]},{"path":"/","description":"","operations":[{"httpMethod":"GET","nickname":"GetAll","type":"","summary":"Get all machines","responseMessages":[{"code":200,"message":"models.Machine","responseModel":"Machine"},{"code":403,"message":"Failed to get all machines","responseModel":""},{"code":401,"message":"Not authorized","responseModel":""}]}]}],"models":{"Machine":{"id":"Machine","properties":{"Available":{"type":"bool","description":"","format":""},"CalcByEnergy":{"type":"bool","description":"","format":""},"CalcByTime":{"type":"bool","description":"","format":""},"CostsPerKwh":{"type":"float32","description":"","format":""},"CostsPerMin":{"type":"float32","description":"","format":""},"Description":{"type":"string","description":"","format":""},"Id":{"type":"int","description":"","format":""},"Name":{"type":"string","description":"","format":""},"UnavailMsg":{"type":"string","description":"","format":""},"UnavailTill":{"type":"\u0026{time Time}","description":"","format":""}}}}},"/user":{"apiVersion":"0.1.0","swaggerVersion":"1.2","basePath":"","resourcePath":"/user","produces":["application/json","application/xml","text/plain","text/html"],"apis":[{"path":"/login","description":"","operations":[{"httpMethod":"POST","nickname":"login","type":"","summary":"Logs user into the system","parameters":[{"paramType":"query","name":"username","description":"\"The username for login\"","dataType":"string","type":"","format":"","allowMultiple":false,"required":true,"minimum":0,"maximum":0},{"paramType":"query","name":"password","description":"\"The password for login\"","dataType":"string","type":"","format":"","allowMultiple":false,"required":true,"minimum":0,"maximum":0}],"responseMessages":[{"code":200,"message":"models.LoginResponse","responseModel":"LoginResponse"},{"code":401,"message":"Failed to authenticate","responseModel":""}]}]},{"path":"/logout","description":"","operations":[{"httpMethod":"GET","nickname":"logout","type":"","summary":"Logs out current logged in user session","responseMessages":[{"code":200,"message":"models.StatusResponse","responseModel":"StatusResponse"}]}]},{"path":"/:uid","description":"","operations":[{"httpMethod":"GET","nickname":"Get","type":"","summary":"get user by uid","parameters":[{"paramType":"path","name":"uid","description":"\"User ID\"","dataType":"int","type":"","format":"","allowMultiple":false,"required":true,"minimum":0,"maximum":0}],"responseMessages":[{"code":200,"message":"models.User","responseModel":"User"},{"code":403,"message":"Variable message","responseModel":""},{"code":401,"message":"Unauthorized","responseModel":""}]}]},{"path":"/","description":"","operations":[{"httpMethod":"GET","nickname":"GetAll","type":"","summary":"Get all users","responseMessages":[{"code":200,"message":"models.User","responseModel":"User"},{"code":403,"message":"Failed to get all users","responseModel":""}]}]},{"path":"/:uid/machines","description":"","operations":[{"httpMethod":"GET","nickname":"GetUserMachines","type":"","summary":"Get user machines","parameters":[{"paramType":"path","name":"uid","description":"\"User ID\"","dataType":"int","type":"","format":"","allowMultiple":false,"required":true,"minimum":0,"maximum":0}],"responseMessages":[{"code":200,"message":"models.Machine","responseModel":"Machine"},{"code":403,"message":"Failed to get user machines","responseModel":""},{"code":401,"message":"Not authorized","responseModel":""}]}]}],"models":{"LoginResponse":{"id":"LoginResponse","properties":{"Status":{"type":"string","description":"","format":""},"UserId":{"type":"int","description":"","format":""}}},"Machine":{"id":"Machine","properties":{"Available":{"type":"bool","description":"","format":""},"CalcByEnergy":{"type":"bool","description":"","format":""},"CalcByTime":{"type":"bool","description":"","format":""},"CostsPerKwh":{"type":"float32","description":"","format":""},"CostsPerMin":{"type":"float32","description":"","format":""},"Description":{"type":"string","description":"","format":""},"Id":{"type":"int","description":"","format":""},"Name":{"type":"string","description":"","format":""},"UnavailMsg":{"type":"string","description":"","format":""},"UnavailTill":{"type":"\u0026{time Time}","description":"","format":""}}},"StatusResponse":{"id":"StatusResponse","properties":{"Status":{"type":"string","description":"","format":""}}},"User":{"id":"User","properties":{"B2b":{"type":"bool","description":"","format":""},"ClientId":{"type":"int","description":"","format":""},"Company":{"type":"string","description":"","format":""},"Email":{"type":"string","description":"","format":""},"FirstName":{"type":"string","description":"","format":""},"Id":{"type":"int","description":"","format":""},"InvoiceAddr":{"type":"int","description":"","format":""},"LastName":{"type":"string","description":"","format":""},"ShipAddr":{"type":"int","description":"","format":""},"Username":{"type":"string","description":"","format":""},"VatRate":{"type":"int","description":"","format":""},"VatUserId":{"type":"string","description":"","format":""}}}}}}`
var rootapi swagger.ResourceListing

var apilist map[string]*swagger.ApiDeclaration

func init() {
	basepath := "/api"
	err := json.Unmarshal([]byte(rootinfo), &rootapi)
	if err != nil {
		beego.Error(err)
	}
	err = json.Unmarshal([]byte(subapi), &apilist)
	if err != nil {
		beego.Error(err)
	}
	beego.GlobalDocApi["Root"] = rootapi
	for k, v := range apilist {
		for i, a := range v.Apis {
			a.Path = urlReplace(k + a.Path)
			v.Apis[i] = a
		}
		v.BasePath = basepath
		beego.GlobalDocApi[strings.Trim(k, "/")] = v
	}
}


func urlReplace(src string) string {
	pt := strings.Split(src, "/")
	for i, p := range pt {
		if len(p) > 0 {
			if p[0] == ':' {
				pt[i] = "{" + p[1:] + "}"
			} else if p[0] == '?' && p[1] == ':' {
				pt[i] = "{" + p[2:] + "}"
			}
		}
	}
	return strings.Join(pt, "/")
}
