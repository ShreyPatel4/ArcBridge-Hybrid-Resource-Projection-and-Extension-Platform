package grpc

// HeartbeatServer is a placeholder for the gRPC service that would stream heartbeats and inventory updates.
// Implementations should translate incoming messages into queue events processed by the REST handlers.
type HeartbeatServer struct{}

// Start launches the mock gRPC server.
func (s *HeartbeatServer) Start() error {
	// TODO: Implement gRPC server using buf-generated protobuf definitions.
	return nil
}
