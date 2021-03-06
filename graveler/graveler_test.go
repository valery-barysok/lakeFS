package graveler_test

import (
	"bytes"
	"context"
	"errors"
	"testing"

	"github.com/go-test/deep"
	"github.com/treeverse/lakefs/graveler"
	"github.com/treeverse/lakefs/graveler/ref"
	"github.com/treeverse/lakefs/graveler/testutil"
	tu "github.com/treeverse/lakefs/testutil"
)

func TestGraveler_List(t *testing.T) {
	conn, _ := tu.GetDB(t, databaseURI)
	branchLocker := ref.NewBranchLocker(conn)
	tests := []struct {
		name        string
		r           *graveler.Graveler
		expectedErr error
		expected    []*graveler.ValueRecord
	}{
		{
			name: "one committed one staged no paths",
			r: graveler.NewGraveler(branchLocker, &testutil.CommittedFake{ValueIterator: testutil.NewValueIteratorFake([]graveler.ValueRecord{{Key: graveler.Key("foo"), Value: &graveler.Value{}}})},
				&testutil.StagingFake{ValueIterator: testutil.NewValueIteratorFake([]graveler.ValueRecord{{Key: graveler.Key("bar"), Value: &graveler.Value{}}})},
				&testutil.RefsFake{RefType: graveler.ReferenceTypeBranch, Commits: map[graveler.CommitID]*graveler.Commit{"": {}}},
			),
			expected: []*graveler.ValueRecord{{Key: graveler.Key("bar"), Value: &graveler.Value{}}, {Key: graveler.Key("foo"), Value: &graveler.Value{}}},
		},
		{
			name: "same path different file",
			r: graveler.NewGraveler(branchLocker, &testutil.CommittedFake{ValueIterator: testutil.NewValueIteratorFake([]graveler.ValueRecord{{Key: graveler.Key("foo"), Value: &graveler.Value{Identity: []byte("original")}}})},
				&testutil.StagingFake{ValueIterator: testutil.NewValueIteratorFake([]graveler.ValueRecord{{Key: graveler.Key("foo"), Value: &graveler.Value{Identity: []byte("other")}}})},
				&testutil.RefsFake{RefType: graveler.ReferenceTypeBranch, Commits: map[graveler.CommitID]*graveler.Commit{"": {}}},
			),
			expected: []*graveler.ValueRecord{{Key: graveler.Key("foo"), Value: &graveler.Value{Identity: []byte("other")}}},
		},
		{
			name: "one committed one staged no paths - with prefix",
			r: graveler.NewGraveler(branchLocker, &testutil.CommittedFake{ValueIterator: testutil.NewValueIteratorFake([]graveler.ValueRecord{{Key: graveler.Key("prefix/foo"), Value: &graveler.Value{}}})},
				&testutil.StagingFake{ValueIterator: testutil.NewValueIteratorFake([]graveler.ValueRecord{{Key: graveler.Key("prefix/bar"), Value: &graveler.Value{}}})},
				&testutil.RefsFake{RefType: graveler.ReferenceTypeBranch, Commits: map[graveler.CommitID]*graveler.Commit{"": {}}},
			),
			expected: []*graveler.ValueRecord{{Key: graveler.Key("prefix/bar"), Value: &graveler.Value{}}, {Key: graveler.Key("prefix/foo"), Value: &graveler.Value{}}},
		},
	}
	ctx := context.Background()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			listing, err := tt.r.List(ctx, "", "")
			if !errors.Is(err, tt.expectedErr) {
				t.Fatalf("wrong error, expected:%s got:%s", tt.expectedErr, err)
			}
			if err != nil {
				return // err == tt.expectedErr
			}
			defer listing.Close()
			var recs []*graveler.ValueRecord
			for listing.Next() {
				recs = append(recs, listing.Value())
			}
			if diff := deep.Equal(recs, tt.expected); diff != nil {
				t.Fatal("List() diff found", diff)
			}
		})
	}
}

func TestGraveler_Get(t *testing.T) {
	conn, _ := tu.GetDB(t, databaseURI)
	branchLocker := ref.NewBranchLocker(conn)
	errTest := errors.New("some kind of err")
	tests := []struct {
		name                string
		r                   *graveler.Graveler
		expectedValueResult graveler.Value
		expectedErr         error
	}{
		{
			name: "commit - exists",
			r: graveler.NewGraveler(branchLocker, &testutil.CommittedFake{ValuesByKey: map[string]*graveler.Value{"key": {Identity: []byte("committed")}}}, nil,
				&testutil.RefsFake{RefType: graveler.ReferenceTypeCommit, Commits: map[graveler.CommitID]*graveler.Commit{"": {}}},
			),
			expectedValueResult: graveler.Value{Identity: []byte("committed")},
		},
		{
			name: "commit - not found",
			r: graveler.NewGraveler(branchLocker, &testutil.CommittedFake{Err: graveler.ErrNotFound}, nil,
				&testutil.RefsFake{RefType: graveler.ReferenceTypeCommit, Commits: map[graveler.CommitID]*graveler.Commit{"": {}}},
			), expectedErr: graveler.ErrNotFound,
		},
		{
			name: "commit - error",
			r: graveler.NewGraveler(branchLocker, &testutil.CommittedFake{Err: errTest}, nil,
				&testutil.RefsFake{RefType: graveler.ReferenceTypeCommit, Commits: map[graveler.CommitID]*graveler.Commit{"": {}}},
			), expectedErr: errTest,
		},
		{
			name: "branch - only staged",
			r: graveler.NewGraveler(branchLocker, &testutil.CommittedFake{Err: graveler.ErrNotFound}, &testutil.StagingFake{Value: &graveler.Value{Identity: []byte("staged")}},
				&testutil.RefsFake{RefType: graveler.ReferenceTypeBranch, Commits: map[graveler.CommitID]*graveler.Commit{"": {}}},
			),
			expectedValueResult: graveler.Value{Identity: []byte("staged")},
		},
		{
			name: "branch - committed and staged",
			r: graveler.NewGraveler(branchLocker, &testutil.CommittedFake{ValuesByKey: map[string]*graveler.Value{"key": {Identity: []byte("committed")}}}, &testutil.StagingFake{Value: &graveler.Value{Identity: []byte("staged")}},
				&testutil.RefsFake{RefType: graveler.ReferenceTypeBranch, Commits: map[graveler.CommitID]*graveler.Commit{"": {}}},
			),
			expectedValueResult: graveler.Value{Identity: []byte("staged")},
		},
		{
			name: "branch - only committed",
			r: graveler.NewGraveler(branchLocker, &testutil.CommittedFake{ValuesByKey: map[string]*graveler.Value{"key": {Identity: []byte("committed")}}}, &testutil.StagingFake{Err: graveler.ErrNotFound},
				&testutil.RefsFake{RefType: graveler.ReferenceTypeBranch, Commits: map[graveler.CommitID]*graveler.Commit{"": {}}},
			),
			expectedValueResult: graveler.Value{Identity: []byte("committed")},
		},
		{
			name: "branch - tombstone",
			r: graveler.NewGraveler(branchLocker, &testutil.CommittedFake{ValuesByKey: map[string]*graveler.Value{"key": {Identity: []byte("committed")}}}, &testutil.StagingFake{Value: nil},
				&testutil.RefsFake{RefType: graveler.ReferenceTypeBranch, Commits: map[graveler.CommitID]*graveler.Commit{"": {}}},
			),
			expectedErr: graveler.ErrNotFound,
		},
		{
			name: "branch - staged return error",
			r: graveler.NewGraveler(branchLocker, &testutil.CommittedFake{}, &testutil.StagingFake{Err: errTest},
				&testutil.RefsFake{RefType: graveler.ReferenceTypeBranch, Commits: map[graveler.CommitID]*graveler.Commit{"": {}}},
			),
			expectedErr: errTest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Value, err := tt.r.Get(context.Background(), "", "", []byte("key"))
			if err != tt.expectedErr {
				t.Fatalf("wrong error, expected:%s got:%s", tt.expectedErr, err)
			}
			if err != nil {
				return // err == tt.expected error
			}
			if string(tt.expectedValueResult.Identity) != string(Value.Identity) {
				t.Errorf("wrong Value address, expected:%s got:%s", tt.expectedValueResult.Identity, Value.Identity)
			}
		})
	}
}

func TestGraveler_DiffUncommitted(t *testing.T) {
	conn, _ := tu.GetDB(t, databaseURI)
	branchLocker := ref.NewBranchLocker(conn)
	tests := []struct {
		name            string
		r               *graveler.Graveler
		amount          int
		expectedErr     error
		expectedHasMore bool
		expectedDiff    graveler.DiffIterator
	}{
		{
			name: "no changes",
			r: graveler.NewGraveler(branchLocker, &testutil.CommittedFake{ValueIterator: testutil.NewValueIteratorFake([]graveler.ValueRecord{{Key: graveler.Key("foo/one"), Value: &graveler.Value{}}}), Err: graveler.ErrNotFound},
				&testutil.StagingFake{ValueIterator: testutil.NewValueIteratorFake([]graveler.ValueRecord{})},
				&testutil.RefsFake{Branch: &graveler.Branch{CommitID: "c1"}, Commits: map[graveler.CommitID]*graveler.Commit{"c1": {MetaRangeID: "mri1"}}},
			),
			amount:       10,
			expectedDiff: testutil.NewDiffIter([]graveler.Diff{}),
		},
		{
			name: "added one",
			r: graveler.NewGraveler(branchLocker, &testutil.CommittedFake{ValueIterator: testutil.NewValueIteratorFake([]graveler.ValueRecord{}), Err: graveler.ErrNotFound},
				&testutil.StagingFake{ValueIterator: testutil.NewValueIteratorFake([]graveler.ValueRecord{{Key: graveler.Key("foo/one"), Value: &graveler.Value{}}})},
				&testutil.RefsFake{Branch: &graveler.Branch{CommitID: "c1"}, Commits: map[graveler.CommitID]*graveler.Commit{"c1": {MetaRangeID: "mri1"}}},
			),
			amount: 10,
			expectedDiff: testutil.NewDiffIter([]graveler.Diff{{
				Key:   graveler.Key("foo/one"),
				Type:  graveler.DiffTypeAdded,
				Value: &graveler.Value{},
			}}),
		},
		{
			name: "changed one",
			r: graveler.NewGraveler(branchLocker, &testutil.CommittedFake{ValueIterator: testutil.NewValueIteratorFake([]graveler.ValueRecord{{Key: graveler.Key("foo/one"), Value: &graveler.Value{}}})},
				&testutil.StagingFake{ValueIterator: testutil.NewValueIteratorFake([]graveler.ValueRecord{{Key: graveler.Key("foo/one"), Value: &graveler.Value{}}})},
				&testutil.RefsFake{Branch: &graveler.Branch{CommitID: "c1"}, Commits: map[graveler.CommitID]*graveler.Commit{"c1": {MetaRangeID: "mri1"}}},
			),
			amount: 10,
			expectedDiff: testutil.NewDiffIter([]graveler.Diff{{
				Key:   graveler.Key("foo/one"),
				Type:  graveler.DiffTypeChanged,
				Value: &graveler.Value{},
			}}),
		},
		{
			name: "removed one",
			r: graveler.NewGraveler(branchLocker, &testutil.CommittedFake{ValueIterator: testutil.NewValueIteratorFake([]graveler.ValueRecord{{Key: graveler.Key("foo/one"), Value: &graveler.Value{}}})},
				&testutil.StagingFake{ValueIterator: testutil.NewValueIteratorFake([]graveler.ValueRecord{{Key: graveler.Key("foo/one"), Value: nil}})},
				&testutil.RefsFake{Branch: &graveler.Branch{CommitID: "c1"}, Commits: map[graveler.CommitID]*graveler.Commit{"c1": {MetaRangeID: "mri1"}}},
			),
			amount: 10,
			expectedDiff: testutil.NewDiffIter([]graveler.Diff{{
				Key:  graveler.Key("foo/one"),
				Type: graveler.DiffTypeRemoved,
			}}),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			diff, err := tt.r.DiffUncommitted(ctx, "repo", "branch")
			if err != tt.expectedErr {
				t.Fatalf("wrong error, expected:%s got:%s", tt.expectedErr, err)
			}
			if err != nil {
				return // err == tt.expectedErr
			}

			// compare iterators
			for diff.Next() {
				if !tt.expectedDiff.Next() {
					t.Fatalf("listing next returned true where expected listing next returned false")
				}
				if diff := deep.Equal(diff.Value(), tt.expectedDiff.Value()); diff != nil {
					t.Errorf("unexpected diff %s", diff)
				}
			}
			if tt.expectedDiff.Next() {
				t.Fatalf("expected listing next returned true where listing next returned false")
			}
		})
	}
}

func TestGraveler_CreateBranch(t *testing.T) {
	conn, _ := tu.GetDB(t, databaseURI)
	branchLocker := ref.NewBranchLocker(conn)
	gravel := graveler.NewGraveler(branchLocker, nil,
		nil,
		&testutil.RefsFake{
			Err:      graveler.ErrNotFound,
			CommitID: "8888888798e3aeface8e62d1c7072a965314b4",
		},
	)
	_, err := gravel.CreateBranch(context.Background(), "", "", "")
	if err != nil {
		t.Fatal("unexpected error on create branch", err)
	}
	// test create branch when branch exists
	gravel = graveler.NewGraveler(branchLocker, nil,
		nil,
		&testutil.RefsFake{
			Branch: &graveler.Branch{},
		},
	)
	_, err = gravel.CreateBranch(context.Background(), "", "", "")
	if !errors.Is(err, graveler.ErrBranchExists) {
		t.Fatal("did not get expected error, expected ErrBranchExists")
	}
}

func TestGraveler_UpdateBranch(t *testing.T) {
	conn, _ := tu.GetDB(t, databaseURI)
	branchLocker := ref.NewBranchLocker(conn)
	gravel := graveler.NewGraveler(branchLocker, nil,
		&testutil.StagingFake{ValueIterator: testutil.NewValueIteratorFake([]graveler.ValueRecord{{Key: graveler.Key("foo/one"), Value: &graveler.Value{}}})},
		&testutil.RefsFake{Branch: &graveler.Branch{}},
	)
	_, err := gravel.UpdateBranch(context.Background(), "", "", "")
	if !errors.Is(err, graveler.ErrConflictFound) {
		t.Fatal("expected update to fail on conflict")
	}
	gravel = graveler.NewGraveler(branchLocker, nil,
		&testutil.StagingFake{ValueIterator: testutil.NewValueIteratorFake([]graveler.ValueRecord{})},
		&testutil.RefsFake{Branch: &graveler.Branch{}},
	)
	_, err = gravel.UpdateBranch(context.Background(), "", "", "")
	if err != nil {
		t.Fatal("did not expect to get error")
	}
}

func TestGraveler_Commit(t *testing.T) {
	conn, _ := tu.GetDB(t, databaseURI)
	branchLocker := ref.NewBranchLocker(conn)
	expectedCommitID := graveler.CommitID("expectedCommitId")
	expectedRangeID := graveler.MetaRangeID("expectedRangeID")
	values := testutil.NewValueIteratorFake([]graveler.ValueRecord{{Key: nil, Value: nil}})
	type fields struct {
		CommittedManager *testutil.CommittedFake
		StagingManager   *testutil.StagingFake
		RefManager       *testutil.RefsFake
	}
	type args struct {
		ctx          context.Context
		repositoryID graveler.RepositoryID
		branchID     graveler.BranchID
		committer    string
		message      string
		metadata     graveler.Metadata
	}
	tests := []struct {
		name        string
		fields      fields
		args        args
		want        graveler.CommitID
		expectedErr error
	}{
		{
			name: "valid commit",
			fields: fields{
				CommittedManager: &testutil.CommittedFake{MetaRangeID: expectedRangeID},
				StagingManager:   &testutil.StagingFake{ValueIterator: values},
				RefManager: &testutil.RefsFake{CommitID: expectedCommitID,
					Branch:  &graveler.Branch{CommitID: expectedCommitID},
					Commits: map[graveler.CommitID]*graveler.Commit{expectedCommitID: {MetaRangeID: expectedRangeID}}},
			},
			args: args{
				ctx:          nil,
				repositoryID: "repo",
				branchID:     "branch",
				committer:    "committer",
				message:      "a message",
				metadata:     graveler.Metadata{},
			},
			want:        expectedCommitID,
			expectedErr: nil,
		},
		{
			name: "fail on staging",
			fields: fields{
				CommittedManager: &testutil.CommittedFake{MetaRangeID: expectedRangeID},
				StagingManager:   &testutil.StagingFake{ValueIterator: values, Err: graveler.ErrNotFound},
				RefManager: &testutil.RefsFake{CommitID: expectedCommitID,
					Branch:  &graveler.Branch{CommitID: expectedCommitID},
					Commits: map[graveler.CommitID]*graveler.Commit{expectedCommitID: {MetaRangeID: expectedRangeID}}},
			},
			args: args{
				ctx:          nil,
				repositoryID: "repo",
				branchID:     "branch",
				committer:    "committer",
				message:      "a message",
				metadata:     nil,
			},
			want:        expectedCommitID,
			expectedErr: graveler.ErrNotFound,
		},
		{
			name: "fail on apply",
			fields: fields{
				CommittedManager: &testutil.CommittedFake{MetaRangeID: expectedRangeID, Err: graveler.ErrConflictFound},
				StagingManager:   &testutil.StagingFake{ValueIterator: values},
				RefManager: &testutil.RefsFake{CommitID: expectedCommitID,
					Branch:  &graveler.Branch{CommitID: expectedCommitID},
					Commits: map[graveler.CommitID]*graveler.Commit{expectedCommitID: {MetaRangeID: expectedRangeID}}},
			},
			args: args{
				ctx:          nil,
				repositoryID: "repo",
				branchID:     "branch",
				committer:    "committer",
				message:      "a message",
				metadata:     nil,
			},
			want:        expectedCommitID,
			expectedErr: graveler.ErrConflictFound,
		},
		{
			name: "fail on add commit",
			fields: fields{
				CommittedManager: &testutil.CommittedFake{MetaRangeID: expectedRangeID},
				StagingManager:   &testutil.StagingFake{ValueIterator: values},
				RefManager: &testutil.RefsFake{CommitID: expectedCommitID,
					Branch:    &graveler.Branch{CommitID: expectedCommitID},
					CommitErr: graveler.ErrConflictFound,
					Commits:   map[graveler.CommitID]*graveler.Commit{expectedCommitID: {MetaRangeID: expectedRangeID}}},
			},
			args: args{
				ctx:          nil,
				repositoryID: "repo",
				branchID:     "branch",
				committer:    "committer",
				message:      "a message",
				metadata:     nil,
			},
			want:        expectedCommitID,
			expectedErr: graveler.ErrConflictFound,
		},
		{
			name: "fail on drop",
			fields: fields{
				CommittedManager: &testutil.CommittedFake{MetaRangeID: expectedRangeID},
				StagingManager:   &testutil.StagingFake{ValueIterator: values, DropErr: graveler.ErrNotFound},
				RefManager: &testutil.RefsFake{CommitID: expectedCommitID,
					Branch:  &graveler.Branch{CommitID: expectedCommitID},
					Commits: map[graveler.CommitID]*graveler.Commit{expectedCommitID: {MetaRangeID: expectedRangeID}}},
			},
			args: args{
				ctx:          nil,
				repositoryID: "repo",
				branchID:     "branch",
				committer:    "committer",
				message:      "a message",
				metadata:     graveler.Metadata{},
			},
			want:        expectedCommitID,
			expectedErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expectedCommitID := graveler.CommitID("expectedCommitId")
			expectedRangeID := graveler.MetaRangeID("expectedRangeID")
			values := testutil.NewValueIteratorFake([]graveler.ValueRecord{{Key: nil, Value: nil}})
			g := graveler.NewGraveler(branchLocker, tt.fields.CommittedManager, tt.fields.StagingManager, tt.fields.RefManager)

			got, err := g.Commit(context.Background(), "", "", tt.args.committer, tt.args.message, tt.args.metadata)
			if !errors.Is(err, tt.expectedErr) {
				t.Fatalf("unexpected err got = %v, wanted = %v", err, tt.expectedErr)
			}
			if err != nil {
				return
			}
			if diff := deep.Equal(tt.fields.CommittedManager.AppliedData, testutil.AppliedData{
				Values:      values,
				MetaRangeID: expectedRangeID,
			}); diff != nil {
				t.Errorf("unexpected apply data %s", diff)
			}

			if diff := deep.Equal(tt.fields.RefManager.AddedCommit, testutil.AddedCommitData{
				Committer:   tt.args.committer,
				Message:     tt.args.message,
				MetaRangeID: expectedRangeID,
				Parents:     graveler.CommitParents{expectedCommitID},
				Metadata:    graveler.Metadata{},
			}); diff != nil {
				t.Errorf("unexpected added commit %s", diff)
			}
			if !tt.fields.StagingManager.DropCalled && tt.fields.StagingManager.DropErr == nil {
				t.Errorf("expected drop to be called")
			}

			if got != expectedCommitID {
				t.Errorf("got wrong commitID, got = %v, want %v", got, expectedCommitID)
			}
		})
	}
}

func TestGraveler_CommitExistingRange(t *testing.T) {
	conn, _ := tu.GetDB(t, databaseURI)
	branchLocker := ref.NewBranchLocker(conn)
	const (
		expectedCommitID       = graveler.CommitID("expectedCommitId")
		expectedParentCommitID = graveler.CommitID("expectedParentCommitId")
		expectedRangeID        = graveler.MetaRangeID("expectedRangeID")
		expectedRepositoryID   = graveler.RepositoryID("expectedRangeID")
	)

	type fields struct {
		CommittedManager *testutil.CommittedFake
		StagingManager   *testutil.StagingFake
		RefManager       *testutil.RefsFake
	}
	type args struct {
		ctx          context.Context
		repositoryID graveler.RepositoryID
		branchID     graveler.BranchID
		committer    string
		message      string
		metadata     graveler.Metadata
	}
	tests := []struct {
		name        string
		fields      fields
		args        args
		want        graveler.CommitID
		expectedErr error
	}{
		{
			name: "valid commit",
			fields: fields{
				CommittedManager: &testutil.CommittedFake{MetaRangeID: expectedRangeID},
				StagingManager:   &testutil.StagingFake{ValueIterator: testutil.NewValueIteratorFake(nil)},
				RefManager: &testutil.RefsFake{CommitID: expectedCommitID, Branch: &graveler.Branch{CommitID: expectedCommitID},
					Commits: map[graveler.CommitID]*graveler.Commit{
						expectedParentCommitID: {},
					},
				}},
			args: args{
				ctx:          nil,
				repositoryID: "repo",
				branchID:     "branch",
				committer:    "committer",
				message:      "a message",
				metadata:     graveler.Metadata{},
			},
			want:        expectedCommitID,
			expectedErr: nil,
		},
		{
			name: "meta range not found",
			fields: fields{
				CommittedManager: &testutil.CommittedFake{Err: graveler.ErrMetaRangeNotFound},
				StagingManager:   &testutil.StagingFake{ValueIterator: testutil.NewValueIteratorFake(nil)},
				RefManager: &testutil.RefsFake{CommitID: expectedCommitID,
					Branch:  &graveler.Branch{CommitID: expectedCommitID},
					Commits: map[graveler.CommitID]*graveler.Commit{expectedParentCommitID: {MetaRangeID: expectedRangeID}}},
			},
			args: args{
				ctx:          nil,
				repositoryID: "repo",
				branchID:     "branch",
				committer:    "committer",
				message:      "a message",
				metadata:     nil,
			},
			want:        expectedCommitID,
			expectedErr: graveler.ErrMetaRangeNotFound,
		},
		{
			name: "fail on add commit",
			fields: fields{
				CommittedManager: &testutil.CommittedFake{MetaRangeID: expectedRangeID},
				StagingManager:   &testutil.StagingFake{ValueIterator: testutil.NewValueIteratorFake(nil)},
				RefManager: &testutil.RefsFake{CommitID: expectedCommitID,
					Branch:    &graveler.Branch{CommitID: expectedCommitID},
					CommitErr: graveler.ErrConflictFound,
					Commits:   map[graveler.CommitID]*graveler.Commit{expectedParentCommitID: {MetaRangeID: expectedRangeID}}},
			},
			args: args{
				ctx:          nil,
				repositoryID: "repo",
				branchID:     "branch",
				committer:    "committer",
				message:      "a message",
				metadata:     nil,
			},
			want:        expectedCommitID,
			expectedErr: graveler.ErrConflictFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := graveler.NewGraveler(branchLocker, tt.fields.CommittedManager, tt.fields.StagingManager, tt.fields.RefManager)
			got, err := g.CommitExistingMetaRange(context.Background(), expectedRepositoryID, expectedParentCommitID, expectedRangeID, tt.args.committer, tt.args.message, tt.args.metadata)
			if !errors.Is(err, tt.expectedErr) {
				t.Fatalf("unexpected err got = %v, wanted = %v", err, tt.expectedErr)
			}
			if err != nil {
				return
			}

			if diff := deep.Equal(tt.fields.RefManager.AddedCommit, testutil.AddedCommitData{
				Committer:   tt.args.committer,
				Message:     tt.args.message,
				MetaRangeID: expectedRangeID,
				Parents:     graveler.CommitParents{expectedParentCommitID},
				Metadata:    graveler.Metadata{},
			}); diff != nil {
				t.Errorf("unexpected added commit %s", diff)
			}
			if tt.fields.StagingManager.DropCalled {
				t.Error("Staging manager drop shouldn't be called")
			}

			if got != expectedCommitID {
				t.Errorf("got wrong commitID, got = %v, want %v", got, expectedCommitID)
			}
		})
	}
}

func TestGraveler_Delete(t *testing.T) {
	conn, _ := tu.GetDB(t, databaseURI)
	branchLocker := ref.NewBranchLocker(conn)
	type fields struct {
		CommittedManager graveler.CommittedManager
		StagingManager   *testutil.StagingFake
		RefManager       graveler.RefManager
	}
	type args struct {
		repositoryID graveler.RepositoryID
		branchID     graveler.BranchID
		key          graveler.Key
	}
	tests := []struct {
		name               string
		fields             fields
		args               args
		expectedSetValue   *graveler.ValueRecord
		expectedRemovedKey graveler.Key
		expectedErr        error
	}{
		{
			name: "exists only in committed",
			fields: fields{
				CommittedManager: &testutil.CommittedFake{
					ValuesByKey: map[string]*graveler.Value{"key": {}},
				},
				StagingManager: &testutil.StagingFake{
					Err: graveler.ErrNotFound,
				},
				RefManager: &testutil.RefsFake{
					Branch:  &graveler.Branch{CommitID: "c1"},
					Commits: map[graveler.CommitID]*graveler.Commit{"c1": {}}},
			},
			args: args{
				key: []byte("key"),
			},
			expectedSetValue: &graveler.ValueRecord{
				Key:   []byte("key"),
				Value: nil,
			},
			expectedErr: nil,
		},
		{
			name: "exists in committed and in staging",
			fields: fields{
				CommittedManager: &testutil.CommittedFake{
					ValuesByKey: map[string]*graveler.Value{"key": {}},
				},
				StagingManager: &testutil.StagingFake{
					Value: &graveler.Value{},
				},
				RefManager: &testutil.RefsFake{
					Branch:  &graveler.Branch{CommitID: "c1"},
					Commits: map[graveler.CommitID]*graveler.Commit{"c1": {}},
				},
			},
			args: args{
				key: []byte("key"),
			},
			expectedSetValue: &graveler.ValueRecord{
				Key:   []byte("key"),
				Value: nil,
			},
			expectedErr: nil,
		},
		{
			name: "exists in committed tombstone in staging",
			fields: fields{
				CommittedManager: &testutil.CommittedFake{
					ValuesByKey: map[string]*graveler.Value{"key": {}},
				},
				StagingManager: &testutil.StagingFake{
					Value: nil,
				},
				RefManager: &testutil.RefsFake{
					Branch:  &graveler.Branch{CommitID: "c1"},
					Commits: map[graveler.CommitID]*graveler.Commit{"c1": {}},
				},
			},
			args:        args{},
			expectedErr: graveler.ErrNotFound,
		},
		{
			name: "exists only in staging - commits",
			fields: fields{
				CommittedManager: &testutil.CommittedFake{
					Err: graveler.ErrNotFound,
				},
				StagingManager: &testutil.StagingFake{
					Value: &graveler.Value{},
				},
				RefManager: &testutil.RefsFake{
					Branch:  &graveler.Branch{CommitID: "c1"},
					Commits: map[graveler.CommitID]*graveler.Commit{"c1": {}},
				},
			},
			args: args{
				key: []byte("key"),
			},
			expectedRemovedKey: []byte("key"),
			expectedErr:        nil,
		},
		{
			name: "exists only in staging - no commits",
			fields: fields{
				CommittedManager: &testutil.CommittedFake{},
				StagingManager: &testutil.StagingFake{
					Value: &graveler.Value{},
				},
				RefManager: &testutil.RefsFake{
					Branch:  &graveler.Branch{},
					Commits: map[graveler.CommitID]*graveler.Commit{},
				},
			},
			args: args{
				key: []byte("key"),
			},
			expectedRemovedKey: []byte("key"),
			expectedErr:        nil,
		},
		{
			name: "not in committed not in staging",
			fields: fields{
				CommittedManager: &testutil.CommittedFake{
					Err: graveler.ErrNotFound,
				},
				StagingManager: &testutil.StagingFake{
					Err: graveler.ErrNotFound,
				},
				RefManager: &testutil.RefsFake{
					Branch:  &graveler.Branch{},
					Commits: map[graveler.CommitID]*graveler.Commit{"": {}},
				},
			},
			args:        args{},
			expectedErr: graveler.ErrNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			g := graveler.NewGraveler(branchLocker, tt.fields.CommittedManager, tt.fields.StagingManager, tt.fields.RefManager)
			if err := g.Delete(ctx, tt.args.repositoryID, tt.args.branchID, tt.args.key); !errors.Is(err, tt.expectedErr) {
				t.Errorf("Delete() returned unexpected error. got = %v, expected %v", err, tt.expectedErr)
			}
			// validate set on staging
			if diff := deep.Equal(tt.fields.StagingManager.LastSetValueRecord, tt.expectedSetValue); diff != nil {
				t.Errorf("unexpected set value %s", diff)
			}
			// validate removed from staging
			if !bytes.Equal(tt.fields.StagingManager.LastRemovedKey, tt.expectedRemovedKey) {
				t.Errorf("unexpected removed key got = %s, expected = %s ", tt.fields.StagingManager.LastRemovedKey, tt.expectedRemovedKey)
			}
		})
	}
}
