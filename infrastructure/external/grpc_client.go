package external

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type AuthInterceptor struct {
	AccessToken string
}

func NewAuthInterceptor() *AuthInterceptor {
	return &AuthInterceptor{}
}

func (i *AuthInterceptor) Unary() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		// log.Printf("--> unary interceptor: %s", method)
		return invoker(i.attachToken(ctx), method, req, reply, cc, opts...)
	}
}

func (i *AuthInterceptor) Stream() grpc.StreamClientInterceptor {
	return func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		// log.Printf("--> stream interceptor: %s", method)
		return streamer(i.attachToken(ctx), desc, cc, method, opts...)
	}
}

func (i *AuthInterceptor) attachToken(ctx context.Context) context.Context {
	return metadata.AppendToOutgoingContext(ctx, "authorization", i.AccessToken)
}

func (i *AuthInterceptor) SetToken(accessToken string) {
	i.AccessToken = accessToken
}

func (i *AuthInterceptor) TransportOpts() []grpc.DialOption {
	return []grpc.DialOption{
		grpc.WithInsecure(),
		grpc.WithUnaryInterceptor(i.Unary()),
		grpc.WithStreamInterceptor(i.Stream()),
	}
}

func GrpcClient(addr string, opts ...grpc.DialOption) (*grpc.ClientConn, error) {

	conn, err := grpc.Dial(
		addr,
		opts...,
	)
	if err != nil {
		return nil, err
	}

	return conn, nil
}
