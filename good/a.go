package good


func CallAnyMethod(v interface{}, method string) interface{} {
	ref := reflect.ValueOf(v)
	refKind := ref.Kind()
	if refKind == reflect.Ptr {
		refKind = ref.Elem().Kind()
	}
	// 如果是结构的话，尝试检索一下他是否有Int、ToInt的函数
	// 不再限制为Struct
	//	if refKind == reflect.Struct {
	fn := ref.MethodByName(method)
	if fn.IsValid() {
		rs := fn.Call(nil)
		if len(rs) > 0 {
			return rs[0].Interface()
		}
	}
	//	}
	return nil

