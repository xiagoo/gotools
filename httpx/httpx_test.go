package httpx

import (
	"io/ioutil"
	"testing"
)

func TestAgent_Get(t *testing.T) {
	agent := NewAgent()
	params := map[string]string{
		"domain": "11111.com",
		"token":  "1",
	}
	resp, err := agent.Get("https://checkapi.aliyun.com/check/checkdomain").AddData(params).GetResponse( nil)
	if err != nil {
		t.Logf("err:%s\n", err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Logf("err:%s\n", err)
	}
	t.Logf("resp:%v\n", string(body))
}

