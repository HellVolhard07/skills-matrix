import grpc
import logging

# Configure logging for handlers
logging.basicConfig(level=logging.INFO, format='%(asctime)s - %(levelname)s - %(message)s')

# Import generated protobuf files (relative import within the package)
from . import user_service_pb2
from . import user_service_pb2_grpc

# Import the mock database
from .mock_db import get_initial_user_profiles

class UserServiceServicer(user_service_pb2_grpc.UserServiceServicer):
    """
    Implements the gRPC methods for the User Service.
    Uses stale, in-memory data provided by mock_db.py.
    """
    def __init__(self):
        # Initialize profiles from the mock database
        self.profiles = get_initial_user_profiles(user_service_pb2)
        logging.info("UserServiceServicer initialized with profiles from mock_db.")

    def CreateUserProfile(self, request, context):
        """
        Creates a new user profile.
        Since we are using stale data, this will simulate adding a new profile
        but it won't persist across service restarts.
        """
        logging.info(f"Handler: Received CreateUserProfile request for user ID: {request.id}")
        if request.id in self.profiles:
            context.set_code(grpc.StatusCode.ALREADY_EXISTS)
            context.set_details(f"Profile for user ID {request.id} already exists.")
            logging.warning(f"Handler: Profile for user ID {request.id} already exists. Returning ALREADY_EXISTS.")
            return user_service_pb2.UserProfile() # Return empty profile

        user_profile = user_service_pb2.UserProfile(
            id=request.id,
            name=request.name,
            bio=request.bio,
            avatar_url=request.avatar_url
        )
        self.profiles[request.id] = user_profile
        logging.info(f"Handler: User profile created for ID: {request.id}, Name: {request.name}")
        return user_profile

    def UpdateUserProfile(self, request, context):
        """
        Updates an existing user profile in the stale data.
        """
        logging.info(f"Handler: Received UpdateUserProfile request for user ID: {request.id}")
        if request.id not in self.profiles:
            context.set_code(grpc.StatusCode.NOT_FOUND)
            context.set_details(f"Profile for user ID {request.id} not found.")
            logging.warning(f"Handler: Profile for user ID {request.id} not found. Returning NOT_FOUND.")
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

        logging.info(f"Handler: User profile updated for ID: {request.id}")
        return current_profile

    def GetUserProfile(self, request, context):
        """
        Retrieves a user profile by ID from the stale data.
        """
        logging.info(f"Handler: Received GetUserProfile request for user ID: {request.id}")
        user_profile = self.profiles.get(request.id)
        if user_profile:
            logging.info(f"Handler: User profile found for ID: {request.id}")
            return user_profile
        else:
            context.set_code(grpc.StatusCode.NOT_FOUND)
            context.set_details(f"Profile for user ID {request.id} not found.")
            logging.warning(f"Handler: Profile for user ID {request.id} not found. Returning NOT_FOUND.")
            return user_service_pb2.UserProfile() # Return empty profile if not found