package session

import (
	"github.com/google/uuid"
	"github.com/vncommunityleague/kazusa/identity"
	"github.com/vncommunityleague/kazusa/internal"
	"net"
	"net/http"
	"strings"
	"time"
)

const Lifetime = time.Hour * 24

type Session struct {
	ID         uuid.UUID          `json:"id" gorm:"primaryKey;type:uuid;default:uuid_generate_v1()"`
	IdentityID uuid.UUID          `json:"-"`
	Identity   *identity.Identity `json:"identity" gorm:"foreignKey:IdentityID"`

	Token string `json:"-"`

	IssuedAt  time.Time `json:"issued_at"`
	ExpiresAt time.Time `json:"expires_at"`

	IPAddress *string `json:"ip_address"`
	UserAgent *string `json:"user_agent"`
	Location  *string `json:"location"`

	Active bool `json:"active"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func NewActiveSession(r *http.Request, identity *identity.Identity) (*Session, error) {
	s := NewInactiveSession()
	s.Active = true

	err := s.Activate(r, identity, time.Now())
	if err != nil {
		return nil, err
	}

	return s, nil
}

func NewInactiveSession() *Session {
	return &Session{
		Token:  "vcl_session_" + internal.RandomString(64),
		Active: false,
	}
}

func (s *Session) IsActive() bool {
	return s.Active && s.ExpiresAt.After(time.Now())
}

func (s *Session) Refresh() *Session {
	s.ExpiresAt = time.Now().Add(Lifetime).UTC()
	return s
}

func (s *Session) Activate(r *http.Request, i *identity.Identity, issuedAt time.Time) error {
	s.Identity = i
	s.IdentityID = i.ID
	s.Active = true
	s.ExpiresAt = issuedAt.Add(Lifetime)
	s.IssuedAt = issuedAt

	s.setMetadata(r)
	return nil
}

func (s *Session) setMetadata(r *http.Request) {
	if ip := r.Header.Get("CF-Connecting-IP"); len(ip) > 0 {
		s.IPAddress = &ip
	} else if ip = r.Header.Get("X-Forwarded-For"); len(ip) > 0 {
		ip = GetClientIPAddressesWithoutInternalIPs(strings.Split(ip, ","))
		s.IPAddress = &ip
	} else {
		s.IPAddress = &r.RemoteAddr
	}

	agent := r.Header["User-Agent"]
	if len(agent) > 0 {
		var ua = strings.Join(agent, " ")
		s.UserAgent = &ua
	}

	var location []string
	if city := r.Header.Get("CF-IPCity"); len(city) > 0 {
		location = append(location, city)
	}
	if country := r.Header.Get("CF-IPCountry"); len(country) > 0 {
		location = append(location, country)
	}

	var loc = strings.Join(location, ", ")
	s.Location = &loc
}

// GetClientIPAddressesWithoutInternalIPs from ory/kratos
func GetClientIPAddressesWithoutInternalIPs(ipAddresses []string) string {
	var res string

	for i := len(ipAddresses) - 1; i >= 0; i-- {
		ip := strings.TrimSpace(ipAddresses[i])

		if !net.ParseIP(ip).IsPrivate() {
			res = ip
			break
		}
	}

	return res
}
