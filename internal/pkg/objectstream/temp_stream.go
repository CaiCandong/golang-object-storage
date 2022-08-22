package objectstream

func NewTempPutStream(server string, name string, size int64) (*TempPutStream, error) {
	t := &TempPutStream{Server: server, Name: name, Size: size}
	err := t.post()
	return t, err
}

// 实现io.Writer口
func (t *TempPutStream) Write(p []byte) (n int, err error) {
	err = t.patch(p)
	if err != nil {
		return 0, err
	}
	return len(p), nil
}

// 根据positive决定是删除临时对象数据还是转正保存到节点中
func (t *TempPutStream) Commit(positive bool) {
	if positive {
		t.put()
	} else {
		t.del()
	}
}

// 临时文件下载
func NewTempGetStream(server, uuid string) (*GetStream, error) {
	t := &TempPutStream{Server: server, UUID: uuid}
	return t.get()
}
