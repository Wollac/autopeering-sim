package main

import (
	"log"
	"math/rand"
	"time"

	"github.com/wollac/autopeering/peer"
	"github.com/wollac/autopeering/salt"
	"go.uber.org/zap"
)

type simPeer struct {
	local *peer.Local
	peer  *peer.Peer
	db    *peer.DB
	log   *zap.SugaredLogger
	rand  *rand.Rand // random number generator
}

func newPeer(name string, i uint16) simPeer {
	var l *zap.Logger
	var err error
	if name == "1" {
		l, err = zap.NewDevelopment()
	} else {
		l, err = zap.NewDevelopment() //zap.NewProduction()
	}
	if err != nil {
		log.Fatalf("cannot initialize logger: %v", err)
	}
	logger := l.Sugar()
	log := logger.Named(name)
	priv, _ := peer.GeneratePrivateKey()
	db := peer.NewMapDB(log.Named("db"))
	local := peer.NewLocal(priv, db)
	s, _ := salt.NewSalt(time.Duration(i) * 1 * time.Second)
	local.SetPrivateSalt(s)
	s, _ = salt.NewSalt(time.Duration(i) * 1 * time.Second)
	local.SetPublicSalt(s)
	p := peer.NewPeer(local.PublicKey(), name)
	return simPeer{local, p, db, log, rand.New(rand.NewSource(time.Now().UnixNano()))}
}
