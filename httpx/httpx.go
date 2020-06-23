package httpx

import (
	"io/ioutil"
	"net/http"
	"strings"
)

type Agent struct {
	client   *http.Client
	url      string
	method   string
	postType string
	header   map[string]string
	data     map[string]string
}

func NewAgent() *Agent {
	return &Agent{
		client: &http.Client{},
	}
}

func (agent *Agent) Get(url string) *Agent {
	agent.url = url
	agent.method = http.MethodGet
	return agent
}

func (agent *Agent) Post(url string) *Agent {
	agent.url = url
	agent.method = http.MethodPost
	return agent
}

func (agent *Agent) AddHeader(header map[string]string) *Agent {
	agent.header = header
	return agent
}
func (agent *Agent) AddData(data map[string]string) *Agent {
	agent.data = data
	return agent
}

func (agent *Agent) getRequest(params map[string]string) (*http.Request, error) {
	var req *http.Request
	var err error
	switch agent.method {
	case http.MethodGet:
		var p []string
		for k, v := range params {
			p = append(p, k+"="+v)
		}
		if len(p) > 0 {
			agent.url = agent.url + "?" + strings.Join(p, "&")
		}
		req, err = http.NewRequest(agent.method, agent.url, nil)
	case http.MethodPost:
		switch agent.postType {
		case "Form":
		}
	default:
	}
	if err != nil {
		return nil, err
	}
	for k, v := range agent.header {
		req.Header.Add(k, v)
	}
	return req, nil
}

func (agent *Agent) GetResponse(callback func(resp *http.Response, err error)) (*http.Response, error) {
	req, err := agent.getRequest(agent.data)
	if err != nil {
		return nil, err
	}
	return agent.client.Do(req)
}
func (agent *Agent) GetResponseBody(callback func(resp *http.Response, err error)) ([]byte, error) {
	req, err := agent.getRequest(agent.data)
	if err != nil {
		return nil, err
	}
	resp, err := agent.client.Do(req)
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}
