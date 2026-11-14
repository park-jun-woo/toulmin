//ff:type feature=engine type=model
//ff:what validateSpecFieldsValueSpec — non-pointer Spec implementation used to exercise the branch where reflect.TypeOf(s).Kind() != reflect.Ptr
package toulmin

type validateSpecFieldsValueSpec struct {
	Value string
}
