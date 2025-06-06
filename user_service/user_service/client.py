import grpc
import logging

# Ensure the generated _pb2.py and _pb2_grpc.py are in the same directory
# or correctly imported based on your project structure.
# If you followed the previous steps, these relative imports should work.
from . import user_service_pb2
from . import user_service_pb2_grpc

# Configure logging for the client
logging.basicConfig(level=logging.INFO, format='%(asctime)s - %(levelname)s - %(message)s')

def run():
    # Connect to the gRPC server
    # The server is running on localhost:50052
    with grpc.insecure_channel('localhost:50052') as channel:
        stub = user_service_pb2_grpc.UserServiceStub(channel)
        logging.info("Connected to User Service at localhost:50052")

        # --- Test 1: Get an existing user profile ---
        user_id_to_get = "user123"
        logging.info(f"\n--- Attempting to get profile for ID: {user_id_to_get} ---")
        try:
            get_request = user_service_pb2.GetUserProfileRequest(id=user_id_to_get)
            user_profile = stub.GetUserProfile(get_request)
            logging.info(f"Successfully retrieved profile: {user_profile.name} (ID: {user_profile.id})")
            logging.info(f"Bio: {user_profile.bio}")
        except grpc.RpcError as e:
            logging.error(f"Error getting profile {user_id_to_get}: {e.code().name} - {e.details()}")

        # --- Test 2: Create a new user profile ---
        new_user_id = "newuser_alpha"
        logging.info(f"\n--- Attempting to create new profile for ID: {new_user_id} ---")
        try:
            create_request = user_service_pb2.CreateUserProfileRequest(
                id=new_user_id,
                name="Alpha Dev",
                bio="Passionate about learning new technologies.",
                avatar_url="http://example.com/avatars/alpha.jpg"
            )
            created_profile = stub.CreateUserProfile(create_request)
            logging.info(f"Successfully created profile: {created_profile.name} (ID: {created_profile.id})")
        except grpc.RpcError as e:
            logging.error(f"Error creating profile {new_user_id}: {e.code().name} - {e.details()}")

        # --- Test 3: Attempt to create an existing user profile (should fail) ---
        logging.info(f"\n--- Attempting to create existing profile for ID: {user_id_to_get} (should fail) ---")
        try:
            create_existing_request = user_service_pb2.CreateUserProfileRequest(
                id=user_id_to_get, # This user already exists
                name="Another Alice",
                bio="Attempt to recreate Alice.",
                avatar_url="http://example.com/avatars/another_alice.jpg"
            )
            stub.CreateUserProfile(create_existing_request)
            logging.warning("Unexpected: Created existing profile successfully.")
        except grpc.RpcError as e:
            if e.code() == grpc.StatusCode.ALREADY_EXISTS:
                logging.info(f"Correctly failed to create existing profile: {e.details()}")
            else:
                logging.error(f"Error creating existing profile: {e.code().name} - {e.details()}")

        # --- Test 4: Update an existing user profile ---
        user_id_to_update = "user456" # Bob The Builder
        logging.info(f"\n--- Attempting to update profile for ID: {user_id_to_update} ---")
        try:
            update_request = user_service_pb2.UpdateUserProfileRequest(
                id=user_id_to_update,
                bio="Updated: DevOps Specialist, now focusing on serverless and AI ops.",
                avatar_url="http://example.com/avatars/bob_new.jpg"
            )
            updated_profile = stub.UpdateUserProfile(update_request)
            logging.info(f"Successfully updated profile for {updated_profile.name} (ID: {updated_profile.id})")
            logging.info(f"New Bio: {updated_profile.bio}")
        except grpc.RpcError as e:
            logging.error(f"Error updating profile {user_id_to_update}: {e.code().name} - {e.details()}")

        # --- Test 5: Get the newly updated profile to verify ---
        logging.info(f"\n--- Verifying updated profile for ID: {user_id_to_update} ---")
        try:
            get_request_updated = user_service_pb2.GetUserProfileRequest(id=user_id_to_update)
            verified_profile = stub.GetUserProfile(get_request_updated)
            logging.info(f"Verified profile Bio: {verified_profile.bio}")
        except grpc.RpcError as e:
            logging.error(f"Error verifying updated profile {user_id_to_update}: {e.code().name} - {e.details()}")

        # --- Test 6: Attempt to get a non-existent user profile (should fail) ---
        non_existent_id = "nonexistent_user"
        logging.info(f"\n--- Attempting to get profile for ID: {non_existent_id} (should fail) ---")
        try:
            get_non_existent_request = user_service_pb2.GetUserProfileRequest(id=non_existent_id)
            stub.GetUserProfile(get_non_existent_request)
            logging.warning(f"Unexpected: Retrieved non-existent profile {non_existent_id}.")
        except grpc.RpcError as e:
            if e.code() == grpc.StatusCode.NOT_FOUND:
                logging.info(f"Correctly failed to get non-existent profile: {e.details()}")
            else:
                logging.error(f"Error getting non-existent profile: {e.code().name} - {e.details()}")

if __name__ == '__main__':
    run()