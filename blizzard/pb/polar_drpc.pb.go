// Code generated by protoc-gen-go-drpc. DO NOT EDIT.
// protoc-gen-go-drpc version: v0.0.32
// source: polar.proto

package pb

import (
	context "context"
	errors "errors"
	protojson "google.golang.org/protobuf/encoding/protojson"
	proto "google.golang.org/protobuf/proto"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
	drpc "storj.io/drpc"
	drpcerr "storj.io/drpc/drpcerr"
)

type drpcEncoding_File_polar_proto struct{}

func (drpcEncoding_File_polar_proto) Marshal(msg drpc.Message) ([]byte, error) {
	return proto.Marshal(msg.(proto.Message))
}

func (drpcEncoding_File_polar_proto) MarshalAppend(buf []byte, msg drpc.Message) ([]byte, error) {
	return proto.MarshalOptions{}.MarshalAppend(buf, msg.(proto.Message))
}

func (drpcEncoding_File_polar_proto) Unmarshal(buf []byte, msg drpc.Message) error {
	return proto.Unmarshal(buf, msg.(proto.Message))
}

func (drpcEncoding_File_polar_proto) JSONMarshal(msg drpc.Message) ([]byte, error) {
	return protojson.Marshal(msg.(proto.Message))
}

func (drpcEncoding_File_polar_proto) JSONUnmarshal(buf []byte, msg drpc.Message) error {
	return protojson.Unmarshal(buf, msg.(proto.Message))
}

type DRPCPolarClient interface {
	DRPCConn() drpc.Conn

	Health(ctx context.Context, in *emptypb.Empty) (*PolarHealth, error)
	Judge(ctx context.Context, in *File) (DRPCPolar_JudgeClient, error)
}

type drpcPolarClient struct {
	cc drpc.Conn
}

func NewDRPCPolarClient(cc drpc.Conn) DRPCPolarClient {
	return &drpcPolarClient{cc}
}

func (c *drpcPolarClient) DRPCConn() drpc.Conn { return c.cc }

func (c *drpcPolarClient) Health(ctx context.Context, in *emptypb.Empty) (*PolarHealth, error) {
	out := new(PolarHealth)
	err := c.cc.Invoke(ctx, "/arctic.Polar/Health", drpcEncoding_File_polar_proto{}, in, out)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *drpcPolarClient) Judge(ctx context.Context, in *File) (DRPCPolar_JudgeClient, error) {
	stream, err := c.cc.NewStream(ctx, "/arctic.Polar/Judge", drpcEncoding_File_polar_proto{})
	if err != nil {
		return nil, err
	}
	x := &drpcPolar_JudgeClient{stream}
	if err := x.MsgSend(in, drpcEncoding_File_polar_proto{}); err != nil {
		return nil, err
	}
	if err := x.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type DRPCPolar_JudgeClient interface {
	drpc.Stream
	Recv() (*JudgeResult, error)
}

type drpcPolar_JudgeClient struct {
	drpc.Stream
}

func (x *drpcPolar_JudgeClient) Recv() (*JudgeResult, error) {
	m := new(JudgeResult)
	if err := x.MsgRecv(m, drpcEncoding_File_polar_proto{}); err != nil {
		return nil, err
	}
	return m, nil
}

func (x *drpcPolar_JudgeClient) RecvMsg(m *JudgeResult) error {
	return x.MsgRecv(m, drpcEncoding_File_polar_proto{})
}

type DRPCPolarServer interface {
	Health(context.Context, *emptypb.Empty) (*PolarHealth, error)
	Judge(*File, DRPCPolar_JudgeStream) error
}

type DRPCPolarUnimplementedServer struct{}

func (s *DRPCPolarUnimplementedServer) Health(context.Context, *emptypb.Empty) (*PolarHealth, error) {
	return nil, drpcerr.WithCode(errors.New("Unimplemented"), drpcerr.Unimplemented)
}

func (s *DRPCPolarUnimplementedServer) Judge(*File, DRPCPolar_JudgeStream) error {
	return drpcerr.WithCode(errors.New("Unimplemented"), drpcerr.Unimplemented)
}

type DRPCPolarDescription struct{}

func (DRPCPolarDescription) NumMethods() int { return 2 }

func (DRPCPolarDescription) Method(n int) (string, drpc.Encoding, drpc.Receiver, interface{}, bool) {
	switch n {
	case 0:
		return "/arctic.Polar/Health", drpcEncoding_File_polar_proto{},
			func(srv interface{}, ctx context.Context, in1, in2 interface{}) (drpc.Message, error) {
				return srv.(DRPCPolarServer).
					Health(
						ctx,
						in1.(*emptypb.Empty),
					)
			}, DRPCPolarServer.Health, true
	case 1:
		return "/arctic.Polar/Judge", drpcEncoding_File_polar_proto{},
			func(srv interface{}, ctx context.Context, in1, in2 interface{}) (drpc.Message, error) {
				return nil, srv.(DRPCPolarServer).
					Judge(
						in1.(*File),
						&drpcPolar_JudgeStream{in2.(drpc.Stream)},
					)
			}, DRPCPolarServer.Judge, true
	default:
		return "", nil, nil, nil, false
	}
}

func DRPCRegisterPolar(mux drpc.Mux, impl DRPCPolarServer) error {
	return mux.Register(impl, DRPCPolarDescription{})
}

type DRPCPolar_HealthStream interface {
	drpc.Stream
	SendAndClose(*PolarHealth) error
}

type drpcPolar_HealthStream struct {
	drpc.Stream
}

func (x *drpcPolar_HealthStream) SendAndClose(m *PolarHealth) error {
	if err := x.MsgSend(m, drpcEncoding_File_polar_proto{}); err != nil {
		return err
	}
	return x.CloseSend()
}

type DRPCPolar_JudgeStream interface {
	drpc.Stream
	Send(*JudgeResult) error
}

type drpcPolar_JudgeStream struct {
	drpc.Stream
}

func (x *drpcPolar_JudgeStream) Send(m *JudgeResult) error {
	return x.MsgSend(m, drpcEncoding_File_polar_proto{})
}