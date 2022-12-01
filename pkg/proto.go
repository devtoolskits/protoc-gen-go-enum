package go_enum

import "google.golang.org/protobuf/reflect/protoreflect"

// EnumMembers returns string names of enum elements/members
func EnumMembers(enum protoreflect.Enum) []string {
	var members []string
	v := enum.Descriptor().Values()
	for i := 0; i < v.Len(); i++ {
		if n := v.Get(i).Name(); n.IsValid() {
			members = append(members, string(n))
		}
	}
	return members
}
