package smtp_server

import (
	"anemone_backend-microservices/internal/mail/config"
	"anemone_backend-microservices/internal/mail/model"
	"anemone_backend-microservices/internal/mail/repository"
	"encoding/base64"
	"errors"
	"io"
	"log"
	"mime"
	"mime/multipart"
	"mime/quotedprintable"
	"net/mail"
	"regexp"
	"strings"
	"time"

	"github.com/emersion/go-smtp"
	"github.com/microcosm-cc/bluemonday"
	"maps"
)

type Server struct {
	cfg  *config.Config
	repo *repository.Repository
}

func NewServer(cfg *config.Config, repo *repository.Repository) *Server {
	return &Server{
		cfg:  cfg,
		repo: repo,
	}
}

func (s *Server) NewSession(c *smtp.Conn) (smtp.Session, error) {
	return &Session{
		repo:   s.repo,
		domain: s.cfg.DomainName,
	}, nil
}

func (s *Server) Start() {
	srv := smtp.NewServer(s)

	srv.Addr = ":" + s.cfg.SMTPPort
	srv.Domain = s.cfg.DomainName
	srv.ReadTimeout = 10 * time.Second
	srv.WriteTimeout = 10 * time.Second
	srv.MaxMessageBytes = 1024 * 1024
	srv.MaxRecipients = 50
	srv.AllowInsecureAuth = true

	log.Printf("INFO: Starting SMTP server at %s for domain %s", srv.Addr, srv.Domain)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("FATAL: Failed to start SMTP server: %v", err)
	}
}

type Session struct {
	repo      *repository.Repository
	domain    string
	from      string
	rcptTo    []string
	addressID int
	rawData   []byte
}

func (s *Session) Mail(from string, opts *smtp.MailOptions) error {
	s.from = from
	return nil
}

func (s *Session) Rcpt(to string, opts *smtp.RcptOptions) error {
	if !strings.HasSuffix(to, "@"+s.domain) {
		return errors.New("invalid recipient domain")
	}

	addr, err := s.repo.FindAddressByString(to)
	if err != nil {
		return errors.New("address does not exist")
	}

	s.rcptTo = append(s.rcptTo, to)
	s.addressID = addr.ID
	return nil
}

func (s *Session) Reset() {
	s.from = ""
	s.rcptTo = nil
	s.addressID = 0
	s.rawData = nil
}

func (s *Session) Logout() error {
	return nil
}

func (s *Session) Data(r io.Reader) error {
	if len(s.rcptTo) == 0 {
		return errors.New("no recipients")
	}

	raw, err := io.ReadAll(r)
	if err != nil {
		return err
	}
	s.rawData = raw

	msg, err := mail.ReadMessage(strings.NewReader(string(raw)))
	if err != nil {
		log.Printf("SMTP DATA: could not read message: %v", err)
		return err
	}

	htmlBody, err := extractAndSanitizeHTML(msg)
	if err != nil {
		log.Printf("SMTP DATA: could not extract/sanitize HTML: %v", err)
		return errors.New("failed to process message body")
	}

	email := &model.Email{
		AddressID:  s.addressID,
		Sender:     s.from,
		Recipients: s.rcptTo,
		Subject:    msg.Header.Get("Subject"),
		Body:       htmlBody,
		RawData:    raw,
	}

	if err := s.repo.SaveEmail(email); err != nil {
		log.Printf("SMTP DATA: failed to save email for address ID %d: %v", s.addressID, err)
		return errors.New("internal server error")
	}

	log.Printf("SMTP DATA: saved email for %s", strings.Join(s.rcptTo, ", "))
	return nil
}

func applyDecoding(r io.Reader, headers mail.Header) io.Reader {
	cte := strings.ToLower(headers.Get("Content-Transfer-Encoding"))

	switch cte {
	case "base64":
		return base64.NewDecoder(base64.StdEncoding, r)
	case "quoted-printable":
		return quotedprintable.NewReader(r)
	default:
		return r
	}
}

func extractAndSanitizeHTML(msg *mail.Message) (string, error) {
	p := bluemonday.UGCPolicy()

	styleRe := regexp.MustCompile(`^[a-zA-Z0-9\s\:\;\#\(\)\-\,\.%]*$`)
	p.AllowAttrs("style").Matching(styleRe).OnElements("p", "span", "div", "td", "th")

	contentType := msg.Header.Get("Content-Type")
	mediaType, params, err := mime.ParseMediaType(contentType)

	bodyReader := applyDecoding(msg.Body, msg.Header)

	if err != nil || (!strings.HasPrefix(mediaType, "multipart/") && mediaType != "text/html") {
		rawBody, _ := io.ReadAll(bodyReader)
		return p.Sanitize(string(rawBody)), nil
	}

	if strings.HasPrefix(mediaType, "multipart/") {
		mr := multipart.NewReader(bodyReader, params["boundary"])
		for {
			part, err := mr.NextPart()
			if err == io.EOF {
				break
			}
			if err != nil {
				continue
			}

			partMsg := &mail.Message{
				Header: make(mail.Header),
				Body:   part,
			}
			maps.Copy(partMsg.Header, part.Header)

			html, err := extractAndSanitizeHTML(partMsg)
			if err == nil && html != "" {
				return html, nil
			}
		}
	}

	if mediaType == "text/html" {
		rawBody, err := io.ReadAll(bodyReader)
		if err != nil {
			return "", err
		}
		return p.Sanitize(string(rawBody)), nil
	}

	return "", nil
}
