# Install script for directory: /workdir/native/src/seal/util

# Set the install prefix
if(NOT DEFINED CMAKE_INSTALL_PREFIX)
  set(CMAKE_INSTALL_PREFIX "/usr/local")
endif()
string(REGEX REPLACE "/$" "" CMAKE_INSTALL_PREFIX "${CMAKE_INSTALL_PREFIX}")

# Set the install configuration name.
if(NOT DEFINED CMAKE_INSTALL_CONFIG_NAME)
  if(BUILD_TYPE)
    string(REGEX REPLACE "^[^A-Za-z0-9_]+" ""
           CMAKE_INSTALL_CONFIG_NAME "${BUILD_TYPE}")
  else()
    set(CMAKE_INSTALL_CONFIG_NAME "Release")
  endif()
  message(STATUS "Install configuration: \"${CMAKE_INSTALL_CONFIG_NAME}\"")
endif()

# Set the component getting installed.
if(NOT CMAKE_INSTALL_COMPONENT)
  if(COMPONENT)
    message(STATUS "Install component: \"${COMPONENT}\"")
    set(CMAKE_INSTALL_COMPONENT "${COMPONENT}")
  else()
    set(CMAKE_INSTALL_COMPONENT)
  endif()
endif()

# Install shared libraries without execute permission?
if(NOT DEFINED CMAKE_INSTALL_SO_NO_EXE)
  set(CMAKE_INSTALL_SO_NO_EXE "1")
endif()

# Is this installation the result of a crosscompile?
if(NOT DEFINED CMAKE_CROSSCOMPILING)
  set(CMAKE_CROSSCOMPILING "TRUE")
endif()

# Set path to fallback-tool for dependency-resolution.
if(NOT DEFINED CMAKE_OBJDUMP)
  set(CMAKE_OBJDUMP "/usr/bin/x86_64-linux-gnu-objdump")
endif()

if(CMAKE_INSTALL_COMPONENT STREQUAL "Unspecified" OR NOT CMAKE_INSTALL_COMPONENT)
  file(INSTALL DESTINATION "${CMAKE_INSTALL_PREFIX}/include/SEAL-4.1/seal/util" TYPE FILE FILES
    "/workdir/native/src/seal/util/blake2.h"
    "/workdir/native/src/seal/util/blake2-impl.h"
    "/workdir/native/src/seal/util/clang.h"
    "/workdir/native/src/seal/util/clipnormal.h"
    "/workdir/native/src/seal/util/common.h"
    "/workdir/native/src/seal/util/croots.h"
    "/workdir/native/src/seal/util/defines.h"
    "/workdir/native/src/seal/util/dwthandler.h"
    "/workdir/native/src/seal/util/fips202.h"
    "/workdir/native/src/seal/util/galois.h"
    "/workdir/native/src/seal/util/gcc.h"
    "/workdir/native/src/seal/util/globals.h"
    "/workdir/native/src/seal/util/hash.h"
    "/workdir/native/src/seal/util/hestdparms.h"
    "/workdir/native/src/seal/util/iterator.h"
    "/workdir/native/src/seal/util/locks.h"
    "/workdir/native/src/seal/util/mempool.h"
    "/workdir/native/src/seal/util/msvc.h"
    "/workdir/native/src/seal/util/numth.h"
    "/workdir/native/src/seal/util/pointer.h"
    "/workdir/native/src/seal/util/polyarithsmallmod.h"
    "/workdir/native/src/seal/util/polycore.h"
    "/workdir/native/src/seal/util/rlwe.h"
    "/workdir/native/src/seal/util/rns.h"
    "/workdir/native/src/seal/util/scalingvariant.h"
    "/workdir/native/src/seal/util/ntt.h"
    "/workdir/native/src/seal/util/streambuf.h"
    "/workdir/native/src/seal/util/uintarith.h"
    "/workdir/native/src/seal/util/uintarithmod.h"
    "/workdir/native/src/seal/util/uintarithsmallmod.h"
    "/workdir/native/src/seal/util/uintcore.h"
    "/workdir/native/src/seal/util/ztools.h"
    )
endif()

string(REPLACE ";" "\n" CMAKE_INSTALL_MANIFEST_CONTENT
       "${CMAKE_INSTALL_MANIFEST_FILES}")
if(CMAKE_INSTALL_LOCAL_ONLY)
  file(WRITE "/workdir/build/native/src/seal/util/install_local_manifest.txt"
     "${CMAKE_INSTALL_MANIFEST_CONTENT}")
endif()
