# -*- coding: utf-8 -*-

import hashlib
from datetime import datetime
from os import mkdir

from sqlalchemy import  Column, Integer, String, DateTime, Float, ForeignKey, Boolean, Unicode, UnicodeText
from sqlalchemy.orm import relationship
from sqlalchemy.ext.declarative import declarative_base
from sqlalchemy import create_engine
from sqlalchemy.orm import sessionmaker
from sqlalchemy.pool import SingletonThreadPool


from globals import basedir, nodeid, newborn, resetted, nuked, FROZEN, PLATFORM, profiledir

#aetherEngine = create_engine('sqlite:///' + basedir + 'Database/aether.db?check_same_thread=False', connect_args={'timeout': 20})
# The modification above allows multithreaded access, however
# it's not thread safe! do not use without implementing queuing or mutex.

if PLATFORM is 'WIN':
    # aetherEngine = create_engine('sqlite:///' + basedir + 'Database/aether.db?check_same_thread=False',
    #                              poolclass=SingletonThreadPool,
    #                              use_threadlocal=True)
    aetherEngine = create_engine('sqlite:///' + profiledir + 'Database/aether.db')
else:
    aetherEngine = create_engine('sqlite:///' + profiledir + 'Database/aether.db')

    #aetherEngine = create_engine('mysql://aether:12345@localhost/aether', encoding='UTF-8', pool_size=50, max_overflow=100)

Base = declarative_base(bind=aetherEngine)
Session = sessionmaker(bind=aetherEngine)

def unicode_(s):
    # I am using this below to coerce None type into empty string and to convert to / keep stuff
    # in unicode. This is important because if I did dumb conversion (discarding non-ascii
    # characters there would be possibility that two arabic posts have the same fingerprint if they
    # line up in terms of time.
    return '.EMPTY.' if s is None else unicode(s)

class Post(Base):
    __tablename__ = 'posts'

    ID = Column(Integer, primary_key=True)
    PostFingerprint = Column(String(64), index=True)
    Subject = Column(Unicode(255))
    Body = Column(UnicodeText)
    OwnerFingerprint = Column(String(64))
    OwnerUsername = Column(Unicode(255))
    CreationDate = Column(DateTime)
    ParentPostFingerprint = Column(String(64), index=True)
    ProtocolVersion = Column(Float)
    Language = Column(String(255))
    # Counters
    UpvoteCount = Column(Integer, default=0, index=True)
    DownvoteCount = Column(Integer, default=0)
    NeutralCount = Column(Integer, default=0)
    ReplyCount = Column(Integer, default=0)
    # Flags
    Upvoted = Column(Boolean, default=False)
    Downvoted = Column(Boolean, default=False)
    Neutral = Column(Boolean, default=False)
    Saved = Column(Boolean, default=False)
    IsReply = Column(Boolean, default=False, index=True)
    Dirty = Column(Boolean, default=True, index=True)
    # Locally set
    LocallyCreated = Column(Boolean, default=False)
    LastVoteDate = Column(DateTime, index=True) #the creation date.
    RankScore = Column(Float, default=0)

    def asDict(self):
       return {c.name: getattr(self, c.name) for c in self.__table__.columns}

    def __init__(self, **kwargs):

        self.Subject = kwargs['Subject'] if 'Subject' in kwargs else None
        self.Body = kwargs['Body'] if 'Body' in kwargs else None
        self.OwnerFingerprint = kwargs['OwnerFingerprint'] if 'OwnerFingerprint' in kwargs else None
        self.OwnerUsername = kwargs['OwnerUsername'] if 'OwnerUsername' in kwargs else None

        self.ParentPostFingerprint = kwargs['ParentPostFingerprint'] if 'ParentPostFingerprint' \
            in kwargs else None
        #self.ProtocolVersion = kwargs['ProtocolVersion'] if 'ProtocolVersion' in kwargs else None # I probably don't need this.
        self.Language = kwargs['Language'] if 'Language' in kwargs else None
        self.UpvoteCount = kwargs['UpvoteCount'] if 'UpvoteCount' in kwargs else 0
        self.DownvoteCount = kwargs['DownvoteCount'] if 'DownvoteCount' in kwargs else 0
        self.NeutralCount = kwargs['NeutralCount'] if 'NeutralCount' in kwargs else 0
        self.ReplyCount = kwargs['ReplyCount'] if 'ReplyCount' in kwargs else 0
        self.Upvoted = kwargs['Upvoted'] if 'Upvoted' in kwargs else None
        self.Downvoted = kwargs['Downvoted'] if 'Downvoted' in kwargs else None
        self.Neutral = kwargs['Neutral'] if 'Neutral' in kwargs else None
        self.Saved = kwargs['Saved'] if 'Saved' in kwargs else None
        self.LastVoteDate = datetime.utcnow()
        self.LocallyCreated = kwargs['LocallyCreated'] if 'LocallyCreated' in kwargs else None
        self.IsReply = kwargs['IsReply'] if 'IsReply' in kwargs else None

        if 'Body' and 'OwnerUsername' and 'OwnerFingerprint' in kwargs is None:
            # If it is a topic
            self.concatInput = unicode_(self.Subject)
            self.CreationDate = None
        else:
            # If it is a Subject or Post
            self.CreationDate = datetime.utcnow().replace(microsecond=0) if 'CreationDate' not in kwargs else kwargs[
                'CreationDate'].replace(microsecond=0)
            self.concatInput = unicode_(self.Subject) \
                               + unicode_(self.Body) \
                               + unicode_(self.CreationDate) \
                               + unicode_(self.ProtocolVersion) \
                               + unicode_(self.OwnerUsername) \
                               + unicode_(self.OwnerFingerprint) \
                               + unicode_(self.ParentPostFingerprint)

        self.PostFingerprint = hashlib.sha256(
                self.concatInput.encode('utf-8')
            ).hexdigest() if 'PostFingerprint' not in kwargs else kwargs['PostFingerprint']

class Node(Base):
    __tablename__ = 'nodes'

    ID = Column(Integer, primary_key=True)
    NodeId = Column(String(64), index=True)
    LastConnectedIP = Column(String(15))
    LastConnectedPort = Column(Integer)
    LastConnectedDate = Column(DateTime, index=True)
    LastRetrievedIP = Column(String(15))
    LastRetrievedPort = Column(String(5))
    LastRetrievedDate = Column(DateTime, index=True)
    LastSyncTimestamp = Column(DateTime) # This timestamp arrives from the remote node node.

    def asDict(self):
       return {c.name: getattr(self, c.name) for c in self.__table__.columns}

class Vote(Base):
    __tablename__ = 'votes'

    ID = Column(Integer, primary_key=True)
    Direction = Column(Integer)
    TargetPostFingerprint = Column(String(64), index=True)
    NodeId = Column(String(64), index=True)

    def asDict(self):
       return {c.name: getattr(self, c.name) for c in self.__table__.columns}

class PostHeader(Base):
    __tablename__ = 'postheaders'
    ID = Column(Integer, primary_key=True)
    PostFingerprint = Column(String(64), index=True)
    ParentPostFingerprint = Column(String(64), index=True)
    Language = Column(String(255))
    #Dirty = Column(Boolean, default=False)

    def asDict(self):
       return {c.name: getattr(self, c.name) for c in self.__table__.columns}

if newborn:
    session = Session()
    if FROZEN:
        try:
            mkdir(profiledir + 'Database')
        except: pass
    Base.metadata.create_all(aetherEngine)
    session.add(Node(NodeId=nodeid, LastConnectedIP='LOCAL'))
    session.commit()
    session.close()

if resetted:
    session = Session()
    node = session.query(Node).filter(Node.LastConnectedIP == 'LOCAL').one()
    node.NodeId = nodeid
    session.add(node)
    session.commit()
    session.close()

if nuked:
    session = Session()
    Base.metadata.create_all(aetherEngine)
    session.add(Node(NodeId=nodeid, LastConnectedIP='LOCAL'))
    session.commit()
    session.close()

