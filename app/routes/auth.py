from flask import request, jsonify
from .. import app
from ..auth import register_user, authenticate_user, logout_user, delete_user, get_user_from_session, is_authenticated

@app.route('/auth/register', methods=['POST'])
def register():
    """
    Register a new user
    ---
    tags:
      - Authentication
    parameters:
      - name: body
        in: body
        required: true
        schema:
          type: object
          properties:
            username:
              type: string
            password:
              type: string
    responses:
      200:
        description: User registered successfully
      400:
        description: Registration failed
    """
    data = request.json
    username = data.get('username')
    password = data.get('password')
    
    if not username or not password:
        return jsonify({"error": "Username and password are required"}), 400
    
    success, message = register_user(username, password)
    
    if success:
        return jsonify({"message": message}), 200
    else:
        return jsonify({"error": message}), 400

@app.route('/auth/login', methods=['POST'])
def login():
    """
    Login a user
    ---
    tags:
      - Authentication
    parameters:
      - name: body
        in: body
        required: true
        schema:
          type: object
          properties:
            username:
              type: string
            password:
              type: string
    responses:
      200:
        description: Login successful
      401:
        description: Login failed
    """
    data = request.json
    username = data.get('username')
    password = data.get('password')
    
    if not username or not password:
        return jsonify({"error": "Username and password are required"}), 400
    
    success, result = authenticate_user(username, password)
    
    if success:
        return jsonify({
            "message": "Login successful",
            "token": result,
            "username": username
        }), 200
    else:
        return jsonify({"error": result}), 401

@app.route('/auth/logout', methods=['POST'])
def logout():
    """
    Logout a user
    ---
    tags:
      - Authentication
    parameters:
      - name: Authorization
        in: header
        type: string
        required: true
        description: Bearer token
    responses:
      200:
        description: Logout successful
    """
    auth_header = request.headers.get('Authorization')
    
    if not auth_header or not auth_header.startswith('Bearer '):
        return jsonify({"message": "Logged out"}), 200
    
    session_id = auth_header.split('Bearer ')[1]
    logout_user(session_id)
    
    return jsonify({"message": "Logged out successfully"}), 200

@app.route('/auth/user', methods=['GET'])
def get_user():
    """
    Get current user information
    ---
    tags:
      - Authentication
    parameters:
      - name: Authorization
        in: header
        type: string
        required: true
        description: Bearer token
    responses:
      200:
        description: User information
      401:
        description: Not authenticated
    """
    auth_header = request.headers.get('Authorization')
    
    if not auth_header or not auth_header.startswith('Bearer '):
        return jsonify({"error": "Not authenticated"}), 401
    
    session_id = auth_header.split('Bearer ')[1]
    user = get_user_from_session(session_id)
    
    if not user:
        return jsonify({"error": "Not authenticated"}), 401
    
    return jsonify({
        "username": user["username"],
        "created_at": user["created_at"]
    }), 200

@app.route('/auth/delete-account', methods=['DELETE'])
def delete_account():
    """
    Delete user account
    ---
    tags:
      - Authentication
    parameters:
      - name: Authorization
        in: header
        type: string
        required: true
        description: Bearer token
    responses:
      200:
        description: Account deleted successfully
      401:
        description: Not authenticated
    """
    auth_header = request.headers.get('Authorization')
    
    if not auth_header or not auth_header.startswith('Bearer '):
        return jsonify({"error": "Not authenticated"}), 401
    
    session_id = auth_header.split('Bearer ')[1]
    user = get_user_from_session(session_id)
    
    if not user:
        return jsonify({"error": "Not authenticated"}), 401
    
    username = user["username"]
    if delete_user(username):
        return jsonify({"message": "Account deleted successfully"}), 200
    else:
        return jsonify({"error": "Failed to delete account"}), 400

@app.route('/auth/check', methods=['GET'])
def check_auth():
    """
    Check if user is authenticated
    ---
    tags:
      - Authentication
    parameters:
      - name: Authorization
        in: header
        type: string
        required: true
        description: Bearer token
    responses:
      200:
        description: Authentication status
    """
    auth_header = request.headers.get('Authorization')
    
    if not auth_header or not auth_header.startswith('Bearer '):
        return jsonify({"authenticated": False}), 200
    
    session_id = auth_header.split('Bearer ')[1]
    authenticated = is_authenticated(session_id)
    
    return jsonify({"authenticated": authenticated}), 200
