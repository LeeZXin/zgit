package apisession

import "sync"

// memStore 内存
type memStore struct {
	sync.RWMutex
	session     map[string]Session
	userSession map[string]Session
}

func newMemStore() Store {
	return &memStore{
		RWMutex:     sync.RWMutex{},
		session:     make(map[string]Session, 8),
		userSession: make(map[string]Session, 8),
	}
}

func (s *memStore) GetBySessionId(sessionId string) (Session, bool, error) {
	s.RLock()
	defer s.RUnlock()
	ret, b := s.session[sessionId]
	return ret, b, nil
}

func (s *memStore) GetByUserId(userId string) (Session, bool, error) {
	s.RLock()
	defer s.RUnlock()
	ret, b := s.userSession[userId]
	return ret, b, nil
}

func (s *memStore) PutSession(session Session) error {
	s.Lock()
	defer s.Unlock()
	s.session[session.SessionId] = session
	s.userSession[session.UserInfo.UserId] = session
	return nil
}

func (s *memStore) DeleteByUserId(userId string) error {
	s.Lock()
	defer s.Unlock()
	session, b := s.userSession[userId]
	if b {
		delete(s.userSession, userId)
		delete(s.session, session.SessionId)
	}
	return nil
}

func (s *memStore) DeleteBySessionId(sessionId string) error {
	s.Lock()
	defer s.Unlock()
	session, b := s.session[sessionId]
	if b {
		delete(s.userSession, session.UserInfo.UserId)
		delete(s.session, sessionId)
	}
	return nil
}

func (s *memStore) RefreshExpiry(sessionId string, expireAt int64) error {
	s.Lock()
	defer s.Unlock()
	session, b := s.session[sessionId]
	if b {
		session.ExpireAt = expireAt
		s.session[sessionId] = session
		s.userSession[session.UserInfo.UserId] = session
	}
	return nil
}
