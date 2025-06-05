# user_service/mock_db.py

import logging

# Configure logging for mock_db
logging.basicConfig(level=logging.INFO, format='%(asctime)s - %(levelname)s - %(message)s')

def get_initial_user_profiles(user_profile_pb2):
    """
    Returns a dictionary of initial, stale user profiles.
    Takes user_profile_pb2 as an argument to create UserProfile objects.
    """
    logging.info("Initializing mock user profiles.")
    return {
        "user123": user_profile_pb2.UserProfile(
            id="user123",
            name="Alice Wonderland",
            bio="Software Engineer with a passion for distributed systems.",
            avatar_url="http://example.com/avatars/alice.jpg"
        ),
        "user456": user_profile_pb2.UserProfile(
            id="user456",
            name="Bob The Builder",
            bio="DevOps Specialist focused on cloud infrastructure.",
            avatar_url="http://example.com/avatars/bob.jpg"
        ),
        "user789": user_profile_pb2.UserProfile(
            id="user789",
            name="Charlie Chaplin",
            bio="Frontend Developer with expertise in React and Vue.",
            avatar_url="http://example.com/avatars/charlie.jpg"
        )
    }