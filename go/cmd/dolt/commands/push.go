// Copyright 2019 Dolthub, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package commands

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/dustin/go-humanize"

	"github.com/dolthub/dolt/go/cmd/dolt/cli"
	"github.com/dolthub/dolt/go/cmd/dolt/errhand"
	eventsapi "github.com/dolthub/dolt/go/gen/proto/dolt/services/eventsapi/v1alpha1"
	"github.com/dolthub/dolt/go/libraries/doltcore/doltdb"
	"github.com/dolthub/dolt/go/libraries/doltcore/env"
	"github.com/dolthub/dolt/go/libraries/doltcore/env/actions"
	"github.com/dolthub/dolt/go/libraries/utils/argparser"
	"github.com/dolthub/dolt/go/libraries/utils/filesys"
	"github.com/dolthub/dolt/go/store/datas"
)

var pushDocs = cli.CommandDocumentationContent{
	ShortDesc: "Update remote refs along with associated objects",
	LongDesc: `Updates remote refs using local refs, while sending objects necessary to complete the given refs.

When the command line does not specify where to push with the {{.LessThan}}remote{{.GreaterThan}} argument, an attempt is made to infer the remote.  If only one remote exists it will be used, if multiple remotes exists, a remote named 'origin' will be attempted.  If there is more than one remote, and none of them are named 'origin' then the command will fail and you will need to specify the correct remote explicitly.

When the command line does not specify what to push with {{.LessThan}}refspec{{.GreaterThan}}... then the current branch will be used.

When neither the command-line does not specify what to push, the default behavior is used, which corresponds to the current branch being pushed to the corresponding upstream branch, but as a safety measure, the push is aborted if the upstream branch does not have the same name as the local one.
`,

	Synopsis: []string{
		"[-u | --set-upstream] [{{.LessThan}}remote{{.GreaterThan}}] [{{.LessThan}}refspec{{.GreaterThan}}]",
	},
}

type PushCmd struct{}

// Name is returns the name of the Dolt cli command. This is what is used on the command line to invoke the command
func (cmd PushCmd) Name() string {
	return "push"
}

// Description returns a description of the command
func (cmd PushCmd) Description() string {
	return "Push to a dolt remote."
}

// CreateMarkdown creates a markdown file containing the helptext for the command at the given path
func (cmd PushCmd) CreateMarkdown(fs filesys.Filesys, path, commandStr string) error {
	ap := cmd.createArgParser()
	return CreateMarkdown(fs, path, cli.GetCommandDocumentation(commandStr, pushDocs, ap))
}

func (cmd PushCmd) createArgParser() *argparser.ArgParser {
	ap := argparser.NewArgParser()
	ap.SupportsFlag(env.SetUpstreamFlag, "u", "For every branch that is up to date or successfully pushed, add upstream (tracking) reference, used by argument-less {{.EmphasisLeft}}dolt pull{{.EmphasisRight}} and other commands.")
	ap.SupportsFlag(env.ForcePushFlag, "f", "Update the remote with local history, overwriting any conflicting history in the remote.")
	return ap
}

// EventType returns the type of the event to log
func (cmd PushCmd) EventType() eventsapi.ClientEventType {
	return eventsapi.ClientEventType_PUSH
}

// Exec executes the command
func (cmd PushCmd) Exec(ctx context.Context, commandStr string, args []string, dEnv *env.DoltEnv) int {
	ap := cmd.createArgParser()
	help, usage := cli.HelpAndUsagePrinters(cli.GetCommandDocumentation(commandStr, pushDocs, ap))
	apr := cli.ParseArgsOrDie(ap, args, help)

	opts, err := env.ParsePushArgs(ctx, apr, dEnv)

	//TODO build verbose error
	//errhand.BuildDError("error: failed to read remotes from config."
	if err != nil {
		return HandleVErrAndExitCode(errhand.VerboseErrorFromError(err), usage)
	}

	err = actions.DoPush(ctx, dEnv, opts, runProgFuncs, stopProgFuncs)

	switch err {
	case doltdb.ErrUpToDate:
		cli.Println("Everything up-to-date")
	default:
	}
	return HandleVErrAndExitCode(errhand.VerboseErrorFromError(err), usage)
}

const minUpdate = 100 * time.Millisecond

var spinnerSeq = []rune{'|', '/', '-', '\\'}

type TextSpinner struct {
	seqPos     int
	lastUpdate time.Time
}

func (ts *TextSpinner) next() string {
	now := time.Now()
	if now.Sub(ts.lastUpdate) > minUpdate {
		ts.seqPos = (ts.seqPos + 1) % len(spinnerSeq)
		ts.lastUpdate = now
	}

	return string([]rune{spinnerSeq[ts.seqPos]})
}

func pullerProgFunc(pullerEventCh chan datas.PullerEvent) {
	var pos int
	var currentTreeLevel int
	var percentBuffered float64
	var tableFilesBuffered int
	var filesUploaded int
	var ts TextSpinner

	uploadRate := ""

	for evt := range pullerEventCh {
		switch evt.EventType {
		case datas.NewLevelTWEvent:
			if evt.TWEventDetails.TreeLevel != 1 {
				currentTreeLevel = evt.TWEventDetails.TreeLevel
				percentBuffered = 0
			}
		case datas.DestDBHasTWEvent:
			if evt.TWEventDetails.TreeLevel != -1 {
				currentTreeLevel = evt.TWEventDetails.TreeLevel
			}

		case datas.LevelUpdateTWEvent:
			if evt.TWEventDetails.TreeLevel != -1 {
				currentTreeLevel = evt.TWEventDetails.TreeLevel
				toBuffer := evt.TWEventDetails.ChunksInLevel - evt.TWEventDetails.ChunksAlreadyHad

				if toBuffer > 0 {
					percentBuffered = 100 * float64(evt.TWEventDetails.ChunksBuffered) / float64(toBuffer)
				}
			}

		case datas.LevelDoneTWEvent:

		case datas.TableFileClosedEvent:
			tableFilesBuffered += 1

		case datas.StartUploadTableFileEvent:

		case datas.UploadTableFileUpdateEvent:
			bps := float64(evt.TFEventDetails.Stats.Read) / evt.TFEventDetails.Stats.Elapsed.Seconds()
			uploadRate = humanize.Bytes(uint64(bps)) + "/s"

		case datas.EndUploadTableFileEvent:
			filesUploaded += 1
		}

		if currentTreeLevel == -1 {
			continue
		}

		var msg string
		if len(uploadRate) > 0 {
			msg = fmt.Sprintf("%s Tree Level: %d, Percent Buffered: %.2f%%, Files Written: %d, Files Uploaded: %d, Current Upload Speed: %s", ts.next(), currentTreeLevel, percentBuffered, tableFilesBuffered, filesUploaded, uploadRate)
		} else {
			msg = fmt.Sprintf("%s Tree Level: %d, Percent Buffered: %.2f%% Files Written: %d, Files Uploaded: %d", ts.next(), currentTreeLevel, percentBuffered, tableFilesBuffered, filesUploaded)
		}

		pos = cli.DeleteAndPrint(pos, msg)
	}
}

func progFunc(progChan chan datas.PullProgress) {
	var latest datas.PullProgress
	last := time.Now().UnixNano() - 1
	lenPrinted := 0
	done := false
	for !done {
		select {
		case progress, ok := <-progChan:
			if !ok {
				done = true
			}
			latest = progress

		case <-time.After(250 * time.Millisecond):
			break
		}

		nowUnix := time.Now().UnixNano()
		deltaTime := time.Duration(nowUnix - last)
		halfSec := 500 * time.Millisecond
		if done || deltaTime > halfSec {
			last = nowUnix
			if latest.KnownCount > 0 {
				progMsg := fmt.Sprintf("Counted chunks: %d, Buffered chunks: %d)", latest.KnownCount, latest.DoneCount)
				lenPrinted = cli.DeleteAndPrint(lenPrinted, progMsg)
			}
		}
	}

	if lenPrinted > 0 {
		cli.Println()
	}
}

func runProgFuncs() (*sync.WaitGroup, chan datas.PullProgress, chan datas.PullerEvent) {
	pullerEventCh := make(chan datas.PullerEvent, 128)
	progChan := make(chan datas.PullProgress, 128)
	wg := &sync.WaitGroup{}

	wg.Add(1)
	go func() {
		defer wg.Done()
		progFunc(progChan)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		pullerProgFunc(pullerEventCh)
	}()

	return wg, progChan, pullerEventCh
}

func stopProgFuncs(wg *sync.WaitGroup, progChan chan datas.PullProgress, pullerEventCh chan datas.PullerEvent) {
	close(progChan)
	close(pullerEventCh)
	wg.Wait()

	cli.Println()
}

func bytesPerSec(bytes uint64, start time.Time) string {
	bps := float64(bytes) / float64(time.Since(start).Seconds())
	return humanize.Bytes(uint64(bps))
}
