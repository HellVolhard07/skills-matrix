# Base image
FROM python:3.9-alpine

# Set working directory
WORKDIR /app

# (Optional) Copy requirements.txt and install dependencies
COPY requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt

# Copy application code
COPY . .

# Expose gRPC port
EXPOSE 50052

# Run gRPC server
CMD ["python", "-m", "user_service.server"]
