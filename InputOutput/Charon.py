from PyQt5.QtCore import *
import ujson

from ORM import Hermes
from globals import nodeid, newborn, appVersion, updateAvailable, onboardingComplete
import globals, aetherProtocol

class Charon(QObject):
    """
        This is the second deepest API level, Charon. This is tasked with converting the
        Hermes' return to JSON and traversal through PyQt, Qt and QtWebkit.

        Charon crosses the River Lethe, to bring its passengers across from the world of
        living to the underworld. Anybody trying to cross the Lethe will forget any and all,
        even if he manages to swim he will have no memory of his life on the world of living.

        Charon is the only way to cross the divide between the backend in Python and the frontend in
        Javascript. This is a two way gate which can carry objects, methods and quite a few other
        items, so this is an entire communication API rather than one single object going through.

        Here you can expose a lot of stuff, but you need to mark the return as int, string, float or
        whatever by explicitly declaring the slot. I need to look a bit more into this slots and
        signals mechanism, but it looks fluid enough.
    """

    # Methods used by front to invoke action in the backend
    @pyqtSlot(str, result=str)
    def getAllPosts(self, fingerprint):
        posts = Hermes.fetchAllDescendantPosts(fingerprint)
        returnList = []
        if posts.__len__() == 0:
            return "[]"
        else:
            for post in posts:
                postJSON = ujson.dumps(post[0].asDict(), ensure_ascii=False)
                returnList.append(postJSON)
            return ujson.dumps(returnList, ensure_ascii=False)
    # Two key things, asDict and this is one single post that goes through.
    # I need to get asDict to the base method and I need to get this system to traverse multiple
    # posts.

    @pyqtSlot(str, result=str)
    def getSinglePost(self, fingerprint):
        post = Hermes.fetchSinglePost(fingerprint)
        if post.__len__() == 0:
            return "[]"
        else:
            # Post[0] because it's all() and not one() in SQLAlchemy, otherwise when not found it
            #  throws exceptions in my face.
            return ujson.dumps(post[0].asDict(), ensure_ascii=False)

    @pyqtSlot(str, int, result=str)
    def getSubjects(self, fingerprint, daysToSubtract):
        posts = Hermes.getSubjects(fingerprint, daysToSubtract)
        returnList = []
        if posts.__len__() == 0:
            return "[]"
        else:
            for post in posts:
                # asDict only converts to dict the values in models table. add the rest by hand.
                postDict = post.asDict()
                postDict['FirstPostBody'] = post.FirstPostBody
                postDict['FirstPostOwner'] = post.FirstPostOwner
                postJSON = ujson.dumps(postDict, ensure_ascii=False)
                returnList.append(postJSON)
            return ujson.dumps(returnList, ensure_ascii=False)

    @pyqtSlot(int, int, result=str)
    def getHomeScreen(self, numberOfTopics, numberOfSubjects):
        topicDictsAsJson = Hermes.getHomeScreen(numberOfTopics, numberOfSubjects)
        if topicDictsAsJson.__len__() == 0:
            return "[]"
        else:
            return ujson.dumps(topicDictsAsJson)



    @pyqtSlot(str, result=bool)
    def getUser(self, fingerprint):
        user = Hermes.fetchUser(fingerprint)
        return user

    @pyqtSlot(str,int, result=str)
    def getSpecificDepthPosts(self, fingerprint, depth):
        posts = Hermes.fetchDescendantPosts(fingerprint, depth)
        returnList = []
        if posts.__len__() == 1:
            return "[]"
        else:
            for post in posts:
                postJSON = ujson.dumps(post[0].asDict(), ensure_ascii=False)
                returnList.append(postJSON)
            return ujson.dumps(returnList, ensure_ascii=False)
    @pyqtSlot(str, result=str)
    def getDirectDescendantPosts(self, fingerprint):
        posts = Hermes.fetchDirectDescendantPosts(fingerprint)
        returnList = []
        if posts.__len__() == 0:
            return "[]"
        else:
            for post in posts:
                postJSON = ujson.dumps(post.asDict(), ensure_ascii=False)
                returnList.append(postJSON)
            return ujson.dumps(returnList, ensure_ascii=False)

    @pyqtSlot(str, int, result=str)
    def getDirectDescendantPostsWithTimespan(self, fingerprint, daysToSubtract):
        posts = Hermes.fetchDirectDescendantPostsWithTimespan(fingerprint, daysToSubtract)
        returnList = []
        if posts.__len__() == 0:
            return "[]"
        else:
            for post in posts:
                postJSON = ujson.dumps(post.asDict(), ensure_ascii=False)
                returnList.append(postJSON)
            return ujson.dumps(returnList, ensure_ascii=False)

    @pyqtSlot(str, result=int)
    def countDirectDescendantPosts(self, fingerprint):
        count = Hermes.countDirectDescendantPosts(fingerprint)
        return count

    @pyqtSlot(str, result=int)
    def countAllDescendantPosts(self, fingerprint):
        count = Hermes.countAllDescendantPosts(fingerprint)
        return count

    @pyqtSlot(result=int)
    def countUppermostTopics(self):
        return Hermes.countTopics()

    @pyqtSlot(result=int)
    def countSubjects(self):
        return Hermes.countSubjects()

    @pyqtSlot(result=str)
    def getUppermostTopics(self):
        posts = Hermes.fetchUppermostTopics()
        returnList = []
        if posts.__len__() == 0:
            return "[]"
        else:
            for post in posts:
                postJSON = ujson.dumps(post.asDict(), ensure_ascii=False)
                returnList.append(postJSON)
            return ujson.dumps(returnList, ensure_ascii=False)

    @pyqtSlot(str, result=str)
    def getTopmostComment(self, fingerprint):
        post = Hermes.getTopmostComment(fingerprint)
        if post is None:
            return "[]"
        else:
            return ujson.dumps(post.asDict(), ensure_ascii=False)

    @pyqtSlot(result=str)
    def getUserProfile(self):
        if Hermes.readUserProfile() != '':
            return Hermes.readUserProfile()
        else:
            return '{}'

    @pyqtSlot(int, str, result=str)
    def getSubjectCounts(self, days, fingerprintsListAsJSON):
        return ujson.dumps(Hermes.multiCountDirectDescendantPostsWithTimespan(days,
                                                                  ujson.loads(
                                                                  fingerprintsListAsJSON)),
                           ensure_ascii=False)

    @pyqtSlot(result=str)
    def getSavedPosts(self):
        posts = Hermes.getSavedPosts()
        returnList = []
        if posts.__len__() == 0:
            return "[]"
        else:
            for post in posts:
                postJSON = ujson.dumps(post.asDict(), ensure_ascii=False)
                returnList.append(postJSON)
            return ujson.dumps(returnList, ensure_ascii=False)

    @pyqtSlot(result=int)
    def countConnectedNodes(self):
        return Hermes.countConnectedNodes()

    @pyqtSlot(result=str)
    def getLocallyCreatedPosts(self):
        posts = Hermes.getLocallyCreatedPosts()
        returnList = []
        if posts.__len__() == 0:
            return "[]"
        else:
            for post in posts:
                postJSON = ujson.dumps(post.asDict(), ensure_ascii=False)
                returnList.append(postJSON)
            return ujson.dumps(returnList, ensure_ascii=False)

    @pyqtSlot(result=int)
    def countReplies(self):
        return Hermes.countReplies()

    @pyqtSlot(result=str)
    def getReplies(self):
        posts = Hermes.getReplies()
        returnList = []
        if posts.__len__() == 0:
            return "[]"
        else:
            for post in posts:
                postJSON = ujson.dumps(post.asDict(), ensure_ascii=False)
                returnList.append(postJSON)
            return ujson.dumps(returnList, ensure_ascii=False)

    @pyqtSlot(str, result=str)
    def getUnregisteredUserPosts(self, name):
        posts = Hermes.getUnregisteredUserPosts(name)
        return multiJsonify(posts)


    @pyqtSlot(str, result=str)
    def getParentSubjectOfGivenPost(self, fingerprint):
        return multiJsonify(Hermes.getParentSubjectOfGivenPost(fingerprint))

    @pyqtSlot(result=str)
    def getLastConnectionTime(self):
        return ujson.dumps(Hermes.getLastConnectionTime())

    # These ones below are write methods.

    @pyqtSlot(str)
    def writeUserProfile(self, userProfileInJSON):
        # I need to do some stuff here to ensure data in json matches the one in profile.
        return Hermes.writeUserProfile(userProfileInJSON)

    @pyqtSlot(str, result=str)
    def createTopic(self, topicName):
        return Hermes.createTopic(topicName)

    @pyqtSlot(str, str, str, str, str, result=str)
    def createPost(self, postSubject, postText, parentPostFingerprint, ownerUsername, postLanguage):
        return Hermes.createPost(postSubject, postText, parentPostFingerprint, ownerUsername, postLanguage)

    @pyqtSlot(str, str, str, str, result=str)
    def createSignedPost(self, postSubject, postText, parentPostFingerprint, ownerFingerprint):
        return Hermes.createSignedPost(postSubject, postText, parentPostFingerprint,
                                     ownerFingerprint)

    @pyqtSlot(str, int, result=bool)
    def votePost(self, fingerprint, voteDirection):
        return Hermes.votePost(fingerprint, voteDirection)

    @pyqtSlot(str, result=bool)
    def savePost(self, fingerprint):
        return Hermes.savePost(fingerprint)

    @pyqtSlot(result=str)
    def nodeId(self):
        return nodeid

    @pyqtSlot(result=bool)
    def newborn(self):
        return newborn

    @pyqtSlot(result=bool)
    def onboardingComplete(self):
        return onboardingComplete

    @pyqtSlot(result=str)
    def appVersion(self):
        return str(appVersion)

    @pyqtSlot(result=bool)
    def updateAvailable(self):
        return updateAvailable

    # It took me so long to figure out how to pass a variable (above) through. I'm an idiot.

    @pyqtSlot(result=bool)
    def markAllRepliesAsRead(self):
        return Hermes.markAllRepliesAsRead()

    @pyqtSlot()
    def quitApp(self):
        from twisted.internet import reactor
        from globals import quitApp
        quitApp(reactor)

    @pyqtSlot()
    def setOnboardingComplete(self):
        globals.setOnboardingComplete(True)

    @pyqtSlot(str, str)
    def connectToNodeWithIP(self, ip, port):
        aetherProtocol.connectWithIP(ip, int(port))

def multiJsonify(posts):
    """
        This method automates the default packaging JSON action used all across the API.
        It also has a Javascript equivalent called multiParsify
    """
    returnList = []
    if posts.__len__() == 0:
        return "[]"
    else:
        for post in posts:
            postJSON = ujson.dumps(post.asDict(), ensure_ascii=False)
            returnList.append(postJSON)
        return ujson.dumps(returnList, ensure_ascii=False)