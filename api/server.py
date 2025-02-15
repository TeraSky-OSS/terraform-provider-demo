from flask import Flask, jsonify, request
import os
import json
import uuid

app = Flask(__name__)

# Configure storage directory
STORAGE_DIR = "cars_data"
if not os.path.exists(STORAGE_DIR):
    os.makedirs(STORAGE_DIR)

def get_car_path(car_id):
    return os.path.join(STORAGE_DIR, f"{car_id}.json")

@app.route('/cars', methods=['POST'])
def create_car():
    # Get request data
    data = request.get_json()
    
    # Generate a unique ID for the car
    car_id = str(uuid.uuid4())
    
    # Create car data
    car_data = {
        "id": car_id,
        "model": data.get("model", ""),
        "year": data.get("year", 0)
    }
    
    # Save to filesystem
    with open(get_car_path(car_id), 'w') as f:
        json.dump(car_data, f)
    
    return jsonify(car_data), 201

@app.route('/cars/<car_id>', methods=['GET'])
def get_car(car_id):
    car_path = get_car_path(car_id)
    
    if not os.path.exists(car_path):
        return '', 404
        
    with open(car_path, 'r') as f:
        car_data = json.load(f)
    
    return jsonify(car_data)

@app.route('/cars/<car_id>', methods=['PUT'])
def update_car(car_id):
    car_path = get_car_path(car_id)
    
    if not os.path.exists(car_path):
        return '', 404
    
    # Get request data    
    data = request.get_json()
        
    with open(car_path, 'r') as f:
        car_data = json.load(f)
    
    # Update fields
    car_data["model"] = data.get("model", car_data.get("model", ""))
    car_data["year"] = data.get("year", car_data.get("year", 0))
    
    # Save updated data
    with open(car_path, 'w') as f:
        json.dump(car_data, f)
    
    return jsonify(car_data)

@app.route('/cars/<car_id>', methods=['DELETE'])
def delete_car(car_id):
    car_path = get_car_path(car_id)
    
    if os.path.exists(car_path):
        os.remove(car_path)
    
    return '', 204

if __name__ == '__main__':
    app.run(port=5000)