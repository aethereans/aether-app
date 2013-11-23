# This file should not import any aether libraries!
# For some reason, if I import stuff here from main the app breaks. The frontend will just show a spinning beachball.
# But imports to all other places work fine.
import sys
import os.path
from os import mkdir
import cPickle as pickle
import hashlib, random
import ujson
import subprocess
from twisted.python import log
from twisted.python.logfile import DailyLogFile
from shutil import copy
from twisted.internet.ssl import ContextFactory, ClientContextFactory
from OpenSSL import SSL, crypto

# Current Running Platform: A lot of things depend on this.. Anything about the filesystem, adding to boot etc.

PLATFORM = 'OSX' # Possible values: OSX, WIN, LNX

# Set this to true if the app received the open signal at boot, so I can hide the splash and the UI.as

appStartedAtBoot = False

# Application version.

appVersion = 100

#Protocol Version of this version of the app. This is separate from app version.

protocolVersion = 100

# If the application is paused.

appIsPaused = False

# The level of logs that will be captured if the logging is enabled.

loggingLevel = 0 # 0: Debug, 1: Warnings 2:Errors (Exceptions will automatically be always logged.)
# Implement this. 0 produces MASSIVE log files.

# These are packet counts (how many items go into one packet) for nodes and headers used in aetherProtocol.

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
    basedir = sys._MEIPASS+'/'
    baseURL = basedir
else: # If running on command line.
    basedir = ''
    baseURL = os.path.dirname(__file__)+'/'
    FROZEN = False


# These items apply to Mercury, the asynchronous database API used by network operations.
maximumRetryCount = 5 # times
retryWaitTime = 5 # seconds

# This applies to the protocol. Maximum allowed time for each connection to complete. This should never be hit, but
# it's an insurance policy if things get borked beyond salvation.

maximumAllowedConnectionTimespan = 5 # minutes

# If database or userSettings does not exist, this is a first run.
newborn = True if not os.path.exists(basedir + 'Database/aether.db') and \
                   not os.path.exists(basedir + 'UserProfile/backendSettings.dat') else False

if newborn:
        print('This is the first run of Aether.')

# If database exists but not settings, then it's a restartRun. In that case,
# I preserve the DB, and update the local node information. This is essentially a reset button.
# This allows the user to get a new nodeid, most of all... Should getting a new node id nuke the DB?
# I'm assuming not, for now.
resetted = True if os.path.exists(basedir + 'Database/aether.db') and \
                 not os.path.exists(basedir + 'UserProfile/backendSettings.dat') else False

if resetted:
    print('Aether has just been resetted. Getting a new Node ID all settings are returned to defaults.')

# Database does not exist, but settings do.
nuked = True if not os.path.exists(basedir + 'Database/aether.db') and \
                  os.path.exists(basedir + 'UserProfile/backendSettings.dat') else False

if nuked:
    print('Aether has just been nuked. I\'m keeping the settings and Node ID, but creating a new database.')


# This gets the stored config details from the file.
try:
    with open(basedir + 'UserProfile/backendSettings.dat', 'rb') as f2:
            nodeid = pickle.load(f2)
            enableWebkitInspector = pickle.load(f2)
            aetherListeningPort = pickle.load(f2)
            updateAvailable = pickle.load(f2)
            onboardingComplete = pickle.load(f2)
except:
    # Setting some defaults.
    nodeid = hashlib.sha256(str(random.getrandbits(256))).hexdigest()
    print('The client picked for itself the Node ID %s' %nodeid)
    enableWebkitInspector = False
    aetherListeningPort = getRandomOpenPort()
    updateAvailable = False
    onboardingComplete = False

    if FROZEN:
        try:
            mkdir(basedir + 'UserProfile')
        except: pass
    with open(basedir + 'UserProfile/backendSettings.dat', 'wb') as f1:
        pickle.dump(nodeid, f1)
        pickle.dump(enableWebkitInspector, f1)
        pickle.dump(aetherListeningPort, f1)
        pickle.dump(updateAvailable, f1)
        pickle.dump(onboardingComplete, f1)


# This sets the config values into the file.

# The three items below, they need to be merged into one function at some point. Just make some enum that points to
# the location in the pickle and the rest is fine.

def setUpdateAvailable(updateAvailable):
    # Do not forget, these arguments are set when the function is initialised. So if there is more than one value change
    # Through the entire app lifetime, the second change is going to overwrite the the first change because at the
    # point of second change the default value in the key.. Fixed this. This is far safer.
    with open(basedir + 'UserProfile/backendSettings.dat', 'wb') as f1:
        pickle.dump(nodeid, f1)
        pickle.dump(enableWebkitInspector, f1)
        pickle.dump(aetherListeningPort, f1)
        pickle.dump(updateAvailable, f1)
        pickle.dump(onboardingComplete, f1)

def setListeningPort(aetherListeningPort):
    with open(basedir + 'UserProfile/backendSettings.dat', 'wb') as f1:
        pickle.dump(nodeid, f1)
        pickle.dump(enableWebkitInspector, f1)
        pickle.dump(aetherListeningPort, f1)
        pickle.dump(updateAvailable, f1)
        pickle.dump(onboardingComplete, f1)

def setOnboardingComplete(onboardingComplete):
    with open(basedir + 'UserProfile/backendSettings.dat', 'wb') as f1:
        pickle.dump(nodeid, f1)
        pickle.dump(enableWebkitInspector, f1)
        pickle.dump(aetherListeningPort, f1)
        pickle.dump(updateAvailable, f1)
        pickle.dump(onboardingComplete, f1)


# These are values gathered from JSON. These values will be updated when frontend requires a JSON change.
userLanguages = []

# Currently selected topics.

selectedTopics = []

# Globals API for notifying about JSON updates.

def refreshBackendValuesGatheredFromJson(userProfileInJSON):
    data = ujson.loads(userProfileInJSON)

    global selectedTopics
    selectedTopics = data['selectedTopics']
    print('Current selected topics: ', selectedTopics)
    global userLanguages
    userLanguages = data['UserDetails']['UserLanguages']
    print('Current user languages is/are %s' %userLanguages)
    startAtBoot = data['UserDetails']['StartAtBoot']
    print('Current Start at Boot preference is: %s' %startAtBoot)
    loggingEnabled = data['UserDetails']['Logging']
    print('Current logging preference is: %s' %loggingEnabled)

    if startAtBoot and FROZEN:
        copy(basedir+'Assets/com.Aether.Aether.plist', os.path.expanduser('~/Library/LaunchAgents/com.Aether.Aether.plist'))
        print('Aether in Boot.')
    else:
        try:
            os.remove('~/Library/LaunchAgents/com.Aether.Aether.plist')
        except: pass
        print('Aether out of Boot.')

    if loggingEnabled:
        print('Logging is ENABLED. Writing to file.')
        #log.startLogging(DailyLogFile.fromFullPath(basedir+'Logs/network.log')) # Use this in real deployment.
        # FIXME logging before launch
        log.startLogging(sys.stdout) # Use this in debug builds.
        print('Logging to file started. NodeId: %s' %nodeid)
        # This logs everything raised in the app. Including
        # everything that is printed to the stdout.
    else:
        print('Logging is DISABLED. Writing to console.')
        log.startLogging(sys.stdout)

    return True

# This is the context factory used to create TLS contexts.

from datetime import datetime
class AetherContextFactory(ContextFactory):
    def getContext(self):
        ctx = SSL.Context(SSL.TLSv1_METHOD)
        # This creates the key pairs and the cert if they do not exist.
        try:
            ctx.use_privatekey_file(basedir+'UserProfile/priv.pem')
            ctx.use_certificate_file(basedir+'UserProfile/cert.pem')
        except:
            # We don't have the requirements, so let's create them.
            print('This machine doesn\'nt seem to have a keypair and a cert. Creating new, at %s' %datetime.utcnow())
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
            newCertFile = open(basedir+'UserProfile/cert.pem', 'wb')
            newCertFile.write(crypto.dump_certificate(crypto.FILETYPE_PEM, cert))
            newCertFile.close()
            newKeyFile = open(basedir+'UserProfile/priv.pem', 'wb')
            newKeyFile.write(crypto.dump_privatekey(crypto.FILETYPE_PEM, k))
            newKeyFile.close()
            print('Key generation finished at %s' %datetime.utcnow())
            ctx.use_privatekey_file(basedir+'UserProfile/priv.pem')
            ctx.use_certificate_file(basedir+'UserProfile/cert.pem')

        return ctx

AetherClientContextFactory = ClientContextFactory


# The global app quit routine.

def quitApp(reactor):
    # This is (probably) buggy...
    if reactor.threadpool is not None:
        reactor.threadpool.stop()
    reactor.stop()
    sys.exit()
