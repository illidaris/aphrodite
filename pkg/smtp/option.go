package smtp

import (
	"fmt"
	"strings"
	"time"
)

type SendOption func(*SendOptions)

func NewSendOptions(opts ...SendOption) SendOptions {
	option := SendOptions{
		ContentType: MAIL_CONTENTTYPE_TEXT, // 默认使用text
		Date:        time.Now(),
	}
	for _, f := range opts {
		f(&option)
	}
	return option
}

type SendOptions struct {
	FromNick    string    // 发件邮箱别称
	From        string    // 发件邮箱
	To          []string  // Target 目标
	Date        time.Time // 日期
	ReplyTo     string    // 回复
	Cc          []string  // CarbonCopy 抄送
	Bcc         []string  // BlindCarbonCopy 密送：BCC 栏中的收件人看不到其他 BCC 栏的收件人，同时 To 和 CC 栏中的收件人也看不到 BCC 栏中的收件人
	Subject     string    // 主题
	ContentType string
	Body        string
}

func (i SendOptions) Targets() []string {
	targets := []string{}
	targets = append(targets, i.To...)
	targets = append(targets, i.Cc...)
	targets = append(targets, i.Bcc...)
	return targets
}
func (i SendOptions) Encode() string {
	from := i.From
	if i.FromNick != "" {
		from = fmt.Sprintf("%s<%s>", i.FromNick, i.From)
	}
	heads := map[string][]string{
		MAIL_HEAD_FROM:        {from},
		MAIL_HEAD_TO:          i.To,
		MAIL_HEAD_DATE:        {i.Date.Format(time.RFC1123Z)},
		MAIL_HEAD_REPLYTO:     {i.ReplyTo},
		MAIL_HEAD_CC:          i.Cc,
		MAIL_HEAD_BCC:         i.Bcc,
		MAIL_HEAD_SUBJECT:     {i.Subject},
		MAIL_HEAD_CONTENTTYPE: {i.ContentType},
	}
	var sb strings.Builder
	for k, v := range heads {
		if len(v) == 0 || v[0] == "" {
			continue
		}
		_, _ = sb.WriteString(fmt.Sprintf("%s: %s\r\n", k, strings.Join(v, ";")))
	}
	_, _ = sb.WriteString("\r\n" + i.Body)
	return sb.String()
}
func WithFromNick(v string) SendOption {
	return func(so *SendOptions) {
		so.FromNick = v
	}
}

func WithFrom(v string) SendOption {
	return func(so *SendOptions) {
		so.From = v
	}
}

func WithTo(vs ...string) SendOption {
	return func(so *SendOptions) {
		so.To = append(so.To, vs...)
	}
}

func WithReplyTo(v string) SendOption {
	return func(so *SendOptions) {
		so.ReplyTo = v
	}
}

func WithDate(v time.Time) SendOption {
	return func(so *SendOptions) {
		so.Date = v
	}
}

func WithCc(vs ...string) SendOption {
	return func(so *SendOptions) {
		so.Cc = append(so.Cc, vs...)
	}
}

func WithBcc(vs ...string) SendOption {
	return func(so *SendOptions) {
		so.Bcc = append(so.Bcc, vs...)
	}
}

func WithSubject(v string) SendOption {
	return func(so *SendOptions) {
		so.Subject = v
	}
}

func WithHtmlBody(v string) SendOption {
	return func(so *SendOptions) {
		so.ContentType = MAIL_CONTENTTYPE_HTML
		so.Body = v
	}
}

func WithBody(v string) SendOption {
	return func(so *SendOptions) {
		so.Body = v
	}
}
