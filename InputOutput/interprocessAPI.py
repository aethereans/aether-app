from twisted.protocols import amp

"""
    This is the protocol definition between processes.
    The main stuff going across are signals, such as boot up, kill app or else.
"""

class commit(amp.Command):
    # When something is added by the user, a new topic, post or subject,
    # fire a commit to globalcommitter and fire a set of new connections.
    arguments = [('PostFingerprint', amp.Unicode())]
    response = []

class killApp(amp.Command):
    response = []
    # When the gui exists, it should also take down the main driver.
    # This signal should only be sent from the GUI to.. hm. Don't know if this is actually the case.

class connectWithIP(amp.Command):
    arguments = [('IP', amp.Unicode()), ('Port', amp.Integer())]
    response = []
    # This fires when the user uses the connect button on the settings.

class thereAreReplies(amp.Command):
    arguments = []
    response = []
