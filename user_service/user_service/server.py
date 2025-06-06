import grpc
from concurrent import futures
import logging

# Configure logging for the Server
logging.basicConfig(level=logging.INFO, format='%(asctime)s - %(levelname)s - %(message)s')

# Import the generated gRPC stub that allows adding the servicer
from . import user_service_pb2_grpc

# Import the handler logic
from .handlers import UserServiceServicer

def serve():
    # Create a gRPC server
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))

    # Add the UserServiceServicer to the server
    # This connects the gRPC methods defined in UserServiceServicer
    # to the incoming gRPC calls.
    user_service_pb2_grpc.add_UserServiceServicer_to_server(UserServiceServicer(), server)

    # Define the server address and port
    server_address = '[::]:50052' # User Service port
    server.add_insecure_port(server_address)

    logging.info(f"User Service starting on {server_address}")
    server.start() # Start the server

    try:
        # Keep the server running until termination signal
        server.wait_for_termination()
    except KeyboardInterrupt:
        logging.info("User Service shutting down.")
        server.stop(0) # Stop the server gracefully

if __name__ == '__main__':
    serve()