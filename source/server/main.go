// Copyright 2019 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package main initializes and starts the Emitto service.
package main

import (
	"context"
	"flag"
	"fmt"
	"net"

	"cloud.google.com/go/storage"
	"github.com/google/emitto/source/filestore"
	"github.com/google/emitto/source/server/fleetspeak"
	"github.com/google/emitto/source/server/service"
	"github.com/google/emitto/source/server/store"
	"google.golang.org/grpc"

	log "github.com/golang/glog"
	pb "github.com/google/emitto/source/server/proto"
	fspb "github.com/google/fleetspeak/fleetspeak/src/server/grpcservice/proto/fleetspeak_grpcservice"
)

var (
	// Server flags.
	port          = flag.Int("port", 4444, "Emitto server port")
	fsAdminAddr   = flag.String("admin_addr", "", "Fleetspeak admin server")
	memoryStorage = flag.Bool("memory_storage", false, "Use memory store and filestore")

	// Google Cloud Project flags.
	projectID     = flag.String("project_id", "", "Google Cloud project ID")
	storageBucket = flag.String("storage_bucket", "", "Google Cloud Storage bucket for storing rule files")
	credFile      = flag.String("cred_file", "", "Path of the JSON application credential file.")

	// Fleetspeak flags.
	certFile = flag.String("cert_file", "", "Path of the Fleetspeak certificate file")
)

func main() {
	ctx := context.Background()

	a, err := fleetspeak.New(*fsAdminAddr, *certFile)
	if err != nil {
		log.Fatalf("unable to connect to the Fleetspeak admin server: %v", err)
	}
	defer a.Close()

	fs, closeFStore := mustGetFileStore(ctx)
	defer closeFStore()

	s, closeStore := mustGetStore(ctx)
	defer closeStore()

	server := grpc.NewServer(nil)
	svc := service.New(s, fs, a)
	pb.RegisterEmittoServer(server, svc)
	fspb.RegisterProcessorServer(server, svc)

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("server failed to listen: %v", err)
	}
	defer l.Close()
	server.Serve(l)
}

func mustGetFileStore(ctx context.Context) (filestore.FileStore, func() error) {
	if *memoryStorage {
		return filestore.NewMemoryFileStore(), func() error { return nil }
	}
	c, err := filestore.NewGCSClient(ctx, *credFile, []string{storage.ScopeFullControl})
	if err != nil {
		log.Fatalf("failed to create Google Cloud Storage client: %v", err)
	}
	return filestore.NewGCSFileStore(*storageBucket, c), c.Close
}

func mustGetStore(ctx context.Context) (store.Store, func() error) {
	if *memoryStorage {
		return store.NewMemoryStore(), func() error { return nil }
	}
	c, err := store.NewGCDClient(ctx, *projectID, *credFile)
	if err != nil {
		log.Fatalf("failed to create Google Cloud Datastore client: %v", err)
	}
	return store.NewDataStore(c), c.Close
}
