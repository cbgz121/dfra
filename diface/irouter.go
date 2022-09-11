package diface

//路由抽象接口， 路由里的数据都是IRequest

type IRouter interface {
	//在处理conn业务之前的钩子方法（Hook）
	PreHandle(request IRequset)
	//在处理conn业务的主方法（Hook）
	Handle(request IRequset)
	//在处理conn业务之后的钩子方法（Hook）
	PostHandle(request IRequset)
}
