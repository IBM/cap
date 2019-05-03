/*
Copyright 2018 The cap Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package commands

import (
	"context"
	"os"
	"runtime"

	"github.com/IBM/cap/go/log"
	"github.com/sirupsen/logrus"
	"golang.org/x/sys/unix"
)

var handledSignals = []os.Signal{
	unix.SIGTERM,
	unix.SIGINT,
	unix.SIGUSR1,
	unix.SIGPIPE,
}

func handleSignals(ctx context.Context, signals chan os.Signal, serverCH chan *Server) chan struct{} {
	done := make(chan struct{}, 1)
	go func() {
		var server *Server
		for {
			select {
			case s := <-serverCH:
				server = s
			case s := <-signals:
				log.G(ctx).WithField("signal", s).Debug("received signal")
				switch s {
				case unix.SIGUSR1:
					dumpStacks()
				case unix.SIGPIPE:
					continue
				default:
					if server == nil {
						close(done)
						return
					}
					// TODO Stop() server method
					close(done)
				}
			}
		}
	}()
	return done
}

// dumpStacks - logs the stack dump of all goroutines
func dumpStacks() {
	buf := make([]byte, 32768)
	for {
		n := runtime.Stack(buf, true)
		if n < len(buf) {
			logrus.Infof("=== BEGIN stack dump of all goroutines ===\n%s\n=== END stack dump ===", buf[:n])
			break
		}
		buf = make([]byte, 2*len(buf))
	}
}
