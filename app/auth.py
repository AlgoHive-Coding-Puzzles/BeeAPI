import uuid
import os
import hashlib
import time
from functools import wraps
from flask import request, abort, session, jsonify

# In-memory user storage
users = {}
sessions = {}

def generate_session_id():
    return str(uuid.uuid4())

def hash_password(password, salt=None):
    if salt is None:
        salt = uuid.uuid4().hex
    hashed = hashlib.sha256((password + salt).encode()).hexdigest()
    return f"{salt}:{hashed}"

def check_password(password, hashed_password):
    salt, hashed = hashed_password.split(":")
    return hashed == hashlib.sha256((password + salt).encode()).hexdigest()

def register_user(username, password):
    if username in users:
        return False, "Username already exists"
    
    if len(password) < 5:
        return False, "Password must be at least 5 characters"
    
    users[username] = {
        "username": username,
        "password": hash_password(password),
        "created_at": time.time()
    }
    return True, "User registered successfully"

def authenticate_user(username, password):
    if username not in users:
        return False, "Invalid username or password"
    
    user = users[username]
    if not check_password(password, user["password"]):
        return False, "Invalid username or password"
    
    session_id = generate_session_id()
    sessions[session_id] = {
        "username": username,
        "created_at": time.time()
    }
    
    return True, session_id

def logout_user(session_id):
    if session_id in sessions:
        del sessions[session_id]
    return True

def delete_user(username):
    if username in users:
        del users[username]
        # Remove any active sessions for this user
        for session_id in list(sessions.keys()):
            if sessions[session_id]["username"] == username:
                del sessions[session_id]
        return True
    return False

def is_authenticated(session_id):
    return session_id in sessions

def get_user_from_session(session_id):
    if session_id in sessions:
        username = sessions[session_id]["username"]
        return users.get(username)
    return None

def login_required(f):
    @wraps(f)
    def decorated_function(*args, **kwargs):
        auth_header = request.headers.get("Authorization")
        
        if not auth_header or not auth_header.startswith("Bearer "):
            return jsonify({"error": "Unauthorized"}), 401
        
        session_id = auth_header.split("Bearer ")[1]
        
        if not is_authenticated(session_id):
            return jsonify({"error": "Unauthorized"}), 401
            
        return f(*args, **kwargs)
    return decorated_function

def create_user_from_env():
    """
    Create users from environment variables.
    Environment variables should be in the format USER_<username>=<password>
    """
    user_count = 0
    for key, value in os.environ.items():
        if key.startswith("USER_"):
            username = key[5:].lower().replace("_", "-")
            password = value
            print(f"Registering user {username} with password {password}")
            register_user(username, password)
            user_count += 1
            
    if user_count == 0:
        print("No users found in environment variables.")
        print("Using a default u:{admin} p:{admin} user.")
        register_user("admin", "admin")    