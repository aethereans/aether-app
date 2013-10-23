"""
    This file decribes the Aether protocol as headers, to describe what goes in an what goes out.
    The particular implementations of the commands are in aetherProtocol.
"""
from twisted.protocols import amp


class Handshake(amp.Command):
    arguments = [('NodeId', amp.Unicode()), ('ListeningPort', amp.Integer()), ('ProtocolVersion', amp.Float())]
    response = [('NodeId', amp.Unicode()), ('ListeningPort', amp.Integer()), ('ProtocolVersion', amp.Float())]

class RequestHeaders(amp.Command):
    arguments = [('LastSyncTimestamp', amp.Unicode()), ('Languages', amp.ListOf(amp.Unicode()))]
    # This is where you send back your remote timestamp (1st time)
    response = []

class ReceiveHeaders(amp.Command):
    arguments = [('PositiveHeaders', amp.ListOf(amp.Unicode())),
                ('NeutralHeaders', amp.ListOf(amp.Unicode())),
                ('NegativeHeaders', amp.ListOf(amp.Unicode())),
                ('TopicHeaders', amp.ListOf(amp.Unicode())),
                ('TotalNumberOfPackets', amp.Integer()),
                ('CurrentPacketNo', amp.Integer())]
    response = []

class RequestPost(amp.Command):
    arguments = [('PostFingerprint', amp.Unicode())]
    response = [('Post', amp.Unicode())]

class RequestNodes(amp.Command):
    arguments = [('LastSyncTimestamp', amp.Unicode())] # This is where you send back your remote timestamp (2nd time)
    response = []

class ReceiveNodes(amp.Command):
    arguments = [('Nodes', amp.ListOf(amp.Unicode())),
                ('TotalNumberOfPackets', amp.Integer()),
                ('CurrentPacketNo', amp.Integer())]
    response = []

class SyncTimestamps(amp.Command):
    arguments = [('NewSyncTimestamp', amp.Unicode())] # just send your damn timestamp when YOU're DONE.
    response = []

class KillConnection(amp.Command):
    arguments = []
    response = []

