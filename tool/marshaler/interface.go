package marshaler

// Marshaler
// @Description: 编码器结构体
type Marshaler interface {
	// 结构体序列化成字节码
	Marshal(interface{}) ([]byte, error)
	// 讲字节码反序列化到结构体
	Unmarshal([]byte, interface{}) error
}
