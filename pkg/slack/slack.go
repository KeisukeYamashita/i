package slack

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/KeisukeYamashita/i/api/v1alpha1"
	"github.com/KeisukeYamashita/i/pkg/pointer"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
)

// Message ...
type Message struct {
	Text        string       `json:"text"`
	Attachments []Attachment `json:"attachments"`
}

// Attachment ...
type Attachment struct {
	Fields []*Field `json:"fields"`
	Title  string   `json:"title"`
}

// Field ...
type Field struct {
	Title string      `json:"title"`
	Value interface{} `json:"value"`
	Short *bool       `json:"short"`
}

// NewField ...
func NewField(title string, value interface{}, short *bool) *Field {
	return &Field{
		Title: title,
		Value: value,
		Short: short,
	}
}

// Client ...
type Client interface {
	PostMessage(*Message) error
}

var _ Client = (*client)(nil)

type client struct {
	HookURL *url.URL
}

// NewClient ...
func NewClient(url *url.URL) Client {
	return &client{
		HookURL: url,
	}
}

// NewMessage ...
func NewMessage(text string, attachments []Attachment) *Message {
	return &Message{
		Text:        text,
		Attachments: attachments,
	}
}

func (c *client) PostMessage(msg *Message) error {
	client := http.DefaultClient
	json, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	req, err := http.NewRequest(
		"POST",
		c.HookURL.String(),
		bytes.NewBuffer(json),
	)
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf("response from slack api was not 200 got %d", resp.StatusCode)
	}

	return nil
}

// NewInvalidPodsMessage ...
func NewInvalidPodsMessage(eye *v1alpha1.Eye, nn types.NamespacedName, pods []corev1.Pod) *Message {
	text := "Deleted old pods"
	attachments := []Attachment{}

	for _, pod := range pods {
		fields := GetPodFields(eye, &pod, nn.Namespace)
		attachment := Attachment{
			Fields: fields,
		}
		attachments = append(attachments, attachment)
	}
	return NewMessage(text, attachments)
}

// GetPodFields ...
func GetPodFields(eye *v1alpha1.Eye, pod *corev1.Pod, ns string) (fields []*Field) {
	eyeField := newShortField("Eye name", eye.Name)
	fields = append(fields, eyeField)
	namespaceField := newShortField("Eye namespace", ns)
	nameField := NewField("Name", pod.ObjectMeta.Name, nil)
	fields = append(fields, nameField)
	fields = append(fields, namespaceField)
	hostIPField := newShortField("Host IP", pod.Status.HostIP)
	fields = append(fields, hostIPField)
	podIPField := newShortField("Pod IP", pod.Status.PodIP)
	fields = append(fields, podIPField)
	podNamespaceField := newShortField("Namespace", pod.Namespace)
	fields = append(fields, podNamespaceField)
	return fields
}

func newShortField(title string, value interface{}) *Field {
	return NewField(title, value, pointer.Bool(true))
}
