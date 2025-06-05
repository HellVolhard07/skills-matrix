import grpc
from concurrent import futures
import logging

from . import user_service_pb2
from . import user_service_pb2_grpc

# Configure logging
logging.basicConfig(level=logging.INFO, format='%(asctime)s - %(levelname)s - %(message)s')

class UserServiceServicer(user_service_pb2_grpc.UserServiceServicer):
    """
    Implements the gRPC methods for the User Service.
    Uses stale, in-memory data instead of a database.
    """
    def __init__(self):
        # Stale, in-memory data for user profiles
        self.profiles = {
            "user123": user_service_pb2.UserProfile(
                id="user123",
                name="Alice Wonderland",
                bio="Software Engineer with a passion for distributed systems.",
                avatar_url="http://example.com/avatars/alice.jpg"
            ),
            "user456": user_service_pb2.UserProfile(
                id="user456",
                name="Bob The Builder",
                bio="DevOps Specialist focused on cloud infrastructure.",
                avatar_url="http://example.com/avatars/bob.jpg"
            ),
            "user789": user_service_pb2.UserProfile(
                id="user789",
                name="Charlie Chaplin",
                bio="Frontend Developer with expertise in React and Vue.",
                avatar_url="http://example.com/avatars/charlie.jpg"
            )
        }
        logging.info("UserServiceServicer initialized with stale profiles.")

    def CreateUserProfile(self, request, context):
        """
        Creates a new user profile.
        Since we are using stale data, this will simulate adding a new profile
        but it won't persist across service restarts.
        """
        logging.info(f"Received CreateUserProfile request for user ID: {request.id}")
        if request.id in self.profiles:
            context.set_code(grpc.StatusCode.ALREADY_EXISTS)
            context.set_details(f"Profile for user ID {request.id} already exists.")
            logging.warning(f"Profile for user ID {request.id} already exists. Returning ALREADY_EXISTS.")
            return user_service_pb2.UserProfile() # Return empty profile

        user_profile = user_service_pb2.UserProfile(
            id=request.id,
            name=request.name,
            bio=request.bio,
            avatar_url=request.avatar_url
        )
        self.profiles[request.id] = user_profile
        logging.info(f"User profile created for ID: {request.id}, Name: {request.name}")
        return user_profile

    def UpdateUserProfile(self, request, context):
        """
        Updates an existing user profile in the stale data.
        """
        logging.info(f"Received UpdateUserProfile request for user ID: {request.id}")
        if request.id not in self.profiles:
            context.set_code(grpc.StatusCode.NOT_FOUND)
            context.set_details(f"Profile for user ID {request.id} not found.")
            logging.warning(f"Profile for user ID {request.id} not found. Returning NOT_FOUND.")
            return user_service_pb2.UserProfile() # Return empty profile

        # Update fields if provided in the request
        current_profile = self.profiles[request.id]
        if request.name:
            current_profile.name = request.name
        if request.bio:
            current_profile.bio = request.bio
        if request.avatar_url:
            current_profile.avatar_url = request.avatar_url
        
        # In-memory update is sufficient for stale data
        # self.profiles[request.id] = current_profile # Not strictly necessary as it's modified in place

        logging.info(f"User profile updated for ID: {request.id}")
        return current_profile

    def GetUserProfile(self, request, context):
        """
        Retrieves a user profile by ID from the stale data.
        """
        logging.info(f"Received GetUserProfile request for user ID: {request.id}")
        user_profile = self.profiles.get(request.id)
        if user_profile:
            logging.info(f"User profile found for ID: {request.id}")
            return user_profile
        else:
            context.set_code(grpc.StatusCode.NOT_FOUND)
            context.set_details(f"Profile for user ID {request.id} not found.")
            logging.warning(f"Profile for user ID {request.id} not found. Returning NOT_FOUND.")
            return user_service_pb2.UserProfile() # Return empty profile if not found

def serve():
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    user_service_pb2_grpc.add_UserServiceServicer_to_server(UserServiceServicer(), server)
    server_address = '[::]:50052' # User Service port
    server.add_insecure_port(server_address)
    logging.info(f"User Service starting on {server_address}")
    server.start()
    try:
        server.wait_for_termination()
    except KeyboardInterrupt:
        logging.info("User Service shutting down.")
        server.stop(0)

if __name__ == '__main__':
    serve()