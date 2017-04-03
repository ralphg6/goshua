package routers

func init() {
	ns := beego.NewNamespace("/v1",

		beego.NSNamespace("/accounts",
			beego.NSInclude(
				&AccountsController{},
			),
		),
	)

	beego.AddNamespace(ns)
}
