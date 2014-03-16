from PyQt5.QtCore import *
import ujson
from datetime import datetime

from ORM import Hermes
from globals import nodeid, newborn, appVersion, updateAvailable, onboardingComplete, PLATFORM
import globals


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

    def __init__(self, hermes):
        QObject.__init__(self)
        self.Hermes = hermes

    # Methods used by front to invoke action in the backend
    @pyqtSlot(str, result=str)
    def getAllPosts(self, fingerprint):
        posts = self.Hermes.fetchAllDescendantPosts(fingerprint)
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
        post = self.Hermes.fetchSinglePost(fingerprint)
        if post.__len__() == 0:
            return "[]"
        else:
            # Post[0] because it's all() and not one() in SQLAlchemy, otherwise when not found it
            #  throws exceptions in my face.
            return ujson.dumps(post[0].asDict(), ensure_ascii=False)

    @pyqtSlot(str, int, result=str)
    def getSubjects(self, fingerprint, daysToSubtract):
        posts = self.Hermes.getSubjects(fingerprint, daysToSubtract)
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
        topicDictsAsJson = self.Hermes.getHomeScreen(numberOfTopics, numberOfSubjects)
        if topicDictsAsJson.__len__() == 0:
            return "[]"
        else:
            return ujson.dumps(topicDictsAsJson)



    @pyqtSlot(str, result=bool)
    def getUser(self, fingerprint):
        user = self.Hermes.fetchUser(fingerprint)
        return user

    @pyqtSlot(str,int, result=str)
    def getSpecificDepthPosts(self, fingerprint, depth):
        posts = self.Hermes.fetchDescendantPosts(fingerprint, depth)
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
        posts = self.Hermes.fetchDirectDescendantPosts(fingerprint)
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
        posts = self.Hermes.fetchDirectDescendantPostsWithTimespan(fingerprint, daysToSubtract)
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
        count = self.Hermes.countDirectDescendantPosts(fingerprint)
        return count

    @pyqtSlot(str, result=int)
    def countAllDescendantPosts(self, fingerprint):
        count = self.Hermes.countAllDescendantPosts(fingerprint)
        return count

    @pyqtSlot(result=int)
    def countUppermostTopics(self):
        return self.Hermes.countTopics()

    @pyqtSlot(result=int)
    def countSubjects(self):
        return self.Hermes.countSubjects()

    @pyqtSlot(result=str)
    def getUppermostTopics(self):
        posts = self.Hermes.fetchUppermostTopics()
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
        post = self.Hermes.getTopmostComment(fingerprint)
        if post is None:
            return "[]"
        else:
            return ujson.dumps(post.asDict(), ensure_ascii=False)

    @pyqtSlot(result=str)
    def getUserProfile(self):
        if self.Hermes.readUserProfile() != '':
            return self.Hermes.readUserProfile()
        else:
            return '{}'

    @pyqtSlot(int, str, result=str)
    def getSubjectCounts(self, days, fingerprintsListAsJSON):
        return ujson.dumps(self.Hermes.multiCountDirectDescendantPostsWithTimespan(days,
                                                                  ujson.loads(
                                                                  fingerprintsListAsJSON)),
                           ensure_ascii=False)

    @pyqtSlot(result=str)
    def getSavedPosts(self):
        posts = self.Hermes.getSavedPosts()
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
        return self.Hermes.countConnectedNodes()

    @pyqtSlot(result=str)
    def getLocallyCreatedPosts(self):
        posts = self.Hermes.getLocallyCreatedPosts()
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
        return self.Hermes.countReplies()

    @pyqtSlot(result=str)
    def getReplies(self):
        posts = self.Hermes.getReplies()
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
        posts = self.Hermes.getUnregisteredUserPosts(name)
        return multiJsonify(posts)


    @pyqtSlot(str, result=str)
    def getParentSubjectOfGivenPost(self, fingerprint):
        return multiJsonify(self.Hermes.getParentSubjectOfGivenPost(fingerprint))

    @pyqtSlot(result=str)
    def getLastConnectionTime(self):
        return ujson.dumps(self.Hermes.getLastConnectionTime())

    # These ones below are write methods.

    @pyqtSlot(str)
    def writeUserProfile(self, userProfileInJSON):
        # I need to do some stuff here to ensure data in json matches the one in profile.
        return self.Hermes.writeUserProfile(userProfileInJSON)

    @pyqtSlot(str, result=str)
    def createTopic(self, topicName):
        return self.Hermes.createTopic(topicName)

    @pyqtSlot(str, str, str, str, str, result=str)
    def createPost(self, postSubject, postText, parentPostFingerprint, ownerUsername, postLanguage):
        return self.Hermes.createPost(postSubject, postText, parentPostFingerprint, ownerUsername, postLanguage)

    @pyqtSlot(str, str, str, str, result=str)
    def createSignedPost(self, postSubject, postText, parentPostFingerprint, ownerFingerprint):
        return self.Hermes.createSignedPost(postSubject, postText, parentPostFingerprint,
                                     ownerFingerprint)

    @pyqtSlot(str, int, result=bool)
    def votePost(self, fingerprint, voteDirection):
        return self.Hermes.votePost(fingerprint, voteDirection)

    @pyqtSlot(str, result=bool)
    def savePost(self, fingerprint):
        return self.Hermes.savePost(fingerprint)

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

    @pyqtSlot(result=str)
    def getOperatingSystem(self):
        return PLATFORM

    # It took me so long to figure out how to pass a variable (above) through. I'm an idiot.

    @pyqtSlot(result=bool)
    def markAllRepliesAsRead(self):
        self.trayIcon.makeIconGoDark()
        return self.Hermes.markAllRepliesAsRead()

    @pyqtSlot(result=bool)
    def markAllSavedsAsNotSaved(self):
        return self.Hermes.markAllSavedsAsNotSaved()

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
        self.Hermes.sendConnectWithIPSignal(ip, port)

    @pyqtSlot(str, result=bool)
    def exportSinglePost(self, post):
        # Here I need to get the post, convert it to dict, and put it into a text file.
        postAsDict = ujson.loads(post)
        if len(postAsDict['Body']) < 25:
            fileName = postAsDict['Body']
        else:
            fileName = postAsDict['Body'][0:25] + '..'
        from PyQt5.QtWidgets import QFileDialog
        dialog = QFileDialog()
        filePath = dialog.getSaveFileName(None, 'Please select a location to save a post.', fileName + '.txt')
        filePath = filePath[0]
        postText = \
u"""%s\n
by %s, at %s\n
Gathered from Aether network. Aether is a distributed network of anonymous forums. Join us at www.getaether.net.
This post is licensed under CC-BY-SA.""" \
        % (postAsDict['Body'], postAsDict['OwnerUsername'], str(datetime.utcfromtimestamp(float(postAsDict['CreationDate']))))
        postText = postText.encode('utf8')
        f = open(filePath, 'wb')
        f.write(postText)
        f.close()
        return True

    @pyqtSlot(str, result=bool)
    def exportAllPosts(self, posts):
        posts = ujson.loads(posts)
        print('posts:', posts)

        fileName = 'Posts compilation at %s' % datetime.now().date()
        from PyQt5.QtWidgets import QFileDialog
        dialog = QFileDialog()
        filePath = dialog.getSaveFileName(None, 'Please select a location to save a post.', fileName + '.txt')
        filePath = filePath[0]
        if not filePath:
            return False

        def produceSinglePostText(postAsDict):
            postText =  \
u"""
%s\n
by %s, at %s\n\n
""" \
            % (postAsDict['Body'], postAsDict['OwnerUsername'], str(datetime.utcfromtimestamp(float(postAsDict['CreationDate']))))
            return postText

        def attachCredits(text):
            creditsText =  \
u"""YOUR SAVED POSTS AT %s
%s
Gathered from Aether network. Aether is a distributed network of anonymous forums. Join us at www.getaether.net.
These posts are licensed under CC-BY-SA.
""" % (datetime.now().date(), text)
            return creditsText

        finalPostText = u""
        for p in posts:
            finalPostText = finalPostText + produceSinglePostText(p)
        f = open(filePath, 'wb')
        f.write(attachCredits(finalPostText).encode('utf8'))# Encode only happens at the end.
        f.close()
        return True





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