#!/bin/bash

# Navigate to the project directory
cd /home/ec2-user/github/simplebank

# Check for updates in the Git repository
git fetch origin main

# Check if there are differences between the local files and the remote repository
if git diff --quiet HEAD origin/main; then
  echo "No changes in the repository. Skipping rebuild and cleanup."
else
  echo "Changes detected. Pulling latest changes."

  # Pull the latest changes
  git pull origin main

  # Stop the existing Docker Compose service
  sudo systemctl stop simplebank.service

  # Remove all stopped containers
  docker container prune -f

  # Remove all unused images, not just dangling ones
  docker image prune -a -f

  # Rebuild and restart the Docker containers
  # The systemd service will handle this
  sudo systemctl start simplebank.service
fi
