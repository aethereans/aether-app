// Persistence > Base
// This file contains the SQL statements that are necessary to handle database action, as well as basic maintenance functions for the database itself.

package persistence

import (
	"aether-core/aether/services/globals"
	"aether-core/aether/services/logging"
	"aether-core/aether/services/toolbox"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	// "github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	// _ "github.com/lib/pq"
	"errors"
	// "github.com/fatih/color"
	// "os"
	"path/filepath"
	"strings"
	"time"
)

// DeleteDatabase removes the existing database in the default location.
func DeleteDatabase() {
	if globals.BackendConfig.GetDbEngine() == "sqlite" {
		dbLoc := filepath.Join(globals.BackendConfig.GetSQLiteDBLocation(), "AetherDB.db")
		toolbox.DeleteFromDisk(dbLoc)
	} else if globals.BackendConfig.GetDbEngine() == "mysql" {
		globals.DbInstance.MustExec("DROP DATABASE `AetherDB`;")
	}
}

// CreateDatabase creates a new database in the default location and places into it the database schema.

func CreateDatabase() {
	err := createDatabase()
	if err != nil {
		if strings.Contains(err.Error(), "Database was locked") {
			logging.Log(1, err)
			if strings.Contains(err.Error(), "Database was locked") {
				logging.Log(1, "This transaction was not committed because database was locked. We'll wait 10 seconds and retry the transaction.")
				time.Sleep(10 * time.Second)
				logging.Log(1, "Retrying the previously failed CreateDatabase transaction.")
				err2 := createDatabase()
				if err2 != nil {
					logging.LogCrash(fmt.Sprintf("The second attempt to commit this data to the database failed. The first attempt had failed because the database was locked. The second attempt failed with the error: %s This database is corrupted. Quitting.", err2))
				} else { // If the reattempted transaction succeeds
					logging.Log(1, "The retry attempt of the failed transaction succeeded.")
				}
			}
		}
	}
}
func createDatabase() error {
	var schemaPrep1 string
	var schemaPrep2 string
	var schemaPrep3 string
	var schemaPrep4 string
	var schema1 string
	// var schema2 string
	var schema3 string
	var schema4 string
	var schema5 string
	var schema6 string
	var schema7 string
	var schema8 string
	var schema9 string
	var schema10 string
	var schema11 string
	var schema12 string
	// var schema13 string
	// var schema14 string
	// var schema15 string
	var schema16 string
	// var schema17 string
	// var schema18 string
	var idxSqlite1 string
	var idxSqlite2 string
	var idxSqlite3 string
	var idxSqlite4 string
	var idxSqlite5 string
	var idxSqlite6 string
	var idxSqlite7 string
	var idxSqlite8 string
	var idxSqlite9 string
	var idxSqlite10 string
	var idxSqlite11 string
	var idxSqlite12 string
	var idxSqlite13 string
	var idxSqlite14 string
	var idxSqlite15 string
	var idxSqlite16 string
	var idxSqlite17 string
	var idxSqlite18 string
	var idxSqlite19 string
	var idxSqlite20 string
	var idxSqlite21 string
	var idxSqlite22 string
	var idxSqlite23 string
	var idxSqlite24 string

	if globals.BackendConfig.DbEngine == "mysql" {
		schemaPrep1 = `
      CREATE DATABASE IF NOT EXISTS AetherDB
      DEFAULT CHARACTER SET utf8mb4
      DEFAULT COLLATE utf8mb4_general_ci;
    `
		schemaPrep2 = `USE AetherDB;`
		schemaPrep3 = `SET sql_mode = 'TRADITIONAL'; `      // strictest
		schemaPrep4 = `SET GLOBAL innodb_file_per_table=1;` // required to enable compression for mysql
		schema1 = `
        CREATE TABLE IF NOT EXISTS BoardOwners (
          BoardFingerprint VARCHAR(64) NOT NULL,
          KeyFingerprint VARCHAR(64) NOT NULL,
          Expiry BIGINT NOT NULL,
          Level SMALLINT NOT NULL,
          PRIMARY KEY(BoardFingerprint, KeyFingerprint)
        )ROW_FORMAT=COMPRESSED;
        `
		schema3 = `
        CREATE TABLE IF NOT EXISTS Boards (
          Fingerprint VARCHAR(64) PRIMARY KEY NOT NULL,
          Name VARCHAR(255) NOT NULL,
          Owner VARCHAR(64) NOT NULL,
          OwnerPublicKey VARCHAR(128) NOT NULL,
          Description MEDIUMTEXT NOT NULL,  -- Converted from varchar(65535) to text, because it doesn't fit into a MYSQL table. Enforce max 65535 chars on the application layer.
          Creation BIGINT NOT NULL,
          ProofOfWork VARCHAR(1024) NOT NULL,
          Signature VARCHAR(512) NOT NULL,
          LastUpdate BIGINT NOT NULL,
          UpdateProofOfWork VARCHAR(1024) NOT NULL,
          UpdateSignature VARCHAR(512) NOT NULL,
          LocalArrival BIGINT NOT NULL,
          LastReferenced BIGINT NOT NULL,
          EntityVersion SMALLINT NOT NULL,
          Language VARCHAR(3) NOT NULL,
          Meta MEDIUMTEXT NOT NULL,
          RealmId VARCHAR(64) NOT NULL,
          EncrContent MEDIUMTEXT NOT NULL,
          INDEX (LastReferenced, LastUpdate, Creation)
        )ROW_FORMAT=COMPRESSED;`
		schema4 = `
        CREATE TABLE IF NOT EXISTS Threads (
          Fingerprint VARCHAR(64) PRIMARY KEY NOT NULL,
          Board VARCHAR(64) NOT NULL,
          Name VARCHAR(255) NOT NULL,
          Body MEDIUMTEXT NOT NULL,
          Link VARCHAR(5000) NOT NULL,
          Owner VARCHAR(64) NOT NULL,
          OwnerPublicKey VARCHAR(128) NOT NULL,
          Creation BIGINT NOT NULL,
          ProofOfWork VARCHAR(1024) NOT NULL,
          Signature VARCHAR(512) NOT NULL,
          LastUpdate BIGINT NOT NULL,
          UpdateProofOfWork VARCHAR(1024) NOT NULL,
          UpdateSignature VARCHAR(512) NOT NULL,
          LocalArrival BIGINT NOT NULL,
          LastReferenced BIGINT NOT NULL,
          EntityVersion SMALLINT NOT NULL,
          Meta MEDIUMTEXT NOT NULL,
          RealmId VARCHAR(64) NOT NULL,
          EncrContent MEDIUMTEXT NOT NULL,
          INDEX (Board, LastReferenced, LastUpdate, Creation)
        )ROW_FORMAT=COMPRESSED;`
		schema5 = `
        CREATE TABLE IF NOT EXISTS Posts (
          Fingerprint VARCHAR(64) PRIMARY KEY NOT NULL,
          Board VARCHAR(64) NOT NULL,
          Thread VARCHAR(64) NOT NULL,
          Parent VARCHAR(64) NOT NULL,
          Body MEDIUMTEXT NOT NULL,
          Owner VARCHAR(64) NOT NULL,
          OwnerPublicKey VARCHAR(128) NOT NULL,
          Creation BIGINT NOT NULL,
          ProofOfWork VARCHAR(1024) NOT NULL,
          Signature VARCHAR(512) NOT NULL,
          LastUpdate BIGINT NOT NULL,
          UpdateProofOfWork VARCHAR(1024) NOT NULL,
          UpdateSignature VARCHAR(512) NOT NULL,
          LocalArrival BIGINT NOT NULL,
          LastReferenced BIGINT NOT NULL,
          EntityVersion SMALLINT NOT NULL,
          Meta MEDIUMTEXT NOT NULL,
          RealmId VARCHAR(64) NOT NULL,
          EncrContent MEDIUMTEXT NOT NULL,
          INDEX (Board, Thread, Parent, LastReferenced, LastUpdate, Creation)
        )ROW_FORMAT=COMPRESSED;`
		schema6 = `
        CREATE TABLE IF NOT EXISTS Votes (
          Fingerprint VARCHAR(64) PRIMARY KEY NOT NULL,
          Board VARCHAR(64) NOT NULL,
          Thread VARCHAR(64) NOT NULL,
          Target VARCHAR(64) NOT NULL,
          Owner VARCHAR(64) NOT NULL,
          OwnerPublicKey VARCHAR(128) NOT NULL,
          TypeClass SMALLINT NOT NULL,
          Type SMALLINT NOT NULL,
          Creation BIGINT NOT NULL,
          ProofOfWork VARCHAR(1024) NOT NULL,
          Signature VARCHAR(512) NOT NULL,
          LastUpdate BIGINT NOT NULL,
          UpdateProofOfWork VARCHAR(1024) NOT NULL,
          UpdateSignature VARCHAR(512) NOT NULL,
          LocalArrival BIGINT NOT NULL,
          LastReferenced BIGINT NOT NULL,
          EntityVersion SMALLINT NOT NULL,
          Meta MEDIUMTEXT NOT NULL,
          RealmId VARCHAR(64) NOT NULL,
          EncrContent MEDIUMTEXT NOT NULL,
          INDEX (Target, LastReferenced, LastUpdate, Creation)
        )ROW_FORMAT=COMPRESSED;`
		schema7 = `
        CREATE TABLE IF NOT EXISTS Addresses (
          Location VARCHAR(256) NOT NULL, -- From 2500
          Sublocation VARCHAR(256) NOT NULL, -- From 2500
          Port INTEGER NOT NULL,
          IPType SMALLINT NOT NULL,
          AddressType SMALLINT NOT NULL,
          LastSuccessfulPing BIGINT NOT NULL,
          LastSuccessfulSync BIGINT NOT NULL,
          ProtocolVersionMajor SMALLINT NOT NULL,
          ProtocolVersionMinor INTEGER NOT NULL,
          ClientVersionMajor SMALLINT NOT NULL,
          ClientVersionMinor INTEGER NOT NULL,
          ClientVersionPatch INTEGER NOT NULL,
          ClientName VARCHAR(255) NOT NULL,
          LocalArrival BIGINT NOT NULL,
          EntityVersion SMALLINT NOT NULL,
          RealmId VARCHAR(64) NOT NULL,
          PRIMARY KEY(Location, Sublocation, Port)
        )ROW_FORMAT=COMPRESSED;`
		schema8 = `
        CREATE TABLE IF NOT EXISTS PublicKeys (
          Fingerprint VARCHAR(64) PRIMARY KEY NOT NULL,
          Type VARCHAR(64) NOT NULL,
          PublicKey VARCHAR(128) NOT NULL,
          Expiry BIGINT NOT NULL,
          Name VARCHAR(64) NOT NULL,
          Info MEDIUMTEXT NOT NULL,
          Creation BIGINT NOT NULL,
          ProofOfWork VARCHAR(1024) NOT NULL,
          Signature VARCHAR(512) NOT NULL,
          LastUpdate BIGINT NOT NULL,
          UpdateProofOfWork VARCHAR(1024) NOT NULL,
          UpdateSignature VARCHAR(512) NOT NULL,
          LocalArrival BIGINT NOT NULL,
          LastReferenced BIGINT NOT NULL,
          EntityVersion SMALLINT NOT NULL,
          Meta MEDIUMTEXT NOT NULL,
          RealmId VARCHAR(64) NOT NULL,
          EncrContent MEDIUMTEXT NOT NULL,
          INDEX (PublicKey, LastReferenced)
        )ROW_FORMAT=COMPRESSED;`
		schema9 = `
        CREATE TABLE IF NOT EXISTS Truststates (
          Fingerprint VARCHAR(64) PRIMARY KEY NOT NULL,
          Target VARCHAR(64) NOT NULL,
          Owner VARCHAR(64) NOT NULL,
          OwnerPublicKey VARCHAR(128) NOT NULL,
          TypeClass SMALLINT NOT NULL,
          Type SMALLINT NOT NULL,
          Domain VARCHAR(64) NOT NULL,
          Expiry BIGINT NOT NULL,
          Creation BIGINT NOT NULL,
          ProofOfWork VARCHAR(1024) NOT NULL,
          Signature VARCHAR(512) NOT NULL,
          LastUpdate BIGINT NOT NULL,
          UpdateProofOfWork VARCHAR(1024) NOT NULL,
          UpdateSignature VARCHAR(512) NOT NULL,
          LocalArrival BIGINT NOT NULL,
          LastReferenced BIGINT NOT NULL,
          EntityVersion SMALLINT NOT NULL,
          Meta MEDIUMTEXT NOT NULL,
          RealmId VARCHAR(64) NOT NULL,
          EncrContent MEDIUMTEXT NOT NULL,
          INDEX (LastReferenced, LastUpdate, Creation)
        )ROW_FORMAT=COMPRESSED;
      `
		schema10 = `
          CREATE TABLE IF NOT EXISTS Nodes (
            Fingerprint VARCHAR(64) PRIMARY KEY NOT NULL,
            BoardsLastCheckin BIGINT NOT NULL,
            ThreadsLastCheckin BIGINT NOT NULL,
            PostsLastCheckin BIGINT NOT NULL,
            VotesLastCheckin BIGINT NOT NULL,
            KeysLastCheckin BIGINT NOT NULL,
            TruststatesLastCheckin BIGINT NOT NULL,
            AddressesLastCheckin BIGINT NOT NULL
          )ROW_FORMAT=COMPRESSED;
        `
		schema11 = `
          CREATE TABLE IF NOT EXISTS Subprotocols (
            Fingerprint VARCHAR(64) PRIMARY KEY NOT NULL,
            Name VARCHAR(64) NOT NULL,
            VersionMajor SMALLINT NOT NULL,
            VersionMinor INTEGER NOT NULL,
            SupportedEntities VARCHAR(5000) NOT NULL
          )ROW_FORMAT=COMPRESSED;`

		schema12 = `
          CREATE TABLE IF NOT EXISTS AddressesSubprotocols (
            AddressLocation VARCHAR(256) NOT NULL,
            AddressSublocation VARCHAR(256) NOT NULL,
            AddressPort INTEGER NOT NULL,
            SubprotocolFingerprint VARCHAR(64) NOT NULL,
            PRIMARY KEY(AddressLocation, AddressSublocation, AddressPort, SubprotocolFingerprint)
          )ROW_FORMAT=COMPRESSED;`
		schema16 = `
          CREATE TABLE IF NOT EXISTS Diagnostics (
            DbRoundtripTestField BIGINT PRIMARY KEY NOT NULL
          )ROW_FORMAT=COMPRESSED;`
	} else if globals.BackendConfig.DbEngine == "sqlite" {
		schemaPrep1 = `
    PRAGMA auto_vacuum = FULL;
    `
		schema1 = `
        CREATE TABLE IF NOT EXISTS "BoardOwners" (
          "BoardFingerprint" varchar(64) NOT NULL
        ,  "KeyFingerprint" varchar(64) NOT NULL
        ,  "Expiry" integer NOT NULL
        ,  "Level" integer NOT NULL
        ,  PRIMARY KEY ("BoardFingerprint","KeyFingerprint")
        );`
		schema3 = `
        CREATE TABLE IF NOT EXISTS "Boards" (
          "Fingerprint" varchar(64) NOT NULL
        ,  "Name" varchar(255) NOT NULL
        ,  "Owner" varchar(64) NOT NULL
        ,  "OwnerPublicKey" varchar(128) NOT NULL
        ,  "Description" text NOT NULL
        ,  "Creation" integer NOT NULL
        ,  "ProofOfWork" varchar(1024) NOT NULL
        ,  "Signature" varchar(512) NOT NULL
        ,  "LastUpdate" integer NOT NULL
        ,  "UpdateProofOfWork" varchar(1024) NOT NULL
        ,  "UpdateSignature" varchar(512) NOT NULL
        ,  "LocalArrival" integer NOT NULL
        ,  "LastReferenced" integer NOT NULL
        ,  "EntityVersion" integer NOT NULL
        ,  "Language" varchar(3) NOT NULL
        ,  "Meta" text NOT NULL
        ,  "RealmId" varchar(64) NOT NULL
        ,  "EncrContent" text NOT NULL
        ,  PRIMARY KEY ("Fingerprint")
        );`
		schema4 = `
        CREATE TABLE IF NOT EXISTS "Threads" (
          "Fingerprint" varchar(64) NOT NULL
        ,  "Board" varchar(64) NOT NULL
        ,  "Name" varchar(255) NOT NULL
        ,  "Body" text NOT NULL
        ,  "Link" varchar(5000) NOT NULL
        ,  "Owner" varchar(64) NOT NULL
        ,  "OwnerPublicKey" varchar(128) NOT NULL
        ,  "Creation" integer NOT NULL
        ,  "ProofOfWork" varchar(1024) NOT NULL
        ,  "Signature" varchar(512) NOT NULL
        ,  "LastUpdate" integer NOT NULL
        ,  "UpdateProofOfWork" varchar(1024) NOT NULL
        ,  "UpdateSignature" varchar(512) NOT NULL
        ,  "LocalArrival" integer NOT NULL
        ,  "LastReferenced" integer NOT NULL
        ,  "EntityVersion" integer NOT NULL
        ,  "Meta" text NOT NULL
        ,  "RealmId" varchar(64) NOT NULL
        ,  "EncrContent" text NOT NULL
        ,  PRIMARY KEY ("Fingerprint")
        );`
		schema5 = `
        CREATE TABLE IF NOT EXISTS "Posts" (
          "Fingerprint" varchar(64) NOT NULL
        ,  "Board" varchar(64) NOT NULL
        ,  "Thread" varchar(64) NOT NULL
        ,  "Parent" varchar(64) NOT NULL
        ,  "Body" text NOT NULL
        ,  "Owner" varchar(64) NOT NULL
        ,  "OwnerPublicKey" varchar(128) NOT NULL
        ,  "Creation" integer NOT NULL
        ,  "ProofOfWork" varchar(1024) NOT NULL
        ,  "Signature" varchar(512) NOT NULL
        ,  "LastUpdate" integer NOT NULL
        ,  "UpdateProofOfWork" varchar(1024) NOT NULL
        ,  "UpdateSignature" varchar(512) NOT NULL
        ,  "LocalArrival" integer NOT NULL
        ,  "LastReferenced" integer NOT NULL
        ,  "EntityVersion" integer NOT NULL
        ,  "Meta" text NOT NULL
        ,  "RealmId" varchar(64) NOT NULL
        ,  "EncrContent" text NOT NULL
        ,  PRIMARY KEY ("Fingerprint")
        );`
		schema6 = `
        CREATE TABLE IF NOT EXISTS "Votes" (
          "Fingerprint" varchar(64) NOT NULL
        ,  "Board" varchar(64) NOT NULL
        ,  "Thread" varchar(64) NOT NULL
        ,  "Target" varchar(64) NOT NULL
        ,  "Owner" varchar(64) NOT NULL
        ,  "OwnerPublicKey" varchar(128) NOT NULL
        ,  "TypeClass" integer NOT NULL
        ,  "Type" integer NOT NULL
        ,  "Creation" integer NOT NULL
        ,  "ProofOfWork" varchar(1024) NOT NULL
        ,  "Signature" varchar(512) NOT NULL
        ,  "LastUpdate" integer NOT NULL
        ,  "UpdateProofOfWork" varchar(1024) NOT NULL
        ,  "UpdateSignature" varchar(512) NOT NULL
        ,  "LocalArrival" integer NOT NULL
        ,  "LastReferenced" integer NOT NULL
        ,  "EntityVersion" integer NOT NULL
        ,  "Meta" text NOT NULL
        ,  "RealmId" varchar(64) NOT NULL
        ,  "EncrContent" text NOT NULL
        ,  PRIMARY KEY ("Fingerprint")
        );`
		schema8 = `
        CREATE TABLE IF NOT EXISTS "PublicKeys" (
          "Fingerprint" varchar(64) NOT NULL
        ,  "Type" varchar(64) NOT NULL
        ,  "PublicKey" text NOT NULL
        ,  "Expiry" integer NOT NULL
        ,  "Name" varchar(64) NOT NULL
        ,  "Info" text NOT NULL
        ,  "Creation" integer NOT NULL
        ,  "ProofOfWork" varchar(1024) NOT NULL
        ,  "Signature" varchar(512) NOT NULL
        ,  "LastUpdate" integer NOT NULL
        ,  "UpdateProofOfWork" varchar(1024) NOT NULL
        ,  "UpdateSignature" varchar(512) NOT NULL
        ,  "LocalArrival" integer NOT NULL
        ,  "LastReferenced" integer NOT NULL
        ,  "EntityVersion" integer NOT NULL
        ,  "Meta" text NOT NULL
        ,  "RealmId" varchar(64) NOT NULL
        ,  "EncrContent" text NOT NULL
        ,  PRIMARY KEY ("Fingerprint")
        );`
		schema9 = `
        CREATE TABLE IF NOT EXISTS "Truststates" (
          "Fingerprint" varchar(64) NOT NULL
        ,  "Target" varchar(64) NOT NULL
        ,  "Owner" varchar(64) NOT NULL
        ,  "OwnerPublicKey" varchar(128) NOT NULL
        ,  "TypeClass" integer NOT NULL
        ,  "Type" integer NOT NULL
        ,  "Domain" varchar(64) NOT NULL
        ,  "Expiry" integer NOT NULL
        ,  "Creation" integer NOT NULL
        ,  "ProofOfWork" varchar(1024) NOT NULL
        ,  "Signature" varchar(512) NOT NULL
        ,  "LastUpdate" integer NOT NULL
        ,  "UpdateProofOfWork" varchar(1024) NOT NULL
        ,  "UpdateSignature" varchar(512) NOT NULL
        ,  "LocalArrival" integer NOT NULL
        ,  "LastReferenced" integer NOT NULL
        ,  "EntityVersion" integer NOT NULL
        ,  "Meta" text NOT NULL
        ,  "RealmId" varchar(64) NOT NULL
        ,  "EncrContent" text NOT NULL
        ,  PRIMARY KEY ("Fingerprint")
        );`
		schema7 = `
        CREATE TABLE IF NOT EXISTS "Addresses" (
          "Location" varchar(256) NOT NULL
        ,  "Sublocation" varchar(256) NOT NULL
        ,  "Port" integer NOT NULL
        ,  "IPType" integer NOT NULL
        ,  "AddressType" integer NOT NULL
        ,  "LastSuccessfulPing" integer NOT NULL
        ,  "LastSuccessfulSync" integer NOT NULL
        ,  "ProtocolVersionMajor" integer NOT NULL
        ,  "ProtocolVersionMinor" integer NOT NULL
        ,  "ClientVersionMajor" integer NOT NULL
        ,  "ClientVersionMinor" integer NOT NULL
        ,  "ClientVersionPatch" integer NOT NULL
        ,  "ClientName" varchar(255) NOT NULL
        ,  "LocalArrival" integer NOT NULL
        ,  "EntityVersion" integer NOT NULL
        ,  "RealmId" varchar(64) NOT NULL
        ,  PRIMARY KEY ("Location","Sublocation","Port")
        );`
		schema10 = `
          CREATE TABLE IF NOT EXISTS "Nodes" (
            "Fingerprint" varchar(64) NOT NULL
          ,  "BoardsLastCheckin" integer NOT NULL
          ,  "ThreadsLastCheckin" integer NOT NULL
          ,  "PostsLastCheckin" integer NOT NULL
          ,  "VotesLastCheckin" integer NOT NULL
          ,  "KeysLastCheckin" integer NOT NULL
          ,  "TruststatesLastCheckin" integer NOT NULL
          ,  "AddressesLastCheckin" integer NOT NULL
          ,  PRIMARY KEY ("Fingerprint")
          );`
		schema11 = `
          CREATE TABLE IF NOT EXISTS "Subprotocols" (
            "Fingerprint" varchar(64) NOT NULL
          ,  "Name" varchar(64) NOT NULL
          ,  "VersionMajor" integer NOT NULL
          ,  "VersionMinor" integer NOT NULL
          ,  "SupportedEntities" varchar(5000) NOT NULL
          ,  PRIMARY KEY ("Fingerprint")
          );`
		schema12 = `
          CREATE TABLE IF NOT EXISTS "AddressesSubprotocols" (
            "AddressLocation" varchar(256) NOT NULL
          ,  "AddressSublocation" varchar(256) NOT NULL
          ,  "AddressPort" integer NOT NULL
          ,  "SubprotocolFingerprint" varchar(64) NOT NULL
          ,  PRIMARY KEY ("AddressLocation","AddressSublocation","AddressPort","SubprotocolFingerprint")
          );`
		schema16 = `
            CREATE TABLE IF NOT EXISTS "Diagnostics" (
              "DbRoundtripTestField" integer NOT NULL
            ,  PRIMARY KEY ("DbRoundtripTestField")
          );`

		idxSqlite1 = `
          CREATE INDEX IF NOT EXISTS "idx_Posts_Thread" ON "Posts" ("Thread");
          `
		idxSqlite2 = `
          CREATE INDEX IF NOT EXISTS "idx_Threads_Board" ON "Threads" ("Board");
          `
		idxSqlite3 = `
          CREATE INDEX IF NOT EXISTS "idx_Votes_Target" ON "Votes" ("Target");
          `
		idxSqlite4 = `
          CREATE INDEX IF NOT EXISTS "idx_PublicKeys_PublicKey" ON "PublicKeys" ("PublicKey");
          `
		idxSqlite5 = `
          CREATE INDEX IF NOT EXISTS "idx_Posts_Parent" ON "Posts" ("Parent");
          `
		// LastReferenced indexes
		idxSqlite6 = `
          CREATE INDEX IF NOT EXISTS "idx_Boards_LastReferenced" ON "Boards" ("LastReferenced");
          `
		idxSqlite7 = `
          CREATE INDEX IF NOT EXISTS "idx_Threads_LastReferenced" ON "Threads" ("LastReferenced");
          `
		idxSqlite8 = `
          CREATE INDEX IF NOT EXISTS "idx_Posts_LastReferenced" ON "Posts" ("LastReferenced");
          `
		idxSqlite9 = `
          CREATE INDEX IF NOT EXISTS "idx_Votes_LastReferenced" ON "Votes" ("LastReferenced");
          `
		idxSqlite10 = `
          CREATE INDEX IF NOT EXISTS "idx_PublicKeys_LastReferenced" ON "PublicKeys" ("LastReferenced");
          `
		idxSqlite11 = `
          CREATE INDEX IF NOT EXISTS "idx_Truststates_LastReferenced" ON "Truststates" ("LastReferenced");
          `
		// LastUpdate indexes
		idxSqlite12 = `
          CREATE INDEX IF NOT EXISTS "idx_Boards_LastUpdate" ON "Boards" ("LastUpdate");
          `
		idxSqlite13 = `
          CREATE INDEX IF NOT EXISTS "idx_Threads_LastUpdate" ON "Threads" ("LastUpdate");
          `
		idxSqlite14 = `
          CREATE INDEX IF NOT EXISTS "idx_Posts_LastUpdate" ON "Posts" ("LastUpdate");
          `
		idxSqlite15 = `
          CREATE INDEX IF NOT EXISTS "idx_Votes_LastUpdate" ON "Votes" ("LastUpdate");
          `
		idxSqlite16 = `
          CREATE INDEX IF NOT EXISTS "idx_PublicKeys_LastUpdate" ON "PublicKeys" ("LastUpdate");
          `
		idxSqlite17 = `
          CREATE INDEX IF NOT EXISTS "idx_Truststates_LastUpdate" ON "Truststates" ("LastUpdate");
          `
		// Creation indexes
		idxSqlite18 = `
          CREATE INDEX IF NOT EXISTS "idx_Boards_Creation" ON "Boards" ("Creation");
          `
		idxSqlite19 = `
          CREATE INDEX IF NOT EXISTS "idx_Threads_Creation" ON "Threads" ("Creation");
          `
		idxSqlite20 = `
          CREATE INDEX IF NOT EXISTS "idx_Posts_Creation" ON "Posts" ("Creation");
          `
		idxSqlite21 = `
          CREATE INDEX IF NOT EXISTS "idx_Votes_Creation" ON "Votes" ("Creation");
          `
		idxSqlite22 = `
          CREATE INDEX IF NOT EXISTS "idx_PublicKeys_Creation" ON "PublicKeys" ("Creation");
          `
		idxSqlite23 = `
          CREATE INDEX IF NOT EXISTS "idx_Truststates_Creation" ON "Truststates" ("Creation");
          `
		// Post's board index (thread and parent already indexed above.)
		idxSqlite24 = `
          CREATE INDEX IF NOT EXISTS "idx_Posts_Board" ON "Posts" ("Board");
          `
	} else {
		logging.LogCrash(fmt.Sprintf("Storage engine you've inputted is not supported. Please change it from the backend user config into something that is supported. You've provided: %s", globals.BackendConfig.GetDbEngine()))
	}

	var creationSchemas []string
	if globals.BackendConfig.GetDbEngine() == "sqlite" {
		creationSchemas = append(creationSchemas, schemaPrep1)
		creationSchemas = append(creationSchemas, schema1)
		creationSchemas = append(creationSchemas, schema3)
		creationSchemas = append(creationSchemas, schema4)
		creationSchemas = append(creationSchemas, schema5)
		creationSchemas = append(creationSchemas, schema6)
		creationSchemas = append(creationSchemas, schema7)
		creationSchemas = append(creationSchemas, schema8)
		creationSchemas = append(creationSchemas, schema9)
		creationSchemas = append(creationSchemas, schema10)
		creationSchemas = append(creationSchemas, schema11)
		creationSchemas = append(creationSchemas, schema12)
		creationSchemas = append(creationSchemas, schema16)
		creationSchemas = append(creationSchemas, idxSqlite1)
		creationSchemas = append(creationSchemas, idxSqlite2)
		creationSchemas = append(creationSchemas, idxSqlite3)
		creationSchemas = append(creationSchemas, idxSqlite4)
		creationSchemas = append(creationSchemas, idxSqlite5)
		creationSchemas = append(creationSchemas, idxSqlite6)
		creationSchemas = append(creationSchemas, idxSqlite7)
		creationSchemas = append(creationSchemas, idxSqlite8)
		creationSchemas = append(creationSchemas, idxSqlite9)
		creationSchemas = append(creationSchemas, idxSqlite10)
		creationSchemas = append(creationSchemas, idxSqlite11)
		creationSchemas = append(creationSchemas, idxSqlite12)
		creationSchemas = append(creationSchemas, idxSqlite13)
		creationSchemas = append(creationSchemas, idxSqlite14)
		creationSchemas = append(creationSchemas, idxSqlite15)
		creationSchemas = append(creationSchemas, idxSqlite16)
		creationSchemas = append(creationSchemas, idxSqlite17)
		creationSchemas = append(creationSchemas, idxSqlite18)
		creationSchemas = append(creationSchemas, idxSqlite19)
		creationSchemas = append(creationSchemas, idxSqlite20)
		creationSchemas = append(creationSchemas, idxSqlite21)
		creationSchemas = append(creationSchemas, idxSqlite22)
		creationSchemas = append(creationSchemas, idxSqlite23)
		creationSchemas = append(creationSchemas, idxSqlite24)
	} else if globals.BackendConfig.GetDbEngine() == "mysql" {
		creationSchemas = append(creationSchemas, schemaPrep1)
		creationSchemas = append(creationSchemas, schemaPrep2)
		creationSchemas = append(creationSchemas, schemaPrep3)
		creationSchemas = append(creationSchemas, schemaPrep4)
		creationSchemas = append(creationSchemas, schema1)
		creationSchemas = append(creationSchemas, schema3)
		creationSchemas = append(creationSchemas, schema4)
		creationSchemas = append(creationSchemas, schema5)
		creationSchemas = append(creationSchemas, schema6)
		creationSchemas = append(creationSchemas, schema7)
		creationSchemas = append(creationSchemas, schema8)
		creationSchemas = append(creationSchemas, schema9)
		creationSchemas = append(creationSchemas, schema10)
		creationSchemas = append(creationSchemas, schema11)
		creationSchemas = append(creationSchemas, schema12)
		creationSchemas = append(creationSchemas, schema16)
	}

	tx, err := globals.DbInstance.Beginx()
	if err != nil {
		logging.LogCrash(err)
	}
	for _, schema := range creationSchemas {
		// fmt.Println(schema)
		_, err2 := tx.Exec(schema)
		if err2 != nil {
			logging.LogCrash(err2)
		}
	}
	err3 := tx.Commit()
	if err3 != nil {
		tx.Rollback()
		logging.Log(1, fmt.Sprintf("CreateDatabase encountered an error when trying to commit to the database. Error is: %s", err3))
		if strings.Contains(err3.Error(), "database is locked") {
			logging.Log(1, fmt.Sprintf("This database seems to be locked. We'll sleep 10 seconds to give it the time it needs to recover. This mostly happens when the app has crashed and there is a hot journal - and SQLite is in the process of repairing the database. THE DATA IN THIS TRANSACTION WAS NOT COMMITTED. PLEASE RETRY."))
			return errors.New("Database was locked. THE DATA IN THIS TRANSACTION WAS NOT COMMITTED. PLEASE RETRY.")
		}
		return err3
	}
	return nil
}

func CheckDatabaseReady() {
	DiagInsert := `REPLACE INTO Diagnostics
  (
    DbRoundtripTestField
  ) VALUES (
    :DbRoundtripTestField
  )`
	DiagDelete := `DELETE FROM Diagnostics`
	// We're using time.now because we don't want DB to optimise out the write and not test the connection that way. Get a random number between 0- 65535 for entry test.
	ss := map[string]interface{}{"DbRoundtripTestField": toolbox.GetInsecureRand(65535)}
	// First, remove everything.
	tx, err := globals.DbInstance.Beginx()
	if err != nil {
		logging.LogCrash(err)
	}
	_, err2 := tx.Exec(DiagDelete)
	if err2 != nil {
		logging.LogCrash(err2)
	}
	tx.Commit()
	// Second, insert a new item.
	tx2, err3 := globals.DbInstance.Beginx()
	if err3 != nil {
		logging.LogCrash(err3)
	}
	_, err4 := tx2.NamedExec(DiagInsert, ss)
	if err4 != nil {
		logging.LogCrash(err4)
	}
	err5 := tx2.Commit()
	if err5 != nil {
		logging.LogCrash(err4)
	}
	logging.Logf(1, "Database is ready. Just verified by removing and inserting data successfully.")
}

// Insertion SQL code used by the writer.

// var nodeInsert = `
// REPLACE INTO Nodes
// (
//   Fingerprint, BoardsLastCheckin, ThreadsLastCheckin, PostsLastCheckin, VotesLastCheckin, KeysLastCheckin, TruststatesLastCheckin, AddressesLastCheckin
// ) VALUES (
//   :Fingerprint,
//   MAX(:BoardsLastCheckin, IFNULL((SELECT BoardsLastCheckin FROM Nodes WHERE Fingerprint = :Fingerprint), 0)),
//   MAX(:ThreadsLastCheckin, IFNULL((SELECT ThreadsLastCheckin FROM Nodes WHERE Fingerprint = :Fingerprint), 0)),
//   MAX(:PostsLastCheckin, IFNULL((SELECT PostsLastCheckin FROM Nodes WHERE Fingerprint = :Fingerprint), 0)),
//   MAX(:VotesLastCheckin, IFNULL((SELECT VotesLastCheckin FROM Nodes WHERE Fingerprint = :Fingerprint), 0)),
//   MAX(:KeysLastCheckin, IFNULL((SELECT KeysLastCheckin FROM Nodes WHERE Fingerprint = :Fingerprint), 0)),
//   MAX(:TruststatesLastCheckin, IFNULL((SELECT TruststatesLastCheckin FROM Nodes WHERE Fingerprint = :Fingerprint), 0)),
//   MAX(:AddressesLastCheckin, IFNULL((SELECT AddressesLastCheckin FROM Nodes WHERE Fingerprint = :Fingerprint), 0))
// )
// `

// NodeInsert just inserts the Node details into the entry. This is mutable.
var nodeInsert = `REPLACE INTO Nodes
(
  Fingerprint, BoardsLastCheckin, ThreadsLastCheckin, PostsLastCheckin,
  VotesLastCheckin, KeysLastCheckin, TruststatesLastCheckin, AddressesLastCheckin
) VALUES (
  :Fingerprint, :BoardsLastCheckin, :ThreadsLastCheckin, :PostsLastCheckin,
  :VotesLastCheckin, :KeysLastCheckin, :TruststatesLastCheckin, :AddressesLastCheckin
)`

/*
In the SQL statements below, the path of the execution is explained.
(Y): LastReferenced insert
BOARD > KEY(Y): Update board's key's lastreferenced timestamp.
*/

// ok, this is an attempted board insert

/*
  whats our exec model?

  attempt to select existing board
  if exists, use case to determine whether it validates
  if it does, insert board, remove board owners, and add new board owners.

*/

/*
   Important - If the board that is being inserted was not in the database before, we do not delete anything from BoardOwners. That is why the null cases are missing from this one.

   THE ORDER OF THE STATEMENTS DO MATTER - The main insert needs to happen the LAST, because after the insert succeeds, all of these IS NULL constraints that we rely on will start to fail, since at that point the candidate is already in, and now the candidate comparison becomes a comparison with itself - which candidate will fail.
*/

var boardInsert_BoardsBoardOwners_DeletePriors = `
WITH ExtantE(Creation, LastUpdate) AS (
  SELECT Creation, LastUpdate
  FROM Boards WHERE Fingerprint = :Fingerprint
)
DELETE FROM BoardOwners
WHERE (
  :LastUpdate > (SELECT LastUpdate FROM ExtantE) AND
  :LastUpdate > (SELECT Creation FROM ExtantE) AND
  :LastUpdate > :Creation AND
  BoardFingerprint = :Fingerprint
);
`
var boardInsert_BoardsKey_LastReferencedUpdate = `
/* Update ORIGINBOARD > KEY(Y) */
WITH ExtantE(Fingerprint, Creation, LastUpdate) AS (
SELECT Fingerprint, Creation, LastUpdate
FROM Boards WHERE Fingerprint = :Fingerprint
)
UPDATE PublicKeys
SET LastReferenced = :LastReferenced
WHERE (
    :LastUpdate > (SELECT LastUpdate FROM ExtantE) AND
    :LastUpdate > (SELECT Creation FROM ExtantE) AND
    :LastUpdate > :Creation AND
    Fingerprint = :Owner AND
    PublicKey = :OwnerPublicKey AND
    PublicKey = :OwnerPublicKey
  OR
    (SELECT Fingerprint FROM ExtantE) IS NULL AND
    :LastUpdate > :Creation AND
    Fingerprint = :Owner AND
    PublicKey = :OwnerPublicKey
  OR
    (SELECT Fingerprint FROM ExtantE) IS NULL AND
    :LastUpdate = 0 AND
    Fingerprint = :Owner AND
    PublicKey = :OwnerPublicKey
);
`
var boardInsert = `
REPLACE INTO Boards
SELECT Candidate.* FROM
(SELECT :Fingerprint AS Fingerprint,
       :Name AS Name,
       :Owner AS Owner,
       :OwnerPublicKey AS OwnerPublicKey,
       :Description AS Description,
       :Creation AS Creation,
       :ProofOfWork AS ProofOfWork,
       :Signature AS Signature,
       :LastUpdate AS LastUpdate,
       :UpdateProofOfWork AS UpdateProofOfWork,
       :UpdateSignature AS UpdateSignature,
       :LocalArrival AS LocalArrival,
       :LastReferenced AS LastReferenced,
       :EntityVersion AS EntityVersion,
       :Language AS Language,
       :Meta AS Meta,
       :RealmId AS RealmId,
       :EncrContent AS EncrContent
       ) AS Candidate
LEFT JOIN Boards ON Candidate.Fingerprint = Boards.Fingerprint
WHERE (
    Candidate.LastUpdate > Boards.LastUpdate AND
    Candidate.LastUpdate > Boards.Creation AND
    Candidate.LastUpdate > Candidate.Creation
  OR
    Boards.Fingerprint IS NULL AND
    Candidate.LastUpdate > Candidate.Creation
  OR
    Boards.Fingerprint IS NULL AND
    Candidate.LastUpdate = 0
);
`

// BoardOwners are mutable, but the condition of mutation is handled in the application layer. The only place the REPLACE could trigger is change of Expiry and level. The BoardFingerprint and KeyFingerprint are identity columns, so anything with different data on those will be committed as a new item.
/* This looks a little complex, but effectively, this is what is happening:
INSERT IF
    a) the board owner's board's candidate last update is newer than extant board's last update AND
    b) the board owner's board's candidate last update is newer than extant board's creation AND
    c) the board owner's board's candidate last update is newer than the candidate board's creation.
  OR
    a) This boardowner does not exist in DB (checking both primary keys) AND
    b) the board owner's board's candidate last update is newer than the candidate board's creation.
  OR
    a) This boardowner does not exist in DB (checking both primary keys) AND
    b) the board owner's board's candidate last update is zero, which means this object has never been updated.
*/
var boardOwnerInsert = `
REPLACE INTO BoardOwners
  WITH ExtantParent(Creation, LastUpdate) AS (
    SELECT Creation, LastUpdate
    FROM Boards WHERE Fingerprint = :BoardFingerprint
  )
  SELECT Candidate.* FROM
  (SELECT :BoardFingerprint AS BoardFingerprint,
          :KeyFingerprint AS KeyFingerprint,
          :Expiry AS Expiry,
          :Level AS Level
          ) AS Candidate
LEFT JOIN BoardOwners ON (
  (Candidate.BoardFingerprint = BoardOwners.BoardFingerprint)
  AND
  (Candidate.KeyFingerprint = BoardOwners.KeyFingerprint)
)
WHERE (
    :ParentBoardLastUpdate > (SELECT LastUpdate FROM ExtantParent) AND
    :ParentBoardLastUpdate > (SELECT Creation FROM ExtantParent) AND
    :ParentBoardLastUpdate > :ParentBoardCreation
  OR
    BoardOwners.BoardFingerprint IS NULL AND
    BoardOwners.KeyFingerprint IS NULL AND
    :ParentBoardLastUpdate > :ParentBoardCreation
  OR
    BoardOwners.BoardFingerprint IS NULL AND
    BoardOwners.KeyFingerprint IS NULL AND
    :ParentBoardLastUpdate = 0
);
`

// Deletion for BoardOwner. This triggers when a person is no longer a moderator, etc. This is the only deletion here because nothing else really gets deleted.

// Deletion for board owner has moved into the board insert. A successful board insertion (both update and creation) will delete all board owners for that board and will reinsert.

var threadInsert_ThreadsBoard_LastReferencedUpdate = `
/* Update ORIGINTHREAD > BOARD(Y) */
WITH ExtantE(Fingerprint, Creation, LastUpdate) AS (
  SELECT Fingerprint, Creation, LastUpdate
  FROM Threads WHERE Fingerprint = :Fingerprint
)
UPDATE Boards
SET LastReferenced = :LastReferenced
WHERE (
    :LastUpdate > (SELECT LastUpdate FROM ExtantE) AND
    :LastUpdate > (SELECT Creation FROM ExtantE) AND
    :LastUpdate > :Creation AND
    Fingerprint = :Board
  OR
    (SELECT Fingerprint FROM ExtantE) IS NULL AND
    :LastUpdate > :Creation AND
    Fingerprint = :Board
  OR
    (SELECT Fingerprint FROM ExtantE) IS NULL AND
    :LastUpdate = 0 AND
    Fingerprint = :Board
);
`
var threadInsert_ThreadsKey_LastReferencedUpdate = `
/* Update ORIGINTHREAD > KEY(Y) */
WITH ExtantE(Fingerprint, Creation, LastUpdate) AS (
  SELECT Fingerprint, Creation, LastUpdate
  FROM Threads WHERE Fingerprint = :Fingerprint
)
UPDATE PublicKeys
SET LastReferenced = :LastReferenced
WHERE (
    :LastUpdate > (SELECT LastUpdate FROM ExtantE) AND
    :LastUpdate > (SELECT Creation FROM ExtantE) AND
    :LastUpdate > :Creation AND
    Fingerprint = :Owner AND
    PublicKey = :OwnerPublicKey
  OR
    (SELECT Fingerprint FROM ExtantE) IS NULL AND
    :LastUpdate > :Creation AND
    Fingerprint = :Owner AND
    PublicKey = :OwnerPublicKey
  OR
    (SELECT Fingerprint FROM ExtantE) IS NULL AND
    :LastUpdate = 0 AND
    Fingerprint = :Owner AND
    PublicKey = :OwnerPublicKey
);
`
var threadInsert_ThreadsBoardsKey_LastReferencedUpdate = `
/* Update ORIGINTHREAD > BOARD > KEY(Y) */
WITH ExtantE(Fingerprint, Creation, LastUpdate) AS (
  SELECT Fingerprint, Creation, LastUpdate
  FROM Threads WHERE Fingerprint = :Fingerprint
),
ParentBoard AS (
  SELECT Owner FROM Boards
  WHERE Fingerprint = :Board
)
UPDATE PublicKeys
SET LastReferenced = :LastReferenced
WHERE (
    :LastUpdate > (SELECT LastUpdate FROM ExtantE) AND
    :LastUpdate > (SELECT Creation FROM ExtantE) AND
    :LastUpdate > :Creation AND
    Fingerprint = (SELECT Owner FROM ParentBoard) AND
    PublicKey = :OwnerPublicKey
  OR
    (SELECT Fingerprint FROM ExtantE) IS NULL AND
    :LastUpdate > :Creation AND
    Fingerprint = (SELECT Owner FROM ParentBoard) AND
    PublicKey = :OwnerPublicKey
  OR
    (SELECT Fingerprint FROM ExtantE) IS NULL AND
    :LastUpdate = 0 AND
    Fingerprint = (SELECT Owner FROM ParentBoard) AND
    PublicKey = :OwnerPublicKey
);
`
var threadInsert = `
REPLACE INTO Threads
  SELECT Candidate.* FROM
  (SELECT :Fingerprint AS Fingerprint,
          :Board AS Board,
          :Name AS Name,
          :Body AS Body,
          :Link AS Link,
          :Owner AS Owner,
          :OwnerPublicKey AS OwnerPublicKey,
          :Creation AS Creation,
          :ProofOfWork AS ProofOfWork,
          :Signature AS Signature,
          :LastUpdate AS LastUpdate,
          :UpdateProofOfWork AS UpdateProofOfWork,
          :UpdateSignature AS UpdateSignature,
          :LocalArrival AS LocalArrival,
          :LastReferenced AS LastReferenced,
          :EntityVersion AS EntityVersion,
          :Meta AS Meta,
          :RealmId AS RealmId,
          :EncrContent AS EncrContent
          ) AS Candidate
LEFT JOIN Threads ON Candidate.Fingerprint = Threads.Fingerprint
WHERE (
    Candidate.LastUpdate > Threads.LastUpdate AND
    Candidate.LastUpdate > Threads.Creation
  OR
    Threads.Fingerprint IS NULL AND
    Candidate.LastUpdate > Candidate.Creation
  OR
    Threads.Fingerprint IS NULL AND
    Candidate.LastUpdate = 0
);
`

var postInsert_PostsBoard_LastReferencedUpdate = `
/* Update ORIGINPOST > BOARD(Y) */
WITH ExtantE(Fingerprint, Creation, LastUpdate) AS (
  SELECT Fingerprint, Creation, LastUpdate
  FROM Posts WHERE Fingerprint = :Fingerprint
)
UPDATE Boards
SET LastReferenced = :LastReferenced
WHERE (
    :LastUpdate > (SELECT LastUpdate FROM ExtantE) AND
    :LastUpdate > (SELECT Creation FROM ExtantE) AND
    :LastUpdate > :Creation AND
    Fingerprint = :Board
  OR
    (SELECT Fingerprint FROM ExtantE) IS NULL AND
    :LastUpdate > :Creation AND
    Fingerprint = :Board
  OR
    (SELECT Fingerprint FROM ExtantE) IS NULL AND
    :LastUpdate = 0 AND
    Fingerprint = :Board
);
`
var postInsert_PostsBoardsKey_LastReferencedUpdate = `
/* Update ORIGINPOST > BOARD > KEY(Y) */
WITH ExtantE(Fingerprint, Creation, LastUpdate) AS (
  SELECT Fingerprint, Creation, LastUpdate
  FROM Posts WHERE Fingerprint = :Fingerprint
),
ParentBoard AS (
  SELECT Owner FROM Boards
  WHERE Fingerprint = :Board
)
UPDATE PublicKeys
SET LastReferenced = :LastReferenced
WHERE (
    :LastUpdate > (SELECT LastUpdate FROM ExtantE) AND
    :LastUpdate > (SELECT Creation FROM ExtantE) AND
    :LastUpdate > :Creation AND
    Fingerprint = (SELECT Owner FROM ParentBoard) AND
    PublicKey = :OwnerPublicKey
  OR
    (SELECT Fingerprint FROM ExtantE) IS NULL AND
    :LastUpdate > :Creation AND
    Fingerprint = (SELECT Owner FROM ParentBoard) AND
    PublicKey = :OwnerPublicKey
  OR
    (SELECT Fingerprint FROM ExtantE) IS NULL AND
    :LastUpdate = 0 AND
    Fingerprint = (SELECT Owner FROM ParentBoard) AND
    PublicKey = :OwnerPublicKey
);
`
var postInsert_PostsThread_LastReferencedUpdate = `
/* Update ORIGINPOST > THREAD(Y) */
WITH ExtantE(Fingerprint, Creation, LastUpdate) AS (
  SELECT Fingerprint, Creation, LastUpdate
  FROM Posts WHERE Fingerprint = :Fingerprint
)
UPDATE Threads
SET LastReferenced = :LastReferenced
WHERE (
    :LastUpdate > (SELECT LastUpdate FROM ExtantE) AND
    :LastUpdate > (SELECT Creation FROM ExtantE) AND
    :LastUpdate > :Creation AND
    Fingerprint = :Thread
  OR
    (SELECT Fingerprint FROM ExtantE) IS NULL AND
    :LastUpdate > :Creation AND
    Fingerprint = :Thread
  OR
    (SELECT Fingerprint FROM ExtantE) IS NULL AND
    :LastUpdate = 0 AND
    Fingerprint = :Thread
);
`
var postInsert_PostsThreadsKey_LastReferencedUpdate = `
/* Update ORIGINPOST > THREAD > KEY(Y) */
WITH ExtantE(Fingerprint, Creation, LastUpdate) AS (
  SELECT Fingerprint, Creation, LastUpdate
  FROM Posts WHERE Fingerprint = :Fingerprint
),
ParentThread AS (
  SELECT Owner FROM Boards
  WHERE Fingerprint = :Thread
)
UPDATE PublicKeys
SET LastReferenced = :LastReferenced
WHERE (
    :LastUpdate > (SELECT LastUpdate FROM ExtantE) AND
    :LastUpdate > (SELECT Creation FROM ExtantE) AND
    :LastUpdate > :Creation AND
    Fingerprint = (SELECT Owner FROM ParentThread) AND
    PublicKey = :OwnerPublicKey
  OR
    (SELECT Fingerprint FROM ExtantE) IS NULL AND
    :LastUpdate > :Creation AND
    Fingerprint = (SELECT Owner FROM ParentThread) AND
    PublicKey = :OwnerPublicKey
  OR
    (SELECT Fingerprint FROM ExtantE) IS NULL AND
    :LastUpdate = 0 AND
    Fingerprint = (SELECT Owner FROM ParentThread) AND
    PublicKey = :OwnerPublicKey
);
`
var postInsert_PostsKey_LastReferencedUpdate = `
/* Update ORIGINPOST > KEY(Y) */
WITH ExtantE(Fingerprint, Creation, LastUpdate) AS (
  SELECT Fingerprint, Creation, LastUpdate
  FROM Posts WHERE Fingerprint = :Fingerprint
)
UPDATE PublicKeys
SET LastReferenced = :LastReferenced
WHERE (
    :LastUpdate > (SELECT LastUpdate FROM ExtantE) AND
    :LastUpdate > (SELECT Creation FROM ExtantE) AND
    :LastUpdate > :Creation AND
    Fingerprint = :Owner AND
    PublicKey = :OwnerPublicKey
  OR
    (SELECT Fingerprint FROM ExtantE) IS NULL AND
    :LastUpdate > :Creation AND
    Fingerprint = :Owner AND
    PublicKey = :OwnerPublicKey
  OR
    (SELECT Fingerprint FROM ExtantE) IS NULL AND
    :LastUpdate = 0 AND
    Fingerprint = :Owner AND
    PublicKey = :OwnerPublicKey
);
`
var postInsert_PostsPosts_Recursive_LastReferencedUpdate = `
/*
Update ORIGINPOST > POST (Y) > ...(Y) > POST(Y)
(post's parent post chain, if any)
*/

WITH RECURSIVE ExtantE(Fingerprint, Creation, LastUpdate) AS (
  SELECT Fingerprint, Creation, LastUpdate
  FROM Posts WHERE Fingerprint = :Fingerprint
),
ParentOf(Fingerprint, ParentPostFingerprint) AS (
  SELECT Fingerprint, Parent FROM Posts
),
AncestorOf(Fingerprint) AS (
  SELECT CASE
  WHEN (SELECT ParentPostFingerprint FROM ParentOf WHERE Fingerprint = :Fingerprint) IS NULL
  THEN (SELECT :Parent)
  ELSE (SELECT ParentPostFingerprint FROM ParentOf WHERE Fingerprint = :Fingerprint)
  END AS ParentPostFingerprint
  UNION ALL
  SELECT ParentPostFingerprint FROM ParentOf
  JOIN AncestorOf
  ON AncestorOf.Fingerprint=ParentOf.Fingerprint
),
FinalTable(Fingerprint) AS (
  SELECT Posts.Fingerprint FROM AncestorOf, Posts
  WHERE
  AncestorOf.Fingerprint=Posts.Fingerprint
)
UPDATE Posts
SET LastReferenced=:LastReferenced
WHERE (
    :LastUpdate > (SELECT LastUpdate FROM ExtantE) AND
    :LastUpdate > (SELECT Creation FROM ExtantE) AND
    :LastUpdate > :Creation AND
    Posts.Fingerprint IN (SELECT Fingerprint FROM FinalTable)
  OR
    (SELECT Fingerprint FROM ExtantE) IS NULL AND
    :LastUpdate > :Creation AND
    Posts.Fingerprint IN (SELECT Fingerprint FROM FinalTable)
  OR
    (SELECT Fingerprint FROM ExtantE) IS NULL AND
    :LastUpdate = 0 AND
    Posts.Fingerprint IN (SELECT Fingerprint FROM FinalTable)
);
`
var postInsert_PostsPostsKeys_Recursive_LastReferencedUpdate = `
/*
Update ORIGINPOST > POST > KEY(Y)
                  > POST > KEY(Y)
                  > ...
                  > POST > KEY(Y)
(post's parent post chain's keys, if any)
*/

WITH RECURSIVE ExtantE(Fingerprint, Creation, LastUpdate) AS (
  SELECT Fingerprint, Creation, LastUpdate
  FROM Posts WHERE Fingerprint = :Fingerprint
),
ParentOf(Fingerprint, ParentPostFingerprint, Owner) AS (
  SELECT Fingerprint, Parent, Owner FROM Posts
),
AncestorOf(Fingerprint) AS (
  SELECT CASE
  WHEN (SELECT ParentPostFingerprint FROM ParentOf WHERE Fingerprint = :Fingerprint) IS NULL
  THEN (SELECT :Parent)
  ELSE (SELECT ParentPostFingerprint FROM ParentOf WHERE Fingerprint = :Fingerprint)
  END AS ParentPostFingerprint
  UNION ALL
  SELECT ParentPostFingerprint FROM ParentOf
  JOIN AncestorOf
  ON AncestorOf.Fingerprint=ParentOf.Fingerprint
),
FinalTable(Fingerprint, Owner) AS (
  SELECT Posts.Fingerprint, Posts.Owner FROM AncestorOf, Posts
  WHERE
  AncestorOf.Fingerprint=Posts.Fingerprint
)
UPDATE PublicKeys
SET LastReferenced=:LastReferenced
WHERE (
    :LastUpdate > (SELECT LastUpdate FROM ExtantE) AND
    :LastUpdate > (SELECT Creation FROM ExtantE) AND
    :LastUpdate > :Creation AND
    PublicKeys.Fingerprint IN (SELECT Owner FROM FinalTable) AND
    PublicKey = :OwnerPublicKey
  OR
    (SELECT Fingerprint FROM ExtantE) IS NULL AND
    :LastUpdate > :Creation AND
    PublicKeys.Fingerprint IN (SELECT Owner FROM FinalTable) AND
    PublicKey = :OwnerPublicKey
  OR
    (SELECT Fingerprint FROM ExtantE) IS NULL AND
    :LastUpdate = 0 AND
    PublicKeys.Fingerprint IN (SELECT Owner FROM FinalTable) AND
    PublicKey = :OwnerPublicKey
);
`
var postInsert = `
REPLACE INTO Posts
SELECT Candidate.* FROM
(SELECT :Fingerprint AS Fingerprint,
        :Board AS Board,
        :Thread AS Thread,
        :Parent AS Parent,
        :Body AS Body,
        :Owner AS Owner,
        :OwnerPublicKey AS OwnerPublicKey,
        :Creation AS Creation,
        :ProofOfWork AS ProofOfWork,
        :Signature AS Signature,
        :LastUpdate AS LastUpdate,
        :UpdateProofOfWork AS UpdateProofOfWork,
        :UpdateSignature AS UpdateSignature,
        :LocalArrival AS LocalArrival,
        :LastReferenced AS LastReferenced,
        :EntityVersion AS EntityVersion,
        :Meta AS Meta,
        :RealmId AS RealmId,
        :EncrContent AS EncrContent
        ) AS Candidate
LEFT JOIN Posts ON Candidate.Fingerprint = Posts.Fingerprint
WHERE (
    Candidate.LastUpdate > Posts.LastUpdate AND
    Candidate.LastUpdate > Posts.Creation
  OR
    Posts.Fingerprint IS NULL AND
    Candidate.LastUpdate > Candidate.Creation
  OR
    Posts.Fingerprint IS NULL AND
    Candidate.LastUpdate = 0
);
`
var voteInsert_VotesKey_LastReferencedUpdate = `
/* Update ORIGINVOTE > KEY(Y) */
WITH ExtantE(Fingerprint, Creation, LastUpdate) AS (
  SELECT Fingerprint, Creation, LastUpdate
  FROM Votes WHERE Fingerprint = :Fingerprint
)
UPDATE PublicKeys
SET LastReferenced = :LastReferenced
WHERE (
    :LastUpdate > (SELECT LastUpdate FROM ExtantE) AND
    :LastUpdate > (SELECT Creation FROM ExtantE) AND
    :LastUpdate > :Creation AND
    Fingerprint = :Owner AND
    PublicKey = :OwnerPublicKey
  OR
    (SELECT Fingerprint FROM ExtantE) IS NULL AND
    :LastUpdate > :Creation AND
    Fingerprint = :Owner AND
    PublicKey = :OwnerPublicKey
  OR
    (SELECT Fingerprint FROM ExtantE) IS NULL AND
    :LastUpdate = 0 AND
    Fingerprint = :Owner AND
    PublicKey = :OwnerPublicKey
);
`
var voteInsert = `
REPLACE INTO Votes
SELECT Candidate.* FROM
(SELECT :Fingerprint AS Fingerprint,
        :Board AS Board,
        :Thread AS Thread,
        :Target AS Target,
        :Owner AS Owner,
        :OwnerPublicKey AS OwnerPublicKey,
        :TypeClass AS TypeClass,
        :Type AS Type,
        :Creation AS Creation,
        :ProofOfWork AS ProofOfWork,
        :Signature AS Signature,
        :LastUpdate AS LastUpdate,
        :UpdateProofOfWork AS UpdateProofOfWork,
        :UpdateSignature AS UpdateSignature,
        :LocalArrival AS LocalArrival,
        :LastReferenced AS LastReferenced,
        :EntityVersion AS EntityVersion,
        :Meta AS Meta,
        :RealmId AS RealmId,
        :EncrContent AS EncrContent
        ) AS Candidate
LEFT JOIN Votes ON Candidate.Fingerprint = Votes.Fingerprint
WHERE (
    Candidate.LastUpdate > Votes.LastUpdate AND
    Candidate.LastUpdate > Votes.Creation AND
    Candidate.LastUpdate > Candidate.Creation
  OR
    Votes.Fingerprint IS NULL AND
    Candidate.LastUpdate > Candidate.Creation
  OR
    Votes.Fingerprint IS NULL AND Candidate.LastUpdate = 0
);
`

// above, you nede to basically do everything that a post does, update board, thread, key, target, board key, thread key, target key (target being a post or a thread, attempt both.)

// actually scratch that. voting on something should only light up the key that created the vote, nothing else.

var keyInsert = `
REPLACE INTO PublicKeys
SELECT Candidate.* FROM
(SELECT :Fingerprint AS Fingerprint,
        :Type AS Type,
        :PublicKey AS PublicKey,
        :Expiry AS Expiry,
        :Name AS Name,
        :Info AS Info,
        :Creation AS Creation,
        :ProofOfWork AS ProofOfWork,
        :Signature AS Signature,
        :LastUpdate AS LastUpdate,
        :UpdateProofOfWork AS UpdateProofOfWork,
        :UpdateSignature AS UpdateSignature,
        :LocalArrival AS LocalArrival,
        :LastReferenced AS LastReferenced,
        :EntityVersion AS EntityVersion,
        :Meta AS Meta,
        :RealmId AS RealmId,
        :EncrContent AS EncrContent
        ) AS Candidate
LEFT JOIN PublicKeys ON Candidate.Fingerprint = PublicKeys.Fingerprint
WHERE (
    Candidate.LastUpdate > PublicKeys.LastUpdate AND
      Candidate.LastUpdate > PublicKeys.Creation
  OR
    PublicKeys.Fingerprint IS NULL AND
    Candidate.LastUpdate > Candidate.Creation
  OR
    PublicKeys.Fingerprint IS NULL AND
    Candidate.LastUpdate = 0
);
`
var truststateInsert_TruststatesKey_LastReferencedUpdate = `
/* Update ORIGINTS > KEY(Y) */
WITH ExtantE(Fingerprint, Creation, LastUpdate) AS (
  SELECT Fingerprint, Creation, LastUpdate
  FROM Truststates WHERE Fingerprint = :Fingerprint
)
UPDATE PublicKeys
SET LastReferenced = :LastReferenced
WHERE (
    :LastUpdate > (SELECT LastUpdate FROM ExtantE) AND
    :LastUpdate > (SELECT Creation FROM ExtantE) AND
    :LastUpdate > :Creation AND
    Fingerprint = :Owner AND
    PublicKey = :OwnerPublicKey
  OR
    (SELECT Fingerprint FROM ExtantE) IS NULL AND
    :LastUpdate > :Creation AND
    Fingerprint = :Owner AND
    PublicKey = :OwnerPublicKey
  OR
    (SELECT Fingerprint FROM ExtantE) IS NULL AND
    :LastUpdate = 0 AND
    Fingerprint = :Owner AND
    PublicKey = :OwnerPublicKey
);
`
var truststateInsert_TruststatesTargetKey_LastReferencedUpdate = `
/* Update ORIGINTS > TARGET(KEY)(Y) */
WITH ExtantE(Fingerprint, Creation, LastUpdate) AS (
  SELECT Fingerprint, Creation, LastUpdate
  FROM Truststates WHERE Fingerprint = :Fingerprint
)
UPDATE PublicKeys
SET LastReferenced = :LastReferenced
WHERE (
    :LastUpdate > (SELECT LastUpdate FROM ExtantE) AND
    :LastUpdate > (SELECT Creation FROM ExtantE) AND
    :LastUpdate > :Creation AND
    Fingerprint = :Target AND
    PublicKey = :OwnerPublicKey
  OR
    (SELECT Fingerprint FROM ExtantE) IS NULL AND
    :LastUpdate > :Creation AND
    Fingerprint = :Target AND
    PublicKey = :OwnerPublicKey
  OR
    (SELECT Fingerprint FROM ExtantE) IS NULL AND
    :LastUpdate = 0 AND
    Fingerprint = :Target AND
    PublicKey = :OwnerPublicKey
);
`
var truststateInsert = `
REPLACE INTO Truststates
SELECT Candidate.* FROM
(SELECT :Fingerprint AS Fingerprint,
        :Target AS Target,
        :Owner AS Owner,
        :OwnerPublicKey AS OwnerPublicKey,
        :TypeClass AS TypeClass,
        :Type AS Type,
        :Domain AS Domain,
        :Expiry AS Expiry,
        :Creation AS Creation,
        :ProofOfWork AS ProofOfWork,
        :Signature AS Signature,
        :LastUpdate AS LastUpdate,
        :UpdateProofOfWork AS UpdateProofOfWork,
        :UpdateSignature AS UpdateSignature,
        :LocalArrival AS LocalArrival,
        :LastReferenced AS LastReferenced,
        :EntityVersion AS EntityVersion,
        :Meta AS Meta,
        :RealmId AS RealmId,
        :EncrContent AS EncrContent
        ) AS Candidate
LEFT JOIN Truststates ON Candidate.Fingerprint = Truststates.Fingerprint
WHERE (
    Candidate.LastUpdate > Truststates.LastUpdate AND
    Candidate.LastUpdate > Truststates.Creation
  OR
    Truststates.Fingerprint IS NULL AND
    Candidate.LastUpdate > Candidate.Creation
  OR
    Truststates.Fingerprint IS NULL AND
    Candidate.LastUpdate = 0
);
`

// Address insert is immutable. This is used for when a node receives data from an address from a node that is not at the aforementioned address. In other words, an address object coming from a third party node not at that address cannot change an existing address saved in the database.
var addressInsertMySQL = `
INSERT IGNORE INTO Addresses
(
  Location,
  Sublocation,
  Port, IPType,
  AddressType,
  LastSuccessfulPing,
  LastSuccessfulSync,
  ProtocolVersionMajor,
  ProtocolVersionMinor,
  ClientVersionMajor,
  ClientVersionMinor,
  ClientVersionPatch,
  ClientName,
  EntityVersion,
  RealmId,
  LocalArrival
) VALUES (
  :Location,
  :Sublocation,
  :Port,
  :IPType,
  :AddressType,
  :LastSuccessfulPing,
  :LastSuccessfulSync,
  :ProtocolVersionMajor,
  :ProtocolVersionMinor,
  :ClientVersionMajor,
  :ClientVersionMinor,
  :ClientVersionPatch,
  :ClientName,
  :EntityVersion,
  :RealmId,
  :LocalArrival
)
  `

var addressInsertSQLite = `
INSERT OR IGNORE INTO Addresses
(
  Location,
  Sublocation,
  Port,
  IPType,
  AddressType,
  LastSuccessfulPing,
  LastSuccessfulSync,
  ProtocolVersionMajor,
  ProtocolVersionMinor,
  ClientVersionMajor,
  ClientVersionMinor,
  ClientVersionPatch,
  ClientName,
  EntityVersion,
  RealmId,
  LocalArrival
) VALUES (
  :Location,
  :Sublocation,
  :Port,
  :IPType,
  :AddressType,
  :LastSuccessfulPing,
  :LastSuccessfulSync,
  :ProtocolVersionMajor,
  :ProtocolVersionMinor,
  :ClientVersionMajor,
  :ClientVersionMinor,
  :ClientVersionPatch,
  :ClientName,
  :EntityVersion,
  :RealmId,
  :LocalArrival
)
  `

// Address update insert is mutable. This is used when the node connects to the address itself. Example: When a node connects to 256.253.231.123:8080, it will update the entry for that address with the data coming from the remote node. This is the only way to mutate an address object.

/* In the WHERE clause below, we do not check for inequivalency of location, sublocation and port because we do the join based on that, they're guaranteed to be the same. We are also explicitly NOT checking for the local arrival being different because this is the entire point of this comparison: if the only thing that changes is the local arrival, we do NOT insert.

Why?

Because a static node's LastSuccessfulSync only gets updated to its timestamp. That means LastSuccessfulSync can actually be in the past. If we generate responses and caches based on LastSuccessfulSync, what happens is that it might actually be placed outside that cache's or responses scope, and in the scope of caches / responses that are already generated. Which effectively means it won't be communicated.

If we do the opposite, and scan based on localarrival, then every time a node is communicated with, the localarrival will update. Even when there are no updates. That means when a static node is hit, it will be placed to the top of the list of addresses that node has connected to and will be served at the top. Considering that static nodes stay online 24/7, this will cause the addresses list that this node serves to be mostly / all made of static nodes.

To prevent that, what we do is that every time LastSuccessfulSync changes in an address we update the localarrival, but ONLY then. If anything changes in an address except localarrival then we update localarrival as well as the changes, but if it's just localarrival change and nothing else, then we do not update.

This solves our problem, because if a static node is updated, the LastSuccessfulSync will be set to its timestamp and the localarrival will be set to now. since we scan by localarrival, it is guaranteed to be included in the response, but so long as the LastSuccessfulSync does not change, that will be it. It won't update every time the node is connected to.
*/
var addressUpdateInsert = `
REPLACE INTO Addresses
SELECT Candidate.* FROM
(SELECT :Location AS Location,
        :Sublocation AS Sublocation,
        :Port AS Port,
        :IPType AS IPType,
        :AddressType AS AddressType,
        MAX(:LastSuccessfulPing,
          IFNULL(
            (SELECT LastSuccessfulPing FROM Addresses WHERE Location = :Location AND Sublocation=:Sublocation AND Port=:Port) , 0
          )
        )
        AS LastSuccessfulPing,
        MAX(:LastSuccessfulSync,
          IFNULL(
            (SELECT LastSuccessfulSync FROM Addresses WHERE Location = :Location AND Sublocation=:Sublocation AND Port=:Port) , 0
          )
        )
        AS LastSuccessfulSync,
        :ProtocolVersionMajor AS ProtocolVersionMajor,
        :ProtocolVersionMinor AS ProtocolVersionMinor,
        :ClientVersionMajor AS ClientVersionMajor,
        :ClientVersionMinor AS ClientVersionMinor,
        :ClientVersionPatch AS ClientVersionPatch,
        :ClientName AS ClientName,
        :LocalArrival AS LocalArrival,
        :EntityVersion AS EntityVersion,
        :RealmId AS RealmId
        ) AS Candidate
LEFT JOIN Addresses ON Candidate.Location = Addresses.Location AND Candidate.Sublocation = Addresses.Sublocation AND Candidate.Port = Addresses.Port
WHERE (
      Candidate.IPType != Addresses.IPType OR
      Candidate.AddressType != Addresses.AddressType OR
      Candidate.LastSuccessfulPing != Addresses.LastSuccessfulPing OR
      Candidate.LastSuccessfulSync != Addresses.LastSuccessfulSync OR
      Candidate.ProtocolVersionMajor != Addresses.ProtocolVersionMajor OR
      Candidate.ProtocolVersionMinor != Addresses.ProtocolVersionMinor OR
      Candidate.ClientVersionMajor != Addresses.ClientVersionMajor OR
      Candidate.ClientVersionMinor != Addresses.ClientVersionMinor OR
      Candidate.ClientVersionPatch != Addresses.ClientVersionPatch OR
      Candidate.ClientName != Addresses.ClientName OR
      Candidate.EntityVersion != Addresses.EntityVersion OR
      Candidate.RealmId != Addresses.RealmId
    OR
      Addresses.Location IS NULL AND
      Addresses.Sublocation IS NULL AND
      Addresses.Port IS NULL
);
`

// Subprotocol insert is the part of address insertion series. This makes it so that we have a list of all subprotocols flying around.
var subprotocolInsert = `
REPLACE INTO Subprotocols
(
  Fingerprint,
  Name,
  VersionMajor,
  VersionMinor,
  SupportedEntities
) VALUES (
  :Fingerprint,
  :Name,
  :VersionMajor,
  :VersionMinor,
  :SupportedEntities
)
  `

// AddressSubprotocolInsert inserts into the many to many junction table so that we can keep track of the subprotocols an address uses.

var addressSubprotocolInsertMySQL = `
INSERT IGNORE INTO AddressesSubprotocols
(
 AddressLocation,
 AddressSublocation,
 AddressPort,
 SubprotocolFingerprint
) VALUES (
 :AddressLocation,
 :AddressSublocation,
 :AddressPort,
 :SubprotocolFingerprint
)
`

var addressSubprotocolInsertSQLite = `
INSERT OR IGNORE INTO AddressesSubprotocols
(
  AddressLocation,
  AddressSublocation,
  AddressPort,
  SubprotocolFingerprint
) VALUES (
  :AddressLocation,
  :AddressSublocation,
  :AddressPort,
  :SubprotocolFingerprint
)
`

// var addressPrune = `
// DELETE FROM Addresses ORDER BY LastSuccessfulPing ASC, LastSuccessfulSync ASC LIMIT (MAX(-(? - (SELECT COUNT(*) FROM Addresses)),0))
// `
var addressPrune = `
DELETE FROM Addresses
WHERE (Location, Sublocation, Port)
NOT IN
  (
  SELECT Location, Sublocation, Port
  FROM Addresses
  ORDER BY LastSuccessfulPing DESC,
           LastSuccessfulSync DESC
  LIMIT ?
  )`
