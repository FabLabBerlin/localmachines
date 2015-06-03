package docs

import (
	"encoding/json"
	"strings"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/swagger"
)

var rootinfo string = `{"apiVersion":"0.1.0","swaggerVersion":"1.2","apis":[{"path":"/users","description":""},{"path":"/memberships","description":""},{"path":"/machines","description":""},{"path":"/activations","description":""},{"path":"/hexabus","description":""},{"path":"/invoices","description":""},{"path":"/urlswitch","description":""}],"info":{"title":"Fabsmith API","description":"Makerspace machine management","contact":"krisjanis.rijnieks@gmail.com"}}`
var subapi string = `{"/activations":{"apiVersion":"0.1.0","swaggerVersion":"1.2","basePath":"","resourcePath":"/activations","produces":["application/json","application/xml","text/plain","text/html"],"apis":[{"path":"/","description":"","operations":[{"httpMethod":"GET","nickname":"Get All","type":"","summary":"Get all activations","parameters":[{"paramType":"query","name":"startDate","description":"\"Period start date\"","dataType":"string","type":"","format":"","allowMultiple":false,"required":true,"minimum":0,"maximum":0},{"paramType":"query","name":"endDate","description":"\"Period end date\"","dataType":"string","type":"","format":"","allowMultiple":false,"required":true,"minimum":0,"maximum":0},{"paramType":"query","name":"userId","description":"\"User ID\"","dataType":"int","type":"","format":"","allowMultiple":false,"required":true,"minimum":0,"maximum":0},{"paramType":"query","name":"includeInvoiced","description":"\"Whether to include already invoiced activations\"","dataType":"bool","type":"","format":"","allowMultiple":false,"required":true,"minimum":0,"maximum":0},{"paramType":"query","name":"itemsPerPage","description":"\"Items per page or max number of items to return\"","dataType":"int","type":"","format":"","allowMultiple":false,"required":true,"minimum":0,"maximum":0},{"paramType":"query","name":"page","description":"\"Current page to show\"","dataType":"int","type":"","format":"","allowMultiple":false,"required":true,"minimum":0,"maximum":0}],"responseMessages":[{"code":200,"message":"models.GetActivationsResponse","responseModel":"GetActivationsResponse"},{"code":403,"message":"Failed to get activations","responseModel":""},{"code":401,"message":"Not authorized","responseModel":""}]}]},{"path":"/:aid","description":"","operations":[{"httpMethod":"GET","nickname":"Get","type":"","summary":"Get activation by activation ID","parameters":[{"paramType":"path","name":"aid","description":"\"Activation ID\"","dataType":"int","type":"","format":"","allowMultiple":false,"required":true,"minimum":0,"maximum":0}],"responseMessages":[{"code":200,"message":"models.Activation","responseModel":"Activation"},{"code":403,"message":"Failed to get activation","responseModel":""},{"code":401,"message":"Not authorized","responseModel":""}]}]},{"path":"/active","description":"","operations":[{"httpMethod":"GET","nickname":"Get Active","type":"","summary":"Get all active activations","responseMessages":[{"code":200,"message":"models.Activation","responseModel":"Activation"},{"code":403,"message":"Failed to get active activations","responseModel":""}]}]},{"path":"/","description":"","operations":[{"httpMethod":"POST","nickname":"Create","type":"","summary":"Create new activation","parameters":[{"paramType":"query","name":"mid","description":"\"Machine ID\"","dataType":"int","type":"","format":"","allowMultiple":false,"required":true,"minimum":0,"maximum":0}],"responseMessages":[{"code":201,"message":"models.ActivationCreateResponse","responseModel":"ActivationCreateResponse"},{"code":403,"message":"Failed to create activation","responseModel":""},{"code":401,"message":"Not authorized","responseModel":""}]}]},{"path":"/:aid","description":"","operations":[{"httpMethod":"PUT","nickname":"Close","type":"","summary":"Close running activation","parameters":[{"paramType":"path","name":"aid","description":"\"Activation ID\"","dataType":"int","type":"","format":"","allowMultiple":false,"required":true,"minimum":0,"maximum":0}],"responseMessages":[{"code":200,"message":"models.StatusResponse","responseModel":"StatusResponse"},{"code":403,"message":"Failed to close activation","responseModel":""},{"code":401,"message":"Not authorized","responseModel":""}]}]},{"path":"/:aid","description":"","operations":[{"httpMethod":"DELETE","nickname":"Delete Activation","type":"","summary":"Delete an activation","parameters":[{"paramType":"path","name":"aid","description":"\"Activation ID\"","dataType":"int","type":"","format":"","allowMultiple":false,"required":true,"minimum":0,"maximum":0}],"responseMessages":[{"code":200,"message":"ok","responseModel":""},{"code":403,"message":"Failed to delete activation","responseModel":""},{"code":401,"message":"Not authorized","responseModel":""}]}]}],"models":{"Activation":{"id":"Activation","properties":{"Active":{"type":"bool","description":"","format":""},"Changed":{"type":"bool","description":"","format":""},"CommentRef":{"type":"string","description":"","format":""},"DiscountFixed":{"type":"float32","description":"","format":""},"DiscountPercents":{"type":"float32","description":"","format":""},"Id":{"type":"int","description":"","format":""},"InvoiceId":{"type":"int","description":"","format":""},"Invoiced":{"type":"bool","description":"","format":""},"MachineId":{"type":"int","description":"","format":""},"TimeEnd":{"type":"\u0026{time Time}","description":"","format":""},"TimeStart":{"type":"\u0026{time Time}","description":"","format":""},"TimeTotal":{"type":"int","description":"","format":""},"UsedKwh":{"type":"float32","description":"","format":""},"UserId":{"type":"int64","description":"","format":""},"VatRate":{"type":"float32","description":"","format":""}}},"ActivationCreateResponse":{"id":"ActivationCreateResponse","properties":{"ActivationId":{"type":"int64","description":"","format":""}}},"GetActivationsResponse":{"id":"GetActivationsResponse","properties":{"ActivationsPage":{"type":"\u0026{787 \u003cnil\u003e Activation}","description":"","format":""},"NumActivations":{"type":"int64","description":"","format":""}}},"StatusResponse":{"id":"StatusResponse","properties":{"Status":{"type":"string","description":"","format":""}}}}},"/hexabus":{"apiVersion":"0.1.0","swaggerVersion":"1.2","basePath":"","resourcePath":"/hexabus","produces":["application/json","application/xml","text/plain","text/html"],"apis":[{"path":"/:mid","description":"","operations":[{"httpMethod":"GET","nickname":"Get","type":"","summary":"Get hexabus mapping by by machine ID","parameters":[{"paramType":"path","name":"mid","description":"\"Machine ID\"","dataType":"int","type":"","format":"","allowMultiple":false,"required":true,"minimum":0,"maximum":0}],"responseMessages":[{"code":200,"message":"models.HexabusMapping","responseModel":"HexabusMapping"},{"code":403,"message":"Failed to get mapping","responseModel":""},{"code":401,"message":"Not authorized","responseModel":""}]}]},{"path":"/","description":"","operations":[{"httpMethod":"POST","nickname":"Create","type":"","summary":"Create hexabus mapping with machine ID","parameters":[{"paramType":"query","name":"mid","description":"\"Machine ID\"","dataType":"int","type":"","format":"","allowMultiple":false,"required":true,"minimum":0,"maximum":0}],"responseMessages":[{"code":200,"message":"ID","responseModel":""},{"code":403,"message":"Failed to create mapping","responseModel":""},{"code":401,"message":"Not authorized","responseModel":""}]}]},{"path":"/:mid","description":"","operations":[{"httpMethod":"DELETE","nickname":"Delete","type":"","summary":"Delete hexabus mapping by by machine ID","parameters":[{"paramType":"path","name":"mid","description":"\"Machine ID\"","dataType":"int","type":"","format":"","allowMultiple":false,"required":true,"minimum":0,"maximum":0}],"responseMessages":[{"code":200,"message":"ok","responseModel":""},{"code":403,"message":"Failed to delete mapping","responseModel":""},{"code":401,"message":"Not authorized","responseModel":""}]}]},{"path":"/:mid","description":"","operations":[{"httpMethod":"PUT","nickname":"Update","type":"","summary":"Update hexabus mapping by by machine ID","parameters":[{"paramType":"path","name":"mid","description":"\"Machine ID\"","dataType":"int","type":"","format":"","allowMultiple":false,"required":true,"minimum":0,"maximum":0},{"paramType":"body","name":"model","description":"\"Hexabus mapping model\"","dataType":"HexabusMapping","type":"","format":"","allowMultiple":false,"required":true,"minimum":0,"maximum":0}],"responseMessages":[{"code":200,"message":"ok","responseModel":""},{"code":403,"message":"Failed to update mapping","responseModel":""},{"code":401,"message":"Not authorized","responseModel":""}]}]}],"models":{"HexabusMapping":{"id":"HexabusMapping","properties":{"Id":{"type":"int64","description":"","format":""},"MachineId":{"type":"int64","description":"","format":""},"SwitchIp":{"type":"string","description":"","format":""}}}}},"/invoices":{"apiVersion":"0.1.0","swaggerVersion":"1.2","basePath":"","resourcePath":"/invoices","produces":["application/json","application/xml","text/plain","text/html"],"apis":[{"path":"/","description":"","operations":[{"httpMethod":"GET","nickname":"Get All Invoices","type":"","summary":"Get all invoices from the database","responseMessages":[{"code":200,"message":"models.Invoice","responseModel":"Invoice"},{"code":403,"message":"Failed to get all invoices","responseModel":""},{"code":401,"message":"Not authorized","responseModel":""}]}]},{"path":"/:iid","description":"","operations":[{"httpMethod":"DELETE","nickname":"Get All Invoices","type":"","summary":"Get all invoices from the database","parameters":[{"paramType":"path","name":"iid","description":"\"Invoice ID\"","dataType":"int","type":"","format":"","allowMultiple":false,"required":true,"minimum":0,"maximum":0}],"responseMessages":[{"code":200,"message":"ok","responseModel":""},{"code":403,"message":"Failed to delete","responseModel":""},{"code":401,"message":"Not authorized","responseModel":""}]}]},{"path":"/","description":"","operations":[{"httpMethod":"POST","nickname":"Create Invoice","type":"","summary":"Create invoice from selection of activations","parameters":[{"paramType":"query","name":"startDate","description":"\"Period start date\"","dataType":"string","type":"","format":"","allowMultiple":false,"required":true,"minimum":0,"maximum":0},{"paramType":"query","name":"endDate","description":"\"Period end date\"","dataType":"string","type":"","format":"","allowMultiple":false,"required":true,"minimum":0,"maximum":0}],"responseMessages":[{"code":200,"message":"models.Invoice","responseModel":"Invoice"},{"code":403,"message":"Failed to create invoice","responseModel":""},{"code":401,"message":"Not authorized","responseModel":""}]}]}],"models":{"Invoice":{"id":"Invoice","properties":{"Activations":{"type":"string","description":"","format":""},"Created":{"type":"\u0026{time Time}","description":"","format":""},"FilePath":{"type":"string","description":"","format":""},"Id":{"type":"int64","description":"","format":""},"PeriodFrom":{"type":"\u0026{time Time}","description":"","format":""},"PeriodTo":{"type":"\u0026{time Time}","description":"","format":""}}}}},"/machines":{"apiVersion":"0.1.0","swaggerVersion":"1.2","basePath":"","resourcePath":"/machines","produces":["application/json","application/xml","text/plain","text/html"],"apis":[{"path":"/","description":"","operations":[{"httpMethod":"GET","nickname":"GetAll","type":"","summary":"Get all machines","responseMessages":[{"code":200,"message":"models.Machine","responseModel":"Machine"},{"code":403,"message":"Failed to get all machines","responseModel":""},{"code":401,"message":"Not authorized","responseModel":""}]}]},{"path":"/:mid","description":"","operations":[{"httpMethod":"GET","nickname":"Get","type":"","summary":"Get machine by machine ID","parameters":[{"paramType":"path","name":"mid","description":"\"Machine ID\"","dataType":"int","type":"","format":"","allowMultiple":false,"required":true,"minimum":0,"maximum":0}],"responseMessages":[{"code":200,"message":"models.Machine","responseModel":"Machine"},{"code":403,"message":"Failed to get machine","responseModel":""},{"code":401,"message":"Not authorized","responseModel":""}]}]},{"path":"/","description":"","operations":[{"httpMethod":"POST","nickname":"Create","type":"","summary":"Create machine","parameters":[{"paramType":"query","name":"mname","description":"\"Machine Name\"","dataType":"string","type":"","format":"","allowMultiple":false,"required":true,"minimum":0,"maximum":0}],"responseMessages":[{"code":200,"message":"models.MachineCreatedResponse","responseModel":"MachineCreatedResponse"},{"code":403,"message":"Failed to create machine","responseModel":""},{"code":401,"message":"Not authorized","responseModel":""}]}]},{"path":"/:mid","description":"","operations":[{"httpMethod":"PUT","nickname":"Update","type":"","summary":"Update machine","parameters":[{"paramType":"path","name":"mid","description":"\"Machine ID\"","dataType":"int","type":"","format":"","allowMultiple":false,"required":true,"minimum":0,"maximum":0},{"paramType":"body","name":"model","description":"\"Machine model\"","dataType":"Machine","type":"","format":"","allowMultiple":false,"required":true,"minimum":0,"maximum":0}],"responseMessages":[{"code":200,"message":"ok","responseModel":""},{"code":403,"message":"Failed to update machine","responseModel":""},{"code":401,"message":"Not authorized","responseModel":""}]}]},{"path":"/:mid","description":"","operations":[{"httpMethod":"DELETE","nickname":"Delete","type":"","summary":"Delete machine","parameters":[{"paramType":"path","name":"mid","description":"\"Machine ID\"","dataType":"int","type":"","format":"","allowMultiple":false,"required":true,"minimum":0,"maximum":0}],"responseMessages":[{"code":200,"message":"ok","responseModel":""},{"code":403,"message":"Failed to delete machine","responseModel":""},{"code":401,"message":"Not authorized","responseModel":""}]}]},{"path":"/:mid/image","description":"","operations":[{"httpMethod":"POST","nickname":"PostImage","type":"","summary":"Post machine image","parameters":[{"paramType":"path","name":"mid","description":"\"Machine ID\"","dataType":"int","type":"","format":"","allowMultiple":false,"required":true,"minimum":0,"maximum":0}],"responseMessages":[{"code":200,"message":"ok","responseModel":""},{"code":400,"message":"Bad Request","responseModel":""},{"code":401,"message":"Not authorized","responseModel":""},{"code":500,"message":"Internal Server Error","responseModel":""}]}]}],"models":{"Machine":{"id":"Machine","properties":{"Available":{"type":"bool","description":"","format":""},"Comments":{"type":"string","description":"","format":""},"Description":{"type":"string","description":"","format":""},"Id":{"type":"int64","description":"","format":""},"Image":{"type":"string","description":"","format":""},"Name":{"type":"string","description":"","format":""},"Price":{"type":"float32","description":"","format":""},"PriceUnit":{"type":"string","description":"","format":""},"Shortname":{"type":"string","description":"","format":""},"UnavailMsg":{"type":"string","description":"","format":""},"UnavailTill":{"type":"\u0026{time Time}","description":"","format":""}}},"MachineCreatedResponse":{"id":"MachineCreatedResponse","properties":{"MachineId":{"type":"int64","description":"","format":""}}}}},"/memberships":{"apiVersion":"0.1.0","swaggerVersion":"1.2","basePath":"","resourcePath":"/memberships","produces":["application/json","application/xml","text/plain","text/html"],"apis":[{"path":"/","description":"","operations":[{"httpMethod":"GET","nickname":"GetAll","type":"","summary":"Get all memberships","responseMessages":[{"code":200,"message":"models.Membership","responseModel":"Membership"},{"code":403,"message":"Failed to get all memberships","responseModel":""}]}]},{"path":"/","description":"","operations":[{"httpMethod":"POST","nickname":"Create","type":"","summary":"Create new membership","parameters":[{"paramType":"query","name":"mname","description":"\"Membership Name\"","dataType":"string","type":"","format":"","allowMultiple":false,"required":true,"minimum":0,"maximum":0}],"responseMessages":[{"code":200,"message":"ID","responseModel":""},{"code":403,"message":"Failed to create membership","responseModel":""},{"code":401,"message":"Not authorized","responseModel":""}]}]},{"path":"/:mid","description":"","operations":[{"httpMethod":"GET","nickname":"Get","type":"","summary":"Get membership by membership ID","parameters":[{"paramType":"path","name":"mid","description":"\"Membership ID\"","dataType":"int","type":"","format":"","allowMultiple":false,"required":true,"minimum":0,"maximum":0}],"responseMessages":[{"code":200,"message":"models.Membership","responseModel":"Membership"},{"code":403,"message":"Failed to get membership","responseModel":""},{"code":401,"message":"Not authorized","responseModel":""}]}]},{"path":"/:mid","description":"","operations":[{"httpMethod":"PUT","nickname":"Update","type":"","summary":"Update membership","parameters":[{"paramType":"path","name":"mid","description":"\"Membership ID\"","dataType":"int","type":"","format":"","allowMultiple":false,"required":true,"minimum":0,"maximum":0},{"paramType":"body","name":"model","description":"\"Membership model\"","dataType":"Membership","type":"","format":"","allowMultiple":false,"required":true,"minimum":0,"maximum":0}],"responseMessages":[{"code":200,"message":"ok","responseModel":""},{"code":403,"message":"Failed to update membership","responseModel":""},{"code":401,"message":"Not authorized","responseModel":""}]}]},{"path":"/:mid","description":"","operations":[{"httpMethod":"DELETE","nickname":"Delete","type":"","summary":"Delete membership","parameters":[{"paramType":"path","name":"mid","description":"\"Membership ID\"","dataType":"int","type":"","format":"","allowMultiple":false,"required":true,"minimum":0,"maximum":0}],"responseMessages":[{"code":200,"message":"ok","responseModel":""},{"code":403,"message":"Failed to delete membership","responseModel":""},{"code":401,"message":"Not authorized","responseModel":""}]}]}],"models":{"Membership":{"id":"Membership","properties":{"AffectedMachines":{"type":"string","description":"","format":""},"Duration":{"type":"int","description":"","format":""},"Id":{"type":"int64","description":"","format":""},"MachinePriceDeduction":{"type":"int","description":"","format":""},"Price":{"type":"float32","description":"","format":""},"ShortName":{"type":"string","description":"","format":""},"Title":{"type":"string","description":"","format":""},"Unit":{"type":"string","description":"","format":""}}}}},"/urlswitch":{"apiVersion":"0.1.0","swaggerVersion":"1.2","basePath":"","resourcePath":"/urlswitch","produces":["application/json","application/xml","text/plain","text/html"],"apis":[{"path":"/","description":"","operations":[{"httpMethod":"POST","nickname":"Create","type":"","summary":"Create UrlSwitch mapping with machine ID","parameters":[{"paramType":"query","name":"mid","description":"\"Machine ID\"","dataType":"int","type":"","format":"","allowMultiple":false,"required":true,"minimum":0,"maximum":0}],"responseMessages":[{"code":200,"message":"ID","responseModel":""},{"code":500,"message":"Internal Server Error","responseModel":""},{"code":401,"message":"Not authorized","responseModel":""}]}]}]},"/users":{"apiVersion":"0.1.0","swaggerVersion":"1.2","basePath":"","resourcePath":"/users","produces":["application/json","application/xml","text/plain","text/html"],"apis":[{"path":"/login","description":"","operations":[{"httpMethod":"POST","nickname":"login","type":"","summary":"Logs user into the system","parameters":[{"paramType":"query","name":"username","description":"\"The username for login\"","dataType":"string","type":"","format":"","allowMultiple":false,"required":true,"minimum":0,"maximum":0},{"paramType":"query","name":"password","description":"\"The password for login\"","dataType":"string","type":"","format":"","allowMultiple":false,"required":true,"minimum":0,"maximum":0}],"responseMessages":[{"code":200,"message":"models.LoginResponse","responseModel":"LoginResponse"},{"code":401,"message":"Failed to authenticate","responseModel":""}]}]},{"path":"/loginuid","description":"","operations":[{"httpMethod":"POST","nickname":"LoginUid","type":"","summary":"Logs user into the system by using NFC UID","parameters":[{"paramType":"query","name":"uid","description":"\"The NFC UID\"","dataType":"string","type":"","format":"","allowMultiple":false,"required":true,"minimum":0,"maximum":0}],"responseMessages":[{"code":200,"message":"models.LoginResponse","responseModel":"LoginResponse"},{"code":401,"message":"Failed to authenticate","responseModel":""}]}]},{"path":"/logout","description":"","operations":[{"httpMethod":"GET","nickname":"logout","type":"","summary":"Logs out current logged in user session","responseMessages":[{"code":200,"message":"models.StatusResponse","responseModel":"StatusResponse"}]}]},{"path":"/","description":"","operations":[{"httpMethod":"GET","nickname":"GetAll","type":"","summary":"Get all users","responseMessages":[{"code":200,"message":"models.User","responseModel":"User"},{"code":403,"message":"Failed to get all users","responseModel":""}]}]},{"path":"/signup","description":"","operations":[{"httpMethod":"POST","nickname":"Signup","type":"","summary":"Accept user signup, create a zombie user with no privileges for later access","parameters":[{"paramType":"body","name":"model","description":"\"User model and password\"","dataType":"UserSignupRequest","type":"","format":"","allowMultiple":false,"required":true,"minimum":0,"maximum":0}],"responseMessages":[{"code":200,"message":"ok","responseModel":""},{"code":500,"message":"Internal Server Error","responseModel":""}]}]},{"path":"/","description":"","operations":[{"httpMethod":"POST","nickname":"Post","type":"","summary":"create user and associated tables","parameters":[{"paramType":"query","name":"email","description":"\"The new user's E-Mail\"","dataType":"string","type":"","format":"","allowMultiple":false,"required":true,"minimum":0,"maximum":0}],"responseMessages":[{"code":201,"message":"models.User","responseModel":"User"},{"code":401,"message":"Unauthorized","responseModel":""},{"code":500,"message":"Internal Server Error","responseModel":""}]}]},{"path":"/:uid","description":"","operations":[{"httpMethod":"GET","nickname":"Get","type":"","summary":"get user by uid","parameters":[{"paramType":"path","name":"uid","description":"\"User ID\"","dataType":"int","type":"","format":"","allowMultiple":false,"required":true,"minimum":0,"maximum":0}],"responseMessages":[{"code":200,"message":"models.User","responseModel":"User"},{"code":403,"message":"Variable message","responseModel":""},{"code":401,"message":"Unauthorized","responseModel":""}]}]},{"path":"/:uid","description":"","operations":[{"httpMethod":"DELETE","nickname":"Delete","type":"","summary":"delete user with uid","parameters":[{"paramType":"path","name":"uid","description":"\"User ID\"","dataType":"int","type":"","format":"","allowMultiple":false,"required":true,"minimum":0,"maximum":0}],"responseMessages":[{"code":0,"message":"200","responseModel":""},{"code":403,"message":"Variable message","responseModel":""},{"code":401,"message":"Unauthorized","responseModel":""}]}]},{"path":"/:uid","description":"","operations":[{"httpMethod":"PUT","nickname":"Put","type":"","summary":"Update user with uid","parameters":[{"paramType":"path","name":"uid","description":"\"User ID\"","dataType":"int","type":"","format":"","allowMultiple":false,"required":true,"minimum":0,"maximum":0}],"responseMessages":[{"code":0,"message":"200","responseModel":""},{"code":400,"message":"Variable message","responseModel":""},{"code":401,"message":"Unauthorized","responseModel":""},{"code":403,"message":"Variable message","responseModel":""}]}]},{"path":"/:uid/machinepermissions","description":"","operations":[{"httpMethod":"GET","nickname":"GetUserMachinePermissions","type":"","summary":"Get current saved machine permissions","parameters":[{"paramType":"path","name":"uid","description":"\"User ID\"","dataType":"int","type":"","format":"","allowMultiple":false,"required":true,"minimum":0,"maximum":0}],"responseMessages":[{"code":200,"message":"models.Machine","responseModel":"Machine"},{"code":500,"message":"Internal Server Error","responseModel":""},{"code":401,"message":"Unauthorized","responseModel":""}]}]},{"path":"/:uid/machines","description":"","operations":[{"httpMethod":"GET","nickname":"GetUserMachines","type":"","summary":"Get user machines, all machines for admin user","parameters":[{"paramType":"path","name":"uid","description":"\"User ID\"","dataType":"int","type":"","format":"","allowMultiple":false,"required":true,"minimum":0,"maximum":0}],"responseMessages":[{"code":200,"message":"models.Machine","responseModel":"Machine"},{"code":500,"message":"Internal Server Error","responseModel":""},{"code":401,"message":"Unauthorized","responseModel":""}]}]},{"path":"/:uid/memberships","description":"","operations":[{"httpMethod":"POST","nickname":"PostUserMemberships","type":"","summary":"Post user membership","parameters":[{"paramType":"path","name":"uid","description":"\"User ID\"","dataType":"int","type":"","format":"","allowMultiple":false,"required":true,"minimum":0,"maximum":0}],"responseMessages":[{"code":200,"message":"models.UserMembership","responseModel":"UserMembership"},{"code":403,"message":"Failed to get user memberships","responseModel":""},{"code":401,"message":"Not authorized","responseModel":""}]}]},{"path":"/:uid/memberships","description":"","operations":[{"httpMethod":"GET","nickname":"GetUserMemberships","type":"","summary":"Get user memberships","parameters":[{"paramType":"path","name":"uid","description":"\"User ID\"","dataType":"int","type":"","format":"","allowMultiple":false,"required":true,"minimum":0,"maximum":0}],"responseMessages":[{"code":200,"message":"models.UserMembership","responseModel":"UserMembership"},{"code":403,"message":"Failed to get user memberships","responseModel":""},{"code":401,"message":"Not authorized","responseModel":""}]}]},{"path":"/:uid/memberships/:umid","description":"","operations":[{"httpMethod":"DELETE","nickname":"DeleteUserMembership","type":"","summary":"Delete user membership","parameters":[{"paramType":"path","name":"uid","description":"\"User ID\"","dataType":"int","type":"","format":"","allowMultiple":false,"required":true,"minimum":0,"maximum":0},{"paramType":"path","name":"umid","description":"\"User Membership ID\"","dataType":"int","type":"","format":"","allowMultiple":false,"required":true,"minimum":0,"maximum":0}],"responseMessages":[{"code":0,"message":"200","responseModel":""},{"code":403,"message":"Failed to get user memberships","responseModel":""},{"code":401,"message":"Not authorized","responseModel":""}]}]},{"path":"/:uid/name","description":"","operations":[{"httpMethod":"GET","nickname":"GetUserName","type":"","summary":"Get user name data only","parameters":[{"paramType":"path","name":"uid","description":"\"User ID\"","dataType":"int","type":"","format":"","allowMultiple":false,"required":true,"minimum":0,"maximum":0}],"responseMessages":[{"code":200,"message":"models.UserNameResponse","responseModel":"UserNameResponse"},{"code":403,"message":"Failed to get user name","responseModel":""},{"code":401,"message":"Not loggen","responseModel":""}]}]},{"path":"/:uid/password","description":"","operations":[{"httpMethod":"POST","nickname":"PostUserPassword","type":"","summary":"Post user password","parameters":[{"paramType":"path","name":"uid","description":"\"User ID\"","dataType":"int","type":"","format":"","allowMultiple":false,"required":true,"minimum":0,"maximum":0}],"responseMessages":[{"code":0,"message":"200","responseModel":""},{"code":403,"message":"Failed to get user","responseModel":""},{"code":401,"message":"Not authorized","responseModel":""}]}]},{"path":"/:uid/nfcuid","description":"","operations":[{"httpMethod":"PUT","nickname":"UpdateNfcUid","type":"","summary":"Update user NFC UID","parameters":[{"paramType":"path","name":"uid","description":"\"User ID\"","dataType":"int","type":"","format":"","allowMultiple":false,"required":true,"minimum":0,"maximum":0},{"paramType":"query","name":"nfcuid","description":"\"NFC UID\"","dataType":"int","type":"","format":"","allowMultiple":false,"required":true,"minimum":0,"maximum":0}],"responseMessages":[{"code":200,"message":"ok","responseModel":""},{"code":403,"message":"Failed to update NFC UID","responseModel":""},{"code":401,"message":"Not authorized","responseModel":""}]}]},{"path":"/:uid/permissions","description":"","operations":[{"httpMethod":"POST","nickname":"CreateUserPermission","type":"","summary":"Create a permission for a user to allow him/her to use a machine","parameters":[{"paramType":"path","name":"uid","description":"\"User ID\"","dataType":"int","type":"","format":"","allowMultiple":false,"required":true,"minimum":0,"maximum":0},{"paramType":"query","name":"mid","description":"\"Machine ID\"","dataType":"int","type":"","format":"","allowMultiple":false,"required":true,"minimum":0,"maximum":0}],"responseMessages":[{"code":200,"message":"ok","responseModel":""},{"code":403,"message":"Failed to create permission","responseModel":""},{"code":401,"message":"Not authorized","responseModel":""}]}]},{"path":"/:uid/permissions","description":"","operations":[{"httpMethod":"DELETE","nickname":"DeleteUserPermission","type":"","summary":"Delete user machine permission","parameters":[{"paramType":"path","name":"uid","description":"\"User ID\"","dataType":"int","type":"","format":"","allowMultiple":false,"required":true,"minimum":0,"maximum":0},{"paramType":"query","name":"mid","description":"\"Machine ID\"","dataType":"int","type":"","format":"","allowMultiple":false,"required":true,"minimum":0,"maximum":0}],"responseMessages":[{"code":200,"message":"ok","responseModel":""},{"code":403,"message":"Failed to delete permission","responseModel":""},{"code":401,"message":"Not authorized","responseModel":""}]}]},{"path":"/:uid/permissions","description":"","operations":[{"httpMethod":"PUT","nickname":"Update User Machine Permissions","type":"","summary":"Update user machine permissions","parameters":[{"paramType":"path","name":"uid","description":"\"User ID\"","dataType":"int","type":"","format":"","allowMultiple":false,"required":true,"minimum":0,"maximum":0},{"paramType":"body","name":"model","description":"\"Permissions Array\"","dataType":"Permission","type":"","format":"","allowMultiple":false,"required":true,"minimum":0,"maximum":0}],"responseMessages":[{"code":200,"message":"ok","responseModel":""},{"code":403,"message":"Failed to update permissions","responseModel":""},{"code":401,"message":"Not authorized","responseModel":""}]}]},{"path":"/:uid/permissions","description":"","operations":[{"httpMethod":"GET","nickname":"Get User Machine Permissions","type":"","summary":"Get user machine permissions","parameters":[{"paramType":"path","name":"uid","description":"\"User ID\"","dataType":"int","type":"","format":"","allowMultiple":false,"required":true,"minimum":0,"maximum":0}],"responseMessages":[{"code":200,"message":"models.Permission","responseModel":""},{"code":403,"message":"Failed to update permissions","responseModel":""},{"code":401,"message":"Not authorized","responseModel":""}]}]}],"models":{"LoginResponse":{"id":"LoginResponse","properties":{"Status":{"type":"string","description":"","format":""},"UserId":{"type":"int64","description":"","format":""}}},"Machine":{"id":"Machine","properties":{"Available":{"type":"bool","description":"","format":""},"Comments":{"type":"string","description":"","format":""},"Description":{"type":"string","description":"","format":""},"Id":{"type":"int64","description":"","format":""},"Image":{"type":"string","description":"","format":""},"Name":{"type":"string","description":"","format":""},"Price":{"type":"float32","description":"","format":""},"PriceUnit":{"type":"string","description":"","format":""},"Shortname":{"type":"string","description":"","format":""},"UnavailMsg":{"type":"string","description":"","format":""},"UnavailTill":{"type":"\u0026{time Time}","description":"","format":""}}},"StatusResponse":{"id":"StatusResponse","properties":{"Status":{"type":"string","description":"","format":""}}},"User":{"id":"User","properties":{"B2b":{"type":"bool","description":"","format":""},"ClientId":{"type":"int","description":"","format":""},"Comments":{"type":"string","description":"","format":""},"Company":{"type":"string","description":"","format":""},"Created":{"type":"\u0026{time Time}","description":"","format":""},"Email":{"type":"string","description":"","format":""},"FirstName":{"type":"string","description":"","format":""},"Id":{"type":"int64","description":"","format":""},"InvoiceAddr":{"type":"int","description":"","format":""},"LastName":{"type":"string","description":"","format":""},"ShipAddr":{"type":"int","description":"","format":""},"UserRole":{"type":"string","description":"","format":""},"Username":{"type":"string","description":"","format":""},"VatRate":{"type":"int","description":"","format":""},"VatUserId":{"type":"string","description":"","format":""}}},"UserMembership":{"id":"UserMembership","properties":{"Id":{"type":"int64","description":"","format":""},"MembershipId":{"type":"int64","description":"","format":""},"StartDate":{"type":"\u0026{time Time}","description":"","format":""},"UserId":{"type":"int64","description":"","format":""}}},"UserNameResponse":{"id":"UserNameResponse","properties":{"FirstName":{"type":"string","description":"","format":""},"LastName":{"type":"string","description":"","format":""},"UserId":{"type":"int64","description":"","format":""}}}}}}`
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
