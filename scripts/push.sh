#!/bin/bash

# Navigate to your git repository directory
# cd /path/to/your/repo

# Add all changes to the staging area
git add .

# Ask for a commit message
echo "Enter commit message: "
read commitMessage

# Commit the changes with the user's message
git commit -m "$commitMessage"

# Push the changes to the remote repository
git push origin main 