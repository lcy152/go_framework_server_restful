package db

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Session struct {
	sess mongo.Session
	ctx  context.Context
}

func (database *Database) StartSession() (*Session, mongo.SessionContext, error) {
	sess, err := database.Client.StartSession(options.Session())
	if err != nil {
		return nil, nil, err
	}

	var session Session
	session.sess = sess
	session.ctx = context.Background()

	var sctx mongo.SessionContext
	err = mongo.WithSession(session.ctx, session.sess, func(sessionContext mongo.SessionContext) error {
		sctx = sessionContext
		return nil
	})
	if err != nil {
		return nil, nil, err
	}
	return &session, sctx, nil
}

func (database *Database) EndSession(session *Session) {
	if session == nil {
		return
	}
	session.sess.EndSession(session.ctx)
}
