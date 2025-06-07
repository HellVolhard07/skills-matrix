# Base image
FROM python:3.9-slim-buster

# Set working directory
WORKDIR /app

# Install system dependencies required for gRPC (for compiling)
RUN apt-get update && apt-get install -y --no-install-recommends \
    gcc \
    libffi-dev \
    python3-dev \
    build-essential \
    && rm -rf /var/lib/apt/lists/*

# Install Python dependencies early to leverage Docker cache
COPY requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt

# Remove unneeded build tools to reduce final image size
RUN apt-get purge -y gcc build-essential libffi-dev python3-dev && \
    apt-get autoremove -y && \
    rm -rf /var/lib/apt/lists/*

# Copy application code
COPY . .

# Expose gRPC port
EXPOSE 50052

# Run gRPC server
CMD ["python", "-m", "user_service.server"]
