// Copyright 2018 Banco Bilbao Vizcaya Argentaria, S.A.

// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at

//     http://www.apache.org/licenses/LICENSE-2.0

// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
// Copyright 2018 Banco Bilbao Vizcaya Argentaria, S.A.

// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at

//     http://www.apache.org/licenses/LICENSE-2.0

// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cli

import (
	"github.com/spf13/cobra"

	"github.com/bbva/qed/balloon/storage"
	"github.com/bbva/qed/log"
	"github.com/bbva/qed/server"
)

func NewServerCommand() *cobra.Command {
	var (
		logLevel, endpoint, dbPath, apiKey, storageName string
		cacheSize                                       uint64
		tamperable                                      bool
	)

	cmd := &cobra.Command{
		Use:   "server",
		Short: "The server for the verifiable log QED",
		Long:  ``,
		// Args:  cobra.NoArgs(),

		Run: func(cmd *cobra.Command, args []string) {

			log.SetLogger("QedServer", logLevel)

			s := server.NewServer(
				endpoint,
				dbPath,
				apiKey,
				cacheSize,
				storageName,
			)

			err := s.ListenAndServe()
			if err != nil {
				log.Errorf("Can't start HTTP Server: ", err)
			}

		},
	}

	cmd.Flags().StringVarP(&apiKey, "apikey", "k", "", "Server api key")
	cmd.Flags().StringVarP(&endpoint, "endpoint", "e", "0.0.0.0:8080", "Endpoint for REST requests on (host:port)")
	cmd.Flags().StringVarP(&dbPath, "path", "p", "/var/tmp/balloon.db", "Set default storage path.")
	cmd.Flags().Uint64VarP(&cacheSize, "cache", "c", storage.SIZE25, "Initialize and reserve custom cache size.")
	cmd.Flags().StringVarP(&storageName, "storage", "s", "badger", "Choose between different storage backends. Eg badge|bolt")
	cmd.Flags().StringVarP(&logLevel, "log", "l", "error", "Choose between log levels: silent, error, info and debug")

	// INFO: testing purposes
	cmd.Flags().BoolVarP(&tamperable, "tamperable", "t", false, "Allow a tamperable api for testing purposes")
	cmd.Flags().MarkHidden("tamperable")

	cmd.MarkFlagRequired("apikey")

	return cmd
}
