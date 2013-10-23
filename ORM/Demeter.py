"""
    Demeter is the goddess of maintenance / sustenance.
    This is the interface slow-cycle maintenance / sustenance loop (Persephone) uses for maintenance
    / sustenance calls such as dirty checking for updates.

    This is a deferred API. Everything returns deferred and runs in threads. This API is used mostly by eventLoop.
"""

from twisted.internet import threads
from sqlalchemy.orm import joinedload
from datetime import datetime, timedelta
from termcolor import cprint
from sqlalchemy.orm import exc
import miniupnpc, random

from ORM.models import *
from globals import aetherListeningPort



def checkUPNPStatus(delay):

    def threadFunction():
        u = miniupnpc.UPnP()
        u.discoverdelay = delay #milliseconds
        def mapUPNP():
            print('Checking UPNP Status... with %d second search...' %(delay/1000))
            numberOfFoundDevices = u.discover()
            print('I have found %s devices.' %numberOfFoundDevices)
            if numberOfFoundDevices > 1:
                print('Multiple UPNP devices found. Don\'t know what to do here. '
                      'Adding mapping to only what selectigd() gives to me.')

            if numberOfFoundDevices == 0 :
                print('There are no UPNP devices in the network or '
                      'this computer is directly connected to the Internet.')
                return False
            i = 0
            mappings = []
            try:
                u.selectigd()
                print('This machine has the IP of %s' %(u.externalipaddress()))
                while True:
                    p = u.getgenericportmapping(i)
                    if p==None:
                        break
                    else:
                        mappings.append(p)
                    i+= 1
                for m in mappings:
                    if m[3][:6] == 'Aether' and \
                            m[2][0] == u.lanaddr and \
                            m[2][1] == aetherListeningPort:
                        print('I have Aether already mapped to this computer at port %s.'
                              %aetherListeningPort)
                        return True

                    elif m[3][:6] == 'Aether' and \
                            m[2][0] == u.lanaddr and \
                            m[2][1] != aetherListeningPort:
                        print('I have Aether already mapped to this computer at port %s, but it is not the port '
                              'Aether is currently using.  Removing the old port and adding a mapping to %s.'
                              %(m[2][1], aetherListeningPort))
                        u.deleteportmapping(m[2][1], 'TCP') # Delete the port mapping. It will add again in the end.

                import random
                nameString = 'Aether.' + str(random.randint(10000000,99999999))
                u.addportmapping(aetherListeningPort, 'TCP', u.lanaddr, aetherListeningPort, nameString, '')
                print('I have mapped TCP port %s on the router to the '
                    'same port of this local machine.' %aetherListeningPort)
                return True

            except Exception, ex:
                print(ex)
                return False
        return mapUPNP()
        # I have no idea how long an UPNP mapping lasts.
        # Should I be removing the mapping request?
        # What happens two machines that are in the same network are using Aether? What's the chance of
        # both being assigned to the same port? I need to randomise the port.

    return threads.deferToThread(threadFunction)


def updatePostStatus():
    """
        This call reads all the dirty headers (those that changed or newly created since the last loop run) and

        a) recounts their votes. (this can be made faster by also flagging the votes, so a post with 1000 upvotes won't
        clog the pipe, only flagged votes would be added.

        b) looks decides on what to neutral broadcast. This can also be made faster by having a vote flagging system
        in place: if having all votes as negative or zeroed could make the post go into neutrality (<10) then do the
        check, otherwise skip. For now, it's pretty dumb, just checks.

        c) a deduplication mechanism. If the dirtied header has two posts corresponding it, before anything else,
        destroy one of those. If more than 2, destroy until one remains.

        QUESTION: Should this be a separate thread? Can I end up with a write lock here?
    """
    def threadFunction():
        s = Session()
        dirtiedHeaders = s.query(PostHeader).filter(PostHeader.Dirty == True)\
            .options(joinedload(PostHeader.Votes)).all()

        for h in dirtiedHeaders: # These headers are all headers that got affected by network or user events.
            if not s.query(Post).filter(Post.PostFingerprint == h.PostFingerprint).count():
                continue # If the post of that header did not arrive yet, leave the header as is and skip to the next.
            while True:
                try:
                    post = s.query(Post).filter(Post.PostFingerprint == h.PostFingerprint).one()
                    break # Get out of the while loop. This will last until only one post remains.
                except exc.MultipleResultsFound:
                    cprint('I have found %s posts with the same fingerprint. '
                           'I\'m starting to kill until I see the last man standing. Hasta la vista baby!'
                           %s.query(Post).filter(Post.PostFingerprint == h.PostFingerprint).count(),
                           'white', 'on_red', attrs=['bold'])
                    postToDelete = s.query(Post).filter(Post.PostFingerprint == h.PostFingerprint).first()
                    s.delete(postToDelete)
                    s.commit()
            upvoteCount = 0
            downvoteCount = 0
            neutralCount = 0

            for v in h.Votes:
                if v.Direction == 1: upvoteCount += 1
                elif v.Direction == -1: downvoteCount += 1
                elif v.Direction == 0: neutralCount += 1

            post.UpvoteCount = upvoteCount
            post.DownvoteCount = downvoteCount
            post.NeutralCount = neutralCount

            # This part below checks for neutral broadcast.
            if post.UpvoteCount + post.DownvoteCount + post.NeutralCount < 10 \
                and post.ParentPostFingerprint is not None:
                # If all votes in total is lesser than 10, and if not a topic.
                # (Does this check for the case all are empty?)
                if post.Upvoted == False and post.Downvoted == False:
                    # If it's not upvoted or downvoted.
                    post.Neutral = True
                else:
                    # If it is actually upvoted or downvoted.
                    post.Neutral = False
                    post.LastVoteDate = datetime.utcnow()
            elif post.UpvoteCount + post.DownvoteCount + post.NeutralCount > 10:
                # If the post actually got more than 10 votes.
                post.Neutral = False
                post.NeutralCount -= 1 # Remove your own.
            s.add(post)
            h.Dirty = False
            s.add(h)
            s.commit()
        s.close()

    return threads.deferToThread(threadFunction)

def getNodesToConnect(amount=10, timeout=10):

    def threadFunction():
        s = Session()

        connecteds = s.query(Node).order_by(Node.LastConnectedDate.desc())\
            .filter(Node.LastConnectedDate < (datetime.utcnow() - timedelta(minutes = timeout)))\
            .filter(Node.LastConnectedIP != 'LOCAL')\
            .limit(amount*4)\
            .all()
        randConnectedIndexes = [int(amount*4*random.random()) for i in xrange(amount/2)] # get random x items from 0-25 range.
        finalConnectedArray = []
        if len(connecteds) < amount*4:
            finalConnectedArray = connecteds[:amount]
        else:
            for rNumber in randConnectedIndexes:
                finalConnectedArray.append(connecteds[rNumber])
        # This prevents connection to nodes connected to in the last [timeout] minutes.



        retrieveds = s.query(Node).order_by(Node.LastRetrievedDate.desc())\
            .filter(Node.LastConnectedIP == None)\
            .limit(amount*4).all()
        randRetrievedIndexes = [int(amount*4*random.random()) for i in xrange(amount/2)] # get random x items from 0-25 range.
        finalRetrievedArray = []
        if len(retrieveds) < amount*4:
            finalRetrievedArray = retrieveds[:amount]
        else:
            for rNumber in randRetrievedIndexes:
                finalRetrievedArray.append(retrieveds[rNumber])

        s.close()
        if (len(connecteds) + len(retrieveds) < amount):
            print('All Connected Nodes: %s' %connecteds)
            print('All Retrieved Nodes: %s' %retrieveds)
            return connecteds + retrieveds # If we don't have enough people, don't filter.
        else:
            print('Randomly Selected Connected Nodes: %s' %finalConnectedArray)
            print('Randomly Selected Retrieved Nodes: %s' %finalRetrievedArray)
            return finalConnectedArray + finalRetrievedArray

    return threads.deferToThread(threadFunction)

def timeit(method):
    import time
    def timed(*args, **kw):
        ts = time.time() *1000
        result = method(*args, **kw)
        te = time.time() *1000

        print '%r (%r, %r) %2.2f ms' % \
              (method.__name__, args, kw, te-ts)
        return result

    return timed

@timeit
def incrementAncestryCommentCount(ephemeralConnectionId):
    def threadFunction():

        def incrementCommentCount(post):
            post.EphemeralConnectionId=None

            if post.OwnerUsername == None and post.OwnerFingerprint == None:
                post.ReplyCount += 1
                s.add(post)
                return True
            elif post.Subject != None:
                # If subject
                post.ReplyCount += 1
                s.add(post)
                try:
                    s.query(Post).filter(Post.PostFingerprint == post.ParentPostFingerprint).one()
                except:
                    while s.query(Post).filter(Post.PostFingerprint == post.ParentPostFingerprint).count() > 1:
                        s.delete(s.query(Post).filter(Post.PostFingerprint == post.ParentPostFingerprint).first())
                        s.commit()
                parentPost = s.query(Post).filter(Post.PostFingerprint == post.ParentPostFingerprint).one()
                return incrementCommentCount(parentPost)
            else:
                #If post
                post.ReplyCount += 1
                s.add(post)
                try:
                    s.query(Post).filter(Post.PostFingerprint == post.ParentPostFingerprint).one()
                except:
                    while s.query(Post).filter(Post.PostFingerprint == post.ParentPostFingerprint).count() > 1:
                        s.delete(s.query(Post).filter(Post.PostFingerprint == post.ParentPostFingerprint).first())
                        s.commit()
                parentPost = s.query(Post).filter(Post.PostFingerprint == post.ParentPostFingerprint).one()
                return incrementCommentCount(parentPost)


        s = Session()
        batch = s.query(Post).filter(Post.EphemeralConnectionId==ephemeralConnectionId).all()
        for post in batch:
            incrementCommentCount(post)
        s.commit()
        s.close()
        return None

    return threads.deferToThread(threadFunction)