# -*- coding: utf-8 -*-

import hashlib
from datetime import datetime
from os import mkdir

from sqlalchemy import  Column, Integer, String, DateTime, Float, ForeignKey, Boolean, Unicode, UnicodeText
from sqlalchemy.orm import relationship
from sqlalchemy.ext.declarative import declarative_base
from sqlalchemy import create_engine
from sqlalchemy.orm import sessionmaker

from globals import basedir, nodeid, newborn, resetted, nuked, FROZEN

#aetherEngine = create_engine('sqlite:///' + basedir + 'Database/aether.db?check_same_thread=False', connect_args={'timeout': 20})
# The modification above allows multithreaded access, however
# it's not thread safe! do not use without implementing queuing or mutex.

aetherEngine = create_engine('sqlite:///' + basedir + 'Database/aether.db', connect_args={'timeout': 20})
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
    PostFingerprint = Column(String, index=True)
    Subject = Column(Unicode)
    Body = Column(UnicodeText)
    OwnerFingerprint = Column(String)
    OwnerUsername = Column(Unicode)
    CreationDate = Column(DateTime)
    ParentPostFingerprint = Column(String, index=True)
    ProtocolVersion = Column(Float)
    Language = Column(String)
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
    # Locally set
    LocallyCreated = Column(Boolean, default=False)
    LastVoteDate = Column(DateTime, index=True) #the creation date.
    # Identifiers
    EphemeralConnectionId = Column(Integer, default=None, index=True)

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
        self.UpvoteCount = kwargs['UpvoteCount'] if 'UpvoteCount' in kwargs else None
        self.DownvoteCount = kwargs['DownvoteCount'] if 'DownvoteCount' in kwargs else None
        self.NeutralCount = kwargs['NeutralCount'] if 'NeutralCount' in kwargs else None
        self.ReplyCount = kwargs['ReplyCount'] if 'ReplyCount' in kwargs else None
        self.Upvoted = kwargs['Upvoted'] if 'Upvoted' in kwargs else None
        self.Downvoted = kwargs['Downvoted'] if 'Downvoted' in kwargs else None
        self.Neutral = kwargs['Neutral'] if 'Neutral' in kwargs else None
        self.Saved = kwargs['Saved'] if 'Saved' in kwargs else None
        self.LastVoteDate = datetime.utcnow()
        self.LocallyCreated = kwargs['LocallyCreated'] if 'LocallyCreated' in kwargs else None
        self.IsReply = kwargs['IsReply'] if 'IsReply' in kwargs else None
        self.EphemeralConnectionId = kwargs['EphemeralConnectionId'] if 'EphemeralConnectionId' in kwargs else None

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

class User(Base):
    __tablename__ = 'users'

    ID = Column(Integer, primary_key=True)
    Fingerprint = Column(String)
    Username = Column(Unicode)
    PublicKey = Column(String)
    # Stuff below are metadata. They are not to be sent across aether but to be fetched.
    FirstName = Column(Unicode)
    LastName = Column(Unicode)
    Mail = Column(Unicode)
    Website = Column(Unicode)
    UserNo = Column(Integer)
    LastRetrieved = Column(DateTime)

class Node(Base):
    __tablename__ = 'nodes'

    ID = Column(Integer, primary_key=True)
    NodeId = Column(String, index=True)
    LastConnectedIP = Column(String)
    LastConnectedPort = Column(Integer)
    LastConnectedDate = Column(DateTime, index=True)
    LastRetrievedIP = Column(String)
    LastRetrievedPort = Column(Integer)
    LastRetrievedDate = Column(DateTime, index=True)
    LastSyncTimestamp = Column(DateTime) # This timestamp arrives from the remote node node.
    Votes = relationship('Vote', back_populates='Node')

    def asDict(self):
       return {c.name: getattr(self, c.name) for c in self.__table__.columns}

class Vote(Base):
    __tablename__ = 'votes'

    ID = Column(Integer, primary_key=True)
    Direction = Column(Integer)
    postheader_id = Column(Integer, ForeignKey('postheaders.ID'))
    node_id = Column(Integer, ForeignKey('nodes.ID'))
    # These relationships below do not work with foreign keys.
    # These are SQLAlchemy constructs and not real database columns.
    PostHeader = relationship('PostHeader', back_populates='Votes')
    Node = relationship('Node', back_populates='Votes')

class PostHeader(Base):
    __tablename__ = 'postheaders'
    ID = Column(Integer, primary_key=True)
    PostFingerprint = Column(String, index=True)
    ParentPostFingerprint = Column(String, index=True)
    Language = Column(String)
    Dirty = Column(Boolean, default=False)
    Votes = relationship('Vote', back_populates='PostHeader')

    def asDict(self):
       return {c.name: getattr(self, c.name) for c in self.__table__.columns}

if newborn:
    session = Session()
    if FROZEN:
        try:
            mkdir(basedir + 'Database')
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

