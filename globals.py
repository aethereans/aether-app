# This file should not import any aether libraries!
# For some reason, if I import stuff here from main the app breaks. The frontend will just show a spinning beachball.
# But imports to all other places work fine.
from __future__ import print_function
import sys
import os.path
from os import mkdir
import cPickle as pickle
import hashlib, random
import ujson
import subprocess
from shutil import copy
from twisted.internet.ssl import ContextFactory, ClientContextFactory
from OpenSSL import SSL, crypto
import os

debugEnabled = True #FIXME

if not debugEnabled:
    def print(*a, **kwargs):
        pass
    def cprint(text, color=None, on_color=None, attrs=None, **kwargs):
        pass

# Current Running Platform: A lot of things depend on this.. Anything about the filesystem, adding to boot etc.

if sys.platform == 'darwin':
    PLATFORM = 'OSX'
elif sys.platform.startswith('win'):
    PLATFORM = 'WIN'
elif sys.platform.startswith('linux'):
    PLATFORM = 'LNX'
else:
    PLATFORM = 'UNKNOWN'
    raise Exception('AetherError: OS Type can not be determined.')

##print('Aether is running on %s' %PLATFORM)

# Set this to true if the app received the open signal at boot, so I can hide the splash and the UI.as

appStartedAtBoot = False

# Application version.

appVersion = 110

#Protocol Version of this version of the app. This is separate from app version.

protocolVersion = 100

# If the application is paused.

appIsPaused = False

# These are packet counts (how many items go into one packet) for nodes and headers used in aetherProtocol.
# These aren't supposed to be user editable. This is something inherent in the protocol and dictated by AMP's 65536 byte
# limit per packet.

nodePacketCount = 10
headerPacketCount = 10

def getRandomOpenPort():
    # Binding to port 0 lets the OS give you an open port.
    import socket
    s = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
    s.bind(("",0))
    s.listen(1)
    port = s.getsockname()[1]
    s.close()
    return port

# These define application status, either running on command line, or packaged into an app bundle.
if getattr(sys, 'frozen', False): # If this is a frozen PyInstaller bundle.
    FROZEN = True
    if PLATFORM == 'OSX':
        basedir = sys._MEIPASS+'/'
        profiledir = os.path.expanduser('~/Library/Application Support/Aether/')
        baseurl = basedir
    elif PLATFORM == 'WIN':
        profiledir = os.environ['ALLUSERSPROFILE']+'\\Aether\\'
        basedir = sys._MEIPASS+'/'
        baseurl = basedir
    elif PLATFORM == 'LNX':
        basedir = sys._MEIPASS+'/'
        profiledir = sys._MEIPASS+'/'
        baseurl = basedir
else: # If running on command line.
    FROZEN = False
    if PLATFORM == 'OSX':
        basedir = ''
        profiledir = os.path.expanduser('~/Library/Application Support/Aether/')
        baseurl = os.path.dirname(__file__)+'/'
    if PLATFORM == 'LNX':
        basedir = ''
        profiledir = ''
        baseurl = os.path.dirname(__file__)+'/'
    elif PLATFORM == 'WIN':
        basedir =  ''
        profiledir = os.environ['ALLUSERSPROFILE']+'\\Aether\\'
        baseurl = os.path.dirname(__file__)+'/'

try:
    mkdir(profiledir)
except: pass


### THE VALUES BELOW GETS OVERWRITTEN WHEN THE APPLICATION STARTS WITH THE VALUES IN USER PROFILE.

# Maximum simultaneous outbound connection count.
maxOutboundCount = 10
# Maximum simultaneous inbound connections per cycle.
maxInboundCount = 3
# How long in minutes the engine waits until a recently connected node is again released to the pool of potential
# connection candidates.
cooldown = 5




# If database or userSettings does not exist, this is a first run.
newborn = True if not os.path.exists(profiledir + 'Database/aether.db') and \
                   not os.path.exists(profiledir + 'UserProfile/backendSettings.dat') else False

if newborn:
       print('This is the first run of Aether.')

# If database exists but not settings, then it's a restartRun. In that case,
# I preserve the DB, and update the local node information. This is essentially a reset button.
# This allows the user to get a new nodeid, most of all... Should getting a new node id nuke the DB?
# I'm assuming not, for now.
resetted = True if os.path.exists(profiledir + 'Database/aether.db') and \
                 not os.path.exists(profiledir + 'UserProfile/backendSettings.dat') else False

if resetted:
    print('Aether has just been resetted. Getting a new Node ID all settings are returned to defaults.')

# Database does not exist, but settings do.
nuked = True if not os.path.exists(profiledir + 'Database/aether.db') and \
                  os.path.exists(profiledir + 'UserProfile/backendSettings.dat') else False

if nuked:
    print('Aether has just ben nuked. I\'m keeping the settings and Node ID, but creating a new database.')


# This gets the stored config details from the file.
try:
    with open(profiledir + 'UserProfile/backendSettings.dat', 'rb') as f2:
            nodeid = pickle.load(f2)
            enableWebkitInspector = pickle.load(f2)
            aetherListeningPort = pickle.load(f2)
            updateAvailable = pickle.load(f2)
            onboardingComplete = pickle.load(f2)

except:
    # Setting some defaults.
    nodeid = hashlib.sha256(str(random.getrandbits(256))).hexdigest()
    ##print('The client picked for itself the Node ID %s' %nodeid)
    enableWebkitInspector = False
    aetherListeningPort = getRandomOpenPort()
    updateAvailable = False
    onboardingComplete = False
    # Max outbound


    try:
        mkdir(profiledir + 'UserProfile')
        mkdir(profiledir + 'Database')
    except: pass
    with open(profiledir + 'UserProfile/backendSettings.dat', 'wb') as f1:
        pickle.dump(nodeid, f1)
        pickle.dump(enableWebkitInspector, f1)
        pickle.dump(aetherListeningPort, f1)
        pickle.dump(updateAvailable, f1)
        pickle.dump(onboardingComplete, f1)
        pickle.dump(maxOutboundCount, f1)
        pickle.dump(maxInboundCount, f1)
        pickle.dump(cooldown, f1)


# This sets the config values into the file.

# The three items below, they need to be merged into one function at some point. Just make some enum that points to
# the location in the pickle and the rest is fine.

def setUpdateAvailable(updateAvailable):
    # Do not forget, these arguments are set when the function is initialised. So if there is more than one value change
    # Through the entire app lifetime, the second change is going to overwrite the the first change because at the
    # point of second change the default value in the key.. Fixed this. This is far safer.
    with open(profiledir + 'UserProfile/backendSettings.dat', 'wb') as f1:
        pickle.dump(nodeid, f1)
        pickle.dump(enableWebkitInspector, f1)
        pickle.dump(aetherListeningPort, f1)
        pickle.dump(updateAvailable, f1)
        pickle.dump(onboardingComplete, f1)

def setListeningPort(aetherListeningPort):
    with open(profiledir + 'UserProfile/backendSettings.dat', 'wb') as f1:
        pickle.dump(nodeid, f1)
        pickle.dump(enableWebkitInspector, f1)
        pickle.dump(aetherListeningPort, f1)
        pickle.dump(updateAvailable, f1)
        pickle.dump(onboardingComplete, f1)

def setOnboardingComplete(onboardingComplete):
    with open(profiledir + 'UserProfile/backendSettings.dat', 'wb') as f1:
        pickle.dump(nodeid, f1)
        pickle.dump(enableWebkitInspector, f1)
        pickle.dump(aetherListeningPort, f1)
        pickle.dump(updateAvailable, f1)
        pickle.dump(onboardingComplete, f1)


# These are values gathered from JSON. These values will be updated when frontend requires a JSON change.
userLanguages = ['English','Turkish','Spanish','French','German','Portuguese','Russian','Chinese','Chineset']

# Currently selected topics.

selectedTopics = []

# Globals API for notifying about JSON updates.

def refreshBackendValuesGatheredFromJson(userProfileInJSON):
    data = ujson.loads(userProfileInJSON)

    global selectedTopics
    selectedTopics = data['selectedTopics']
    ##print('Current selected topics: ', selectedTopics)
    global userLanguages
    userLanguages = data['UserDetails']['UserLanguages']
    ##print('Current user languages is/are %s' %userLanguages)
    startAtBoot = data['UserDetails']['StartAtBoot']
    ##print('Current Start at Boot preference is: %s' %startAtBoot)
    global maxInboundCount
    global maxOutboundCount
    global cooldown
    maxOutboundCount = int(data['UserDetails']['maxOutboundCount'])
    print('max outbound it set:', maxOutboundCount)
    print('max inbound it set:', maxInboundCount)
    print('cooldown it set:', cooldown)
    maxInboundCount = int(data['UserDetails']['maxInboundCount'])
    cooldown = int(data['UserDetails']['cooldown'])

    if startAtBoot and FROZEN and PLATFORM == 'OSX':
        try:
            copy(basedir+'Assets/com.Aether.Aether.plist', os.path.expanduser('~/Library/LaunchAgents/com.Aether.Aether.plist'))
        except:
            mkdir(os.path.expanduser('~/Library/LaunchAgents/'))

            # FIXME
            copy(basedir+'Assets/com.Aether.Aether.plist', os.path.expanduser('~/Library/LaunchAgents/com.Aether.Aether.plist'))
        ##print('Aether in Boot.')
    else:
        try:
            os.remove(os.path.expanduser('~/Library/LaunchAgents/com.Aether.Aether.plist'))
        except: pass
        ##print('Aether out of Boot.')

    return True

if not newborn and not resetted:
    try:
        f = open(profiledir+'UserProfile/UserProfile.json', 'rb')
        refreshBackendValuesGatheredFromJson(f.read())
        f.close()
    except:
        pass

# This is the context factory used to create TLS contexts.

class AetherContextFactory(ContextFactory):
    def getContext(self):
        ctx = SSL.Context(SSL.TLSv1_METHOD)
        # This creates the key pairs and the cert if they do not exist.
        try:
            ctx.use_privatekey_file(profiledir+'UserProfile/priv.pem')
            ctx.use_certificate_file(profiledir+'UserProfile/cert.pem')
        except:
            # We don't have the requirements, so let's create them.
            ##print('This machine doesn\'nt seem to have a keypair and a cert. Creating new, at %s' %datetime.utcnow())
            k = crypto.PKey()
            k.generate_key(crypto.TYPE_RSA, 2048)
            cert = crypto.X509()
            cert.get_subject().countryName = 'XI'
            cert.get_subject().stateOrProvinceName = 'The Internet'
            cert.get_subject().localityName = 'Aether'
            cert.set_serial_number(1000)
            cert.gmtime_adj_notBefore(0)
            cert.gmtime_adj_notAfter(10*365*24*60*60)
            cert.set_issuer(cert.get_subject())
            cert.set_pubkey(k)
            cert.sign(k, 'sha1')
            newCertFile = open(profiledir+'UserProfile/cert.pem', 'wb')
            newCertFile.write(crypto.dump_certificate(crypto.FILETYPE_PEM, cert))
            newCertFile.close()
            newKeyFile = open(profiledir+'UserProfile/priv.pem', 'wb')
            newKeyFile.write(crypto.dump_privatekey(crypto.FILETYPE_PEM, k))
            newKeyFile.close()
            ##print('Key generation finished at %s' %datetime.utcnow())
            ctx.use_privatekey_file(profiledir+'UserProfile/priv.pem')
            ctx.use_certificate_file(profiledir+'UserProfile/cert.pem')

        return ctx

AetherClientContextFactory = ClientContextFactory

aetherClientContextFactoryInstance = AetherClientContextFactory()

aetherContextFactoryInstance = AetherContextFactory()

# The global app quit routine.

def quitApp(reactor):
    # This is (probably) buggy...
    if reactor.threadpool is not None:
        reactor.threadpool.stop()
    reactor.stop()
    sys.exit()

def raiseAndFocusApp():
    raiseWindowCmd = \
    '''osascript<<END
    tell application "Aether"
    activate
    end tell
    END'''
    import signal, time
    p = subprocess.Popen(raiseWindowCmd, shell=True)
    time.sleep(0.01)
    d = os.kill(p.pid, signal.SIGKILL) # so as to stop errors emanating from this.

    # HA. This is funny. For the dear reader who's appalled by what's going on here: There is no way
    # to get raise behaviour on Mac otherwise. Because when I try to raise through normal ways, it tries to
    # raise PyInstaller's bootloader, which is hidden and frozen. If you try to not kill the subprocess,
    # Applescript will take ownership of the process and it will just wait until the end of time.
    # Every time you execute this you would be creating a new process which does absolutely nothing, yet
    # impossible to quit, as they're also frozen because of PyInstaller.


