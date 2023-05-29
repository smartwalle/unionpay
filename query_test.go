package unionpay_test

import "testing"

func TestClient_Query(t *testing.T) {
	var rsp, err = client.Query("testssss")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(rsp)
}
