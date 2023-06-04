# Aether makefile for building and/or cross-compiling.

# This file assumes that you are running it on a mainstream UNIX. (OS X, Linux). Other OSes might work, but are untested.

# Usage

# To build the backend server:
# make compile-linux-x64-extverify-be
# (See {{ COMPILE INSTRUCTIONS }})

# To build based on os:

# make mac
# make win
# make linux

# To build server
# make linux-server-extv (for external verify)
# make linux-server

all: remove_workdir create_workdir compile-mac compile-win compile-linux  compile-linux-server compile-linux-extv-server package move-updater-metadata dircleanup

# Desktop
mac: remove_workdir create_workdir compile-mac package-mac-only move-updater-metadata-mac dircleanup-mac
win: remove_workdir create_workdir compile-win package-win-only move-updater-metadata-win dircleanup-win
linux: remove_workdir create_workdir compile-linux package-linux-only move-updater-metadata-linux dircleanup-linux
linux-native: remove_workdir create_workdir compile-linux-native package-linux-native-only move-updater-metadata-linux-native dircleanup-linux-native

# Server
linux-server-extv: remove_workdir create_workdir package-linux-extv-server-only
linux-server: remove_workdir create_workdir package-linux-server-only

# Go parameters
GOCMD=go
EXTERNAL_VERIFY_ENABLED=-tags extvenabled
EXTERNAL_VERIFY_DISABLED=-tags

#
# {{ COMPILE INSTRUCTIONS }}
# To build, set the external verify state to disabled
#

EXTERNAL_VERIFY_STATE=$(EXTERNAL_VERIFY_ENABLED)
# 																			^^^^^^^ to DISABLED
GOBUILD=$(GOCMD) build
GOBUILD_WITH_EXTVERIFY=$(GOCMD) build $(EXTERNAL_VERIFY_STATE)
GOCLEAN=$(GOCMD) clean
UPX= upx -1
BE_BINARY_NAME=aether-backend
FE_BINARY_NAME=aether-frontend
CGO=CGO_ENABLED=0
# CGO=
# Workdir prep

WORKDIR=client/MAKE_BINARIES

remove_workdir:
		rm -rf $(WORKDIR); rm -rf packaged

create_workdir: remove_workdir
		mkdir $(WORKDIR)

package: compile-mac compile-win compile-linux compile-linux-server compile-linux-extv-server
	cd client; npm run dist;

package-mac-only: compile-mac
	cd client; npm run dist-mac;

package-win-only: compile-win
	cd client; npm run dist-win;

package-linux-only: compile-linux
	cd client; npm run dist-linux;

package-linux-native-only: compile-linux-native
	cd client; npm run dist-linux-native;

package-linux-server-only: compile-linux-server
	cd client; npm run dist-linux-server

package-linux-extv-server-only: compile-linux-extv-server
	cd client; npm run dist-linux-server

# Workdir cleanup - electron builder leaves empty .git folders around, which confuses Sublime's git implementation by showing all files as new. This uses rmdir instead of rm -rf because rmdir will not remove a non-empty directory. Just a safeguard against blowing up the real git directory.

dircleanup-mac: move-updater-metadata-mac
	cd client; rmdir .git 2> /dev/null || true; cd ../../../; rmdir .git 2> /dev/null || true;
dircleanup-win: move-updater-metadata-win
	cd client; rmdir .git 2> /dev/null || true; cd ../../../; rmdir .git 2> /dev/null || true;
dircleanup-linux: move-updater-metadata-linux
	cd client; rmdir .git 2> /dev/null || true; cd ../../../; rmdir .git 2> /dev/null || true;
dircleanup-linux-native: move-updater-metadata-linux-native
	cd client; rmdir .git 2> /dev/null || true; cd ../../../; rmdir .git 2> /dev/null || true;

dircleanup: move-updater-metadata
	cd client; rmdir .git 2> /dev/null || true; cd ../../../; rmdir .git 2> /dev/null || true;

# Carry over updater metadata

# This moves the update.yml, update-mac.yml, update-linux.yml files to the appropriate folders. Why is this needed? Because those files are generated after the afterAllArtifactBuild hook of electron-builder, therefore our afterpack script that runs on that hook is not able to get to those files. This one just looks at the latest build by timestamp, and moves the files in those places.

move-updater-metadata-mac: package-mac-only
	LATESTVERSION=$$(ls -td ../../ReleaseArchive/*/ | head -1); LATESTBUILD=$$(ls -td $$LATESTVERSION*/ | head -1); cd $$LATESTBUILD; mv ../../../BundledReleases/latest-mac.yml ./mac/latest-mac.yml 2> /dev/null || true;

move-updater-metadata-win: package-win-only
	LATESTVERSION=$$(ls -td ../../ReleaseArchive/*/ | head -1); LATESTBUILD=$$(ls -td $$LATESTVERSION*/ | head -1); cd $$LATESTBUILD; mv ../../../BundledReleases/latest.yml ./win/latest.yml 2> /dev/null || true;

move-updater-metadata-linux: package-linux-only
	LATESTVERSION=$$(ls -td ../../ReleaseArchive/*/ | head -1); LATESTBUILD=$$(ls -td $$LATESTVERSION*/ | head -1); cd $$LATESTBUILD; mv ../../../BundledReleases/latest-linux.yml ./linux/latest-linux.yml 2> /dev/null || true;

move-updater-metadata-linux-native: package-linux-native-only
	LATESTVERSION=$$(ls -td ../../ReleaseArchive/*/ | head -1); LATESTBUILD=$$(ls -td $$LATESTVERSION*/ | head -1); cd $$LATESTBUILD; mv ../../../BundledReleases/latest-linux.yml ./linux/latest-linux.yml 2> /dev/null || true;

move-updater-metadata: package
	LATESTVERSION=$$(ls -td ../../ReleaseArchive/*/ | head -1); LATESTBUILD=$$(ls -td $$LATESTVERSION*/ | head -1); cd $$LATESTBUILD; mv ../../../BundledReleases/latest-mac.yml ./mac/latest-mac.yml 2> /dev/null || true; mv ../../../BundledReleases/latest-linux.yml ./linux/latest-linux.yml 2> /dev/null || true; mv ../../../BundledReleases/latest.yml ./win/latest.yml 2> /dev/null || true;

# 2> /dev/null - route stderr to dev null
# || true - make the return code always return successful

# Go linker flags

FLAGS_LINKER_MAC_NATIVE= -ldflags="-s -w"
FLAGS_LINKER_LINUX_NATIVE= -ldflags="-s -w"

FLAGS_LINKER_LINUX_CROSS= -ldflags="-s -w -linkmode external -extldflags -static"
FLAGS_LINKER_WIN_CROSS= -ldflags="-s -w -extldflags -static"

# Go Compiler flags

# OS
MAC=GOOS=darwin
WIN=GOOS=windows
LINUX=GOOS=linux

# ARCH

# X86
X64=GOARCH=amd64
X86=GOARCH=386

# ARM
ARM64=GOARCH=arm64
ARM32=GOARCH=arm

# Native-compile flags

# Mac
# None - Mac only supports native compile, since you cannot build Mac OS binaries on non-Mac platforms, so the default is native build.

# Linux
CCFLAGS_LINUX_NATIVE_X64=CC=gcc CXX=g++

# Windows
# None - No one has attempted this yet.

# Cross-compile flags
CCFLAGS_LINUX_X64=CC=x86_64-linux-musl-gcc CXX=x86_64-linux-musl-g++
CCFLAGS_LINUX_X86=CC=i486-linux-musl-gcc CXX=i486-linux-musl-g++
CCFLAGS_LINUX_ARM64=CC=aarch64-linux-musl-gcc CXX=aarch64-linux-musl-g++
CCFLAGS_LINUX_ARM32=CC=arm-linux-musleabi-gcc CXX=arm-linux-musleabi-g++

CCFLAGS_WIN_X64=CC=x86_64-w64-mingw32-gcc CXX=x86_64-w64-mingw32-g++
CCFLAGS_WIN_X86=CC=i686-w64-mingw32-gcc CXX=i686-w64-mingw32-g++

# Flag sets

# Mac flagset (native)
MAC_X64=$(CGO) $(MAC) $(X64)
MAC_X86=$(CGO) $(MAC) $(X86)

# Linux flags (native)
LINUX_X64_NATIVE=$(CGO) $(LINUX) $(X64) $(CCFLAGS_LINUX_NATIVE_X64)

# Windows flags (native)
# None - No one has attempted this yet.


# Windows flagset (cross-compile)
WIN_X64=$(CGO) $(WIN) $(X64) $(CCFLAGS_WIN_X64)
WIN_X86=$(CGO) $(WIN) $(X86)  $(CCFLAGS_WIN_X86)

# Linux flagset (cross-compile)
LINUX_X64=$(CGO) $(LINUX) $(X64) $(CCFLAGS_LINUX_X64)
LINUX_X86=$(CGO) $(LINUX) $(X86) $(CCFLAGS_LINUX_X86)
LINUX_ARM64=$(CGO) $(LINUX) $(ARM64) $(CCFLAGS_LINUX_ARM64)
LINUX_ARM32=$(CGO) $(LINUX) $(ARM32) $(CCFLAGS_LINUX_ARM32)



# Names

MAC_X64_NAME=-mac-x64
MAC_X86_NAME=-mac-ia32

WIN_X64_NAME=-win-x64.exe
WIN_X86_NAME=-win-ia32.exe

LINUX_X64_NAME=-linux-x64
LINUX_X64_EXTVERIFY_NAME=-linux-x64-extverify
LINUX_X86_NAME=-linux-ia32
LINUX_ARM64_NAME=-linux-arm64
LINUX_ARM32_NAME=-linux-arm32

# Compile paths

## Desktop
compile-mac: compile-mac-x64 move-mac-x64-be move-mac-x64-fe
compile-win: compile-win-x64 compile-win-x86 move-win-x64-be move-win-x64-fe move-win-x86-be move-win-x86-fe
compile-linux: compile-linux-x64 compile-linux-x86 compile-linux-arm64 compile-linux-arm32 move-linux-x64-be move-linux-x64-fe move-linux-x86-be move-linux-x86-fe move-linux-arm64-be move-linux-arm64-fe move-linux-arm32-be move-linux-arm32-fe
compile-linux-native: compile-linux-native-x64 move-linux-native-x64-be move-linux-native-x64-fe
# For Linux native, we don't have x86, arm64 and arm32 yet - those we still cross-compile from Mac OS.

## Server
# compile-linux-extv-server: compile-linux-x64-extverify-be move-linux-x64-extverify-be compress-linux-x64-extverify-be version-linux-x64-extverify-be
# compile-linux-server: compile-linux-x64-be move-linux-x64-be compress-linux-x64-be version-linux-x64-be

## Server
compile-linux-extv-server: compile-linux-x64-extverify-be move-linux-x64-extverify-be compress-linux-x64-extverify-be
compile-linux-server: compile-linux-x64-be move-linux-x64-be compress-linux-x64-be



#########
# MAC X64 (NATIVE)
#########

compile-mac-x64: compile-mac-x64-be move-mac-x64-be compress-mac-x64-be compile-mac-x64-fe move-mac-x64-fe compress-mac-x64-fe

compile-mac-x64-be: create_workdir
	cd backend; $(MAC_X64) $(GOBUILD) $(FLAGS_LINKER_MAC_NATIVE) -o $(BE_BINARY_NAME)$(MAC_X64_NAME);

move-mac-x64-be: compile-mac-x64-be
	mv backend/$(BE_BINARY_NAME)$(MAC_X64_NAME) $(WORKDIR)

compress-mac-x64-be: move-mac-x64-be
	$(UPX) $(WORKDIR)/$(BE_BINARY_NAME)$(MAC_X64_NAME)

compile-mac-x64-fe:
	cd frontend; $(MAC_X64) $(GOBUILD) $(FLAGS_LINKER_MAC_NATIVE) -o $(FE_BINARY_NAME)$(MAC_X64_NAME);

move-mac-x64-fe: compile-mac-x64-fe
	mv frontend/$(FE_BINARY_NAME)$(MAC_X64_NAME) $(WORKDIR)

compress-mac-x64-fe: move-mac-x64-fe
	$(UPX) $(WORKDIR)/$(FE_BINARY_NAME)$(MAC_X64_NAME)



#########
# MAC X86 (NATIVE)
#########

compile-mac-x86: compile-mac-x86-be move-mac-x86-be compress-mac-x86-be compile-mac-x86-fe move-mac-x86-fe compress-mac-x86-fe

compile-mac-x86-be: create_workdir
	cd backend; $(MAC_X86) $(GOBUILD) $(FLAGS_LINKER_MAC_NATIVE) -o $(BE_BINARY_NAME)$(MAC_X86_NAME);

move-mac-x86-be: compile-mac-x86-be
	mv backend/$(BE_BINARY_NAME)$(MAC_X86_NAME) $(WORKDIR)

compress-mac-x86-be: move-mac-x86-be
	$(UPX) $(WORKDIR)/$(BE_BINARY_NAME)$(MAC_X86_NAME)


compile-mac-x86-fe: create_workdir
	cd frontend; $(MAC_X86) $(GOBUILD) $(FLAGS_LINKER_MAC_NATIVE) -o $(FE_BINARY_NAME)$(MAC_X86_NAME);

move-mac-x86-fe: compile-mac-x86-fe
	mv frontend/$(FE_BINARY_NAME)$(MAC_X86_NAME) $(WORKDIR)

compress-mac-x86-fe: move-mac-x86-fe
	$(UPX) $(WORKDIR)/$(FE_BINARY_NAME)$(MAC_X86_NAME)



#########
# MAC X64 (CROSS-COMPILE)
#########

# Not possible - Can't target Mac on non-Mac OSes.



#########
# MAC X86 (CROSS-COMPILE)
#########

# Not possible - Can't target Mac on non-Mac OSes.



#########
# WIN X64 (NATIVE)
#########

# Nobody tried to compile on Windows yet.



#########
# WIN X86 (NATIVE)
#########

# Nobody tried to compile on Windows yet.



#########
# WIN X64 (CROSS-COMPILE)
#########

compile-win-x64: compile-win-x64-be move-win-x64-be compress-win-x64-be compile-win-x64-fe move-win-x64-fe compress-win-x64-fe

compile-win-x64-be: create_workdir
	cd backend; $(WIN_X64) $(GOBUILD) $(FLAGS_LINKER_WIN_CROSS) -o $(BE_BINARY_NAME)$(WIN_X64_NAME);

move-win-x64-be: compile-win-x64-be
	mv backend/$(BE_BINARY_NAME)$(WIN_X64_NAME) $(WORKDIR)

compress-win-x64-be: move-win-x64-be
	$(UPX) $(WORKDIR)/$(BE_BINARY_NAME)$(WIN_X64_NAME)

compile-win-x64-fe: create_workdir
	cd frontend; $(WIN_X64) $(GOBUILD) $(FLAGS_LINKER_WIN_CROSS) -o $(FE_BINARY_NAME)$(WIN_X64_NAME);

move-win-x64-fe: compile-win-x64-fe
	mv frontend/$(FE_BINARY_NAME)$(WIN_X64_NAME) $(WORKDIR)

compress-win-x64-fe: move-win-x64-fe
	$(UPX) $(WORKDIR)/$(FE_BINARY_NAME)$(WIN_X64_NAME)



#########
# WIN X86 (CROSS-COMPILE)
#########

compile-win-x86: compile-win-x86-be move-win-x86-be compress-win-x86-be compile-win-x86-fe move-win-x86-fe compress-win-x86-fe

compile-win-x86-be: create_workdir
	cd backend; $(WIN_X86) $(GOBUILD) $(FLAGS_LINKER_WIN_CROSS) -o $(BE_BINARY_NAME)$(WIN_X86_NAME);

move-win-x86-be: compile-win-x86-be
	mv backend/$(BE_BINARY_NAME)$(WIN_X86_NAME) $(WORKDIR)

compress-win-x86-be: move-win-x86-be
	$(UPX) $(WORKDIR)/$(BE_BINARY_NAME)$(WIN_X86_NAME)

compile-win-x86-fe: create_workdir
	cd frontend; $(WIN_X86) $(GOBUILD) $(FLAGS_LINKER_WIN_CROSS) -o $(FE_BINARY_NAME)$(WIN_X86_NAME);

move-win-x86-fe: compile-win-x86-fe
	mv frontend/$(FE_BINARY_NAME)$(WIN_X86_NAME) $(WORKDIR)

compress-win-x86-fe: move-win-x86-fe
	$(UPX) $(WORKDIR)/$(FE_BINARY_NAME)$(WIN_X86_NAME)



#########
# LNX X64 (NATIVE)
#########

compile-linux-native-x64: compile-linux-native-x64-be move-linux-native-x64-be compress-linux-native-x64-be compile-linux-native-x64-fe move-linux-native-x64-fe compress-linux-native-x64-fe

compile-linux-native-x64-be: create_workdir
	cd backend; $(LINUX_X64_NATIVE) $(GOBUILD) $(FLAGS_LINKER_LINUX_NATIVE) -o $(BE_BINARY_NAME)$(LINUX_X64_NAME)

move-linux-native-x64-be: compile-linux-native-x64-be
	mv backend/$(BE_BINARY_NAME)$(LINUX_X64_NAME) $(WORKDIR)

compress-linux-native-x64-be: move-linux-native-x64-be
	$(UPX) $(WORKDIR)/$(BE_BINARY_NAME)$(LINUX_X64_NAME)

# We handle versioning in package.json.
# version-linux-x64-be: compress-linux-x64-be
# 	NAME=$$(node client/buildresources/get-version-from-git.js print); mv $(WORKDIR)/$(BE_BINARY_NAME)$(LINUX_X64_NAME) $(WORKDIR)/$(BE_BINARY_NAME)$(LINUX_X64_NAME)-$$NAME

compile-linux-native-x64-fe: create_workdir
	cd frontend; $(LINUX_X64_NATIVE) $(GOBUILD) $(FLAGS_LINKER_LINUX_NATIVE) -o $(FE_BINARY_NAME)$(LINUX_X64_NAME);

move-linux-native-x64-fe: compile-linux-native-x64-fe
	mv frontend/$(FE_BINARY_NAME)$(LINUX_X64_NAME) $(WORKDIR)

compress-linux-native-x64-fe: move-linux-native-x64-fe
	$(UPX) $(WORKDIR)/$(FE_BINARY_NAME)$(LINUX_X64_NAME)

#########
# LNX X86 (NATIVE)
#########

# Nobody tried this yet.



###########
# LNX ARM64 (NATIVE)
###########

# Nobody tried this yet.



###########
# LNX ARM32 (NATIVE)
###########

# Nobody tried this yet.



#########
# LNX X64 (CROSS-COMPILE)
#########

compile-linux-x64: compile-linux-x64-be move-linux-x64-be compress-linux-x64-be compile-linux-x64-fe move-linux-x64-fe compress-linux-x64-fe

compile-linux-x64-be: create_workdir
	cd backend; $(LINUX_X64) $(GOBUILD) $(FLAGS_LINKER_LINUX_CROSS) -o $(BE_BINARY_NAME)$(LINUX_X64_NAME)

move-linux-x64-be: compile-linux-x64-be
	mv backend/$(BE_BINARY_NAME)$(LINUX_X64_NAME) $(WORKDIR)

compress-linux-x64-be: move-linux-x64-be
	$(UPX) $(WORKDIR)/$(BE_BINARY_NAME)$(LINUX_X64_NAME)

# We handle versioning in package.json.
# version-linux-x64-be: compress-linux-x64-be
# 	NAME=$$(node client/buildresources/get-version-from-git.js print); mv $(WORKDIR)/$(BE_BINARY_NAME)$(LINUX_X64_NAME) $(WORKDIR)/$(BE_BINARY_NAME)$(LINUX_X64_NAME)-$$NAME

compile-linux-x64-fe: create_workdir
	cd frontend; $(LINUX_X64) $(GOBUILD) $(FLAGS_LINKER_LINUX_CROSS) -o $(FE_BINARY_NAME)$(LINUX_X64_NAME);

move-linux-x64-fe: compile-linux-x64-fe
	mv frontend/$(FE_BINARY_NAME)$(LINUX_X64_NAME) $(WORKDIR)

compress-linux-x64-fe: move-linux-x64-fe
	$(UPX) $(WORKDIR)/$(FE_BINARY_NAME)$(LINUX_X64_NAME)


# Linux X64 extverify backend
# (for running it behind a reverse proxy - you normally don't need this, since Aether doesn't need / is not designed to have a reverse proxy to sit in front. If you do, you need to implement Aether into your reverse proxy, and this provides the interface for it)

compile-linux-x64-extverify-be: create_workdir
	cd backend; $(LINUX_X64) $(GOBUILD_WITH_EXTVERIFY) $(FLAGS_LINKER_LINUX_CROSS) -o $(BE_BINARY_NAME)$(LINUX_X64_EXTVERIFY_NAME)

move-linux-x64-extverify-be: compile-linux-x64-extverify-be
	mv backend/$(BE_BINARY_NAME)$(LINUX_X64_EXTVERIFY_NAME) $(WORKDIR)

compress-linux-x64-extverify-be: move-linux-x64-extverify-be
	$(UPX) $(WORKDIR)/$(BE_BINARY_NAME)$(LINUX_X64_EXTVERIFY_NAME)

# We handle versioning in package.json.
# version-linux-x64-extverify-be: compress-linux-x64-extverify-be
# 	NAME=$$(node client/buildresources/get-version-from-git.js print); mv $(WORKDIR)/$(BE_BINARY_NAME)$(LINUX_X64_EXTVERIFY_NAME) $(WORKDIR)/$(BE_BINARY_NAME)$(LINUX_X64_EXTVERIFY_NAME)-$$NAME



#########
# LNX X86 (CROSS-COMPILE)
#########

compile-linux-x86: compile-linux-x86-be move-linux-x86-be compress-linux-x86-be compile-linux-x86-fe move-linux-x86-fe compress-linux-x86-fe

compile-linux-x86-be: create_workdir
	cd backend; $(LINUX_X86) $(GOBUILD) $(FLAGS_LINKER_LINUX_CROSS) -o $(BE_BINARY_NAME)$(LINUX_X86_NAME);

move-linux-x86-be: compile-linux-x86-be
	mv backend/$(BE_BINARY_NAME)$(LINUX_X86_NAME) $(WORKDIR)

compress-linux-x86-be: move-linux-x86-be
	$(UPX) $(WORKDIR)/$(BE_BINARY_NAME)$(LINUX_X86_NAME)

compile-linux-x86-fe: create_workdir
	cd frontend; $(LINUX_X86) $(GOBUILD) $(FLAGS_LINKER_LINUX_CROSS) -o $(FE_BINARY_NAME)$(LINUX_X86_NAME);

move-linux-x86-fe: compile-linux-x86-fe
	mv frontend/$(FE_BINARY_NAME)$(LINUX_X86_NAME) $(WORKDIR)

compress-linux-x86-fe: move-linux-x86-fe
	$(UPX) $(WORKDIR)/$(FE_BINARY_NAME)$(LINUX_X86_NAME)



###########
# LNX ARM64 (CROSS-COMPILE)
###########

compile-linux-arm64: compile-linux-arm64-be move-linux-arm64-be compress-linux-arm64-be compile-linux-arm64-fe move-linux-arm64-fe compress-linux-arm64-fe

compile-linux-arm64-be: create_workdir
	cd backend; $(LINUX_ARM64) $(GOBUILD) $(FLAGS_LINKER_LINUX_CROSS) -o $(BE_BINARY_NAME)$(LINUX_ARM64_NAME);

move-linux-arm64-be: compile-linux-arm64-be
	mv backend/$(BE_BINARY_NAME)$(LINUX_ARM64_NAME) $(WORKDIR)

compress-linux-arm64-be: move-linux-arm64-be
	$(UPX) $(WORKDIR)/$(BE_BINARY_NAME)$(LINUX_ARM64_NAME)

compile-linux-arm64-fe: create_workdir
	cd frontend; $(LINUX_ARM64) $(GOBUILD) $(FLAGS_LINKER_LINUX_CROSS) -o $(FE_BINARY_NAME)$(LINUX_ARM64_NAME);

move-linux-arm64-fe: compile-linux-arm64-fe
	mv frontend/$(FE_BINARY_NAME)$(LINUX_ARM64_NAME) $(WORKDIR)

compress-linux-arm64-fe: move-linux-arm64-fe
	$(UPX) $(WORKDIR)/$(FE_BINARY_NAME)$(LINUX_ARM64_NAME)



###########
# LNX ARM32 (CROSS-COMPILE)
###########

compile-linux-arm32: compile-linux-arm32-be move-linux-arm32-be compress-linux-arm32-be compile-linux-arm32-fe move-linux-arm32-fe compress-linux-arm32-fe

compile-linux-arm32-be: create_workdir
	cd backend; $(LINUX_ARM32) $(GOBUILD) $(FLAGS_LINKER_LINUX_CROSS) -o $(BE_BINARY_NAME)$(LINUX_ARM32_NAME);

move-linux-arm32-be: compile-linux-arm32-be
	mv backend/$(BE_BINARY_NAME)$(LINUX_ARM32_NAME) $(WORKDIR)

compress-linux-arm32-be: move-linux-arm32-be
	$(UPX) $(WORKDIR)/$(BE_BINARY_NAME)$(LINUX_ARM32_NAME)

compile-linux-arm32-fe: create_workdir
	cd frontend; $(LINUX_ARM32) $(GOBUILD) $(FLAGS_LINKER_LINUX_CROSS) -o $(FE_BINARY_NAME)$(LINUX_ARM32_NAME);

move-linux-arm32-fe: compile-linux-arm32-fe
	mv frontend/$(FE_BINARY_NAME)$(LINUX_ARM32_NAME) $(WORKDIR)

compress-linux-arm32-fe: move-linux-arm32-fe
	$(UPX) $(WORKDIR)/$(FE_BINARY_NAME)$(LINUX_ARM32_NAME)

# bincompress:
# 		cd backend; $(UPX) $(BE_BINARY_NAME)
# 		cd frontend; $(UPX) $(FE_BINARY_NAME)
# test:
# 		$(GOTEST) -v ./...
# clean:
# 		cd backend; $(GOCLEAN); rm -f $(BE_BINARY_NAME)
# 		cd frontend; $(GOCLEAN); rm -f $(FE_BINARY_NAME)

# # Cross compilation


# all: remove_workdir create_workdir clean build bincompress
