# user_service/Dockerfile

# Use an official Python runtime as a parent image
FROM python:3.9-slim-buster

# Set the working directory in the container
WORKDIR /app

# Install necessary gRPC dependencies
# grpcio requires some build tools, so install them temporarily
RUN apt-get update && \
    apt-get install -y --no-install-recommends \
    gcc \
    build-essential \
    libffi-dev \
    python3-dev \
    && rm -rf /var/lib/apt/lists/*

# Copy the requirements file and install dependencies
# This step is often done separately to leverage Docker caching
COPY requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt

# Remove build dependencies to make the image smaller
RUN apt-get purge -y gcc build-essential libffi-dev python3-dev && \
    apt-get autoremove -y && \
    rm -rf /var/lib/apt/lists/*

# Copy the entire user_service application directory into the container
# This copies server.py, handlers.py, mock_db.py, and the generated pb2 files
COPY . .

# Expose the port your gRPC server will listen on
EXPOSE 50052

# Command to run the gRPC server when the container starts
# Use python -m to ensure package structure is respected
CMD ["python", "-m", "user_service.server"]