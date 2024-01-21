package gomongo

func Int64Ptr(v *int) *int64 {
	if v == nil {
		return nil
	}
	v64 := int64(*v)
	return &v64
}
