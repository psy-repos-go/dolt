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

package actions

import (
	"context"
	"errors"
	"github.com/dolthub/dolt/go/cmd/dolt/cli"
	"github.com/dolthub/dolt/go/cmd/dolt/errhand"
	eventsapi "github.com/dolthub/dolt/go/gen/proto/dolt/services/eventsapi/v1alpha1"
	"github.com/dolthub/dolt/go/libraries/doltcore/doltdb"
	"github.com/dolthub/dolt/go/libraries/doltcore/env"
	"github.com/dolthub/dolt/go/libraries/doltcore/ref"
	"github.com/dolthub/dolt/go/libraries/doltcore/remotestorage"
	"github.com/dolthub/dolt/go/libraries/events"
	"github.com/dolthub/dolt/go/libraries/utils/earl"
	"github.com/dolthub/dolt/go/store/datas"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"sync"
)

var ErrCantFF = errors.New("can't fast forward merge")

// Push will update a destination branch, in a given destination database if it can be done as a fast forward merge.
// This is accomplished first by verifying that the remote tracking reference for the source database can be updated to
// the given commit via a fast forward merge.  If this is the case, an attempt will be made to update the branch in the
// destination db to the given commit via fast forward move.  If that succeeds the tracking branch is updated in the
// source db.
func Push(ctx context.Context, dEnv *env.DoltEnv, mode ref.UpdateMode, destRef ref.BranchRef, remoteRef ref.RemoteRef, srcDB, destDB *doltdb.DoltDB, commit *doltdb.Commit, progChan chan datas.PullProgress, pullerEventCh chan datas.PullerEvent) error {
	var err error
	if mode == ref.FastForwardOnly {
		canFF, err := srcDB.CanFastForward(ctx, remoteRef, commit)

		if err != nil {
			return err
		} else if !canFF {
			return ErrCantFF
		}
	}

	rf, err := commit.GetStRef()

	if err != nil {
		return err
	}

	err = destDB.PushChunks(ctx, dEnv.TempTableFilesDir(), srcDB, rf, progChan, pullerEventCh)

	if err != nil {
		return err
	}

	switch mode {
	case ref.ForceUpdate:
		err = destDB.SetHeadToCommit(ctx, destRef, commit)
		if err != nil {
			return err
		}
		err = srcDB.SetHeadToCommit(ctx, remoteRef, commit)
	case ref.FastForwardOnly:
		err = destDB.FastForward(ctx, destRef, commit)
		if err != nil {
			return err
		}
		err = srcDB.FastForward(ctx, remoteRef, commit)
	}

	return err
}

func DoPush(ctx context.Context, dEnv *env.DoltEnv, opts *env.PushOpts, progStarter func () (*sync.WaitGroup, chan datas.PullProgress, chan datas.PullerEvent), progStopper func (wg *sync.WaitGroup, progChan chan datas.PullProgress, pullerEventCh chan datas.PullerEvent)) (verr errhand.VerboseError) {
	destDB, err := opts.Remote.GetRemoteDB(ctx, dEnv.DoltDB.ValueReadWriter().Format())

	if err != nil {
		bdr := errhand.BuildDError("error: failed to get remote db").AddCause(err)

		if err == remotestorage.ErrInvalidDoltSpecPath {
			urlObj, _ := earl.Parse(opts.Remote.Url)
			bdr.AddDetails("For the remote: %s %s", opts.Remote.Name, opts.Remote.Url)

			path := urlObj.Path
			if path[0] == '/' {
				path = path[1:]
			}

			bdr.AddDetails("'%s' should be in the format 'organization/repo'", path)
		}

		return bdr.Build()
	}

	switch opts.SrcRef.GetType() {
	case ref.BranchRefType:
		if opts.SrcRef == ref.EmptyBranchRef {
			verr = deleteRemoteBranch(ctx, opts.DestRef, opts.RemoteRef, dEnv.DoltDB, destDB, opts.Remote)
		} else {
			verr = PushToRemoteBranch(ctx, dEnv, opts.Mode, opts.SrcRef, opts.DestRef, opts.RemoteRef, dEnv.DoltDB, destDB, opts.Remote, progStarter, progStopper)
		}
	case ref.TagRefType:
		verr = pushTagToRemote(ctx, dEnv, opts.SrcRef, opts.DestRef, dEnv.DoltDB, destDB, progStarter, progStopper)
	default:
		verr = errhand.BuildDError("cannot push ref %s of type %s", opts.SrcRef.String(), opts.SrcRef.GetType()).Build()
	}

	if verr != nil {
		return verr
	}

	if opts.SetUpstream {
		dEnv.RepoState.Branches[opts.SrcRef.GetPath()] = env.BranchConfig{
			Merge: ref.MarshalableRef{
				Ref: opts.DestRef,
			},
			Remote: opts.Remote.Name,
		}

		err := dEnv.RepoState.Save(dEnv.FS)

		if err != nil {
			verr = errhand.BuildDError("error: failed to save repo state").AddCause(err).Build()
		}
	}

	return verr
}

// PushTag pushes a commit tag and all underlying data from a local source database to a remote destination database.
func PushTag(ctx context.Context, dEnv *env.DoltEnv, destRef ref.TagRef, srcDB, destDB *doltdb.DoltDB, tag *doltdb.Tag, progChan chan datas.PullProgress, pullerEventCh chan datas.PullerEvent) error {
	var err error

	rf, err := tag.GetStRef()

	if err != nil {
		return err
	}

	err = destDB.PushChunks(ctx, dEnv.TempTableFilesDir(), srcDB, rf, progChan, pullerEventCh)

	if err != nil {
		return err
	}

	return destDB.SetHead(ctx, destRef, rf)
}


func deleteRemoteBranch(ctx context.Context, toDelete, remoteRef ref.DoltRef, localDB, remoteDB *doltdb.DoltDB, remote env.Remote) errhand.VerboseError {
	err := DeleteRemoteBranch(ctx, toDelete.(ref.BranchRef), remoteRef.(ref.RemoteRef), localDB, remoteDB)

	if err != nil {
		return errhand.BuildDError("error: failed to delete '%s' from remote '%s'", toDelete.String(), remote.Name).Build()
	}

	return nil
}

func PushToRemoteBranch(ctx context.Context, dEnv *env.DoltEnv, mode ref.UpdateMode, srcRef, destRef, remoteRef ref.DoltRef, localDB, remoteDB *doltdb.DoltDB, remote env.Remote, progStarter func () (*sync.WaitGroup, chan datas.PullProgress, chan datas.PullerEvent), progStopper func (wg *sync.WaitGroup, progChan chan datas.PullProgress, pullerEventCh chan datas.PullerEvent)) errhand.VerboseError {
	evt := events.GetEventFromContext(ctx)

	u, err := earl.Parse(remote.Url)

	if err == nil {
		if u.Scheme != "" {
			evt.SetAttribute(eventsapi.AttributeID_REMOTE_URL_SCHEME, u.Scheme)
		}
	}

	cs, _ := doltdb.NewCommitSpec(srcRef.GetPath())
	cm, err := localDB.Resolve(ctx, cs, dEnv.RepoStateReader().CWBHeadRef())

	if err != nil {
		return errhand.BuildDError("error: refspec '%v' not found.", srcRef.GetPath()).Build()
	} else {
		wg, progChan, pullerEventCh := progStarter()
		err = Push(ctx, dEnv, mode, destRef.(ref.BranchRef), remoteRef.(ref.RemoteRef), localDB, remoteDB, cm, progChan, pullerEventCh)
		progStopper(wg, progChan, pullerEventCh)

		if err != nil {
			if err == doltdb.ErrUpToDate {
				cli.Println("Everything up-to-date")
			} else if err == doltdb.ErrIsAhead || err == ErrCantFF || err == datas.ErrMergeNeeded {
				cli.Printf("To %s\n", remote.Url)
				cli.Printf("! [rejected]          %s -> %s (non-fast-forward)\n", destRef.String(), remoteRef.String())
				cli.Printf("error: failed to push some refs to '%s'\n", remote.Url)
				cli.Println("hint: Updates were rejected because the tip of your current branch is behind")
				cli.Println("hint: its remote counterpart. Integrate the remote changes (e.g.")
				cli.Println("hint: 'dolt pull ...') before pushing again.")
				return errhand.BuildDError("").Build()
			} else {
				status, ok := status.FromError(err)
				if ok && status.Code() == codes.PermissionDenied {
					cli.Println("hint: have you logged into DoltHub using 'dolt login'?")
					cli.Println("hint: check that user.email in 'dolt config --list' has write perms to DoltHub repo")
				}
				if rpcErr, ok := err.(*remotestorage.RpcError); ok {
					return errhand.BuildDError("error: push failed").AddCause(err).AddDetails(rpcErr.FullDetails()).Build()
				} else {
					return errhand.BuildDError("error: push failed").AddCause(err).Build()
				}
			}
		}
	}

	return nil
}

func pushTagToRemote(ctx context.Context, dEnv *env.DoltEnv, srcRef, destRef ref.DoltRef, localDB, remoteDB *doltdb.DoltDB, progStarter func () (*sync.WaitGroup, chan datas.PullProgress, chan datas.PullerEvent), progStopper func (wg *sync.WaitGroup, progChan chan datas.PullProgress, pullerEventCh chan datas.PullerEvent)) errhand.VerboseError {
	tg, err := localDB.ResolveTag(ctx, srcRef.(ref.TagRef))

	if err != nil {
		return errhand.VerboseErrorFromError(err)
	}

	wg, progChan, pullerEventCh := progStarter()

	err = PushTag(ctx, dEnv, destRef.(ref.TagRef), localDB, remoteDB, tg, progChan, pullerEventCh)

	progStopper(wg, progChan, pullerEventCh)

	if err != nil {
		if err == doltdb.ErrUpToDate {
			cli.Println("Everything up-to-date")
		} else {
			return errhand.BuildDError("error: push failed").AddCause(err).Build()
		}
	}

	return nil
}

// DeleteRemoteBranch validates targetRef is a branch on the remote database, and then deletes it, then deletes the
// remote tracking branch from the local database.
func DeleteRemoteBranch(ctx context.Context, targetRef ref.BranchRef, remoteRef ref.RemoteRef, localDB, remoteDB *doltdb.DoltDB) error {
	hasRef, err := remoteDB.HasRef(ctx, targetRef)

	if err != nil {
		return err
	}

	if hasRef {
		err = remoteDB.DeleteBranch(ctx, targetRef)
	}

	if err != nil {
		return err
	}

	err = localDB.DeleteBranch(ctx, remoteRef)

	if err != nil {
		return err
	}

	return nil
}

// FetchCommit takes a fetches a commit and all underlying data from a remote source database to the local destination database.
func FetchCommit(ctx context.Context, dEnv *env.DoltEnv, srcDB, destDB *doltdb.DoltDB, srcDBCommit *doltdb.Commit, progChan chan datas.PullProgress, pullerEventCh chan datas.PullerEvent) error {
	stRef, err := srcDBCommit.GetStRef()

	if err != nil {
		return err
	}

	return destDB.PullChunks(ctx, dEnv.TempTableFilesDir(), srcDB, stRef, progChan, pullerEventCh)
}

// FetchCommit takes a fetches a commit tag and all underlying data from a remote source database to the local destination database.
func FetchTag(ctx context.Context, dEnv *env.DoltEnv, srcDB, destDB *doltdb.DoltDB, srcDBTag *doltdb.Tag, progChan chan datas.PullProgress, pullerEventCh chan datas.PullerEvent) error {
	stRef, err := srcDBTag.GetStRef()

	if err != nil {
		return err
	}

	return destDB.PullChunks(ctx, dEnv.TempTableFilesDir(), srcDB, stRef, progChan, pullerEventCh)
}

// Clone pulls all data from a remote source database to a local destination database.
func Clone(ctx context.Context, srcDB, destDB *doltdb.DoltDB, eventCh chan<- datas.TableFileEvent) error {
	return srcDB.Clone(ctx, destDB, eventCh)
}
