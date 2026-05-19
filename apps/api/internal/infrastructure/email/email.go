package email

import (
	"context"
	"fmt"
	"html/template"
	"strings"

	"github.com/rs/zerolog"
)

type Service interface {
	SendVerification(ctx context.Context, to, name, link string) error
	SendPasswordReset(ctx context.Context, to, name, link string) error
}

type LogService struct {
	log zerolog.Logger
}

func NewLogService(log zerolog.Logger) *LogService {
	return &LogService{log: log.With().Str("component", "email").Logger()}
}

func (s *LogService) SendVerification(ctx context.Context, to, name, link string) error {
	subject, html, text := verificationEmail(name, link)
	s.log.Info().
		Str("to", to).
		Str("subject", subject).
		Str("html", html).
		Str("text", text).
		Msg("verification email queued")
	return nil
}

func (s *LogService) SendPasswordReset(ctx context.Context, to, name, link string) error {
	subject, html, text := passwordResetEmail(name, link)
	s.log.Info().
		Str("to", to).
		Str("subject", subject).
		Str("html", html).
		Str("text", text).
		Msg("password reset email queued")
	return nil
}

func verificationEmail(name, link string) (string, string, string) {
	subject := "Verify your Meetext workspace"
	html := renderEmail("Verify your email", greeting(name), "Your meeting intelligence workspace is ready. Confirm this email to protect your account and keep client conversations secure.", "Verify email", link)
	text := fmt.Sprintf("Verify your Meetext account: %s", link)
	return subject, html, text
}

func passwordResetEmail(name, link string) (string, string, string) {
	subject := "Reset your Meetext password"
	html := renderEmail("Reset your password", greeting(name), "Use this secure link to choose a new password. The link expires soon, and your old sessions will be retired after reset.", "Reset password", link)
	text := fmt.Sprintf("Reset your Meetext password: %s", link)
	return subject, html, text
}

func greeting(name string) string {
	first := strings.Fields(strings.TrimSpace(name))
	if len(first) == 0 {
		return "Hi,"
	}
	return "Hi " + first[0] + ","
}

func renderEmail(title, greeting, body, cta, link string) string {
	escapedTitle := template.HTMLEscapeString(title)
	escapedGreeting := template.HTMLEscapeString(greeting)
	escapedBody := template.HTMLEscapeString(body)
	escapedCTA := template.HTMLEscapeString(cta)
	escapedLink := template.HTMLEscapeString(link)
	return fmt.Sprintf(`<!doctype html>
<html>
  <body style="margin:0;background:#f6f7fb;font-family:Inter,Arial,sans-serif;color:#18181b;">
    <table role="presentation" width="100%%" cellpadding="0" cellspacing="0" style="padding:32px 16px;background:#f6f7fb;">
      <tr>
        <td align="center">
          <table role="presentation" width="100%%" cellpadding="0" cellspacing="0" style="max-width:560px;background:#ffffff;border:1px solid #e4e4e7;border-radius:14px;overflow:hidden;">
            <tr><td style="padding:28px 32px 12px;font-size:13px;font-weight:700;letter-spacing:.14em;text-transform:uppercase;color:#52525b;">Meetext</td></tr>
            <tr><td style="padding:0 32px 8px;font-size:28px;line-height:1.2;font-weight:700;color:#09090b;">%s</td></tr>
            <tr><td style="padding:8px 32px 0;font-size:16px;line-height:1.7;color:#3f3f46;">%s</td></tr>
            <tr><td style="padding:8px 32px 26px;font-size:15px;line-height:1.7;color:#52525b;">%s</td></tr>
            <tr><td style="padding:0 32px 30px;"><a href="%s" style="display:inline-block;background:#09090b;color:#ffffff;text-decoration:none;border-radius:8px;padding:13px 18px;font-size:14px;font-weight:700;">%s</a></td></tr>
            <tr><td style="padding:22px 32px;background:#fafafa;border-top:1px solid #e4e4e7;font-size:12px;line-height:1.6;color:#71717a;">If the button does not work, paste this link into your browser:<br><span style="word-break:break-all;color:#3f3f46;">%s</span></td></tr>
          </table>
        </td>
      </tr>
    </table>
  </body>
</html>`, escapedTitle, escapedGreeting, escapedBody, escapedLink, escapedCTA, escapedLink)
}
