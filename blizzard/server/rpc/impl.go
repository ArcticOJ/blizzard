package rpc

import (
	"blizzard/blizzard/pb/blizzard"
	"context"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

type Blizzard struct {
}

func (*Blizzard) FetchTestCases(stream blizzard.DRPCBlizzard_FetchTestCasesStream) error {
	// Get the problem id first, and exit on error
	problem, e := stream.Recv()
	if e != nil {
		_ = stream.Close()
		return e
	}
	// Try parsing the problem id, and continue if it is valid
	if p, ok := problem.Id.(*blizzard.Case_Problem); ok && p != nil {
		for {
			testCase, e := stream.Recv()
			if e != nil {
				_ = stream.Close()
				break
			}
			if c, ok := testCase.Id.(*blizzard.Case_Index); ok && c != nil {
				stream.Send(&blizzard.CaseData{Input: []byte("test"), Output: []byte("uwu")})
			} else {
				_ = stream.Close()
				break
			}
		}
	} else {
		_ = stream.Close()
	}
	return nil
}

func (*Blizzard) Alive(context.Context, *emptypb.Empty) (*wrapperspb.BoolValue, error) {
	return wrapperspb.Bool(true), nil
}
