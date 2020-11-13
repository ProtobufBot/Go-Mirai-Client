package util

import "testing"

func TestGetBytes(t *testing.T) {
	bytes, err := GetBytes("http://tnoodle.lz1998.xin/view/222.png?scramble=U2+R+U%27+F%27+U2+R%27+U%27+R+U2+R2+F%27")
	if err != nil {
		t.Error(err)
	}
	t.Logf("%+v", bytes)
}
