package berrt

// type MockAdvisor struct{}

// func (advisor MockAdvisor) Static(errs ring.Ring) bool {
// 	return true
// }

// func (advisor MockAdvisor) Dynamic(err <-chan errnie.Error) bool {
// 	return true
// }

// func TestNewAdvisor(t *testing.T) {
// 	Convey("Given an struct type", t, func() {
// 		r := ring.New(1)
// 		r.Value = errnie.NewError(errors.New("test"))

// 		So(
// 			NewAdvisor(MockAdvisor{}).Static(*r),
// 			ShouldEqual,
// 			NewAdvisor(ProtoAdvisor{}).Static(*r),
// 		)
// 	})
// }
